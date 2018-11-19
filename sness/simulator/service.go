// Package simulator defines the simulation utility to test the sness.
package simulator

import (
	"context"

	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/mattn/go-colorable"

	pb "github.com/ovcharovvladimir/Prysm/proto/beacon/p2p/v1"
	"github.com/ovcharovvladimir/Prysm/shared/event"
	"github.com/ovcharovvladimir/Prysm/shared/p2p"
	"github.com/ovcharovvladimir/Prysm/shared/slotticker"
	"github.com/ovcharovvladimir/Prysm/sness/params"
	"github.com/ovcharovvladimir/Prysm/sness/types"
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	"github.com/ovcharovvladimir/essentiaHybrid/log"
)

type p2pAPI interface {
	Subscribe(msg proto.Message, channel chan p2p.Message) event.Subscription
	Send(msg proto.Message, peer p2p.Peer)
	Broadcast(msg proto.Message)
}

type powChainService interface {
	LatestBlockHash() common.Hash
}

// Simulator struct.
type Simulator struct {
	ctx              context.Context
	cancel           context.CancelFunc
	p2p              p2pAPI
	web3Service      powChainService
	beaconDB         beaconDB
	enablePOWChain   bool
	blockRequestChan chan p2p.Message
}

// Config options for the simulator service.
type Config struct {
	BlockRequestBuf int
	P2P             p2pAPI
	Web3Service     powChainService
	BeaconDB        beaconDB
	EnablePOWChain  bool
}

type beaconDB interface {
	GetChainHead() (*types.Block, error)
	GetGenesisTime() (time.Time, error)
	GetActiveState() (*types.ActiveState, error)
	GetCrystallizedState() (*types.CrystallizedState, error)
}

// DefaultConfig options for the simulator.
func DefaultConfig() *Config {
	return &Config{
		BlockRequestBuf: 100,
	}
}

// NewSimulator creates a simulator instance for a syncer to consume fake, generated blocks.
func NewSimulator(ctx context.Context, cfg *Config) *Simulator {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(colorable.NewColorableStdout(), log.TerminalFormat(true))))

	ctx, cancel := context.WithCancel(ctx)
	return &Simulator{
		ctx:              ctx,
		cancel:           cancel,
		p2p:              cfg.P2P,
		web3Service:      cfg.Web3Service,
		beaconDB:         cfg.BeaconDB,
		enablePOWChain:   cfg.EnablePOWChain,
		blockRequestChan: make(chan p2p.Message, cfg.BlockRequestBuf),
	}
}

// Start the sim.
func (sim *Simulator) Start() {
	log.Info("Starting service")
	genesisTime, err := sim.beaconDB.GetGenesisTime()
	if err != nil {
		log.Crit(err.Error())
		return
	}

	slotTicker := slotticker.GetSlotTicker(genesisTime, params.GetConfig().SlotDuration)
	go func() {
		sim.run(slotTicker.C(), sim.blockRequestChan)
		close(sim.blockRequestChan)
		slotTicker.Done()
	}()
}

// Stop the sim.
func (sim *Simulator) Stop() error {
	defer sim.cancel()
	log.Info("Stopping service")
	return nil
}

func (sim *Simulator) run(slotInterval <-chan uint64, requestChan <-chan p2p.Message) {
	blockReqSub := sim.p2p.Subscribe(&pb.BeaconBlockRequest{}, sim.blockRequestChan)
	defer blockReqSub.Unsubscribe()

	lastBlock, err := sim.beaconDB.GetChainHead()
	if err != nil {
		log.Error("Could not fetch latest block", "err", err)
		return
	}

	lastHash, err := lastBlock.Hash()
	if err != nil {
		log.Error("Could not get hash of the latest block", "err", err.Error())
	}
	broadcastedBlocks := map[[32]byte]*types.Block{}

	for {
		select {
		case <-sim.ctx.Done():
			log.Debug("Simulator context closed, exiting goroutine")
			return
		case slot := <-slotInterval:
			aState, err := sim.beaconDB.GetActiveState()
			if err != nil {
				log.Error("Failed to get active state", "err", err.Error())
				continue
			}
			cState, err := sim.beaconDB.GetCrystallizedState()
			if err != nil {
				log.Error("Failed to get crystallized state", "err", err.Error())
				continue
			}

			aStateHash, err := aState.Hash()
			if err != nil {
				log.Error("Failed to hash active state", "err", err.Error())
				continue
			}

			cStateHash, err := cState.Hash()
			if err != nil {
				log.Error("Failed to hash crystallized state", "err", err.Error())
				continue
			}

			var powChainRef []byte
			if sim.enablePOWChain {
				powChainRef = sim.web3Service.LatestBlockHash().Bytes()
			} else {
				powChainRef = []byte{byte(slot)}
			}

			parentHash := make([]byte, 32)
			copy(parentHash, lastHash[:])
			block := types.NewBlock(&pb.BeaconBlock{
				Slot:                  slot,
				Timestamp:             ptypes.TimestampNow(),
				PowChainRef:           powChainRef,
				ActiveStateRoot:       aStateHash[:],
				CrystallizedStateRoot: cStateHash[:],
				AncestorHashes:        [][]byte{parentHash},
				RandaoReveal:          params.GetConfig().SimulatedBlockRandao[:],
				Attestations: []*pb.AggregatedAttestation{
					{Slot: slot - 1, AttesterBitfield: []byte{byte(255)}},
				},
			})

			hash, err := block.Hash()
			if err != nil {
				log.Error("Could not hash simulated block", "err", err.Error())
				continue
			}
			sim.p2p.Broadcast(&pb.BeaconBlockHashAnnounce{
				Hash: hash[:],
			})
			log.Info("Broadcast block hash", "hash", hash, "slot", slot)

			broadcastedBlocks[hash] = block
			lastHash = hash
		case msg := <-requestChan:
			data := msg.Data.(*pb.BeaconBlockRequest)
			var hash [32]byte
			copy(hash[:], data.Hash)

			block := broadcastedBlocks[hash]
			if block == nil {
				log.Error("Requested block not found", "hash", hash)
				continue
			}
			log.Info("Responding to full block request", "hash", hash)

			// Sends the full block body to the requester.
			res := &pb.BeaconBlockResponse{Block: block.Proto(), Attestation: &pb.AggregatedAttestation{
				Slot:             block.SlotNumber(),
				AttesterBitfield: []byte{byte(255)},
			}}
			sim.p2p.Send(res, msg.Peer)

			delete(broadcastedBlocks, hash)
		}
	}
}
