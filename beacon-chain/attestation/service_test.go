package attestation

import (
	"context"
	"testing"

	btestutil "github.com/ovcharovvladimir/Prysm/beacon-chain/testutil"
	"github.com/ovcharovvladimir/Prysm/beacon-chain/types"
	"github.com/ovcharovvladimir/Prysm/shared/testutil"
	"github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestIncomingAttestations(t *testing.T) {
	hook := logTest.NewGlobal()
	beaconDB := btestutil.SetupDB(t)
	defer btestutil.TeardownDB(t, beaconDB)
	service := NewAttestationService(context.Background(), &Config{BeaconDB: beaconDB})

	exitRoutine := make(chan bool)
	go func() {
		service.aggregateAttestations()
		<-exitRoutine
	}()

	service.incomingChan <- types.NewAttestation(nil)
	service.cancel()
	exitRoutine <- true

	testutil.AssertLogsContain(t, hook, "Forwarding aggregated attestation")
}
