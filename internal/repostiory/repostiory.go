package repostiory

import "githup.com/Mo277210/booking/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation)(int,error)
	InsertRoomRestriction(r models.RoomRestrictions)error
}