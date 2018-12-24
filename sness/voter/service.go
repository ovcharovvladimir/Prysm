// Package voter defines all relevant functionality for a Notary actor
// within a sharded Ethereum blockchain.
package voter

import (
	"github.com/ovcharovvladimir/Prysm/sness/database"
	"github.com/ovcharovvladimir/Prysm/sness/mainchain"
	"github.com/ovcharovvladimir/Prysm/sness/params"
	"github.com/ovcharovvladimir/Prysm/shared/p2p"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "voter")

// Notary holds functionality required to run a collation voter
// in a sharded system. Must satisfy the Service interface defined in
// sharding/service.go.
type Notary struct {
	config    *params.Config
	vrcClient *mainchain.VRCClient
	p2p       *p2p.Server
	dbService *database.ShardDB
}

// NewNotary creates a new voter instance.
func NewNotary(config *params.Config, vrcClient *mainchain.VRCClient, p2p *p2p.Server, dbService *database.ShardDB) (*Notary, error) {
	return &Notary{config, vrcClient, p2p, dbService}, nil
}

// Start the main routine for a voter.
func (n *Notary) Start() {
	log.Info("Starting voter service")
	go n.notarizeCollations()
}

// Stop the main loop for notarizing collations.
func (n *Notary) Stop() error {
	log.Info("Stopping voter service")
	return nil
}

// notarizeCollations checks incoming block headers and determines if
// we are an eligible voter for collations.
func (n *Notary) notarizeCollations() {

	// TODO: handle this better through goroutines. Right now, these methods
	// are blocking.
	if n.vrcClient.DepositFlag() {
		if err := joinNotaryPool(n.vrcClient, n.vrcClient); err != nil {
			log.Errorf("Could not fetch current block number: %v", err)
			return
		}
	}

	if err := subscribeBlockHeaders(n.vrcClient.ChainReader(), n.vrcClient, n.vrcClient.Account()); err != nil {
		log.Errorf("Could not fetch current block number: %v", err)
		return
	}
}
