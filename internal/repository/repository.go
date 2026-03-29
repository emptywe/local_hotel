package repository

import (
	"time"

	"github.com/emptywe/local_hotel/model"
)

// DataBaseRepo describes database methods and behaviour
type DataBaseRepo interface {
	InsertReservation(res model.Reservation) (int, error)
	InsertRoomRestriction(res model.RoomRestriction) error
	SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]model.Room, error)
	GetRoomById(id int) (model.Room, error)

	GetUserById(id int) (model.User, error)
	UpdateUser(user model.User) error
	Authenticate(email, password string) (int, error)

	AllReservations() (reservations []model.Reservation, err error)
	NewReservations() (reservations []model.Reservation, err error)
	GetReservationById(id int) (reservation model.Reservation, err error)
	UpdateReservation(res model.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedReservation(id int, processed bool) error
	AllRooms() (rooms []model.Room, err error)
	GetRestrictionsForRoomByDate(roomId int, startDate, endDate time.Time) (restrictions []model.RoomRestriction, err error)
	InsertBlockForRoom(roomId int, startDate time.Time) error
	DeleteBlockById(id int) error
}
