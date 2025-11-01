package models

import (
	// "text/template"
	"time"
	"html/template"
)

// Users holds the user data modle
type User struct{

	ID int
	FirstName string
	LastName string
	Email string
	Password string
	AccessLevel int
	CreateAt time.Time
	UpdateAt time.Time



}

// Rooms holds the room data
type Room struct {
	ID int
	RoomName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions holds the restriction data
type Restriction struct {
	ID int
	RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
// Reservations holds the reservation data
type Reservation struct {
	ID int
	FirstName string
	LastName string
	Email string
	Phone string
	StartDate time.Time
	EndDate time.Time
	RoomID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room

	
}

//RoomRestrictions holds the room restriction
type RoomRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

// MailData holds the data about an email
type MailData struct {
	To      string
	From    string
	Subject string
	Content template.HTML
	Template string
}