package initialsync

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	pb "github.com/ovcharovvladimir/Prysm/proto/beacon/p2p/v1"
	"github.com/ovcharovvladimir/Prysm/shared/event"
	"github.com/ovcharovvladimir/Prysm/shared/p2p"
	"github.com/ovcharovvladimir/Prysm/shared/testutil"
	"github.com/ovcharovvladimir/Prysm/sness/types"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

type mockP2P struct {
}

func (mp *mockP2P) Subscribe(msg proto.Message, channel chan p2p.Message) event.Subscription {
	return new(event.Feed).Subscribe(channel)
}

func (mp *mockP2P) Broadcast(msg proto.Message) {}

func (mp *mockP2P) Send(msg proto.Message, peer p2p.Peer) {
}

type mockSyncService struct {
	hasStarted bool
	isSynced   bool
}

func (ms *mockSyncService) Start() {
	ms.hasStarted = true
}

func (ms *mockSyncService) IsSyncedWithNetwork() bool {
	return ms.isSynced
}

type mockDB struct{}

func (m *mockDB) SaveBlock(*types.Block) error {
	return nil
}

func TestSetBlockForInitialSync(t *testing.T) {
	hook := logTest.NewGlobal()

	cfg := Config{
		P2P:         &mockP2P{},
		SyncService: &mockSyncService{},
		BeaconDB:    &mockDB{},
	}

	ss := NewInitialSyncService(context.Background(), cfg)

	exitRoutine := make(chan bool)
	delayChan := make(chan time.Time)
	defer func() {
		close(exitRoutine)
		close(delayChan)
	}()

	go func() {
		ss.run(delayChan)
		exitRoutine <- true
	}()

	genericHash := make([]byte, 32)
	genericHash[0] = 'a'

	block := &pb.BeaconBlock{
		PowChainRef:           []byte{1, 2, 3},
		AncestorHashes:        [][]byte{genericHash},
		Slot:                  uint64(20),
		CrystallizedStateRoot: genericHash,
	}

	blockResponse := &pb.BeaconBlockResponse{Block: block}
	msg1 := p2p.Message{
		Peer: p2p.Peer{},
		Data: blockResponse,
	}

	ss.blockBuf <- msg1

	ss.cancel()
	<-exitRoutine

	var hash [32]byte
	copy(hash[:], blockResponse.Block.CrystallizedStateRoot)

	if hash != ss.initialCrystallizedStateRoot {
		t.Fatalf("Crystallized state hash not updated: %#x", blockResponse.Block.CrystallizedStateRoot)
	}

	hook.Reset()
}

func TestSavingBlocksInSync(t *testing.T) {
	hook := logTest.NewGlobal()

	cfg := Config{
		P2P:         &mockP2P{},
		SyncService: &mockSyncService{},
		BeaconDB:    &mockDB{},
	}
	ss := NewInitialSyncService(context.Background(), cfg)

	exitRoutine := make(chan bool)
	delayChan := make(chan time.Time)

	defer func() {
		close(exitRoutine)
		close(delayChan)
	}()

	go func() {
		ss.run(delayChan)
		exitRoutine <- true
	}()

	genericHash := make([]byte, 32)
	genericHash[0] = 'a'

	crystallizedState := &pb.CrystallizedState{
		LastFinalizedSlot: 99,
	}

	stateResponse := &pb.CrystallizedStateResponse{
		CrystallizedState: crystallizedState,
	}

	incorrectState := &pb.CrystallizedState{
		LastFinalizedSlot: 9,
		LastJustifiedSlot: 20,
	}

	incorrectStateResponse := &pb.CrystallizedStateResponse{
		CrystallizedState: incorrectState,
	}

	crystallizedStateRoot, err := types.NewCrystallizedState(crystallizedState).Hash()
	if err != nil {
		t.Fatalf("unable to get hash of crystallized state: %v", err)
	}

	getBlockResponseMsg := func(Slot uint64) p2p.Message {
		block := &pb.BeaconBlock{
			PowChainRef:           []byte{1, 2, 3},
			AncestorHashes:        [][]byte{genericHash},
			Slot:                  Slot,
			CrystallizedStateRoot: crystallizedStateRoot[:],
		}

		blockResponse := &pb.BeaconBlockResponse{
			Block: block,
		}

		return p2p.Message{
			Peer: p2p.Peer{},
			Data: blockResponse,
		}
	}

	msg1 := getBlockResponseMsg(0)

	msg2 := p2p.Message{
		Peer: p2p.Peer{},
		Data: incorrectStateResponse,
	}

	ss.blockBuf <- msg1
	ss.crystallizedStateBuf <- msg2

	if ss.currentSlot == incorrectStateResponse.CrystallizedState.LastFinalizedSlot {
		t.Fatalf("Crystallized state updated incorrectly: %#x", ss.currentSlot)
	}

	msg2.Data = stateResponse

	ss.crystallizedStateBuf <- msg2

	if crystallizedStateRoot != ss.initialCrystallizedStateRoot {
		br := msg1.Data.(*pb.BeaconBlockResponse)
		t.Fatalf("Crystallized state hash not updated: %#x", br.Block.CrystallizedStateRoot)
	}

	msg1 = getBlockResponseMsg(30)
	ss.blockBuf <- msg1

	if stateResponse.CrystallizedState.GetLastFinalizedSlot() != ss.currentSlot {
		t.Fatalf("Slot saved when it was not supposed too: %v", stateResponse.CrystallizedState.GetLastFinalizedSlot())
	}

	msg1 = getBlockResponseMsg(100)
	ss.blockBuf <- msg1

	ss.cancel()
	<-exitRoutine

	br := msg1.Data.(*pb.BeaconBlockResponse)
	if br.Block.GetSlot() != ss.currentSlot {
		t.Fatalf("Slot not updated despite receiving a valid block: %v", ss.currentSlot)
	}

	hook.Reset()
}

func TestDelayChan(t *testing.T) {
	hook := logTest.NewGlobal()
	cfg := Config{
		P2P:         &mockP2P{},
		SyncService: &mockSyncService{},
		BeaconDB:    &mockDB{},
	}
	ss := NewInitialSyncService(context.Background(), cfg)

	exitRoutine := make(chan bool)
	delayChan := make(chan time.Time)

	defer func() {
		close(exitRoutine)
		close(delayChan)
	}()

	go func() {
		ss.run(delayChan)
		exitRoutine <- true
	}()

	genericHash := make([]byte, 32)
	genericHash[0] = 'a'

	crystallizedstate := &pb.CrystallizedState{
		LastFinalizedSlot: 99,
	}

	stateResponse := &pb.CrystallizedStateResponse{
		CrystallizedState: crystallizedstate,
	}

	crystallizedStateRoot, err := types.NewCrystallizedState(stateResponse.CrystallizedState).Hash()
	if err != nil {
		t.Fatalf("unable to get hash of crystallized state: %v", err)
	}

	block := &pb.BeaconBlock{
		PowChainRef:           []byte{1, 2, 3},
		AncestorHashes:        [][]byte{genericHash},
		Slot:                  uint64(20),
		CrystallizedStateRoot: crystallizedStateRoot[:],
	}

	blockResponse := &pb.BeaconBlockResponse{
		Block: block,
	}

	msg1 := p2p.Message{
		Peer: p2p.Peer{},
		Data: blockResponse,
	}

	msg2 := p2p.Message{
		Peer: p2p.Peer{},
		Data: stateResponse,
	}

	ss.blockBuf <- msg1

	ss.crystallizedStateBuf <- msg2

	blockResponse.Block.Slot = 100
	msg1.Data = blockResponse

	ss.blockBuf <- msg1

	delayChan <- time.Time{}

	ss.cancel()
	<-exitRoutine

	testutil.AssertLogsContain(t, hook, "Exiting initial sync and starting normal sync")

	hook.Reset()
}

func TestIsSyncedWithNetwork(t *testing.T) {
	hook := logTest.NewGlobal()
	mockSync := &mockSyncService{}
	cfg := Config{
		P2P:                 &mockP2P{},
		SyncService:         mockSync,
		BeaconDB:            &mockDB{},
		SyncPollingInterval: 1,
	}
	ss := NewInitialSyncService(context.Background(), cfg)

	mockSync.isSynced = true
	ss.Start()
	ss.Stop()

	testutil.AssertLogsContain(t, hook, "Chain state detected, exiting initial sync")
	testutil.AssertLogsContain(t, hook, "Stopping service")

	hook.Reset()
}

func TestIsNotSyncedWithNetwork(t *testing.T) {
	hook := logTest.NewGlobal()
	mockSync := &mockSyncService{}
	cfg := Config{
		P2P:                 &mockP2P{},
		SyncService:         mockSync,
		BeaconDB:            &mockDB{},
		SyncPollingInterval: 1,
	}
	ss := NewInitialSyncService(context.Background(), cfg)

	mockSync.isSynced = false
	ss.Start()
	ss.Stop()

	testutil.AssertLogsDoNotContain(t, hook, "Chain state detected, exiting initial sync")
	testutil.AssertLogsContain(t, hook, "Stopping service")

	hook.Reset()
}