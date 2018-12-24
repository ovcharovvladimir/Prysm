package mainchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ovcharovvladimir/Prysm/sness/contracts"
	"github.com/ovcharovvladimir/Prysm/sness/params"
	"github.com/ovcharovvladimir/essentiaHybrid/node"
	"github.com/ovcharovvladimir/essentiaHybrid/rpc"
)

// dialRPC endpoint to node.
func dialRPC(endpoint string) (*rpc.Client, error) {
	if endpoint == "" {
		endpoint = node.DefaultIPCEndpoint(ClientIdentifier)
	}
	return rpc.Dial(endpoint)
}

// initVRC initializes the sharding manager contract bindings.
// If the VRC does not exist, it will be deployed.
func initVRC(s *VRCClient) (*contracts.VRC, error) {
	b, err := s.client.CodeAt(context.Background(), params.DefaultConfig.VRCAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get contract code at %s: %v", params.DefaultConfig.VRCAddress.Hex(), err)
	}

	// Deploy VRC for development only.
	// TODO: Separate contract deployment from the sharding node. It would only need to be deployed
	// once on the mainnet, so this code would not need to ship with the node.
	if len(b) == 0 {
		log.Infof("No sharding manager contract found at %s, deploying new contract", params.DefaultConfig.VRCAddress.Hex())

		txOps, err := s.CreateTXOpts(big.NewInt(0))
		if err != nil {
			return nil, fmt.Errorf("unable to initiate the transaction: %v", err)
		}

		addr, tx, contract, err := contracts.DeployVRC(txOps, s.client)
		if err != nil {
			return nil, fmt.Errorf("unable to deploy sharding manager contract: %v", err)
		}

		for pending := true; pending; _, pending, err = s.client.TransactionByHash(context.Background(), tx.Hash()) {
			if err != nil {
				return nil, fmt.Errorf("unable to get transaction by hash: %v", err)
			}
			time.Sleep(1 * time.Second)
		}

		log.Infof("New contract deployed at %s", addr.Hex())
		return contract, nil
	}

	return contracts.NewVRC(params.DefaultConfig.VRCAddress, s.client)
}
