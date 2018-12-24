package voter

import (
	"math/big"
	"testing"

	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/Prysm/sness/internal"
	shardparams "github.com/ovcharovvladimir/Prysm/sness/params"
	"github.com/ovcharovvladimir/Prysm/sness/types"
)

// Verifies that Notary implements the Actor interface.
var _ = types.Actor(&Notary{})

func TestHasAccountBeenDeregistered(t *testing.T) {
	backend, vrc := internal.SetupMockClient(t)
	client := &internal.MockClient{VRC: vrc, T: t, Backend: backend, BlockNumber: 1}

	client.SetDepositFlag(true)
	err := joinNotaryPool(client, client)
	if err != nil {
		t.Error(err)
	}
}

func TestIsAccountInNotaryPool(t *testing.T) {
	backend, vrc := internal.SetupMockClient(t)
	client := &internal.MockClient{VRC: vrc, T: t, Backend: backend}

	// address should not be in pool initially.
	b, err := isAccountInNotaryPool(client, client.Account())
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal("account unexpectedly in voter pool")
	}

	txOpts, _ := client.CreateTXOpts(shardparams.DefaultConfig.NotaryDeposit)
	if _, err := vrc.RegisterNotary(txOpts); err != nil {
		t.Fatalf("Failed to deposit: %v", err)
	}
	client.CommitWithBlock()
	b, err = isAccountInNotaryPool(client, client.Account())
	if err != nil {
		t.Error(err)
	}
	if !b {
		t.Error("account not in voter pool when expected to be")
	}
}

func TestJoinNotaryPool(t *testing.T) {
	backend, vrc := internal.SetupMockClient(t)
	client := &internal.MockClient{VRC: vrc, T: t, Backend: backend}

	// There should be no voter initially.
	numNotaries, err := vrc.NotaryPoolLength(&bind.CallOpts{})
	if err != nil {
		t.Error(err)
	}
	if big.NewInt(0).Cmp(numNotaries) != 0 {
		t.Errorf("unexpected number of notaries. Got %d, wanted 0.", numNotaries)
	}

	client.SetDepositFlag(false)
	err = joinNotaryPool(client, client)
	if err == nil {
		t.Error("joined voter pool while --deposit was not present")
	}

	client.SetDepositFlag(true)
	err = joinNotaryPool(client, client)
	if err != nil {
		t.Error(err)
	}

	// Now there should be one voter.
	numNotaries, err = vrc.NotaryPoolLength(&bind.CallOpts{})
	if err != nil {
		t.Error(err)
	}
	if big.NewInt(1).Cmp(numNotaries) != 0 {
		t.Errorf("unexpected number of notaries. Got %d, wanted 1", numNotaries)
	}

	// Join while deposited should do nothing.
	err = joinNotaryPool(client, client)
	if err != nil {
		t.Error(err)
	}

	numNotaries, err = vrc.NotaryPoolLength(&bind.CallOpts{})
	if err != nil {
		t.Error(err)
	}
	if big.NewInt(1).Cmp(numNotaries) != 0 {
		t.Errorf("unexpected number of notaries. Got %d, wanted 1", numNotaries)
	}
}
