package state

import (
	"fmt"

	acm "github.com/hyperledger/burrow/account"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
)

type MemoryState struct {
	Accounts map[crypto.Address]acm.Account
	Storage  map[crypto.Address]map[binary.Word256]binary.Word256
}

var _ IterableWriter = &MemoryState{}

// Get an in-memory state Iterable
func NewMemoryState() *MemoryState {
	return &MemoryState{
		Accounts: make(map[crypto.Address]acm.Account),
		Storage:  make(map[crypto.Address]map[binary.Word256]binary.Word256),
	}
}

func (ms *MemoryState) GetAccount(address crypto.Address) (acm.Account, error) {
	return ms.Accounts[address], nil
}

func (ms *MemoryState) UpdateAccount(updatedAccount acm.Account) error {
	if updatedAccount == nil {
		return fmt.Errorf("UpdateAccount passed nil account in MemoryState")
	}
	ms.Accounts[updatedAccount.Address()] = updatedAccount
	return nil
}

func (ms *MemoryState) RemoveAccount(address crypto.Address) error {
	delete(ms.Accounts, address)
	return nil
}

func (ms *MemoryState) GetStorage(address crypto.Address, key binary.Word256) (binary.Word256, error) {
	storage, ok := ms.Storage[address]
	if !ok {
		return binary.Zero256, fmt.Errorf("could not find storage for account %s", address)
	}
	value, ok := storage[key]
	if !ok {
		return binary.Zero256, fmt.Errorf("could not find key %x for account %s", key, address)
	}
	return value, nil
}

func (ms *MemoryState) SetStorage(address crypto.Address, key, value binary.Word256) error {
	storage, ok := ms.Storage[address]
	if !ok {
		storage = make(map[binary.Word256]binary.Word256)
		ms.Storage[address] = storage
	}
	storage[key] = value
	return nil
}

func (ms *MemoryState) IterateAccounts(consumer func(acm.Account) (stop bool)) (stopped bool, err error) {
	for _, acc := range ms.Accounts {
		if consumer(acc) {
			return true, nil
		}
	}
	return false, nil
}

func (ms *MemoryState) IterateStorage(address crypto.Address, consumer func(key, value binary.Word256) (stop bool)) (stopped bool, err error) {
	for key, value := range ms.Storage[address] {
		if consumer(key, value) {
			return true, nil
		}
	}
	return false, nil
}
