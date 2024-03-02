package database

import (
	"time"
)

type User struct {
	ID          int64 `gorm:"primaryKey"`
	FirstName   string
	LastName    string
	Role        string
	PhoneNumber string
	IsBoy       bool

	IsFumStudent  bool
	StudentNumber string
	FumFaculty    string

	IsStudent        bool
	IsMashhadStudent bool

	UniversityName string

	EntranceYear string
	IsMastPhd    bool
	StudentMajor string

	IsGraduateStudent bool
	IsStudentRelative bool
}

type Place struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"unique"`
	Address  string
	Capacity uint
}

type Event struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	Name           string `gorm:"unique"`
	Description    string
	IsPaid         bool
	MaxTicketBatch uint

	PlaceID uint

	IsActive bool
}

type Ticket struct {
	ID           uint `gorm:"primaryKey"`
	PurchaseTime time.Time
	UserCount    uint

	UserID  int64
	EventID uint
}
