package models

import "time"

// Reservation holds the reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}
// Users holds the user data modle
type Users struct{

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
type Rooms struct {
	ID int
	RoomName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions holds the restriction data
type Restrictions struct {
	ID int
	RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
// Reservations holds the reservation data
type Reservations struct {
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
	Room Rooms

	
}

//RoomRestrictions holds the room restriction
type RoomRestrictions struct {
	ID int
	RoomID int
	ReservationID int
	RestrictionID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Rooms
	Restriction Restrictions
	Reservation Reservations
}