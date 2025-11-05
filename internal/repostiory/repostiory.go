package repostiory

import (
	"time"

	"githup.com/Mo277210/booking/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestrictions) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error) 
	GetReservationByID(id int) (models.Reservation, error)
    UpdateResertvation(u models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedReservation(id int, processed int) error
	 AllRooms () ([]models.Room, error)
	 GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestrictions, error)
}


