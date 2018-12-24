package voter

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ovcharovvladimir/Prysm/sness/contracts"
	"github.com/ovcharovvladimir/Prysm/sness/mainchain"
	shardparams "github.com/ovcharovvladimir/Prysm/sness/params"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	gethTypes "github.com/ovcharovvladimir/essentiaHybrid/core/types"
	"github.com/ovcharovvladimir/essentiaHybrid/params"
	"github.com/sirupsen/logrus"
)

// subscribeBlockHeaders checks incoming block headers and determines if
// we are an eligible voter for collations. Then, it finds the pending tx's
// from the running geth node and sorts them by descending order of gas price,
// eliminates those that ask for too much gas, and routes them over
// to the VRC to create a collation.
func subscribeBlockHeaders(reader mainchain.Reader, caller mainchain.ContractCaller, account *accounts.Account) error {
	headerChan := make(chan *gethTypes.Header, 16)

	_, err := reader.SubscribeNewHead(context.Background(), headerChan)
	if err != nil {
		return fmt.Errorf("unable to subscribe to incoming headers. %v", err)
	}

	log.Info("Listening for new headers...")

	for {
		// TODO: Error handling for getting disconnected from the client.
		head := <-headerChan
		// Query the current state to see if we are an eligible voter.
		log.WithFields(logrus.Fields{
			"number": head.Number.String(),
		}).Info("Received new header")

		// Check if we are in the voter pool before checking if we are an eligible voter.
		v, err := isAccountInNotaryPool(caller, account)
		if err != nil {
			return fmt.Errorf("unable to verify client in voter pool. %v", err)
		}

		if v {
			if err := checkVRCForNotary(caller, account); err != nil {
				return fmt.Errorf("unable to watch shards. %v", err)
			}
		}
	}
}

// checkVRCForNotary checks if we are an eligible voter for
// collation for the available shards in the VRC. The function calls
// getEligibleNotary from the VRC and voter a collation if
// conditions are met.
func checkVRCForNotary(caller mainchain.ContractCaller, account *accounts.Account) error {
	log.Info("Checking if we are an eligible collation voter for a shard...")
	shardCount, err := caller.GetShardCount()
	if err != nil {
		return fmt.Errorf("can't get shard count from vrc: %v", err)
	}
	for s := int64(0); s < shardCount; s++ {
		// Checks if we are an eligible voter according to the VRC.
		addr, err := caller.VRCCaller().GetNotaryInCommittee(&bind.CallOpts{}, big.NewInt(s))

		if err != nil {
			return err
		}

		if addr == account.Address {
			log.Infof("Selected as voter on shard: %d", s)
		}

	}

	return nil
}

// getNotaryRegistry retrieves the registry of the registered account.
func getNotaryRegistry(caller mainchain.ContractCaller, account *accounts.Account) (*contracts.Registry, error) {

	var nreg contracts.Registry
	nreg, err := caller.VRCCaller().NotaryRegistry(&bind.CallOpts{}, account.Address)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve voter registry: %v", err)
	}

	return &nreg, nil
}

// isAccountInNotaryPool checks if the user is in the voter pool because
// we can't guarantee our tx for deposit will be in the next block header we receive.
// The function calls IsNotaryDeposited from the VRC and returns true if
// the user is in the voter pool.
func isAccountInNotaryPool(caller mainchain.ContractCaller, account *accounts.Account) (bool, error) {

	nreg, err := getNotaryRegistry(caller, account)
	if err != nil {
		return false, err
	}

	if !nreg.Deposited {
		log.Warnf("Account %s not in voter pool.", account.Address.Hex())
	}

	return nreg.Deposited, nil
}

// joinNotaryPool checks if the deposit flag is true and the account is a
// voter in the VRC. If the account is not in the set, it will deposit ESS
// into contract.
func joinNotaryPool(manager mainchain.ContractManager, client mainchain.FullClient) error {
	if !client.DepositFlag() {
		return errors.New("joinNotaryPool called when deposit flag was not set")
	}

	if b, err := isAccountInNotaryPool(manager, client.Account()); b || err != nil {
		if b {
			log.Info("Already joined voter pool")
			return nil
		}
		return err
	}

	log.Info("Joining voter pool")
	txOps, err := manager.CreateTXOpts(shardparams.DefaultConfig.NotaryDeposit)

	if err != nil {
		return fmt.Errorf("unable to initiate the deposit transaction: %v", err)
	}

	tx, err := manager.VRCTransactor().RegisterNotary(txOps)
	if err != nil {
		return fmt.Errorf("unable to deposit ESS and become a voter: %v", err)
	}

	err = client.WaitForTransaction(context.Background(), tx.Hash(), 400)
	if err != nil {
		return err
	}

	receipt, err := client.TransactionReceipt(tx.Hash())
	if err != nil {
		return err
	}
	log.Info("receipt.Status", "st", receipt)
	if receipt.Status == gethTypes.ReceiptStatusFailed {
		return errors.New("transaction was not successful, unable to deposit ESS and become a voter")
	}

	if inPool, err := isAccountInNotaryPool(manager, client.Account()); !inPool || err != nil {
		if err != nil {
			return err
		}
		return errors.New("account has not been able to be deposited in voter pool")
	}

	log.Infof("Deposited %dESS into contract with transaction hash: %s", new(big.Int).Div(shardparams.DefaultConfig.NotaryDeposit, big.NewInt(params.Ether)), tx.Hash().Hex())

	return nil
}
