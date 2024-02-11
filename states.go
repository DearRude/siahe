package main

import (
	"context"
	"sync"

	"gorm.io/gorm"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type UserState uint

const (
	CommandState UserState = iota
	SignUpAskFirstName
	SignUpAskLastName
	SignUpAskPhoneNumber
	SignUpAskGender
)

var (
	StateMap = &UserStateMap{data: make(map[int64]UserState)}
	UserMap  = &UserDataMap{data: make(map[int64]User)}

	sender *message.Sender
	db     *gorm.DB
)

type UpdateMessage struct {
	Ctx context.Context
	Ent tg.Entities
	Unm *tg.UpdateNewMessage

	SenderID int64
	Message  *tg.Message
}

type UpdateCallback struct {
	Ctx context.Context
	Ent tg.Entities
	Ubc *tg.UpdateBotCallbackQuery

	PeerUser *tg.InputPeerUser
}

type UserStateMap struct {
	data map[int64]UserState
	mu   sync.RWMutex
}

// Set adds or updates a user state in the map
func (m *UserStateMap) Set(userID int64, state UserState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[userID] = state
}

// Get retrieves the user state associated with the given user ID
func (m *UserStateMap) Get(userID int64) (UserState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	state, ok := m.data[userID]
	return state, ok
}

// Delete removes the user state associated with the given user ID
func (m *UserStateMap) Delete(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, userID)
}

// UserDataMap is a concurrent-safe map for user data
type UserDataMap struct {
	data map[int64]User
	mu   sync.RWMutex
}

// Set adds or updates a user in the map
func (m *UserDataMap) Set(userID int64, userData User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[userID] = userData
}

// Get retrieves the user data associated with the given user ID
func (m *UserDataMap) Get(userID int64) (User, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	userData, ok := m.data[userID]
	return userData, ok
}

// Delete removes the user associated with the given user ID
func (m *UserDataMap) Delete(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, userID)
}
