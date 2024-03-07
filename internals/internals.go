package internals

import (
	"context"
	"sync"

	"github.com/gotd/td/tg"

	db "github.com/DearRude/siahe/database"
)

type UserState uint

const (
	CommandState UserState = iota
	SignUpAskFirstName
	SignUpAskLastName
	SignUpAskPhoneNumber
	SignUpAskGender
	SignUpAskIsFumStudent
	SignUpAskStudentNumber
	SignUpAskFumFaculty
	SignUpAskIsStudent
	SignUpAskIsMashhadStudent
	SignUpAskIsMastPhd
	SignUpAskUniversityName
	SignUpAskEntraceYear
	SignUpAskStudentMajor
	SignUpAskIsGraduate
	SignUpAskIsStudentRelative
	SignUpCheckInfo
	GetTicketInit
	GetTicketCount
	GetTicketPayment
)

type UpdateMessage struct {
	Ctx context.Context
	Ent tg.Entities
	Unm *tg.UpdateNewMessage

	PeerUser *tg.InputPeerUser
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

// Constructor
func NewUserStateMap() UserStateMap {
	return UserStateMap{data: make(map[int64]UserState)}
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
	data map[int64]db.User
	mu   sync.RWMutex
}

// Constructor
func NewUserDataMap() UserDataMap {
	return UserDataMap{data: make(map[int64]db.User)}
}

// Set adds or updates a user in the map
func (m *UserDataMap) Set(userID int64, userData db.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[userID] = userData
}

// Get retrieves the user data associated with the given user ID
func (m *UserDataMap) Get(userID int64) (db.User, bool) {
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

// EventDataMap is a concurrent-safe map for event data
type EventDataMap struct {
	data map[int64]db.Event
	mu   sync.RWMutex
}

// Constructor
func NewEventDataMap() EventDataMap {
	return EventDataMap{data: make(map[int64]db.Event)}
}

// Set adds or updates a event in the map
func (m *EventDataMap) Set(eventID int64, eventData db.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[eventID] = eventData
}

// Get retrieves the event data associated with the given event ID
func (m *EventDataMap) Get(eventID int64) (db.Event, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	eventData, ok := m.data[eventID]
	return eventData, ok
}

// Delete removes the event associated with the given event ID
func (m *EventDataMap) Delete(eventID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, eventID)
}
