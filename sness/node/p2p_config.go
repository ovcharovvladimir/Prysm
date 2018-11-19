package node

import (
	"github.com/golang/protobuf/proto"
	"github.com/ovcharovvladimir/Prysm/shared/cmd"
	"github.com/ovcharovvladimir/Prysm/shared/p2p"
	"github.com/ovcharovvladimir/Prysm/shared/p2p/adapter/tracer"

	pb "github.com/ovcharovvladimir/Prysm/proto/beacon/p2p/v1"
	"gopkg.in/urfave/cli.v1"
)

var topicMappings = map[pb.Topic]proto.Message{
	pb.Topic_BEACON_BLOCK_HASH_ANNOUNCE:          &pb.BeaconBlockHashAnnounce{},
	pb.Topic_BEACON_BLOCK_REQUEST:                &pb.BeaconBlockRequest{},
	pb.Topic_BEACON_BLOCK_REQUEST_BY_SLOT_NUMBER: &pb.BeaconBlockRequestBySlotNumber{},
	pb.Topic_BEACON_BLOCK_RESPONSE:               &pb.BeaconBlockResponse{},
	pb.Topic_CRYSTALLIZED_STATE_HASH_ANNOUNCE:    &pb.CrystallizedStateHashAnnounce{},
	pb.Topic_CRYSTALLIZED_STATE_REQUEST:          &pb.CrystallizedStateRequest{},
	pb.Topic_CRYSTALLIZED_STATE_RESPONSE:         &pb.CrystallizedStateResponse{},
	pb.Topic_ACTIVE_STATE_HASH_ANNOUNCE:          &pb.ActiveStateHashAnnounce{},
	pb.Topic_ACTIVE_STATE_REQUEST:                &pb.ActiveStateRequest{},
	pb.Topic_ACTIVE_STATE_RESPONSE:               &pb.ActiveStateResponse{},
}

func configureP2P(ctx *cli.Context) (*p2p.Server, error) {
	s, err := p2p.NewServer()
	if err != nil {
		return nil, err
	}

	traceAdapter, err := tracer.New("beacon-chain",
		ctx.GlobalString(cmd.TracingEndpointFlag.Name),
		ctx.GlobalFloat64(cmd.TraceSampleFractionFlag.Name),
		ctx.GlobalBool(cmd.EnableTracingFlag.Name))
	if err != nil {
		return nil, err
	}

	// TODO(437): Define default adapters for logging, monitoring, etc.
	adapters := []p2p.Adapter{traceAdapter}
	for k, v := range topicMappings {
		s.RegisterTopic(k.String(), v, adapters...)
	}

	return s, nil
}
