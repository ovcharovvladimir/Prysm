// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	//	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ovcharovvladimir/essentiaHybrid"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi"
	"github.com/ovcharovvladimir/essentiaHybrid/accounts/abi/bind"
	"github.com/ovcharovvladimir/essentiaHybrid/common"
	"github.com/ovcharovvladimir/essentiaHybrid/core/types"
	"github.com/ovcharovvladimir/essentiaHybrid/event"
)

// VRCABI is the input ABI used to generate the binding from.
const VRCABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"shardCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"currentVote\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_shardId\",\"type\":\"uint256\"},{\"name\":\"_period\",\"type\":\"uint256\"},{\"name\":\"_index\",\"type\":\"uint256\"},{\"name\":\"_chunkRoot\",\"type\":\"bytes32\"}],\"name\":\"submitVote\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deregisterNotary\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_shardId\",\"type\":\"uint256\"},{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"hasVoted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_shardId\",\"type\":\"uint256\"}],\"name\":\"getNotaryInCommittee\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"registerNotary\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"notaryRegistry\",\"outputs\":[{\"name\":\"deregisteredPeriod\",\"type\":\"uint256\"},{\"name\":\"poolIndex\",\"type\":\"uint256\"},{\"name\":\"balance\",\"type\":\"uint256\"},{\"name\":\"deposited\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastSubmittedCollation\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastApprovedCollation\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"releaseNotary\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"notaryPool\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_shardId\",\"type\":\"uint256\"}],\"name\":\"getVoteCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CHALLENGE_PERIOD\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_shardId\",\"type\":\"uint256\"},{\"name\":\"_period\",\"type\":\"uint256\"},{\"name\":\"_chunkRoot\",\"type\":\"bytes32\"},{\"name\":\"_signature\",\"type\":\"bytes32\"}],\"name\":\"addHeader\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"collationRecords\",\"outputs\":[{\"name\":\"chunkRoot\",\"type\":\"bytes32\"},{\"name\":\"proposer\",\"type\":\"address\"},{\"name\":\"isElected\",\"type\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"notaryPoolLength\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"shardId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"chunkRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"period\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"proposerAddress\",\"type\":\"address\"}],\"name\":\"HeaderAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"notary\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"poolIndex\",\"type\":\"uint256\"}],\"name\":\"NotaryRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"notary\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"poolIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"deregisteredPeriod\",\"type\":\"uint256\"}],\"name\":\"NotaryDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"notary\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"poolIndex\",\"type\":\"uint256\"}],\"name\":\"NotaryReleased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"shardId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"chunkRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"period\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"notaryAddress\",\"type\":\"address\"}],\"name\":\"VoteSubmitted\",\"type\":\"event\"}]"

// VRCBin is the compiled bytecode used for deploying new contracts.
const VRCBin = `0x60806040526064600c5534801561001557600080fd5b50610d58806100256000396000f3006080604052600436106100f05763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166304e9c77a81146100f55780630c8da4cc1461011c5780634f33ffa01461013457806358377bd11461015757806364390ff11461016c578063673221af1461019b57806368e9513e146101cf5780636bdd3271146101d757806383ceeabe1461022057806397d369a2146102385780639910851d14610250578063a81f451014610265578063b2c2f2e81461027d578063c3a079ed14610295578063c4d5f198146102aa578063e9e0b683146102cb578063f6f67d3614610314575b600080fd5b34801561010157600080fd5b5061010a610329565b60408051918252519081900360200190f35b34801561012857600080fd5b5061010a60043561032f565b34801561014057600080fd5b50610155600435602435604435606435610341565b005b34801561016357600080fd5b506101556104e7565b34801561017857600080fd5b50610187600435602435610612565b604080519115158252519081900360200190f35b3480156101a757600080fd5b506101b3600435610635565b60408051600160a060020a039092168252519081900360200190f35b6101556106ee565b3480156101e357600080fd5b506101f8600160a060020a03600435166108b1565b6040805194855260208501939093528383019190915215156060830152519081900360800190f35b34801561022c57600080fd5b5061010a6004356108dc565b34801561024457600080fd5b5061010a6004356108ee565b34801561025c57600080fd5b50610155610900565b34801561027157600080fd5b506101b3600435610a35565b34801561028957600080fd5b5061010a600435610a5d565b3480156102a157600080fd5b5061010a610a72565b3480156102b657600080fd5b50610155600435602435604435606435610a77565b3480156102d757600080fd5b506102e6600435602435610bd6565b60408051948552600160a060020a039093166020850152901515838301526060830152519081900360800190f35b34801561032057600080fd5b5061010a610c29565b600c5481565b60036020526000908152604090205481565b60008085101580156103545750600c5485105b151561035f57600080fd5b60054304841461036e57600080fd5b600085815260056020526040902054841461038857600080fd5b6087831061039557600080fd5b600085815260046020908152604080832087845290915290205482146103ba57600080fd5b600160a060020a03331660009081526001602052604090206003015460ff1615156103e457600080fd5b6103ee8584610612565b156103f857600080fd5b33600160a060020a031661040b86610635565b600160a060020a03161461041e57600080fd5b6104288584610c2f565b61043185610a5d565b9050605a8110610495576000858152600660209081526040808320879055600482528083208784529091529020600101805474ff00000000000000000000000000000000000000001916740100000000000000000000000000000000000000001790555b6040805183815260208101869052600160a060020a03331681830152905186917fc99370212b708f699fb6945a17eb34d0fc1ccd5b45d88f4d9682593a45d6e833919081900360600190a25050505050565b33600160a060020a038116600090815260016020819052604082209081015460039091015490919060ff16151561051d57600080fd5b82600160a060020a031660008381548110151561053657fe5b600091825260209091200154600160a060020a03161461055557600080fd5b61055d610c53565b50600160a060020a0382166000908152600160205260409020600543049081905561058782610c76565b600080548390811061059557fe5b600091825260209182902001805473ffffffffffffffffffffffffffffffffffffffff191690556002805460001901905560408051600160a060020a0386168152918201849052818101839052517f90e5afdc8fd31453dcf6e37154fa117ddf3b0324c96c65015563df9d5e4b5a759181900360600190a1505050565b60009182526003602052604090912054600160ff9290920360020a900481161490565b6000600543048180808080610648610c53565b600b5486111561065c57600a549450610662565b60095494505b600160a060020a03331660009081526001602081815260409283902090910154825160001960058b020180408083529382018390528185018d90529351908190036060019020909650919450925085908115156106bb57fe5b0690506000818154811015156106cd57fe5b600091825260209091200154600160a060020a031698975050505050505050565b33600160a060020a03811660009081526001602052604081206003015460ff161561071857600080fd5b34683635c9adc5dea000001461072d57600080fd5b610735610c53565b61073d610ce7565b156107a05750600254600080546001810182559080527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e56301805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0384161790556107e9565b6107a8610cee565b9050816000828154811015156107ba57fe5b9060005260206000200160006101000a815481600160a060020a030219169083600160a060020a031602179055505b600280546001908101825560408051608081018252600080825260208083018781523484860190815260608501878152600160a060020a038b168552928790529490922092518355905193820193909355905192810192909255516003909101805460ff1916911515919091179055600a5481106108695760018101600a555b60408051600160a060020a03841681526020810183905281517fa4fe15c53db34d35a5117acc26c27a2653dc68e2dadfc21ed211e38b7864d7a7929181900390910190a15050565b6001602081905260009182526040909120805491810154600282015460039092015490919060ff1684565b60056020526000908152604090205481565b60066020526000908152604090205481565b33600160a060020a038116600090815260016020819052604082208082015460039091015490929160ff90911615151461093957600080fd5b600160a060020a038316600090815260016020526040902054151561095d57600080fd5b600160a060020a038316600090815260016020526040902054613f0001600543041161098857600080fd5b50600160a060020a0382166000818152600160208190526040808320600281018054858355938201859055849055600301805460ff191690555190929183156108fc02918491818181858888f193505050501580156109eb573d6000803e3d6000fd5b5060408051600160a060020a03851681526020810184905281517faee20171b64b7f3360a142659094ce929970d6963dcea8c34a9bf1ece8033680929181900390910190a1505050565b6000805482908110610a4357fe5b600091825260209091200154600160a060020a0316905081565b60009081526003602052604090205460ff1690565b601981565b60008410158015610a895750600c5484105b1515610a9457600080fd5b600543048314610aa357600080fd5b6000848152600560205260409020548311610abd57600080fd5b610ac5610c53565b60408051608081018252838152600160a060020a033381166020808401828152600085870181815260608088018a81528d8452600486528984208d8552865289842098518955935160018901805493511515740100000000000000000000000000000000000000000274ff0000000000000000000000000000000000000000199290991673ffffffffffffffffffffffffffffffffffffffff19909416939093171696909617905590516002909501949094558884526005808252858520439190910490556003815284842093909355835186815292830187905282840152915186927f2d0a86178d2fd307b47be157a766e6bee19bc26161c32f9781ee0e818636f09c928290030190a250505050565b60046020908152600092835260408084209091529082529020805460018201546002909201549091600160a060020a038116917401000000000000000000000000000000000000000090910460ff169084565b60025481565b600091825260036020526040909120805460ff9290920360020a9091176001019055565b600b546005430490811015610c6757610c73565b600a54600955600b8190555b50565b6008546007541415610cbc57600780546001810182556000919091527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c68801819055610cdb565b806007600854815481101515610cce57fe5b6000918252602090912001555b50600880546001019055565b6008541590565b60006001600854111515610d0157600080fd5b600880546000190190819055600780549091908110610d1c57fe5b90600052602060002001549050905600a165627a7a7230582073009df4f7ec65cf53714d6f4b6b07ed4a0804bda664a944c998e5d97bbe1b520029`

// DeployVRC deploys a new Ethereum contract, binding an instance of VRC to it.
func DeployVRC(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRC, error) {
	parsed, err := abi.JSON(strings.NewReader(VRCABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRCBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRC{VRCCaller: VRCCaller{contract: contract}, VRCTransactor: VRCTransactor{contract: contract}, VRCFilterer: VRCFilterer{contract: contract}}, nil
}

// VRC is an auto generated Go binding around an Ethereum contract.
type VRC struct {
	VRCCaller     // Read-only binding to the contract
	VRCTransactor // Write-only binding to the contract
	VRCFilterer   // Log filterer for contract events
}

// VRCCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRCSession struct {
	Contract     *VRC              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRCCallerSession struct {
	Contract *VRCCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VRCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRCTransactorSession struct {
	Contract     *VRCTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRCRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRCRaw struct {
	Contract *VRC // Generic contract binding to access the raw methods on
}

// VRCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRCCallerRaw struct {
	Contract *VRCCaller // Generic read-only contract binding to access the raw methods on
}

// VRCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRCTransactorRaw struct {
	Contract *VRCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRC creates a new instance of VRC, bound to a specific deployed contract.
func NewVRC(address common.Address, backend bind.ContractBackend) (*VRC, error) {
	contract, err := bindVRC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRC{VRCCaller: VRCCaller{contract: contract}, VRCTransactor: VRCTransactor{contract: contract}, VRCFilterer: VRCFilterer{contract: contract}}, nil
}

// NewVRCCaller creates a new read-only instance of VRC, bound to a specific deployed contract.
func NewVRCCaller(address common.Address, caller bind.ContractCaller) (*VRCCaller, error) {
	contract, err := bindVRC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRCCaller{contract: contract}, nil
}

// NewVRCTransactor creates a new write-only instance of VRC, bound to a specific deployed contract.
func NewVRCTransactor(address common.Address, transactor bind.ContractTransactor) (*VRCTransactor, error) {
	contract, err := bindVRC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRCTransactor{contract: contract}, nil
}

// NewVRCFilterer creates a new log filterer instance of VRC, bound to a specific deployed contract.
func NewVRCFilterer(address common.Address, filterer bind.ContractFilterer) (*VRCFilterer, error) {
	contract, err := bindVRC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRCFilterer{contract: contract}, nil
}

// bindVRC binds a generic wrapper to an already deployed contract.
func bindVRC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRC *VRCRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRC.Contract.VRCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRC *VRCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRC.Contract.VRCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRC *VRCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRC.Contract.VRCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRC *VRCCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRC *VRCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRC *VRCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRC.Contract.contract.Transact(opts, method, params...)
}

// CHALLENGEPERIOD is a free data retrieval call binding the contract method 0xc3a079ed.
//
// Solidity: function CHALLENGE_PERIOD() constant returns(uint256)
func (_VRC *VRCCaller) CHALLENGEPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "CHALLENGE_PERIOD")
	return *ret0, err
}

// CHALLENGEPERIOD is a free data retrieval call binding the contract method 0xc3a079ed.
//
// Solidity: function CHALLENGE_PERIOD() constant returns(uint256)
func (_VRC *VRCSession) CHALLENGEPERIOD() (*big.Int, error) {
	return _VRC.Contract.CHALLENGEPERIOD(&_VRC.CallOpts)
}

// CHALLENGEPERIOD is a free data retrieval call binding the contract method 0xc3a079ed.
//
// Solidity: function CHALLENGE_PERIOD() constant returns(uint256)
func (_VRC *VRCCallerSession) CHALLENGEPERIOD() (*big.Int, error) {
	return _VRC.Contract.CHALLENGEPERIOD(&_VRC.CallOpts)
}

// CollationRecords is a free data retrieval call binding the contract method 0xe9e0b683.
//
// Solidity: function collationRecords( uint256,  uint256) constant returns(chunkRoot bytes32, proposer address, isElected bool, signature bytes32)
func (_VRC *VRCCaller) CollationRecords(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	ChunkRoot [32]byte
	Proposer  common.Address
	IsElected bool
	Signature [32]byte
}, error) {
	ret := new(struct {
		ChunkRoot [32]byte
		Proposer  common.Address
		IsElected bool
		Signature [32]byte
	})
	out := ret
	err := _VRC.contract.Call(opts, out, "collationRecords", arg0, arg1)
	return *ret, err
}

// CollationRecords is a free data retrieval call binding the contract method 0xe9e0b683.
//
// Solidity: function collationRecords( uint256,  uint256) constant returns(chunkRoot bytes32, proposer address, isElected bool, signature bytes32)
func (_VRC *VRCSession) CollationRecords(arg0 *big.Int, arg1 *big.Int) (struct {
	ChunkRoot [32]byte
	Proposer  common.Address
	IsElected bool
	Signature [32]byte
}, error) {
	return _VRC.Contract.CollationRecords(&_VRC.CallOpts, arg0, arg1)
}

// CollationRecords is a free data retrieval call binding the contract method 0xe9e0b683.
//
// Solidity: function collationRecords( uint256,  uint256) constant returns(chunkRoot bytes32, proposer address, isElected bool, signature bytes32)
func (_VRC *VRCCallerSession) CollationRecords(arg0 *big.Int, arg1 *big.Int) (struct {
	ChunkRoot [32]byte
	Proposer  common.Address
	IsElected bool
	Signature [32]byte
}, error) {
	return _VRC.Contract.CollationRecords(&_VRC.CallOpts, arg0, arg1)
}

// CurrentVote is a free data retrieval call binding the contract method 0x0c8da4cc.
//
// Solidity: function currentVote( uint256) constant returns(bytes32)
func (_VRC *VRCCaller) CurrentVote(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "currentVote", arg0)
	return *ret0, err
}

// CurrentVote is a free data retrieval call binding the contract method 0x0c8da4cc.
//
// Solidity: function currentVote( uint256) constant returns(bytes32)
func (_VRC *VRCSession) CurrentVote(arg0 *big.Int) ([32]byte, error) {
	return _VRC.Contract.CurrentVote(&_VRC.CallOpts, arg0)
}

// CurrentVote is a free data retrieval call binding the contract method 0x0c8da4cc.
//
// Solidity: function currentVote( uint256) constant returns(bytes32)
func (_VRC *VRCCallerSession) CurrentVote(arg0 *big.Int) ([32]byte, error) {
	return _VRC.Contract.CurrentVote(&_VRC.CallOpts, arg0)
}

// GetNotaryInCommittee is a free data retrieval call binding the contract method 0x673221af.
//
// Solidity: function getNotaryInCommittee(_shardId uint256) constant returns(address)
func (_VRC *VRCCaller) GetNotaryInCommittee(opts *bind.CallOpts, _shardId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "getNotaryInCommittee", _shardId)
	return *ret0, err
}

// GetNotaryInCommittee is a free data retrieval call binding the contract method 0x673221af.
//
// Solidity: function getNotaryInCommittee(_shardId uint256) constant returns(address)
func (_VRC *VRCSession) GetNotaryInCommittee(_shardId *big.Int) (common.Address, error) {
	return _VRC.Contract.GetNotaryInCommittee(&_VRC.CallOpts, _shardId)
}

// GetNotaryInCommittee is a free data retrieval call binding the contract method 0x673221af.
//
// Solidity: function getNotaryInCommittee(_shardId uint256) constant returns(address)
func (_VRC *VRCCallerSession) GetNotaryInCommittee(_shardId *big.Int) (common.Address, error) {
	return _VRC.Contract.GetNotaryInCommittee(&_VRC.CallOpts, _shardId)
}

// GetVoteCount is a free data retrieval call binding the contract method 0xb2c2f2e8.
//
// Solidity: function getVoteCount(_shardId uint256) constant returns(uint256)
func (_VRC *VRCCaller) GetVoteCount(opts *bind.CallOpts, _shardId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "getVoteCount", _shardId)
	return *ret0, err
}

// GetVoteCount is a free data retrieval call binding the contract method 0xb2c2f2e8.
//
// Solidity: function getVoteCount(_shardId uint256) constant returns(uint256)
func (_VRC *VRCSession) GetVoteCount(_shardId *big.Int) (*big.Int, error) {
	return _VRC.Contract.GetVoteCount(&_VRC.CallOpts, _shardId)
}

// GetVoteCount is a free data retrieval call binding the contract method 0xb2c2f2e8.
//
// Solidity: function getVoteCount(_shardId uint256) constant returns(uint256)
func (_VRC *VRCCallerSession) GetVoteCount(_shardId *big.Int) (*big.Int, error) {
	return _VRC.Contract.GetVoteCount(&_VRC.CallOpts, _shardId)
}

// HasVoted is a free data retrieval call binding the contract method 0x64390ff1.
//
// Solidity: function hasVoted(_shardId uint256, _index uint256) constant returns(bool)
func (_VRC *VRCCaller) HasVoted(opts *bind.CallOpts, _shardId *big.Int, _index *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "hasVoted", _shardId, _index)
	return *ret0, err
}

// HasVoted is a free data retrieval call binding the contract method 0x64390ff1.
//
// Solidity: function hasVoted(_shardId uint256, _index uint256) constant returns(bool)
func (_VRC *VRCSession) HasVoted(_shardId *big.Int, _index *big.Int) (bool, error) {
	return _VRC.Contract.HasVoted(&_VRC.CallOpts, _shardId, _index)
}

// HasVoted is a free data retrieval call binding the contract method 0x64390ff1.
//
// Solidity: function hasVoted(_shardId uint256, _index uint256) constant returns(bool)
func (_VRC *VRCCallerSession) HasVoted(_shardId *big.Int, _index *big.Int) (bool, error) {
	return _VRC.Contract.HasVoted(&_VRC.CallOpts, _shardId, _index)
}

// LastApprovedCollation is a free data retrieval call binding the contract method 0x97d369a2.
//
// Solidity: function lastApprovedCollation( uint256) constant returns(uint256)
func (_VRC *VRCCaller) LastApprovedCollation(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "lastApprovedCollation", arg0)
	return *ret0, err
}

// LastApprovedCollation is a free data retrieval call binding the contract method 0x97d369a2.
//
// Solidity: function lastApprovedCollation( uint256) constant returns(uint256)
func (_VRC *VRCSession) LastApprovedCollation(arg0 *big.Int) (*big.Int, error) {
	return _VRC.Contract.LastApprovedCollation(&_VRC.CallOpts, arg0)
}

// LastApprovedCollation is a free data retrieval call binding the contract method 0x97d369a2.
//
// Solidity: function lastApprovedCollation( uint256) constant returns(uint256)
func (_VRC *VRCCallerSession) LastApprovedCollation(arg0 *big.Int) (*big.Int, error) {
	return _VRC.Contract.LastApprovedCollation(&_VRC.CallOpts, arg0)
}

// LastSubmittedCollation is a free data retrieval call binding the contract method 0x83ceeabe.
//
// Solidity: function lastSubmittedCollation( uint256) constant returns(uint256)
func (_VRC *VRCCaller) LastSubmittedCollation(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "lastSubmittedCollation", arg0)
	return *ret0, err
}

// LastSubmittedCollation is a free data retrieval call binding the contract method 0x83ceeabe.
//
// Solidity: function lastSubmittedCollation( uint256) constant returns(uint256)
func (_VRC *VRCSession) LastSubmittedCollation(arg0 *big.Int) (*big.Int, error) {
	return _VRC.Contract.LastSubmittedCollation(&_VRC.CallOpts, arg0)
}

// LastSubmittedCollation is a free data retrieval call binding the contract method 0x83ceeabe.
//
// Solidity: function lastSubmittedCollation( uint256) constant returns(uint256)
func (_VRC *VRCCallerSession) LastSubmittedCollation(arg0 *big.Int) (*big.Int, error) {
	return _VRC.Contract.LastSubmittedCollation(&_VRC.CallOpts, arg0)
}

// NotaryPool is a free data retrieval call binding the contract method 0xa81f4510.
//
// Solidity: function voterPool( uint256) constant returns(address)
func (_VRC *VRCCaller) NotaryPool(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "notaryPool", arg0)
	return *ret0, err
}

// NotaryPool is a free data retrieval call binding the contract method 0xa81f4510.
//
// Solidity: function voterPool( uint256) constant returns(address)
func (_VRC *VRCSession) NotaryPool(arg0 *big.Int) (common.Address, error) {
	return _VRC.Contract.NotaryPool(&_VRC.CallOpts, arg0)
}

// NotaryPool is a free data retrieval call binding the contract method 0xa81f4510.
//
// Solidity: function voterPool( uint256) constant returns(address)
func (_VRC *VRCCallerSession) NotaryPool(arg0 *big.Int) (common.Address, error) {
	return _VRC.Contract.NotaryPool(&_VRC.CallOpts, arg0)
}

// NotaryPoolLength is a free data retrieval call binding the contract method 0xf6f67d36.
//
// Solidity: function voterPoolLength() constant returns(uint256)
func (_VRC *VRCCaller) NotaryPoolLength(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "notaryPoolLength")
	return *ret0, err
}

// NotaryPoolLength is a free data retrieval call binding the contract method 0xf6f67d36.
//
// Solidity: function voterPoolLength() constant returns(uint256)
func (_VRC *VRCSession) NotaryPoolLength() (*big.Int, error) {
	return _VRC.Contract.NotaryPoolLength(&_VRC.CallOpts)
}

// NotaryPoolLength is a free data retrieval call binding the contract method 0xf6f67d36.
//
// Solidity: function voterPoolLength() constant returns(uint256)
func (_VRC *VRCCallerSession) NotaryPoolLength() (*big.Int, error) {
	return _VRC.Contract.NotaryPoolLength(&_VRC.CallOpts)
}

// NotaryRegistry is a free data retrieval call binding the contract method 0x6bdd3271.
//
// Solidity: function voterRegistry( address) constant returns(deregisteredPeriod uint256, poolIndex uint256, balance uint256, deposited bool)
func (_VRC *VRCCaller) NotaryRegistry(opts *bind.CallOpts, arg0 common.Address) (struct {
	DeregisteredPeriod *big.Int
	PoolIndex          *big.Int
	Balance            *big.Int
	Deposited          bool
}, error) {
	ret := new(struct {
		DeregisteredPeriod *big.Int
		PoolIndex          *big.Int
		Balance            *big.Int
		Deposited          bool
	})
	out := ret
	err := _VRC.contract.Call(opts, out, "notaryRegistry", arg0)
	return *ret, err
}

// NotaryRegistry is a free data retrieval call binding the contract method 0x6bdd3271.
//
// Solidity: function voterRegistry( address) constant returns(deregisteredPeriod uint256, poolIndex uint256, balance uint256, deposited bool)
func (_VRC *VRCSession) NotaryRegistry(arg0 common.Address) (struct {
	DeregisteredPeriod *big.Int
	PoolIndex          *big.Int
	Balance            *big.Int
	Deposited          bool
}, error) {
	return _VRC.Contract.NotaryRegistry(&_VRC.CallOpts, arg0)
}

// NotaryRegistry is a free data retrieval call binding the contract method 0x6bdd3271.
//
// Solidity: function voterRegistry( address) constant returns(deregisteredPeriod uint256, poolIndex uint256, balance uint256, deposited bool)
func (_VRC *VRCCallerSession) NotaryRegistry(arg0 common.Address) (struct {
	DeregisteredPeriod *big.Int
	PoolIndex          *big.Int
	Balance            *big.Int
	Deposited          bool
}, error) {
	return _VRC.Contract.NotaryRegistry(&_VRC.CallOpts, arg0)
}

// ShardCount is a free data retrieval call binding the contract method 0x04e9c77a.
//
// Solidity: function shardCount() constant returns(uint256)
func (_VRC *VRCCaller) ShardCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRC.contract.Call(opts, out, "shardCount")
	return *ret0, err
}

// ShardCount is a free data retrieval call binding the contract method 0x04e9c77a.
//
// Solidity: function shardCount() constant returns(uint256)
func (_VRC *VRCSession) ShardCount() (*big.Int, error) {
	return _VRC.Contract.ShardCount(&_VRC.CallOpts)
}

// ShardCount is a free data retrieval call binding the contract method 0x04e9c77a.
//
// Solidity: function shardCount() constant returns(uint256)
func (_VRC *VRCCallerSession) ShardCount() (*big.Int, error) {
	return _VRC.Contract.ShardCount(&_VRC.CallOpts)
}

// AddHeader is a paid mutator transaction binding the contract method 0xc4d5f198.
//
// Solidity: function addHeader(_shardId uint256, _period uint256, _chunkRoot bytes32, _signature bytes32) returns()
func (_VRC *VRCTransactor) AddHeader(opts *bind.TransactOpts, _shardId *big.Int, _period *big.Int, _chunkRoot [32]byte, _signature [32]byte) (*types.Transaction, error) {
	return _VRC.contract.Transact(opts, "addHeader", _shardId, _period, _chunkRoot, _signature)
}

// AddHeader is a paid mutator transaction binding the contract method 0xc4d5f198.
//
// Solidity: function addHeader(_shardId uint256, _period uint256, _chunkRoot bytes32, _signature bytes32) returns()
func (_VRC *VRCSession) AddHeader(_shardId *big.Int, _period *big.Int, _chunkRoot [32]byte, _signature [32]byte) (*types.Transaction, error) {
	return _VRC.Contract.AddHeader(&_VRC.TransactOpts, _shardId, _period, _chunkRoot, _signature)
}

// AddHeader is a paid mutator transaction binding the contract method 0xc4d5f198.
//
// Solidity: function addHeader(_shardId uint256, _period uint256, _chunkRoot bytes32, _signature bytes32) returns()
func (_VRC *VRCTransactorSession) AddHeader(_shardId *big.Int, _period *big.Int, _chunkRoot [32]byte, _signature [32]byte) (*types.Transaction, error) {
	return _VRC.Contract.AddHeader(&_VRC.TransactOpts, _shardId, _period, _chunkRoot, _signature)
}

// DeregisterNotary is a paid mutator transaction binding the contract method 0x58377bd1.
//
// Solidity: function deregisterNotary() returns()
func (_VRC *VRCTransactor) DeregisterNotary(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRC.contract.Transact(opts, "deregisterNotary")
}

// DeregisterNotary is a paid mutator transaction binding the contract method 0x58377bd1.
//
// Solidity: function deregisterNotary() returns()
func (_VRC *VRCSession) DeregisterNotary() (*types.Transaction, error) {
	return _VRC.Contract.DeregisterNotary(&_VRC.TransactOpts)
}

// DeregisterNotary is a paid mutator transaction binding the contract method 0x58377bd1.
//
// Solidity: function deregisterNotary() returns()
func (_VRC *VRCTransactorSession) DeregisterNotary() (*types.Transaction, error) {
	return _VRC.Contract.DeregisterNotary(&_VRC.TransactOpts)
}

// RegisterNotary is a paid mutator transaction binding the contract method 0x68e9513e.
//
// Solidity: function registerNotary() returns()
func (_VRC *VRCTransactor) RegisterNotary(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRC.contract.Transact(opts, "registerNotary")
}

// RegisterNotary is a paid mutator transaction binding the contract method 0x68e9513e.
//
// Solidity: function registerNotary() returns()
func (_VRC *VRCSession) RegisterNotary() (*types.Transaction, error) {
	return _VRC.Contract.RegisterNotary(&_VRC.TransactOpts)
}

// RegisterNotary is a paid mutator transaction binding the contract method 0x68e9513e.
//
// Solidity: function registerNotary() returns()
func (_VRC *VRCTransactorSession) RegisterNotary() (*types.Transaction, error) {
	return _VRC.Contract.RegisterNotary(&_VRC.TransactOpts)
}

// ReleaseNotary is a paid mutator transaction binding the contract method 0x9910851d.
//
// Solidity: function releaseNotary() returns()
func (_VRC *VRCTransactor) ReleaseNotary(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRC.contract.Transact(opts, "releaseNotary")
}

// ReleaseNotary is a paid mutator transaction binding the contract method 0x9910851d.
//
// Solidity: function releaseNotary() returns()
func (_VRC *VRCSession) ReleaseNotary() (*types.Transaction, error) {
	return _VRC.Contract.ReleaseNotary(&_VRC.TransactOpts)
}

// ReleaseNotary is a paid mutator transaction binding the contract method 0x9910851d.
//
// Solidity: function releaseNotary() returns()
func (_VRC *VRCTransactorSession) ReleaseNotary() (*types.Transaction, error) {
	return _VRC.Contract.ReleaseNotary(&_VRC.TransactOpts)
}

// SubmitVote is a paid mutator transaction binding the contract method 0x4f33ffa0.
//
// Solidity: function submitVote(_shardId uint256, _period uint256, _index uint256, _chunkRoot bytes32) returns()
func (_VRC *VRCTransactor) SubmitVote(opts *bind.TransactOpts, _shardId *big.Int, _period *big.Int, _index *big.Int, _chunkRoot [32]byte) (*types.Transaction, error) {
	return _VRC.contract.Transact(opts, "submitVote", _shardId, _period, _index, _chunkRoot)
}

// SubmitVote is a paid mutator transaction binding the contract method 0x4f33ffa0.
//
// Solidity: function submitVote(_shardId uint256, _period uint256, _index uint256, _chunkRoot bytes32) returns()
func (_VRC *VRCSession) SubmitVote(_shardId *big.Int, _period *big.Int, _index *big.Int, _chunkRoot [32]byte) (*types.Transaction, error) {
	return _VRC.Contract.SubmitVote(&_VRC.TransactOpts, _shardId, _period, _index, _chunkRoot)
}

// SubmitVote is a paid mutator transaction binding the contract method 0x4f33ffa0.
//
// Solidity: function submitVote(_shardId uint256, _period uint256, _index uint256, _chunkRoot bytes32) returns()
func (_VRC *VRCTransactorSession) SubmitVote(_shardId *big.Int, _period *big.Int, _index *big.Int, _chunkRoot [32]byte) (*types.Transaction, error) {
	return _VRC.Contract.SubmitVote(&_VRC.TransactOpts, _shardId, _period, _index, _chunkRoot)
}

// VRCHeaderAddedIterator is returned from FilterHeaderAdded and is used to iterate over the raw logs and unpacked data for HeaderAdded events raised by the VRC contract.
type VRCHeaderAddedIterator struct {
	Event *VRCHeaderAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRCHeaderAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRCHeaderAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRCHeaderAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRCHeaderAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRCHeaderAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRCHeaderAdded represents a HeaderAdded event raised by the VRC contract.
type VRCHeaderAdded struct {
	ShardId         *big.Int
	ChunkRoot       [32]byte
	Period          *big.Int
	ProposerAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterHeaderAdded is a free log retrieval operation binding the contract event 0x2d0a86178d2fd307b47be157a766e6bee19bc26161c32f9781ee0e818636f09c.
//
// Solidity: event HeaderAdded(shardId indexed uint256, chunkRoot bytes32, period uint256, proposerAddress address)
func (_VRC *VRCFilterer) FilterHeaderAdded(opts *bind.FilterOpts, shardId []*big.Int) (*VRCHeaderAddedIterator, error) {

	var shardIdRule []interface{}
	for _, shardIdItem := range shardId {
		shardIdRule = append(shardIdRule, shardIdItem)
	}

	logs, sub, err := _VRC.contract.FilterLogs(opts, "HeaderAdded", shardIdRule)
	if err != nil {
		return nil, err
	}
	return &VRCHeaderAddedIterator{contract: _VRC.contract, event: "HeaderAdded", logs: logs, sub: sub}, nil
}

// WatchHeaderAdded is a free log subscription operation binding the contract event 0x2d0a86178d2fd307b47be157a766e6bee19bc26161c32f9781ee0e818636f09c.
//
// Solidity: event HeaderAdded(shardId indexed uint256, chunkRoot bytes32, period uint256, proposerAddress address)
func (_VRC *VRCFilterer) WatchHeaderAdded(opts *bind.WatchOpts, sink chan<- *VRCHeaderAdded, shardId []*big.Int) (event.Subscription, error) {

	var shardIdRule []interface{}
	for _, shardIdItem := range shardId {
		shardIdRule = append(shardIdRule, shardIdItem)
	}

	logs, sub, err := _VRC.contract.WatchLogs(opts, "HeaderAdded", shardIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRCHeaderAdded)
				if err := _VRC.contract.UnpackLog(event, "HeaderAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// VRCNotaryDeregisteredIterator is returned from FilterNotaryDeregistered and is used to iterate over the raw logs and unpacked data for NotaryDeregistered events raised by the VRC contract.
type VRCNotaryDeregisteredIterator struct {
	Event *VRCNotaryDeregistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRCNotaryDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRCNotaryDeregistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRCNotaryDeregistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRCNotaryDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRCNotaryDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRCNotaryDeregistered represents a NotaryDeregistered event raised by the VRC contract.
type VRCNotaryDeregistered struct {
	Notary             common.Address
	PoolIndex          *big.Int
	DeregisteredPeriod *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterNotaryDeregistered is a free log retrieval operation binding the contract event 0x90e5afdc8fd31453dcf6e37154fa117ddf3b0324c96c65015563df9d5e4b5a75.
//
// Solidity: event NotaryDeregistered(voter address, poolIndex uint256, deregisteredPeriod uint256)
func (_VRC *VRCFilterer) FilterNotaryDeregistered(opts *bind.FilterOpts) (*VRCNotaryDeregisteredIterator, error) {

	logs, sub, err := _VRC.contract.FilterLogs(opts, "NotaryDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRCNotaryDeregisteredIterator{contract: _VRC.contract, event: "NotaryDeregistered", logs: logs, sub: sub}, nil
}

// WatchNotaryDeregistered is a free log subscription operation binding the contract event 0x90e5afdc8fd31453dcf6e37154fa117ddf3b0324c96c65015563df9d5e4b5a75.
//
// Solidity: event NotaryDeregistered(voter address, poolIndex uint256, deregisteredPeriod uint256)
func (_VRC *VRCFilterer) WatchNotaryDeregistered(opts *bind.WatchOpts, sink chan<- *VRCNotaryDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRC.contract.WatchLogs(opts, "NotaryDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRCNotaryDeregistered)
				if err := _VRC.contract.UnpackLog(event, "NotaryDeregistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// VRCNotaryRegisteredIterator is returned from FilterNotaryRegistered and is used to iterate over the raw logs and unpacked data for NotaryRegistered events raised by the VRC contract.
type VRCNotaryRegisteredIterator struct {
	Event *VRCNotaryRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRCNotaryRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRCNotaryRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRCNotaryRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRCNotaryRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRCNotaryRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRCNotaryRegistered represents a NotaryRegistered event raised by the VRC contract.
type VRCNotaryRegistered struct {
	Notary    common.Address
	PoolIndex *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNotaryRegistered is a free log retrieval operation binding the contract event 0xa4fe15c53db34d35a5117acc26c27a2653dc68e2dadfc21ed211e38b7864d7a7.
//
// Solidity: event NotaryRegistered(voter address, poolIndex uint256)
func (_VRC *VRCFilterer) FilterNotaryRegistered(opts *bind.FilterOpts) (*VRCNotaryRegisteredIterator, error) {

	logs, sub, err := _VRC.contract.FilterLogs(opts, "NotaryRegistered")
	if err != nil {
		return nil, err
	}
	return &VRCNotaryRegisteredIterator{contract: _VRC.contract, event: "NotaryRegistered", logs: logs, sub: sub}, nil
}

// WatchNotaryRegistered is a free log subscription operation binding the contract event 0xa4fe15c53db34d35a5117acc26c27a2653dc68e2dadfc21ed211e38b7864d7a7.
//
// Solidity: event NotaryRegistered(voter address, poolIndex uint256)
func (_VRC *VRCFilterer) WatchNotaryRegistered(opts *bind.WatchOpts, sink chan<- *VRCNotaryRegistered) (event.Subscription, error) {

	logs, sub, err := _VRC.contract.WatchLogs(opts, "NotaryRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRCNotaryRegistered)
				if err := _VRC.contract.UnpackLog(event, "NotaryRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// VRCNotaryReleasedIterator is returned from FilterNotaryReleased and is used to iterate over the raw logs and unpacked data for NotaryReleased events raised by the VRC contract.
type VRCNotaryReleasedIterator struct {
	Event *VRCNotaryReleased // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRCNotaryReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRCNotaryReleased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRCNotaryReleased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRCNotaryReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRCNotaryReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRCNotaryReleased represents a NotaryReleased event raised by the VRC contract.
type VRCNotaryReleased struct {
	Notary    common.Address
	PoolIndex *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNotaryReleased is a free log retrieval operation binding the contract event 0xaee20171b64b7f3360a142659094ce929970d6963dcea8c34a9bf1ece8033680.
//
// Solidity: event NotaryReleased(voter address, poolIndex uint256)
func (_VRC *VRCFilterer) FilterNotaryReleased(opts *bind.FilterOpts) (*VRCNotaryReleasedIterator, error) {

	logs, sub, err := _VRC.contract.FilterLogs(opts, "NotaryReleased")
	if err != nil {
		return nil, err
	}
	return &VRCNotaryReleasedIterator{contract: _VRC.contract, event: "NotaryReleased", logs: logs, sub: sub}, nil
}

// WatchNotaryReleased is a free log subscription operation binding the contract event 0xaee20171b64b7f3360a142659094ce929970d6963dcea8c34a9bf1ece8033680.
//
// Solidity: event NotaryReleased(voter address, poolIndex uint256)
func (_VRC *VRCFilterer) WatchNotaryReleased(opts *bind.WatchOpts, sink chan<- *VRCNotaryReleased) (event.Subscription, error) {

	logs, sub, err := _VRC.contract.WatchLogs(opts, "NotaryReleased")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRCNotaryReleased)
				if err := _VRC.contract.UnpackLog(event, "NotaryReleased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// VRCVoteSubmittedIterator is returned from FilterVoteSubmitted and is used to iterate over the raw logs and unpacked data for VoteSubmitted events raised by the VRC contract.
type VRCVoteSubmittedIterator struct {
	Event *VRCVoteSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRCVoteSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRCVoteSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRCVoteSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRCVoteSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRCVoteSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRCVoteSubmitted represents a VoteSubmitted event raised by the VRC contract.
type VRCVoteSubmitted struct {
	ShardId       *big.Int
	ChunkRoot     [32]byte
	Period        *big.Int
	NotaryAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterVoteSubmitted is a free log retrieval operation binding the contract event 0xc99370212b708f699fb6945a17eb34d0fc1ccd5b45d88f4d9682593a45d6e833.
//
// Solidity: event VoteSubmitted(shardId indexed uint256, chunkRoot bytes32, period uint256, voterAddress address)
func (_VRC *VRCFilterer) FilterVoteSubmitted(opts *bind.FilterOpts, shardId []*big.Int) (*VRCVoteSubmittedIterator, error) {

	var shardIdRule []interface{}
	for _, shardIdItem := range shardId {
		shardIdRule = append(shardIdRule, shardIdItem)
	}

	logs, sub, err := _VRC.contract.FilterLogs(opts, "VoteSubmitted", shardIdRule)
	if err != nil {
		return nil, err
	}
	return &VRCVoteSubmittedIterator{contract: _VRC.contract, event: "VoteSubmitted", logs: logs, sub: sub}, nil
}

// WatchVoteSubmitted is a free log subscription operation binding the contract event 0xc99370212b708f699fb6945a17eb34d0fc1ccd5b45d88f4d9682593a45d6e833.
//
// Solidity: event VoteSubmitted(shardId indexed uint256, chunkRoot bytes32, period uint256, voterAddress address)
func (_VRC *VRCFilterer) WatchVoteSubmitted(opts *bind.WatchOpts, sink chan<- *VRCVoteSubmitted, shardId []*big.Int) (event.Subscription, error) {

	var shardIdRule []interface{}
	for _, shardIdItem := range shardId {
		shardIdRule = append(shardIdRule, shardIdItem)
	}

	logs, sub, err := _VRC.contract.WatchLogs(opts, "VoteSubmitted", shardIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRCVoteSubmitted)
				if err := _VRC.contract.UnpackLog(event, "VoteSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
