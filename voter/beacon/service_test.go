package beacon

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	pbp2p "github.com/ovcharovvladimir/Prysm/proto/beacon/p2p/v1"
	pb "github.com/ovcharovvladimir/Prysm/proto/beacon/rpc/v1"
	"github.com/ovcharovvladimir/Prysm/shared/testutil"
	"github.com/ovcharovvladimir/Prysm/voter/internal"
	"github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(ioutil.Discard)
}

type mockClient struct {
	ctrl *gomock.Controller
}

func (fc *mockClient) BeaconServiceClient() pb.BeaconServiceClient {
	mockServiceClient := internal.NewMockBeaconServiceClient(fc.ctrl)

	blockStream := internal.NewMockBeaconService_LatestBeaconBlockClient(fc.ctrl)
	blockStream.EXPECT().Recv().Return(&pbp2p.BeaconBlock{}, io.EOF)
	stateStream := internal.NewMockBeaconService_LatestCrystallizedStateClient(fc.ctrl)
	stateStream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{}, io.EOF)

	mockServiceClient.EXPECT().LatestBeaconBlock(
		gomock.Any(),
		&empty.Empty{},
	).Return(blockStream, nil)
	mockServiceClient.EXPECT().LatestCrystallizedState(
		gomock.Any(),
		&empty.Empty{},
	).Return(stateStream, nil)
	return mockServiceClient
}

func TestChannelGetters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	b := NewBeaconVoter(context.Background(), &Config{AttesterChanBuf: 1, ProposerChanBuf: 1}, &mockClient{ctrl})
	b.proposerChan <- false
	proposerVal := <-b.ProposerAssignment()
	if proposerVal {
		t.Error("Expected false value from channel, received true")
	}
	b.attesterChan <- false
	attesterVal := <-b.AttesterAssignment()
	if attesterVal {
		t.Error("Expected false value from channel, received true")
	}
}

func TestLifecycle(t *testing.T) {
	hook := logTest.NewGlobal()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	b := NewBeaconVoter(context.Background(), &Config{AttesterChanBuf: 0, ProposerChanBuf: 0}, &mockClient{ctrl})

	// Testing default config values.
	cfg := DefaultConfig()
	if cfg.AttesterChanBuf != 5 && cfg.ProposerChanBuf != 5 {
		t.Error("Default config values incorrect")
	}

	b.Start()
	// TODO: find a better way to test this. The problem is that start is nonblocking, and it ends
	// before the for loops of its inner goroutines begin.
	time.Sleep(time.Millisecond * 10)
	testutil.AssertLogsContain(t, hook, "Starting service")
	b.Stop()
	testutil.AssertLogsContain(t, hook, "Stopping service")
}

func TestFetchBeaconBlocks(t *testing.T) {
	hook := logTest.NewGlobal()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := NewBeaconVoter(context.Background(), &Config{AttesterChanBuf: 1, ProposerChanBuf: 1}, &mockClient{ctrl})

	// Create mock for the stream returned by LatestBeaconBlock.
	stream := internal.NewMockBeaconService_LatestBeaconBlockClient(ctrl)

	// If the block's slot number from the stream matches the assigned attestation height,
	// trigger a log.
	stream.EXPECT().Recv().Return(&pbp2p.BeaconBlock{SlotNumber: 10}, nil)
	stream.EXPECT().Recv().Return(&pbp2p.BeaconBlock{}, io.EOF)
	b.assignedHeight = 10
	b.responsibility = "attester"

	mockServiceClient := internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestBeaconBlock(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, nil)

	b.fetchBeaconBlocks(mockServiceClient)
	<-b.attesterChan

	testutil.AssertLogsContain(t, hook, "Latest beacon block slot number")
	testutil.AssertLogsContain(t, hook, "Assigned attestation slot number reached")

	// If the voter is assigned to be a proposer, trigger a log upon next
	// SlotNumber being reached.
	stream = internal.NewMockBeaconService_LatestBeaconBlockClient(ctrl)

	stream.EXPECT().Recv().Return(&pbp2p.BeaconBlock{SlotNumber: 1}, nil)
	stream.EXPECT().Recv().Return(&pbp2p.BeaconBlock{}, io.EOF)
	b.responsibility = "proposer"

	mockServiceClient = internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestBeaconBlock(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, nil)

	b.fetchBeaconBlocks(mockServiceClient)
	<-b.proposerChan

	testutil.AssertLogsContain(t, hook, "Latest beacon block slot number")
	testutil.AssertLogsContain(t, hook, "Assigned proposal slot number reached")

	// Testing an error coming from the stream.
	stream = internal.NewMockBeaconService_LatestBeaconBlockClient(ctrl)
	stream.EXPECT().Recv().Return(&pbp2p.BeaconBlock{}, errors.New("stream error"))
	stream.EXPECT().Recv().Return(&pbp2p.BeaconBlock{}, io.EOF)

	mockServiceClient = internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestBeaconBlock(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, nil)

	b.fetchBeaconBlocks(mockServiceClient)

	testutil.AssertLogsContain(t, hook, "stream error")

	// Creating a faulty stream will trigger error.
	mockServiceClient = internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestBeaconBlock(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, errors.New("stream creation failed"))

	b.fetchBeaconBlocks(mockServiceClient)
	testutil.AssertLogsContain(t, hook, "stream creation failed")
}

func TestFetchCrystallizedState(t *testing.T) {
	hook := logTest.NewGlobal()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := NewBeaconVoter(context.Background(), &Config{AttesterChanBuf: 0, ProposerChanBuf: 0}, &mockClient{ctrl})

	// Creating a faulty stream will trigger error.
	stream := internal.NewMockBeaconService_LatestCrystallizedStateClient(ctrl)
	mockServiceClient := internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestCrystallizedState(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, errors.New("stream creation failed"))

	b.fetchCrystallizedState(mockServiceClient)

	testutil.AssertLogsContain(t, hook, "stream creation failed")

	// Stream recv error should trigger error log.
	stream = internal.NewMockBeaconService_LatestCrystallizedStateClient(ctrl)
	stream.EXPECT().Recv().Return(nil, errors.New("recv error"))
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{}, io.EOF)

	mockServiceClient = internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestCrystallizedState(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, nil)

	b.fetchCrystallizedState(mockServiceClient)

	testutil.AssertLogsContain(t, hook, "recv error")

	// Being unable to marshal the received crystallized state should log an error.
	stream = internal.NewMockBeaconService_LatestCrystallizedStateClient(ctrl)
	stream.EXPECT().Recv().Return(nil, nil)
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{}, io.EOF)

	mockServiceClient = internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestCrystallizedState(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, nil)

	b.fetchCrystallizedState(mockServiceClient)

	testutil.AssertLogsContain(t, hook, "Could not marshal crystallized state proto")

	// If the current voter is not found within the active validators list, log a debug message.
	voter := &pbp2p.ValidatorRecord{WithdrawalAddress: []byte("0x01"), StartDynasty: 1, EndDynasty: 10}
	stream = internal.NewMockBeaconService_LatestCrystallizedStateClient(ctrl)
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{Validators: []*pbp2p.ValidatorRecord{voter}, CurrentDynasty: 5}, nil)
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{}, io.EOF)

	mockServiceClient = internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestCrystallizedState(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, nil)

	b.fetchCrystallizedState(mockServiceClient)

	testutil.AssertLogsContain(t, hook, "Validator index not found in latest crystallized state's active voter list")

	// A faulty client.ShuffleValidators should log error.
	voter = &pbp2p.ValidatorRecord{WithdrawalAddress: []byte{}, StartDynasty: 1, EndDynasty: 10}
	stream = internal.NewMockBeaconService_LatestCrystallizedStateClient(ctrl)
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{Validators: []*pbp2p.ValidatorRecord{voter}, CurrentDynasty: 5}, nil)
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{}, io.EOF)

	mockServiceClient = internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestCrystallizedState(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, nil)
	mockServiceClient.EXPECT().FetchShuffledValidatorIndices(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, errors.New("something went wrong"))

	b.fetchCrystallizedState(mockServiceClient)

	testutil.AssertLogsContain(t, hook, "Could not fetch shuffled voter indices: something went wrong")

	// Height should be assigned based on the result of ShuffleValidators.
	validator1 := &pbp2p.ValidatorRecord{WithdrawalAddress: []byte("0x0"), StartDynasty: 1, EndDynasty: 10}
	validator2 := &pbp2p.ValidatorRecord{WithdrawalAddress: []byte("0x1"), StartDynasty: 1, EndDynasty: 10}
	validator3 := &pbp2p.ValidatorRecord{WithdrawalAddress: []byte{}, StartDynasty: 1, EndDynasty: 10}
	stream = internal.NewMockBeaconService_LatestCrystallizedStateClient(ctrl)
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{Validators: []*pbp2p.ValidatorRecord{validator1, validator2, validator3}, CurrentDynasty: 5}, nil)
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{}, io.EOF)

	mockServiceClient = internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestCrystallizedState(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, nil)
	mockServiceClient.EXPECT().FetchShuffledValidatorIndices(
		gomock.Any(),
		gomock.Any(),
	).Return(&pb.ShuffleResponse{
		AssignedAttestationHeights: []uint64{0, 1, 2},
		CutoffIndices:              []uint64{0, 1, 2},
		ShuffledValidatorIndices:   []uint64{2, 1, 0},
	}, nil)

	b.fetchCrystallizedState(mockServiceClient)

	testutil.AssertLogsContain(t, hook, "Validator selected as attester")

	// If the voter is the last index in the shuffled voter indices, it should be assigned
	// to be a proposer.
	validator1 = &pbp2p.ValidatorRecord{WithdrawalAddress: []byte("0x0"), StartDynasty: 1, EndDynasty: 10}
	validator2 = &pbp2p.ValidatorRecord{WithdrawalAddress: []byte("0x1"), StartDynasty: 1, EndDynasty: 10}
	validator3 = &pbp2p.ValidatorRecord{WithdrawalAddress: []byte{}, StartDynasty: 1, EndDynasty: 10}
	stream = internal.NewMockBeaconService_LatestCrystallizedStateClient(ctrl)
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{Validators: []*pbp2p.ValidatorRecord{validator1, validator2, validator3}, CurrentDynasty: 5}, nil)
	stream.EXPECT().Recv().Return(&pbp2p.CrystallizedState{}, io.EOF)

	mockServiceClient = internal.NewMockBeaconServiceClient(ctrl)
	mockServiceClient.EXPECT().LatestCrystallizedState(
		gomock.Any(),
		gomock.Any(),
	).Return(stream, nil)
	mockServiceClient.EXPECT().FetchShuffledValidatorIndices(
		gomock.Any(),
		gomock.Any(),
	).Return(&pb.ShuffleResponse{
		ShuffledValidatorIndices: []uint64{0, 1, 2},
	}, nil)

	b.fetchCrystallizedState(mockServiceClient)

	testutil.AssertLogsContain(t, hook, "Validator selected as proposer of the next slot")
}
