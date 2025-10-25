package repostiory

import (
	"time"

	"githup.com/Mo277210/booking/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation)(int,error)
	InsertRoomRestriction(r models.RoomRestrictions)error
    SearchAvailabilityByDatesByRoomID(start,end time.Time,roomID int)(bool,error)
	SearchAvailabilityForAllRooms(start,end time.Time)([]models.Room, error)
	GetRoomByID(id int) (models.Room,error)
}