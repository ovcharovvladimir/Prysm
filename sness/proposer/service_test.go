package proposer

import (
	"crypto/rand"
	//"math"
	"math/big"
	"testing"

	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	gethTypes "github.com/ovcharovvladimir/essentiaHybrid/core/types"
	"github.com/ovcharovvladimir/Prysm/sness/internal"
	"github.com/ovcharovvladimir/Prysm/sness/params"
)

// TODO: Fix the tests so that the following tests can be tested as a package and stop leaking into each other,
// find another way rather than using logs to do this

/*
func settingUpProposer(t *testing.T) (*Proposer, *internal.MockClient) {
	backend, vrc := internal.SetupMockClient(t)
	node := &internal.MockClient{VRC: vrc, T: t, Backend: backend}

	server, err := p2p.NewServer()
	if err != nil {
		t.Fatalf("Failed to start server %v", err)
	}
	server.Start()

	pool, err := txpool.NewTXPool(server)
	if err != nil {
		t.Fatalf("Failed to start txpool %v", err)
	}
	pool.Start()

	config := &database.ShardDBConfig{DataDir: "", Name: "", InMemory: true}

	db, err := database.NewShardDB(config)
	if err != nil {
		t.Fatalf("Failed create shardDB %v", err)
	}

	db.Start()

	fakeSyncer, err := syncer.NewSyncer(params.DefaultConfig, &mainchain.VRCClient{}, server, db, 1)
	if err != nil {
		t.Fatalf("Failed to start syncer %v", err)
	}

	fakeProposer, err := NewProposer(params.DefaultConfig, node, server, pool, db, 1, fakeSyncer)
	if err != nil {
		t.Fatalf("Failed to create proposer %v", err)
	}
	fakeProposer.config.CollationSizeLimit = int64(math.Pow(float64(2), float64(10)))

	return fakeProposer, node

}

func TestProposerRoundTrip(t *testing.T) {
	hook := logTest.NewGlobal()
	fakeProposer, node := settingUpProposer(t)

	input := make([]byte, 0, 2000)
	for len(input) < int(fakeProposer.config.CollationSizeLimit/4) {
		input = append(input, []byte{'t', 'e', 's', 't', 'i', 'n', 'g'}...)
	}
	tx := pb.Transaction{Input: input}

	for i := 0; i < 5; i++ {
		node.CommitWithBlock()
	}
	fakeProposer.Start()

	for i := 0; i < 7; i++ {
		fakeProposer.p2p.Broadcast(&tx)
		<-fakeProposer.msgChan
	}

	want := "Collation created"
	length := len(hook.AllEntries())
	for length < 5 {
		length = len(hook.AllEntries())
	}

	msg := hook.AllEntries()[4]

	if msg.Message != want {
		t.Errorf("Incorrect log, wanted %v but got %v", want, msg.Message)
	}

	fakeProposer.cancel()
	fakeProposer.Stop()
	fakeProposer.dbService.Stop()
	fakeProposer.txpool.Stop()
	fakeProposer.p2p.Stop()

	hook.Reset()

}

func TestIncompleteCollation(t *testing.T) {
	hook := logTest.NewGlobal()
	fakeProposer, node := settingUpProposer(t)

	input := make([]byte, 0, 2000)
	for int64(len(input)) < (fakeProposer.config.CollationSizeLimit)/4 {
		input = append(input, []byte{'t', 'e', 's', 't', 'i', 'n', 'g'}...)
	}
	tx := pb.Transaction{Input: input}

	for i := 0; i < 5; i++ {
		node.CommitWithBlock()
	}
	fakeProposer.Start()

	for i := 0; i < 3; i++ {
		fakeProposer.p2p.Broadcast(&tx)
		<-fakeProposer.msgChan
	}

	want := "Starting proposer service"

	msg := hook.LastEntry()
	if msg.Message != want {
		t.Errorf("Incorrect log, wanted %v but got %v", want, msg.Message)
	}

	length := len(hook.AllEntries())

	if length != 4 {
		t.Errorf("Number of logs was supposed to be 4 but is %v", length)
	}

	fakeProposer.cancel()
	fakeProposer.Stop()
	fakeProposer.dbService.Stop()
	fakeProposer.txpool.Stop()
	fakeProposer.p2p.Stop()
	hook.Reset()
}

func TestCollationWitInDiffPeriod(t *testing.T) {
	hook := logTest.NewGlobal()
	fakeProposer, node := settingUpProposer(t)

	defer fakeProposer.p2p.Stop()
	defer fakeProposer.txpool.Stop()
	defer fakeProposer.dbService.Stop()
	defer fakeProposer.Stop()

	input := make([]byte, 0, 2000)
	for int64(len(input)) < (fakeProposer.config.CollationSizeLimit)/4 {
		input = append(input, []byte{'t', 'e', 's', 't', 'i', 'n', 'g'}...)
	}
	tx := pb.Transaction{Input: input}

	for i := 0; i < 5; i++ {
		node.CommitWithBlock()
	}
	fakeProposer.Start()

	fakeProposer.p2p.Broadcast(&tx)
	<-fakeProposer.msgChan

	for i := 0; i < 5; i++ {
		node.CommitWithBlock()
	}

	fakeProposer.p2p.Broadcast(&tx)
	<-fakeProposer.msgChan

	want := "Collation created"
	length := len(hook.AllEntries())
	for length < 5 {
		length = len(hook.AllEntries())
	}

	msg := hook.AllEntries()[4]

	if msg.Message != want {
		t.Errorf("Incorrect log, wanted %v but got %v", want, msg.Message)
	}

	fakeProposer.cancel()
	fakeProposer.Stop()
	fakeProposer.dbService.Stop()
	fakeProposer.txpool.Stop()
	fakeProposer.p2p.Stop()
	hook.Reset()
}
*/
func TestCreateCollation(t *testing.T) {
	backend, vrc := internal.SetupMockClient(t)
	node := &internal.MockClient{VRC: vrc, T: t, Backend: backend}
	var txs []*gethTypes.Transaction
	for i := 0; i < 10; i++ {
		data := make([]byte, 1024)
		rand.Read(data)
		txs = append(txs, gethTypes.NewTransaction(0, common.HexToAddress("0x0"),
			nil, 0, nil, data))
	}

	_, err := createCollation(node, node.Account(), node, big.NewInt(0), big.NewInt(1), txs)
	if err != nil {
		t.Fatalf("Create collation failed: %v", err)
	}

	// fast forward to 2nd period.
	for i := 0; i < 2*int(params.DefaultConfig.PeriodLength); i++ {
		backend.Commit()
	}

	// negative test case #1: create collation with shard > shardCount.
	_, err = createCollation(node, node.Account(), node, big.NewInt(101), big.NewInt(2), txs)
	if err == nil {
		t.Errorf("Create collation should have failed with invalid shard number")
	}
	// negative test case #2, create collation with blob size > collationBodySizeLimit.
	var badTxs []*gethTypes.Transaction
	for i := 0; i <= 1024; i++ {
		data := make([]byte, 1024)
		rand.Read(data)
		badTxs = append(badTxs, gethTypes.NewTransaction(0, common.HexToAddress("0x0"),
			nil, 0, nil, data))
	}
	_, err = createCollation(node, node.Account(), node, big.NewInt(0), big.NewInt(2), badTxs)
	if err == nil {
		t.Errorf("Create collation should have failed with Txs longer than collation body limit")
	}

	// normal test case #1 create collation with correct parameters.
	collation, err := createCollation(node, node.Account(), node, big.NewInt(5), big.NewInt(5), txs)
	if err != nil {
		t.Errorf("Create collation failed: %v", err)
	}
	if collation.Header().Period().Cmp(big.NewInt(5)) != 0 {
		t.Errorf("Incorrect collation period, want 5, got %v ", collation.Header().Period())
	}
	if collation.Header().ShardID().Cmp(big.NewInt(5)) != 0 {
		t.Errorf("Incorrect shard id, want 5, got %v ", collation.Header().ShardID())
	}
	if *collation.ProposerAddress() != node.Account().Address {
		t.Errorf("Incorrect proposer address, got %v", *collation.ProposerAddress())
	}
	if collation.Header().Sig() != [32]byte{} {
		t.Errorf("Proposer signaure can not be empty")
	}
}

func TestAddCollation(t *testing.T) {
	backend, vrc := internal.SetupMockClient(t)
	node := &internal.MockClient{VRC: vrc, T: t, Backend: backend}
	var txs []*gethTypes.Transaction
	for i := 0; i < 10; i++ {
		data := make([]byte, 1024)
		rand.Read(data)
		txs = append(txs, gethTypes.NewTransaction(0, common.HexToAddress("0x0"),
			nil, 0, nil, data))
	}

	collation, err := createCollation(node, node.Account(), node, big.NewInt(0), big.NewInt(1), txs)
	if err != nil {
		t.Errorf("Create collation failed: %v", err)
	}

	// fast forward to next period.
	for i := 0; i < int(params.DefaultConfig.PeriodLength); i++ {
		backend.Commit()
	}

	// normal test case #1 create collation with normal parameters.
	err = AddHeader(node, node, collation)
	if err != nil {
		t.Errorf("%v", err)
	}
	backend.Commit()

	// verify collation was correctly added from VRC.
	collationFromVRC, err := vrc.CollationRecords(&bind.CallOpts{}, big.NewInt(0), big.NewInt(1))
	if err != nil {
		t.Errorf("Failed to get collation record")
	}
	if collationFromVRC.Proposer != node.Account().Address {
		t.Errorf("Incorrect proposer address, got %v", *collation.ProposerAddress())
	}
	if common.BytesToHash(collationFromVRC.ChunkRoot[:]) != *collation.Header().ChunkRoot() {
		t.Errorf("Incorrect chunk root, got %v", collationFromVRC.ChunkRoot)
	}

	// negative test case #1 create the same collation that just got added to VRC.
	_, err = createCollation(node, node.Account(), node, big.NewInt(0), big.NewInt(1), txs)
	if err == nil {
		t.Errorf("Create collation should fail due to same collation in VRC")
	}
}

func TestCheckCollation(t *testing.T) {
	backend, vrc := internal.SetupMockClient(t)
	node := &internal.MockClient{VRC: vrc, T: t, Backend: backend}
	var txs []*gethTypes.Transaction
	for i := 0; i < 10; i++ {
		data := make([]byte, 1024)
		rand.Read(data)
		txs = append(txs, gethTypes.NewTransaction(0, common.HexToAddress("0x0"),
			nil, 0, nil, data))
	}

	collation, err := createCollation(node, node.Account(), node, big.NewInt(0), big.NewInt(1), txs)
	if err != nil {
		t.Errorf("Create collation failed: %v", err)
	}

	for i := 0; i < int(params.DefaultConfig.PeriodLength); i++ {
		backend.Commit()
	}

	err = AddHeader(node, node, collation)
	if err != nil {
		t.Errorf("%v", err)
	}
	backend.Commit()

	// normal test case 1: check if we can still add header for period 1, should return false.
	a, err := checkHeaderAdded(node, big.NewInt(0), big.NewInt(1))
	if err != nil {
		t.Errorf("Can not check header submitted: %v", err)
	}
	if a {
		t.Errorf("Check header submitted shouldn't return: %v", a)
	}
	// normal test case 2: check if we can add header for period 2, should return true.
	a, err = checkHeaderAdded(node, big.NewInt(0), big.NewInt(2))
	if err != nil {
		t.Error(err)
	}
	if !a {
		t.Errorf("Check header submitted shouldn't return: %v", a)
	}
}
