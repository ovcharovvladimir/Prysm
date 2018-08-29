// Package legacyutil exists to convert Ethererum 2.0 types to go-ethereum or
// Ethereum 1.0 types.
package legacyutil

import (
	"math/big"

	pb "github.com/ovcharovvladimir/Prysm/proto/sharding/p2p/v1"
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	gethTypes "github.com/ovcharovvladimir/essentiaHybrid/core/types"
)

// TransformTransaction of proto transaction to gess's transction.
func TransformTransaction(t *pb.Transaction) *gethTypes.Transaction {
	return gethTypes.NewTransaction(
		t.Nonce,
		common.BytesToAddress(t.Recipient),
		big.NewInt(0).SetUint64(t.Value),
		t.GasLimit,
		big.NewInt(0).SetUint64(t.GasPrice),
		t.Input,
	)
}
