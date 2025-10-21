package models

import "time"

// Reservation holds the reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

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