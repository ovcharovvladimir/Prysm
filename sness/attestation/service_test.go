package attestation

import (
	"context"
	"testing"

	"github.com/ovcharovvladimir/Prysm/shared/testutil"
	btestutil "github.com/ovcharovvladimir/Prysm/sness/testutil"
	"github.com/ovcharovvladimir/Prysm/sness/types"
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
