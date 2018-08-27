package types

import (
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	pb "github.com/ovcharovvladimir/Prysm/proto/beacon/rpc/v1"
)

// BeaconValidator defines a service that interacts with a beacon node via RPC to determine
// attestation/proposal responsibilities.
type BeaconValidator interface {
	AttesterAssignment() <-chan bool
	ProposerAssignment() <-chan bool
}

// CollationFetcher defines functionality for a struct that is able to extract
// respond with collation information to the caller. Shard implements this interface.
type CollationFetcher interface {
	CollationByHeaderHash(headerHash *common.Hash) (*Collation, error)
}

// RPCClient defines a struct that opens up RPC client services via gRPC.
type RPCClient interface {
	BeaconServiceClient() pb.BeaconServiceClient
}
