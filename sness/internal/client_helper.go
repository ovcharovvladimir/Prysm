package internal

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ovcharovvladimir/Prysm/sness/contracts"
	shardparams "github.com/ovcharovvladimir/Prysm/sness/params"
	ethereum "github.com/ovcharovvladimir/essentiaHybrid"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind/backends"
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	"github.com/ovcharovvladimir/essentiaHybrid/core"
	gethTypes "github.com/ovcharovvladimir/essentiaHybrid/core/types"
	"github.com/ovcharovvladimir/essentiaHybrid/crypto"
)

var (
	key, _            = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	addr              = crypto.PubkeyToAddress(key.PublicKey)
	accountBalance, _ = new(big.Int).SetString("1001000000000000000000", 10)
)

// MockClient for testing proposer.
type MockClient struct {
	VRC         *contracts.VRC
	T           *testing.T
	depositFlag bool
	Backend     *backends.SimulatedBackend
	BlockNumber int64
}

func (m *MockClient) Account() *accounts.Account {
	return &accounts.Account{Address: addr}
}

func (m *MockClient) VRCCaller() *contracts.VRCCaller {
	return &m.VRC.VRCCaller
}

func (m *MockClient) ChainReader() ethereum.ChainReader {
	return nil
}

func (m *MockClient) VRCTransactor() *contracts.VRCTransactor {
	return &m.VRC.VRCTransactor
}

func (m *MockClient) VRCFilterer() *contracts.VRCFilterer {
	return &m.VRC.VRCFilterer
}

func (m *MockClient) WaitForTransaction(ctx context.Context, hash common.Hash, durationInSeconds time.Duration) error {
	m.CommitWithBlock()
	m.FastForward(1)
	return nil
}

func (m *MockClient) TransactionReceipt(hash common.Hash) (*gethTypes.Receipt, error) {
	return m.Backend.TransactionReceipt(context.Background(), hash)
}

func (m *MockClient) CreateTXOpts(value *big.Int) (*bind.TransactOpts, error) {
	txOpts := TransactOpts()
	txOpts.Value = value
	return txOpts, nil
}

func (m *MockClient) SetDepositFlag(value bool) {
	m.depositFlag = value
}

func (m *MockClient) DepositFlag() bool {
	return m.depositFlag
}

func (m *MockClient) Sign(hash common.Hash) ([]byte, error) {
	return nil, nil
}

func (m *MockClient) GetShardCount() (int64, error) {
	return 100, nil
}

func (m *MockClient) CommitWithBlock() {
	m.Backend.Commit()
	m.BlockNumber = m.BlockNumber + 1
}

func (m *MockClient) FastForward(p int) {
	for i := 0; i < p*int(shardparams.DefaultConfig.PeriodLength); i++ {
		m.CommitWithBlock()
	}
}

func (m *MockClient) SubscribeNewHead(ctx context.Context, ch chan<- *gethTypes.Header) (ethereum.Subscription, error) {
	return nil, nil
}

func (m *MockClient) BlockByNumber(ctx context.Context, number *big.Int) (*gethTypes.Block, error) {
	return gethTypes.NewBlockWithHeader(&gethTypes.Header{Number: big.NewInt(m.BlockNumber)}), nil
}

func TransactOpts() *bind.TransactOpts {
	return bind.NewKeyedTransactor(key)
}

func SetupMockClient(t *testing.T) (*backends.SimulatedBackend, *contracts.VRC) {
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: accountBalance}}, 2100)
	_, _, VRC, err := contracts.DeployVRC(TransactOpts(), backend)
	if err != nil {
		t.Fatalf("Failed to deploy VRC contract: %v", err)
	}
	backend.Commit()
	return backend, VRC
}
