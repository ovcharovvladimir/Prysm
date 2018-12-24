package mainchain

import (
	"context"
	"math/big"
	"time"

	"github.com/ovcharovvladimir/Prysm/sness/contracts"
	ethereum "github.com/ovcharovvladimir/essentiaHybrid"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	gethTypes "github.com/ovcharovvladimir/essentiaHybrid/core/types"
)

// Signer defines an interface that can read from the Ethereum mainchain as well call
// read-only methods and functions from the Sharding Manager Contract.
type Signer interface {
	Sign(hash common.Hash) ([]byte, error)
}

// ContractManager specifies an interface that defines both read/write
// operations on a contract in the Ethereum mainchain.
type ContractManager interface {
	ContractCaller
	ContractTransactor
}

// ContractCaller defines an interface that can read from a contract on the
// Ethereum mainchain as well as call its read-only methods and functions.
type ContractCaller interface {
	VRCCaller() *contracts.VRCCaller
	GetShardCount() (int64, error)
}

// ContractTransactor defines an interface that can transact with a contract on the
// Ethereum mainchain as well as call its methods and functions.
type ContractTransactor interface {
	VRCTransactor() *contracts.VRCTransactor
	CreateTXOpts(value *big.Int) (*bind.TransactOpts, error)
}

// essclient defines the methods that will be used to perform rpc calls
// to the main geth node, and be responsible for other user-specific data
type ess_client interface {
	Account() *accounts.Account
	WaitForTransaction(ctx context.Context, hash common.Hash, durationInSeconds time.Duration) error
	TransactionReceipt(hash common.Hash) (*gethTypes.Receipt, error)
	DepositFlag() bool
}

type FullClient interface {
	ess_client
	Reader
	ContractCaller
	Signer
	ContractTransactor
}

// Reader defines an interface for a struct that can read mainchain information
// such as blocks, transactions, receipts, and more. Useful for testing.
type Reader interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*gethTypes.Block, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *gethTypes.Header) (ethereum.Subscription, error)
}

// RecordFetcher serves as an interface for a struct that can fetch collation information
// from a sharding manager contract on the Ethereum mainchain.
type RecordFetcher interface {
	CollationRecords(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
		ChunkRoot [32]byte
		Proposer  common.Address
		IsElected bool
		Signature [32]byte
	}, error)
}
