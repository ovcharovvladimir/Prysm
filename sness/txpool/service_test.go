package txpool

import "github.com/ovcharovvladimir/Prysm/shared"

// Verifies that TXPool implements the Service interface.
var _ = shared.Service(&TXPool{})
