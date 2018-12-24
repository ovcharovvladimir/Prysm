package contracts

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind/backends"
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	"github.com/ovcharovvladimir/essentiaHybrid/core"
	"github.com/ovcharovvladimir/essentiaHybrid/core/types"
	"github.com/ovcharovvladimir/essentiaHybrid/crypto"
	"github.com/ovcharovvladimir/Prysm/sness/params"
)

type vrcTestHelper struct {
	testAccounts []testAccount
	backend      *backends.SimulatedBackend
	vrc          *VRC
}

type testAccount struct {
	addr    common.Address
	privKey *ecdsa.PrivateKey
	txOpts  *bind.TransactOpts
}

var (
	accountBalance2000Eth, _     = new(big.Int).SetString("2000000000000000000000", 10)
	voterDepositInsufficient, _ = new(big.Int).SetString("999000000000000000000", 10)
	voterDeposit, _             = new(big.Int).SetString("1000000000000000000000", 10)
	ctx                          = context.Background()
)

// newVRCTestHelper is a helper function to initialize backend with the n test accounts,
// a vrc struct to interact with vrc, and a simulated back end to for deployment.
func newVRCTestHelper(n int) (*vrcTestHelper, error) {
	genesis := make(core.GenesisAlloc)
	var testAccounts []testAccount

	for i := 0; i < n; i++ {
		privKey, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(privKey.PublicKey)
		txOpts := bind.NewKeyedTransactor(privKey)
		testAccounts = append(testAccounts, testAccount{addr, privKey, txOpts})
		genesis[addr] = core.GenesisAccount{
			Balance: accountBalance2000Eth,
		}
	}
	backend := backends.NewSimulatedBackend(genesis)
	_, _, vrc, err := deployVRCContract(backend, testAccounts[0].privKey)
	if err != nil {
		return nil, err
	}
	return &vrcTestHelper{testAccounts: testAccounts, backend: backend, vrc: vrc}, nil
}

// deployVRCContract is a helper function for deploying VRC.
func deployVRCContract(backend *backends.SimulatedBackend, key *ecdsa.PrivateKey) (common.Address, *types.Transaction, *VRC, error) {
	transactOpts := bind.NewKeyedTransactor(key)
	defer backend.Commit()
	return DeployVRC(transactOpts, backend)
}

// fastForward is a helper function to skip through n period.
func (s *vrcTestHelper) fastForward(p int) {
	for i := 0; i < p*int(params.DefaultConfig.PeriodLength); i++ {
		s.backend.Commit()
	}
}

// registerNotaries is a helper function register notaries in batch.
func (s *vrcTestHelper) registerNotaries(deposit *big.Int, params ...int) error {
	for i := params[0]; i < params[1]; i++ {
		s.testAccounts[i].txOpts.Value = deposit
		_, err := s.vrc.RegisterNotary(s.testAccounts[i].txOpts)
		if err != nil {
			return err
		}
		s.backend.Commit()

		voter, _ := s.vrc.NotaryRegistry(&bind.CallOpts{}, s.testAccounts[i].addr)
		if !voter.Deposited ||
			voter.PoolIndex.Cmp(big.NewInt(int64(i))) != 0 ||
			voter.DeregisteredPeriod.Cmp(big.NewInt(0)) != 0 {
			return fmt.Errorf("Incorrect voter registry. Want - deposited:true, index:%v, period:0"+
				"Got - deposited:%v, index:%v, period:%v ", i, voter.Deposited, voter.PoolIndex, voter.DeregisteredPeriod)
		}
	}
	// Filter VRC logs by voterRegistered.
	log, err := s.vrc.FilterNotaryRegistered(&bind.FilterOpts{})
	if err != nil {
		return err
	}
	// Iterate voterRegistered logs, compare each address and poolIndex.
	for i := 0; i < params[1]; i++ {
		log.Next()
		if log.Event.Notary != s.testAccounts[i].addr {
			return fmt.Errorf("incorrect address in voterRegistered log. Want: %v Got: %v", s.testAccounts[i].addr, log.Event.Notary)
		}
		// Verify voterPoolIndex is incremental starting from 1st registered Notary.
		if log.Event.PoolIndex.Cmp(big.NewInt(int64(i))) != 0 {
			return fmt.Errorf("incorrect index in voterRegistered log. Want: %v Got: %v", i, log.Event.Notary)
		}
	}
	return nil
}

// deregisterNotaries is a helper function that deregister notaries in batch.
func (s *vrcTestHelper) deregisterNotaries(params ...int) error {
	for i := params[0]; i < params[1]; i++ {
		s.testAccounts[i].txOpts.Value = big.NewInt(0)
		_, err := s.vrc.DeregisterNotary(s.testAccounts[i].txOpts)
		if err != nil {
			return fmt.Errorf("Failed to deregister voter: %v", err)
		}
		s.backend.Commit()
		voter, _ := s.vrc.NotaryRegistry(&bind.CallOpts{}, s.testAccounts[i].addr)
		if voter.DeregisteredPeriod.Cmp(big.NewInt(0)) == 0 {
			return fmt.Errorf("Degistered period can not be 0 right after deregistration")
		}
	}
	// Filter VRC logs by voterDeregistered.
	log, err := s.vrc.FilterNotaryDeregistered(&bind.FilterOpts{})
	if err != nil {
		return err
	}
	// Iterate voterDeregistered logs, compare each address, poolIndex and verify period is set.
	for i := 0; i < params[1]; i++ {
		log.Next()
		if log.Event.Notary != s.testAccounts[i].addr {
			return fmt.Errorf("incorrect address in voterDeregistered log. Want: %v Got: %v", s.testAccounts[i].addr, log.Event.Notary)
		}
		if log.Event.PoolIndex.Cmp(big.NewInt(int64(i))) != 0 {
			return fmt.Errorf("incorrect index in voterDeregistered log. Want: %v Got: %v", i, log.Event.Notary)
		}
		if log.Event.DeregisteredPeriod.Cmp(big.NewInt(0)) == 0 {
			return fmt.Errorf("incorrect period in voterDeregistered log. Got: %v", log.Event.DeregisteredPeriod)
		}
	}
	return nil
}

// addHeader is a helper function to add header to vrc.
func (s *vrcTestHelper) addHeader(a *testAccount, shard *big.Int, period *big.Int, chunkRoot uint8) error {
	sig := [32]byte{'S', 'I', 'G', 'N', 'A', 'T', 'U', 'R', 'E'}
	_, err := s.vrc.AddHeader(a.txOpts, shard, period, [32]byte{chunkRoot}, sig)
	if err != nil {
		return err
	}
	s.backend.Commit()

	p, err := s.vrc.LastSubmittedCollation(&bind.CallOpts{}, shard)
	if err != nil {
		return fmt.Errorf("Can't get last submitted collation's period number: %v", err)
	}
	if p.Cmp(period) != 0 {
		return fmt.Errorf("Incorrect last period, when header was added. Got: %v", p)
	}

	cr, err := s.vrc.CollationRecords(&bind.CallOpts{}, shard, period)
	if err != nil {
		return err
	}
	if cr.ChunkRoot != [32]byte{chunkRoot} {
		return fmt.Errorf("Chunkroot mismatched. Want: %v, Got: %v", chunkRoot, cr)
	}

	// Filter VRC logs by headerAdded.
	shardIndex := []*big.Int{shard}
	logPeriod := uint64(period.Int64() * params.DefaultConfig.PeriodLength)
	log, err := s.vrc.FilterHeaderAdded(&bind.FilterOpts{Start: logPeriod}, shardIndex)
	if err != nil {
		return err
	}
	log.Next()
	if log.Event.ProposerAddress != s.testAccounts[0].addr {
		return fmt.Errorf("incorrect proposer address in addHeader log. Want: %v Got: %v", s.testAccounts[0].addr, log.Event.ProposerAddress)
	}
	if log.Event.ChunkRoot != [32]byte{chunkRoot} {
		return fmt.Errorf("chunk root missmatch in addHeader log. Want: %v Got: %v", [32]byte{chunkRoot}, log.Event.ChunkRoot)
	}
	return nil
}

// submitVote is a helper function for voter to submit vote on a given header.
func (s *vrcTestHelper) submitVote(a *testAccount, shard *big.Int, period *big.Int, index *big.Int, chunkRoot uint8) error {
	_, err := s.vrc.SubmitVote(a.txOpts, shard, period, index, [32]byte{chunkRoot})
	if err != nil {
		return fmt.Errorf("Notary submit vote failed: %v", err)
	}
	s.backend.Commit()

	v, err := s.vrc.HasVoted(&bind.CallOpts{}, shard, index)
	if err != nil {
		return fmt.Errorf("Check voter's vote failed: %v", err)
	}
	if !v {
		return fmt.Errorf("Notary's indexd bit did not cast to 1 in index %v", index)
	}
	// Filter VRC logs by submitVote.
	shardIndex := []*big.Int{shard}
	logPeriod := uint64(period.Int64() * params.DefaultConfig.PeriodLength)
	log, err := s.vrc.FilterVoteSubmitted(&bind.FilterOpts{Start: logPeriod}, shardIndex)
	if err != nil {
		return err
	}
	log.Next()
	if log.Event.NotaryAddress != a.addr {
		return fmt.Errorf("incorrect voter address in submitVote log. Want: %v Got: %v", s.testAccounts[0].addr, a.addr)
	}
	if log.Event.ChunkRoot != [32]byte{chunkRoot} {
		return fmt.Errorf("chunk root missmatch in submitVote log. Want: %v Got: %v", common.BytesToHash([]byte{chunkRoot}), common.BytesToHash(log.Event.ChunkRoot[:]))
	}
	return nil
}

// checkNotaryPoolLength is a helper function to verify current voter pool
// length is equal to n.
func checkNotaryPoolLength(vrc *VRC, n *big.Int) error {
	numNotaries, err := vrc.NotaryPoolLength(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("Failed to get voter pool length: %v", err)
	}
	if numNotaries.Cmp(n) != 0 {
		return fmt.Errorf("Incorrect count from voter pool. Want: %v, Got: %v", n, numNotaries)
	}
	return nil
}

// TestContractCreation tests VRC smart contract can successfully be deployed.
func TestContractCreation(t *testing.T) {
	_, err := newVRCTestHelper(1)
	if err != nil {
		t.Fatalf("can't deploy VRC: %v", err)
	}
}

// TestNotaryRegister tests voter registers in a normal condition.
func TestNotaryRegister(t *testing.T) {
	// Initializes 3 accounts to register as notaries.
	const voterCount = 3
	s, _ := newVRCTestHelper(voterCount)

	// Verify voter 0 has not registered.
	voter, err := s.vrc.NotaryRegistry(&bind.CallOpts{}, s.testAccounts[0].addr)
	if err != nil {
		t.Errorf("Can't get voter registry info: %v", err)
	}
	if voter.Deposited {
		t.Errorf("Notary has not registered. Got deposited flag: %v", voter.Deposited)
	}

	// Test voter 0 has registered.
	err = s.registerNotaries(voterDeposit, 0, 1)
	if err != nil {
		t.Errorf("Register voter failed: %v", err)
	}
	// Test voter 1 and 2 have registered.
	err = s.registerNotaries(voterDeposit, 1, 3)
	if err != nil {
		t.Errorf("Register voter failed: %v", err)
	}
	// Check total numbers of notaries in pool, should be 3
	err = checkNotaryPoolLength(s.vrc, big.NewInt(voterCount))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}
}

// TestNotaryRegisterInsufficientEther tests voter registers with insufficient deposit.
func TestNotaryRegisterInsufficientEther(t *testing.T) {
	s, _ := newVRCTestHelper(1)
	if err := s.registerNotaries(voterDepositInsufficient, 0, 1); err == nil {
		t.Errorf("Notary register should have failed with insufficient deposit")
	}
}

// TestNotaryDoubleRegisters tests voter registers twice.
func TestNotaryDoubleRegisters(t *testing.T) {
	s, _ := newVRCTestHelper(1)

	// Notary 0 registers.
	err := s.registerNotaries(voterDeposit, 0, 1)
	if err != nil {
		t.Errorf("Register voter failed: %v", err)
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(1))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}

	// Notary 0 registers again, This time should fail.
	if err = s.registerNotaries(big.NewInt(0), 0, 1); err == nil {
		t.Errorf("Notary register should have failed with double registers")
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(1))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}
}

// TestNotaryDeregister tests voter deregisters in a normal condition.
func TestNotaryDeregister(t *testing.T) {
	s, _ := newVRCTestHelper(1)

	// Notary 0 registers.
	err := s.registerNotaries(voterDeposit, 0, 1)
	if err != nil {
		t.Errorf("Failed to release voter: %v", err)
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(1))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}
	// Fast forward 20 periods to check voter's deregistered period field is set correctly.
	s.fastForward(20)

	// Notary 0 deregisters.
	s.deregisterNotaries(0, 1)
	err = checkNotaryPoolLength(s.vrc, big.NewInt(0))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}
}

// TestNotaryDeregisterThenRegister tests voter deregisters then registers before lock up ends.
func TestNotaryDeregisterThenRegister(t *testing.T) {
	s, _ := newVRCTestHelper(1)

	// Notary 0 registers.
	err := s.registerNotaries(voterDeposit, 0, 1)
	if err != nil {
		t.Errorf("Failed to register voter: %v", err)
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(1))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}
	s.fastForward(1)

	// Notary 0 deregisters.
	s.deregisterNotaries(0, 1)
	err = checkNotaryPoolLength(s.vrc, big.NewInt(0))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}

	// Notary 0 re-registers again.
	err = s.registerNotaries(voterDeposit, 0, 1)
	if err == nil {
		t.Error("Expected re-registration to fail")
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(0))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}
}

// TestNotaryRelease tests voter releases in a normal condition.
func TestNotaryRelease(t *testing.T) {
	s, _ := newVRCTestHelper(1)

	// Notary 0 registers.
	err := s.registerNotaries(voterDeposit, 0, 1)
	if err != nil {
		t.Errorf("Failed to register voter: %v", err)
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(1))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}
	s.fastForward(1)

	// Notary 0 deregisters.
	s.deregisterNotaries(0, 1)
	err = checkNotaryPoolLength(s.vrc, big.NewInt(0))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}

	// Fast forward until lockup ends.
	s.fastForward(int(params.DefaultConfig.NotaryLockupLength + 1))

	// Notary 0 releases.
	_, err = s.vrc.ReleaseNotary(s.testAccounts[0].txOpts)
	if err != nil {
		t.Errorf("Failed to release voter: %v", err)
	}
	s.backend.Commit()
	voter, err := s.vrc.NotaryRegistry(&bind.CallOpts{}, s.testAccounts[0].addr)
	if err != nil {
		t.Errorf("Can't get voter registry info: %v", err)
	}
	if voter.Deposited {
		t.Errorf("Notary deposit flag should be false after released")
	}
	balance, err := s.backend.BalanceAt(ctx, s.testAccounts[0].addr, nil)
	if err != nil {
		t.Errorf("Can't get account balance, err: %s", err)
	}
	if balance.Cmp(voterDeposit) < 0 {
		t.Errorf("Notary did not receive deposit after lock up ends")
	}
}

// TestNotaryInstantRelease tests voter releases before lockup ends.
func TestNotaryInstantRelease(t *testing.T) {
	s, _ := newVRCTestHelper(1)

	// Notary 0 registers.
	if err := s.registerNotaries(voterDeposit, 0, 1); err != nil {
		t.Error(err)
	}
	if err := checkNotaryPoolLength(s.vrc, big.NewInt(1)); err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}
	s.fastForward(1)

	// Notary 0 deregisters.
	s.deregisterNotaries(0, 1)
	if err := checkNotaryPoolLength(s.vrc, big.NewInt(0)); err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}

	// Notary 0 tries to release before lockup ends.
	if _, err := s.vrc.ReleaseNotary(s.testAccounts[0].txOpts); err == nil {
		t.Error("Expected release voter to fail")
	}
	s.backend.Commit()
	voter, err := s.vrc.NotaryRegistry(&bind.CallOpts{}, s.testAccounts[0].addr)
	if err != nil {
		t.Errorf("Can't get voter registry info: %v", err)
	}
	if !voter.Deposited {
		t.Errorf("Notary deposit flag should be true before released")
	}
	balance, err := s.backend.BalanceAt(ctx, s.testAccounts[0].addr, nil)
	if err != nil {
		t.Error(err)
	}
	if balance.Cmp(voterDeposit) > 0 {
		t.Errorf("Notary received deposit before lockup ends")
	}
}

// TestCommitteeListsAreDifferent tests different shards have different voter committee.
func TestCommitteeListsAreDifferent(t *testing.T) {
	const voterCount = 1000
	s, _ := newVRCTestHelper(voterCount)

	// Register 1000 notaries to s.vrc.
	err := s.registerNotaries(voterDeposit, 0, 1000)
	if err != nil {
		t.Errorf("Failed to release voter: %v", err)
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(1000))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}

	// Compare sampled first 5 notaries of shard 0 to shard 1, they should not be identical.
	for i := 0; i < 5; i++ {
		addr0, _ := s.vrc.GetNotaryInCommittee(&bind.CallOpts{}, big.NewInt(0))
		addr1, _ := s.vrc.GetNotaryInCommittee(&bind.CallOpts{}, big.NewInt(1))
		if addr0 == addr1 {
			t.Errorf("Shard 0 committee list is identical to shard 1's committee list")
		}
	}
}

// TestGetCommitteeWithNonMember tests unregistered voter tries to be in the committee.
func TestGetCommitteeWithNonMember(t *testing.T) {
	const voterCount = 11
	s, _ := newVRCTestHelper(voterCount)

	// Register 10 notaries to s.vrc, leave 1 address free.
	err := s.registerNotaries(voterDeposit, 0, 10)
	if err != nil {
		t.Errorf("Failed to release voter: %v", err)
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(10))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}

	// Verify the unregistered account is not in the voter pool list.
	for i := 0; i < 10; i++ {
		addr, _ := s.vrc.GetNotaryInCommittee(&bind.CallOpts{}, big.NewInt(0))
		if s.testAccounts[10].addr == addr {
			t.Errorf("Account %s is not a voter", s.testAccounts[10].addr.String())
		}
	}
}

// TestGetCommitteeWithinSamePeriod tests voter registers and samples within the same period.
func TestGetCommitteeWithinSamePeriod(t *testing.T) {
	s, _ := newVRCTestHelper(1)

	// Notary 0 registers.
	err := s.registerNotaries(voterDeposit, 0, 1)
	if err != nil {
		t.Errorf("Failed to release voter: %v", err)
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(1))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}

	// Notary 0 samples for itself within the same period after registration.
	sampledAddr, _ := s.vrc.GetNotaryInCommittee(&bind.CallOpts{}, big.NewInt(0))
	if s.testAccounts[0].addr != sampledAddr {
		t.Errorf("Unable to sample voter address within same period of registration, got addr: %v", sampledAddr)
	}
}

// TestGetCommitteeAfterDeregisters tests voter tries to be in committee after deregistered.
func TestGetCommitteeAfterDeregisters(t *testing.T) {
	const voterCount = 10
	s, _ := newVRCTestHelper(voterCount)

	// Register 10 notaries to s.vrc.
	err := s.registerNotaries(voterDeposit, 0, voterCount)
	if err != nil {
		t.Errorf("Failed to release voter: %v", err)
	}
	err = checkNotaryPoolLength(s.vrc, big.NewInt(10))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}

	// Deregister voter 0 from s.vrc.
	s.deregisterNotaries(0, 1)
	err = checkNotaryPoolLength(s.vrc, big.NewInt(9))
	if err != nil {
		t.Errorf("Notary pool length mismatched: %v", err)
	}

	// Verify degistered voter 0 is not in the voter pool list.
	for i := 0; i < 10; i++ {
		addr, _ := s.vrc.GetNotaryInCommittee(&bind.CallOpts{}, big.NewInt(0))
		if s.testAccounts[0].addr == addr {
			t.Errorf("Account %s is not a voter", s.testAccounts[0].addr.String())
		}
	}
}

// TestNormalAddHeader tests proposer add header in normal condition.
func TestNormalAddHeader(t *testing.T) {
	s, _ := newVRCTestHelper(1)
	s.fastForward(1)

	// Proposer adds header consists shard 0, period 1 and chunkroot 0xA.
	err := s.addHeader(&s.testAccounts[0], big.NewInt(0), big.NewInt(1), 'A')
	if err != nil {
		t.Errorf("Proposer adds header failed: %v", err)
	}
	s.fastForward(1)

	// Proposer adds header consists shard 0, period 2 and chunkroot 0xB.
	err = s.addHeader(&s.testAccounts[0], big.NewInt(0), big.NewInt(2), 'B')
	if err != nil {
		t.Errorf("Proposer adds header failed: %v", err)
	}
	// Proposer adds header consists shard 1, period 2 and chunkroot 0xC.
	err = s.addHeader(&s.testAccounts[0], big.NewInt(1), big.NewInt(2), 'C')
	if err != nil {
		t.Errorf("Proposer adds header failed: %v", err)
	}
}

// TestAddTwoHeadersAtSamePeriod tests we can't add two headers within the same period.
func TestAddTwoHeadersAtSamePeriod(t *testing.T) {
	s, _ := newVRCTestHelper(1)
	s.fastForward(1)

	// Proposer adds header consists shard 0, period 1 and chunkroot 0xA.
	err := s.addHeader(&s.testAccounts[0], big.NewInt(0), big.NewInt(1), 'A')
	if err != nil {
		t.Errorf("Proposer adds header failed: %v", err)
	}

	// Proposer attempts to add another header chunkroot 0xB on the same period for the same shard.
	err = s.addHeader(&s.testAccounts[0], big.NewInt(0), big.NewInt(1), 'B')
	if err == nil {
		t.Errorf("Proposer is not allowed to add 2 headers within same period")
	}
}

// TestAddHeadersAtWrongPeriod tests proposer adds header in the wrong period.
func TestAddHeadersAtWrongPeriod(t *testing.T) {
	s, _ := newVRCTestHelper(1)
	s.fastForward(1)

	// Proposer adds header at wrong period, shard 0, period 0 and chunkroot 0xA.
	err := s.addHeader(&s.testAccounts[0], big.NewInt(0), big.NewInt(0), 'A')
	if err == nil {
		t.Errorf("Proposer adds header at wrong period should have failed")
	}
	// Proposer adds header at wrong period, shard 0, period 2 and chunkroot 0xA.
	err = s.addHeader(&s.testAccounts[0], big.NewInt(0), big.NewInt(2), 'A')
	if err == nil {
		t.Errorf("Proposer adds header at wrong period should have failed")
	}
	// Proposer adds header at correct period, shard 0, period 1 and chunkroot 0xA.
	err = s.addHeader(&s.testAccounts[0], big.NewInt(0), big.NewInt(1), 'A')
	if err != nil {
		t.Errorf("Proposer adds header failed: %v", err)
	}
}

// TestSubmitVote tests voter submit votes in normal condition.
func TestSubmitVote(t *testing.T) {
	s, _ := newVRCTestHelper(1)
	// Notary 0 registers.
	if err := s.registerNotaries(voterDeposit, 0, 1); err != nil {
		t.Error(err)
	}
	s.fastForward(1)

	// Proposer adds header consists shard 0, period 1 and chunkroot 0xA.
	period1 := big.NewInt(1)
	shard0 := big.NewInt(0)
	index0 := big.NewInt(0)
	s.testAccounts[0].txOpts.Value = big.NewInt(0)
	if err := s.addHeader(&s.testAccounts[0], shard0, period1, 'A'); err != nil {
		t.Errorf("Proposer adds header failed: %v", err)
	}

	// Notary 0 votes on header.
	c, err := s.vrc.GetVoteCount(&bind.CallOpts{}, shard0)
	if err != nil {
		t.Errorf("Get voter vote count failed: %v", err)
	}
	if c.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Incorrect voter vote count, want: 0, got: %v", c)
	}

	// Notary votes on the header that was submitted.
	err = s.submitVote(&s.testAccounts[0], shard0, period1, index0, 'A')
	if err != nil {
		t.Fatalf("Notary submits vote failed: %v", err)
	}
	c, err = s.vrc.GetVoteCount(&bind.CallOpts{}, shard0)
	if err != nil {
		t.Error(err)
	}
	if c.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Incorrect voter vote count, want: 1, got: %v", c)
	}

	// Check header's approved with the current period, should be period 0.
	p, err := s.vrc.LastApprovedCollation(&bind.CallOpts{}, shard0)
	if err != nil {
		t.Fatalf("Get period of last approved header failed: %v", err)
	}
	if p.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Incorrect period submitted, want: 0, got: %v", p)
	}

}

// TestSubmitVoteTwice tests voter tries to submit same vote twice.
func TestSubmitVoteTwice(t *testing.T) {
	s, _ := newVRCTestHelper(1)
	// Notary 0 registers.
	if err := s.registerNotaries(voterDeposit, 0, 1); err != nil {
		t.Error(err)
	}
	s.fastForward(1)

	// Proposer adds header consists shard 0, period 1 and chunkroot 0xA.
	period1 := big.NewInt(1)
	shard0 := big.NewInt(0)
	index0 := big.NewInt(0)
	s.testAccounts[0].txOpts.Value = big.NewInt(0)
	if err := s.addHeader(&s.testAccounts[0], shard0, period1, 'A'); err != nil {
		t.Errorf("Proposer adds header failed: %v", err)
	}

	// Notary 0 votes on header.
	if err := s.submitVote(&s.testAccounts[0], shard0, period1, index0, 'A'); err != nil {
		t.Errorf("Notary submits vote failed: %v", err)
	}

	// Notary 0 votes on header again, it should fail.
	if err := s.submitVote(&s.testAccounts[0], shard0, period1, index0, 'A'); err == nil {
		t.Errorf("voter voting twice should have failed")
	}

	// Check voter's vote count is correct in shard.
	c, _ := s.vrc.GetVoteCount(&bind.CallOpts{}, shard0)
	if c.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Incorrect voter vote count, want: 1, got: %v", c)
	}
}

// TestSubmitVoteByNonEligibleNotary tests a non-eligible voter tries to submit vote.
func TestSubmitVoteByNonEligibleNotary(t *testing.T) {
	s, _ := newVRCTestHelper(1)
	s.fastForward(1)

	// Proposer adds header consists shard 0, period 1 and chunkroot 0xA.
	period1 := big.NewInt(1)
	shard0 := big.NewInt(0)
	index0 := big.NewInt(0)
	err := s.addHeader(&s.testAccounts[0], shard0, period1, 'A')
	if err != nil {
		t.Errorf("Proposer adds header failed: %v", err)
	}

	// Unregistered Notary 0 votes on header, it should fail.
	err = s.submitVote(&s.testAccounts[0], shard0, period1, index0, 'A')
	if err == nil {
		t.Errorf("Non registered voter submits vote should have failed")
	}

	// Check voter's vote count is correct in shard.
	c, _ := s.vrc.GetVoteCount(&bind.CallOpts{}, shard0)
	if c.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Incorrect voter vote count, want: 0, got: %v", c)
	}
}

// TestSubmitVoteWithOutAHeader tests a voter tries to submit vote before header gets added.
func TestSubmitVoteWithOutAHeader(t *testing.T) {
	s, _ := newVRCTestHelper(1)
	// Notary 0 registers.
	if err := s.registerNotaries(voterDeposit, 0, 1); err != nil {
		t.Error(err)
	}
	s.fastForward(1)

	// Proposer adds header consists shard 0, period 1 and chunkroot 0xA.
	period1 := big.NewInt(1)
	shard0 := big.NewInt(0)
	index0 := big.NewInt(0)
	s.testAccounts[0].txOpts.Value = big.NewInt(0)

	// Notary 0 votes on header, it should fail because no header has added.
	if err := s.submitVote(&s.testAccounts[0], shard0, period1, index0, 'A'); err == nil {
		t.Errorf("Notary votes should have failed due to missing header")
	}

	// Check voter's vote count is correct in shard
	c, _ := s.vrc.GetVoteCount(&bind.CallOpts{}, shard0)
	if c.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Incorrect voter vote count, want: 1, got: %v", c)
	}
}

// TestSubmitVoteWithInvalidArgs tests voter submits vote using wrong chunkroot and period.
func TestSubmitVoteWithInvalidArgs(t *testing.T) {
	s, _ := newVRCTestHelper(1)
	// Notary 0 registers.
	if err := s.registerNotaries(voterDeposit, 0, 1); err != nil {
		t.Error(err)
	}
	s.fastForward(1)

	// Proposer adds header consists shard 0, period 1 and chunkroot 0xA.
	period1 := big.NewInt(1)
	shard0 := big.NewInt(0)
	index0 := big.NewInt(0)
	s.testAccounts[0].txOpts.Value = big.NewInt(0)
	if err := s.addHeader(&s.testAccounts[0], shard0, period1, 'A'); err != nil {
		t.Errorf("Proposer adds header failed: %v", err)
	}

	// Notary voting with incorrect period.
	period2 := big.NewInt(2)
	if err := s.submitVote(&s.testAccounts[0], shard0, period2, index0, 'A'); err == nil {
		t.Errorf("Notary votes should have failed due to incorrect period")
	}

	// Notary voting with incorrect chunk root.
	if err := s.submitVote(&s.testAccounts[0], shard0, period2, index0, 'B'); err == nil {
		t.Errorf("Notary votes should have failed due to incorrect chunk root")
	}
}
