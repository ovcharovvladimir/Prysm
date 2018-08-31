package casper

import (
	pb "github.com/ovcharovvladimir/Prysm/proto/beacon/p2p/v1"
	   "github.com/ovcharovvladimir/Prysm/beacon-chain/params"
	   "github.com/ovcharovvladimir/Prysm/beacon-chain/utils"
   	   "github.com/sirupsen/logrus"
)
var log = logrus.WithField("prefix", "casper")

// 
// CalculateRewards adjusts validators 
// balances by applying rewards or penalties
// based on FFG incentive structure.
// 
// Changed : AS
// Date    : 31.09.2018 16:20
// Title   : Calculate Rewards
// Note    : 
// Note    : Changed coefficient
// Descr.  : 1.25 = 25% (attester factor)
//           3.75 = 75% (total factor)
// TODO    : Waiting test
// 
func CalculateRewards(attestations []*pb.AttestationRecord, validators []*pb.ValidatorRecord, dynasty uint64, totalDeposit uint64) ([]*pb.ValidatorRecord, error) {
	
	activeValidators := ActiveValidatorIndices(validators, dynasty)
	attesterDeposits := GetAttestersTotalDeposit(attestations)

     // Changed
     // Replace parameters 
     // Сoefficient attester 
	attesterFactor   := attesterDeposits    * 1.25      // Old = 3
	totalFactor      := uint64(totalDeposit * 3.75)     // Old = 2 


    // Chek conditionl
    // Old  conditional (<=)
    // TODO: Need test
	if attesterFactor >= totalFactor {
	   log.Debug("Applying rewards or penalties for the validators from last cycle.")
	   
	   for i, attesterIndex := range activeValidators {
		          voted := utils.CheckBit(attestations[len(attestations)-1].AttesterBitfield, int(attesterIndex))
		
			  if voted {
			     validators[i].Balance += params.AttesterReward
			  }  else {
			     validators[i].Balance -= params.AttesterReward
			  }
		}
	}

	return validators, nil
}
