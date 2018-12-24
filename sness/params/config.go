// Package params defines important configuration options to be used when instantiating
// objects within the sharding package. For example, it defines objects such as a
// Config that will be useful when creating new shard instances.
package params

import (
	"math"
	"math/big"

	"github.com/ovcharovvladimir/essentiaHybrid/common"
)

// DefaultConfig contains default configs for node to use in the sharded universe.
var DefaultConfig = &Config{
	VRCAddress:            common.HexToAddress("0x93d6E7BbC08FB9FD9Ad07744B5C597c9f45b9eA7"),
	PeriodLength:          5,
	NotaryDeposit:         new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil), // 1000 ESS
	NotaryLockupLength:    16128,
	ProposerLockupLength:  48,
	NotaryCommitteeSize:   135,
	NotaryQuorumSize:      90,
	NotaryChallengePeriod: 25,
	CollationSizeLimit:    int64(math.Pow(float64(2), float64(20))),
}

// DefaultChainConfig contains default chain configs of an individual shard.
var DefaultChainConfig = &ChainConfig{}

// Config contains configs for node to participate in the sharded universe.
type Config struct {
	VRCAddress            common.Address // VRCAddress is the address of VRC in mainchain.
	PeriodLength          int64          // PeriodLength is num of blocks in period.
	NotaryDeposit         *big.Int       // NotaryDeposit is a required deposit size in wei.
	NotaryLockupLength    int64          // NotaryLockupLength to lockup voter deposit from time of deregistration.
	ProposerLockupLength  int64          // ProposerLockupLength to lockup proposer deposit from time of deregistration.
	NotaryCommitteeSize   int64          // NotaryCommitSize sampled per block from the notaries pool per period per shard.
	NotaryQuorumSize      int64          // NotaryQuorumSize votes the collation needs to get accepted to the canonical chain.
	NotaryChallengePeriod int64          // NotaryChallengePeriod is the duration a voter has to store collations for.
	CollationSizeLimit    int64          // CollationSizeLimit is the maximum size the serialized blobs in a collation can take.
}

// ChainConfig contains chain config of an individual shard. Still to be designed.
type ChainConfig struct{}
