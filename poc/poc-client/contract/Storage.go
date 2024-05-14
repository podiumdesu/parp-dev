// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package store

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PaychanPayChan is an auto generated low-level Go binding around an user-defined struct.
type PaychanPayChan struct {
	Id               [32]byte
	Sender           common.Address
	Recipient        common.Address
	SenderDeposit    *big.Int
	StartTime        *big.Int
	Status           *big.Int
	Fee              *big.Int
	DisputeStartTime *big.Int
	DisputeDuration  *big.Int
	SenderConfirm    bool
	RecipientConfirm bool
}

// PaychanMetaData contains all meta data concerning the Paychan contract.
var PaychanMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"ChannelOpened\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"balance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"bal\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"closeChan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"confirmClosure\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"greeting\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"senderDeposit\",\"type\":\"uint256\"}],\"name\":\"openChan\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"paychanCheck\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"addresspayable\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"senderDeposit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"disputeStartTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"disputeDuration\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"senderConfirm\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"recipientConfirm\",\"type\":\"bool\"}],\"internalType\":\"structpaychan.PayChan\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"paychanSelectedArguments\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rec\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"senderB\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PaychanABI is the input ABI used to generate the binding from.
// Deprecated: Use PaychanMetaData.ABI instead.
var PaychanABI = PaychanMetaData.ABI

// Paychan is an auto generated Go binding around an Ethereum contract.
type Paychan struct {
	PaychanCaller     // Read-only binding to the contract
	PaychanTransactor // Write-only binding to the contract
	PaychanFilterer   // Log filterer for contract events
}

// PaychanCaller is an auto generated read-only Go binding around an Ethereum contract.
type PaychanCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaychanTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PaychanTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaychanFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PaychanFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaychanSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PaychanSession struct {
	Contract     *Paychan          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PaychanCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PaychanCallerSession struct {
	Contract *PaychanCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// PaychanTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PaychanTransactorSession struct {
	Contract     *PaychanTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// PaychanRaw is an auto generated low-level Go binding around an Ethereum contract.
type PaychanRaw struct {
	Contract *Paychan // Generic contract binding to access the raw methods on
}

// PaychanCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PaychanCallerRaw struct {
	Contract *PaychanCaller // Generic read-only contract binding to access the raw methods on
}

// PaychanTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PaychanTransactorRaw struct {
	Contract *PaychanTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPaychan creates a new instance of Paychan, bound to a specific deployed contract.
func NewPaychan(address common.Address, backend bind.ContractBackend) (*Paychan, error) {
	contract, err := bindPaychan(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Paychan{PaychanCaller: PaychanCaller{contract: contract}, PaychanTransactor: PaychanTransactor{contract: contract}, PaychanFilterer: PaychanFilterer{contract: contract}}, nil
}

// NewPaychanCaller creates a new read-only instance of Paychan, bound to a specific deployed contract.
func NewPaychanCaller(address common.Address, caller bind.ContractCaller) (*PaychanCaller, error) {
	contract, err := bindPaychan(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PaychanCaller{contract: contract}, nil
}

// NewPaychanTransactor creates a new write-only instance of Paychan, bound to a specific deployed contract.
func NewPaychanTransactor(address common.Address, transactor bind.ContractTransactor) (*PaychanTransactor, error) {
	contract, err := bindPaychan(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PaychanTransactor{contract: contract}, nil
}

// NewPaychanFilterer creates a new log filterer instance of Paychan, bound to a specific deployed contract.
func NewPaychanFilterer(address common.Address, filterer bind.ContractFilterer) (*PaychanFilterer, error) {
	contract, err := bindPaychan(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PaychanFilterer{contract: contract}, nil
}

// bindPaychan binds a generic wrapper to an already deployed contract.
func bindPaychan(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PaychanMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Paychan *PaychanRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Paychan.Contract.PaychanCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Paychan *PaychanRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paychan.Contract.PaychanTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Paychan *PaychanRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Paychan.Contract.PaychanTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Paychan *PaychanCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Paychan.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Paychan *PaychanTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paychan.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Paychan *PaychanTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Paychan.Contract.contract.Transact(opts, method, params...)
}

// Balance is a free data retrieval call binding the contract method 0xe3d670d7.
//
// Solidity: function balance(address addr) view returns(uint256 bal)
func (_Paychan *PaychanCaller) Balance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Paychan.contract.Call(opts, &out, "balance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balance is a free data retrieval call binding the contract method 0xe3d670d7.
//
// Solidity: function balance(address addr) view returns(uint256 bal)
func (_Paychan *PaychanSession) Balance(addr common.Address) (*big.Int, error) {
	return _Paychan.Contract.Balance(&_Paychan.CallOpts, addr)
}

// Balance is a free data retrieval call binding the contract method 0xe3d670d7.
//
// Solidity: function balance(address addr) view returns(uint256 bal)
func (_Paychan *PaychanCallerSession) Balance(addr common.Address) (*big.Int, error) {
	return _Paychan.Contract.Balance(&_Paychan.CallOpts, addr)
}

// Greeting is a free data retrieval call binding the contract method 0xa21d05d9.
//
// Solidity: function greeting(address from) pure returns(string)
func (_Paychan *PaychanCaller) Greeting(opts *bind.CallOpts, from common.Address) (string, error) {
	var out []interface{}
	err := _Paychan.contract.Call(opts, &out, "greeting", from)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Greeting is a free data retrieval call binding the contract method 0xa21d05d9.
//
// Solidity: function greeting(address from) pure returns(string)
func (_Paychan *PaychanSession) Greeting(from common.Address) (string, error) {
	return _Paychan.Contract.Greeting(&_Paychan.CallOpts, from)
}

// Greeting is a free data retrieval call binding the contract method 0xa21d05d9.
//
// Solidity: function greeting(address from) pure returns(string)
func (_Paychan *PaychanCallerSession) Greeting(from common.Address) (string, error) {
	return _Paychan.Contract.Greeting(&_Paychan.CallOpts, from)
}

// PaychanCheck is a free data retrieval call binding the contract method 0xcfd5aabd.
//
// Solidity: function paychanCheck(bytes32 channelId) view returns((bytes32,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bool,bool))
func (_Paychan *PaychanCaller) PaychanCheck(opts *bind.CallOpts, channelId [32]byte) (PaychanPayChan, error) {
	var out []interface{}
	err := _Paychan.contract.Call(opts, &out, "paychanCheck", channelId)

	if err != nil {
		return *new(PaychanPayChan), err
	}

	out0 := *abi.ConvertType(out[0], new(PaychanPayChan)).(*PaychanPayChan)

	return out0, err

}

// PaychanCheck is a free data retrieval call binding the contract method 0xcfd5aabd.
//
// Solidity: function paychanCheck(bytes32 channelId) view returns((bytes32,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bool,bool))
func (_Paychan *PaychanSession) PaychanCheck(channelId [32]byte) (PaychanPayChan, error) {
	return _Paychan.Contract.PaychanCheck(&_Paychan.CallOpts, channelId)
}

// PaychanCheck is a free data retrieval call binding the contract method 0xcfd5aabd.
//
// Solidity: function paychanCheck(bytes32 channelId) view returns((bytes32,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bool,bool))
func (_Paychan *PaychanCallerSession) PaychanCheck(channelId [32]byte) (PaychanPayChan, error) {
	return _Paychan.Contract.PaychanCheck(&_Paychan.CallOpts, channelId)
}

// PaychanSelectedArguments is a free data retrieval call binding the contract method 0x2aefff5e.
//
// Solidity: function paychanSelectedArguments(bytes32 channelId) view returns(address sender, address rec, uint256 status, uint256 senderB, uint256 fee)
func (_Paychan *PaychanCaller) PaychanSelectedArguments(opts *bind.CallOpts, channelId [32]byte) (struct {
	Sender  common.Address
	Rec     common.Address
	Status  *big.Int
	SenderB *big.Int
	Fee     *big.Int
}, error) {
	var out []interface{}
	err := _Paychan.contract.Call(opts, &out, "paychanSelectedArguments", channelId)

	outstruct := new(struct {
		Sender  common.Address
		Rec     common.Address
		Status  *big.Int
		SenderB *big.Int
		Fee     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Sender = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Rec = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Status = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.SenderB = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Fee = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PaychanSelectedArguments is a free data retrieval call binding the contract method 0x2aefff5e.
//
// Solidity: function paychanSelectedArguments(bytes32 channelId) view returns(address sender, address rec, uint256 status, uint256 senderB, uint256 fee)
func (_Paychan *PaychanSession) PaychanSelectedArguments(channelId [32]byte) (struct {
	Sender  common.Address
	Rec     common.Address
	Status  *big.Int
	SenderB *big.Int
	Fee     *big.Int
}, error) {
	return _Paychan.Contract.PaychanSelectedArguments(&_Paychan.CallOpts, channelId)
}

// PaychanSelectedArguments is a free data retrieval call binding the contract method 0x2aefff5e.
//
// Solidity: function paychanSelectedArguments(bytes32 channelId) view returns(address sender, address rec, uint256 status, uint256 senderB, uint256 fee)
func (_Paychan *PaychanCallerSession) PaychanSelectedArguments(channelId [32]byte) (struct {
	Sender  common.Address
	Rec     common.Address
	Status  *big.Int
	SenderB *big.Int
	Fee     *big.Int
}, error) {
	return _Paychan.Contract.PaychanSelectedArguments(&_Paychan.CallOpts, channelId)
}

// CloseChan is a paid mutator transaction binding the contract method 0x91bc9212.
//
// Solidity: function closeChan(bytes32 channelId, uint256 value) returns()
func (_Paychan *PaychanTransactor) CloseChan(opts *bind.TransactOpts, channelId [32]byte, value *big.Int) (*types.Transaction, error) {
	return _Paychan.contract.Transact(opts, "closeChan", channelId, value)
}

// CloseChan is a paid mutator transaction binding the contract method 0x91bc9212.
//
// Solidity: function closeChan(bytes32 channelId, uint256 value) returns()
func (_Paychan *PaychanSession) CloseChan(channelId [32]byte, value *big.Int) (*types.Transaction, error) {
	return _Paychan.Contract.CloseChan(&_Paychan.TransactOpts, channelId, value)
}

// CloseChan is a paid mutator transaction binding the contract method 0x91bc9212.
//
// Solidity: function closeChan(bytes32 channelId, uint256 value) returns()
func (_Paychan *PaychanTransactorSession) CloseChan(channelId [32]byte, value *big.Int) (*types.Transaction, error) {
	return _Paychan.Contract.CloseChan(&_Paychan.TransactOpts, channelId, value)
}

// ConfirmClosure is a paid mutator transaction binding the contract method 0x298c88d2.
//
// Solidity: function confirmClosure(bytes32 channelId, uint256 value) payable returns()
func (_Paychan *PaychanTransactor) ConfirmClosure(opts *bind.TransactOpts, channelId [32]byte, value *big.Int) (*types.Transaction, error) {
	return _Paychan.contract.Transact(opts, "confirmClosure", channelId, value)
}

// ConfirmClosure is a paid mutator transaction binding the contract method 0x298c88d2.
//
// Solidity: function confirmClosure(bytes32 channelId, uint256 value) payable returns()
func (_Paychan *PaychanSession) ConfirmClosure(channelId [32]byte, value *big.Int) (*types.Transaction, error) {
	return _Paychan.Contract.ConfirmClosure(&_Paychan.TransactOpts, channelId, value)
}

// ConfirmClosure is a paid mutator transaction binding the contract method 0x298c88d2.
//
// Solidity: function confirmClosure(bytes32 channelId, uint256 value) payable returns()
func (_Paychan *PaychanTransactorSession) ConfirmClosure(channelId [32]byte, value *big.Int) (*types.Transaction, error) {
	return _Paychan.Contract.ConfirmClosure(&_Paychan.TransactOpts, channelId, value)
}

// OpenChan is a paid mutator transaction binding the contract method 0x63bbe7a3.
//
// Solidity: function openChan(address to, uint256 senderDeposit) payable returns(bytes32 channelId)
func (_Paychan *PaychanTransactor) OpenChan(opts *bind.TransactOpts, to common.Address, senderDeposit *big.Int) (*types.Transaction, error) {
	return _Paychan.contract.Transact(opts, "openChan", to, senderDeposit)
}

// OpenChan is a paid mutator transaction binding the contract method 0x63bbe7a3.
//
// Solidity: function openChan(address to, uint256 senderDeposit) payable returns(bytes32 channelId)
func (_Paychan *PaychanSession) OpenChan(to common.Address, senderDeposit *big.Int) (*types.Transaction, error) {
	return _Paychan.Contract.OpenChan(&_Paychan.TransactOpts, to, senderDeposit)
}

// OpenChan is a paid mutator transaction binding the contract method 0x63bbe7a3.
//
// Solidity: function openChan(address to, uint256 senderDeposit) payable returns(bytes32 channelId)
func (_Paychan *PaychanTransactorSession) OpenChan(to common.Address, senderDeposit *big.Int) (*types.Transaction, error) {
	return _Paychan.Contract.OpenChan(&_Paychan.TransactOpts, to, senderDeposit)
}

// PaychanChannelOpenedIterator is returned from FilterChannelOpened and is used to iterate over the raw logs and unpacked data for ChannelOpened events raised by the Paychan contract.
type PaychanChannelOpenedIterator struct {
	Event *PaychanChannelOpened // Event containing the contract specifics and raw log

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
func (it *PaychanChannelOpenedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaychanChannelOpened)
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
		it.Event = new(PaychanChannelOpened)
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
func (it *PaychanChannelOpenedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaychanChannelOpenedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaychanChannelOpened represents a ChannelOpened event raised by the Paychan contract.
type PaychanChannelOpened struct {
	ChannelId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelOpened is a free log retrieval operation binding the contract event 0x7ffc675d721b8768e035a323722842743bb523487b535906abc97e6b3e2095d0.
//
// Solidity: event ChannelOpened(bytes32 indexed channelId)
func (_Paychan *PaychanFilterer) FilterChannelOpened(opts *bind.FilterOpts, channelId [][32]byte) (*PaychanChannelOpenedIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _Paychan.contract.FilterLogs(opts, "ChannelOpened", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &PaychanChannelOpenedIterator{contract: _Paychan.contract, event: "ChannelOpened", logs: logs, sub: sub}, nil
}

// WatchChannelOpened is a free log subscription operation binding the contract event 0x7ffc675d721b8768e035a323722842743bb523487b535906abc97e6b3e2095d0.
//
// Solidity: event ChannelOpened(bytes32 indexed channelId)
func (_Paychan *PaychanFilterer) WatchChannelOpened(opts *bind.WatchOpts, sink chan<- *PaychanChannelOpened, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _Paychan.contract.WatchLogs(opts, "ChannelOpened", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaychanChannelOpened)
				if err := _Paychan.contract.UnpackLog(event, "ChannelOpened", log); err != nil {
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

// ParseChannelOpened is a log parse operation binding the contract event 0x7ffc675d721b8768e035a323722842743bb523487b535906abc97e6b3e2095d0.
//
// Solidity: event ChannelOpened(bytes32 indexed channelId)
func (_Paychan *PaychanFilterer) ParseChannelOpened(log types.Log) (*PaychanChannelOpened, error) {
	event := new(PaychanChannelOpened)
	if err := _Paychan.contract.UnpackLog(event, "ChannelOpened", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
