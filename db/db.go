package db

import (
	"errors"
	"fmt"
	"sync"

	"github.com/tusupov/gousers/model"
)

var (
	errAmountZero      = errors.New("Amount must be more than zero")
	errAmountNotEnough = errors.New("Amount is not enough")
)

type db struct {
	sync.RWMutex
	lastId uint64
	list   map[uint64]*model.User
}

// New create user with amount
func (store *db) New(amount uint64) (model.User, error) {

	// Check amount
	if amount == 0 {
		return model.User{}, errAmountZero
	}

	store.Lock()
	defer store.Unlock()

	store.lastId++
	store.list[store.lastId] = &model.User{
		Id:     store.lastId,
		Amount: amount,
	}

	return *store.list[store.lastId], nil
}

// Get user by id
func (store *db) Get(id uint64) (model.User, bool) {
	store.RLock()
	defer store.RUnlock()

	usr, ok := store.list[id]
	if !ok {
		return model.User{}, false
	}

	return *usr, true
}

// Transfer amount from user to user
func (store *db) Transfer(fromUserId, toUserId uint64, amount uint64) error {

	// Get from user
	store.RLock()
	fromUser, ok := store.list[fromUserId]
	store.RUnlock()
	if !ok {
		return errors.New(fmt.Sprintf("User [%d] not found", fromUserId))
	}

	// Get to user
	store.RLock()
	toUser, ok := store.list[toUserId]
	store.RUnlock()
	if !ok {
		return errors.New(fmt.Sprintf("User [%d] not found", toUserId))
	}

	store.Lock()
	defer store.Unlock()

	// Check amount
	if fromUser.Amount < amount {
		return errAmountNotEnough
	}

	fromUser.Amount -= amount
	toUser.Amount += amount

	return nil
}

// default store for app
var defaultStore = NewDB()

// New user for default store
func New(amount uint64) (model.User, error) {
	return defaultStore.New(amount)
}

// New user user by id for default store
func Get(id uint64) (model.User, bool) {
	return defaultStore.Get(id)
}

// Transfer user user by id for default store
func Transfer(fromUserId, toUserId uint64, amount uint64) error {
	return defaultStore.Transfer(fromUserId, toUserId, amount)
}

func NewDB() *db {
	return &db{
		lastId: 0,
		list:   make(map[uint64]*model.User),
	}
}
