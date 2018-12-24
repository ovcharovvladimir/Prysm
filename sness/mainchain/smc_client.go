// Package mainchain defines services that interacts with a Geth node via RPC.
// This package is useful for an actor in a sharded system to interact with
// a Sharding Manager Contract.
package mainchain

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/ovcharovvladimir/Prysm/sness/contracts"
	ethereum "github.com/ovcharovvladimir/essentiaHybrid"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/keystore"
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	gethTypes "github.com/ovcharovvladimir/essentiaHybrid/core/types"
	"github.com/ovcharovvladimir/essentiaHybrid/essclient"
	"github.com/ovcharovvladimir/essentiaHybrid/node"
	"github.com/ovcharovvladimir/essentiaHybrid/rpc"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "mainchain")

// ClientIdentifier tells us what client the node we interact with over RPC is running.
const ClientIdentifier = "gess"

// VRCClient defines a struct that interacts with a
// mainchain node via RPC. Specifically, it aids in VRC bindings that are useful
// to other sharding services.
type VRCClient struct {
	endpoint     string             // Endpoint to JSON RPC.
	dataDirPath  string             // Path to the data directory.
	depositFlag  bool               // Keeps track of the deposit option passed in as via CLI flags.
	passwordFile string             // Path to the account password file.
	client       *essclient.Client  // Ethereum RPC client.
	keystore     *keystore.KeyStore // Keystore containing the single signer.
	vrc          *contracts.VRC     // The deployed sharding management contract.
	rpcClient    *rpc.Client        // The RPC client connection to the main geth node.
}

// NewVRCClient constructs a new instance of an VRCClient.
func NewVRCClient(endpoint string, dataDirPath string, depositFlag bool, passwordFile string) (*VRCClient, error) {
	config := &node.Config{
		DataDir: dataDirPath,
	}

	scryptN, scryptP, keydir, err := config.AccountConfig()
	if err != nil {
		return nil, err
	}

	ks := keystore.NewKeyStore(keydir, scryptN, scryptP)

	vrcClient := &VRCClient{
		keystore:     ks,
		endpoint:     endpoint,
		depositFlag:  depositFlag,
		dataDirPath:  dataDirPath,
		passwordFile: passwordFile,
	}

	return vrcClient, nil
}

// Start the VRC Client and connect to running geth node.
func (s *VRCClient) Start() {
	// Sets up a connection to a Geth node via RPC.
	rpcClient, err := dialRPC(s.endpoint)
	if err != nil {
		log.Panicf("Cannot start rpc client: %v", err)
		return
	}

	s.rpcClient = rpcClient
	s.client = essclient.NewClient(rpcClient)

	// Check account existence and unlock account before starting.
	accounts := s.keystore.Accounts()
	if len(accounts) == 0 {
		log.Panic("No accounts found")
		return
	}
	if _, err := s.unlockAccount(accounts[0]); err != nil {
		log.Panicf("Cannot unlock account: %v", err)
		return
	}

	// Initializes bindings to VRC.
	vrc, err := initVRC(s)
	if err != nil {
		log.Panicf("Failed to initialize VRC: %v", err)
		return
	}

	s.vrc = vrc
}

// Stop VRCClient immediately. This cancels any pending RPC connections.
func (s *VRCClient) Stop() error {
	s.rpcClient.Close()
	return nil
}

// CreateTXOpts creates a *TransactOpts with a signer using the default account on the keystore.
func (s *VRCClient) CreateTXOpts(value *big.Int) (*bind.TransactOpts, error) {
	account := s.Account()
	keyJSON, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Info("Error when read utc file")
		return nil, err
	}

	privKey, err := keystore.DecryptKey(keyJSON, "123")
	if err != nil {
		return nil, fmt.Errorf("unable to fetch privKey: %v", err)
	}
	txOps := bind.NewKeyedTransactor(privKey.PrivateKey)
	txOps.Value = value
	log.Info("txOps", "value", value)
	return txOps, nil
	// return &bind.TransactOpts{
	// 	From:  account.Address,
	// 	Value: value,
	// 	Signer: func(signer gethTypes.Signer, addr common.Address, tx *gethTypes.Transaction) (*gethTypes.Transaction, error) {
	// 		networkID, err := s.client.NetworkID(context.Background())
	// 		if err != nil {
	// 			return nil, fmt.Errorf("unable to fetch networkID: %v", err)
	// 		}

	// 		return s.keystore.SignTx(*account, tx, networkID /* chainID */)
	// 	},
	// }, nil
}

// Account to use for sharding transactions.
func (s *VRCClient) Account() *accounts.Account {
	accounts := s.keystore.Accounts()
	return &accounts[0]
}

// ChainReader for interacting with the chain.
func (s *VRCClient) ChainReader() ethereum.ChainReader {
	return ethereum.ChainReader(s.client)
}

func (s *VRCClient) BlockByNumber(ctx context.Context, number *big.Int) (*gethTypes.Block, error) {
	return s.ChainReader().BlockByNumber(ctx, number)
}

func (s *VRCClient) SubscribeNewHead(ctx context.Context, ch chan<- *gethTypes.Header) (ethereum.Subscription, error) {
	return s.ChainReader().SubscribeNewHead(ctx, ch)
}

// VRCCaller to interact with the sharding manager contract.
func (s *VRCClient) VRCCaller() *contracts.VRCCaller {
	if s.vrc == nil {
		return nil
	}
	return &s.vrc.VRCCaller
}

// VRCTransactor allows us to send tx's to the VRC programmatically.
func (s *VRCClient) VRCTransactor() *contracts.VRCTransactor {
	if s.vrc == nil {
		return nil
	}
	return &s.vrc.VRCTransactor
}

// VRCFilterer allows for easy filtering of events from the Sharding Manager Contract.
func (s *VRCClient) VRCFilterer() *contracts.VRCFilterer {
	if s.vrc == nil {
		return nil
	}
	return &s.vrc.VRCFilterer
}

// WaitForTransaction waits for transaction to be mined and returns an error if it takes
// too long.
func (s *VRCClient) WaitForTransaction(ctx context.Context, hash common.Hash, durationInSeconds time.Duration) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, durationInSeconds*time.Second)

	for pending, err := true, error(nil); pending; _, pending, err = s.client.TransactionByHash(ctxTimeout, hash) {
		if err != nil {
			cancel()
			return fmt.Errorf("unable to retrieve transaction: %v", err)
		}
		if ctxTimeout.Err() != nil {
			cancel()
			return fmt.Errorf("transaction timed out, transaction was not able to be mined in the duration: %v", ctxTimeout.Err())
		}
	}
	cancel()
	ctxTimeout.Done()
	log.Infof("Transaction: %s has been mined", hash.Hex())
	return nil
}

// TransactionReceipt allows an VRCClient to retrieve transaction receipts on
// the mainchain by hash.
func (s *VRCClient) TransactionReceipt(hash common.Hash) (*gethTypes.Receipt, error) {

	receipt, err := s.client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return nil, err
	}

	return receipt, err
}

// DepositFlag returns true for cli flag --deposit.
func (s *VRCClient) DepositFlag() bool {
	return s.depositFlag
}

// SetDepositFlag updates the deposit flag property of VRCClient.
func (s *VRCClient) SetDepositFlag(deposit bool) {
	s.depositFlag = deposit
}

// DataDirPath returns the datadir flag as a string.
func (s *VRCClient) DataDirPath() string {
	return s.dataDirPath
}

// unlockAccount will unlock the specified account using utils.PasswordFileFlag
// or empty string if unset.
func (s *VRCClient) unlockAccount(account accounts.Account) (bool, error) {
	pass := ""

	if s.passwordFile != "" {
		file, err := os.Open(s.passwordFile)
		if err != nil {
			return false, fmt.Errorf("unable to open file containing account password %s. %v", s.passwordFile, err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
		if !scanner.Scan() {
			err = scanner.Err()
			if err != nil {
				return false, fmt.Errorf("unable to read contents of file %v", err)
			}
			return false, errors.New("password not found in file")
		}

		pass = scanner.Text()
	}

	return s.keystore.Unlock(account, pass)
}

// Sign signs the hash of collationHeader contents by
// using default account on keystore and returns signed signature.
func (s *VRCClient) Sign(hash common.Hash) ([]byte, error) {
	account := s.Account()
	return s.keystore.SignHash(*account, hash.Bytes())
}

// GetShardCount gets the count of the total shards
// currently operating in the sharded universe.
func (s *VRCClient) GetShardCount() (int64, error) {
	shardCount, err := s.VRCCaller().ShardCount(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}
	return shardCount.Int64(), nil
}
