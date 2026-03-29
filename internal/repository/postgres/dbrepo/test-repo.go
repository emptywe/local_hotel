package dbrepo

import (
	"errors"
	"time"

	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/internal/repository"
	"github.com/emptywe/local_hotel/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// testDBRepo test repo struct
type testDBRepo struct {
	App *config.AppConfig
	DB  *pgxpool.Pool
}

// NewTestRepo creating new test repo
func NewTestRepo(app *config.AppConfig) repository.DataBaseRepo {
	return &testDBRepo{
		App: app,
	}
}

// InsertReservation inserts reservation into database
func (r *testDBRepo) InsertReservation(res model.Reservation) (int, error) {
	if time.Since(res.StartDate) > 0 || time.Since(res.EndDate) > 0 {
		return 0, errors.New("invalid dates time")
	}
	if res.Room.Id == 0 {
		return 0, errors.New("invalid room struct")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into database
func (r *testDBRepo) InsertRoomRestriction(res model.RoomRestriction) error {
	if res.RoomId == 1999 {
		return errors.New("invalid room id")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomId returns true if passed dates are available for specific room
func (r *testDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {
	if roomId == 1999 {
		return false, errors.New("sql: no rows")
	}
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms for given date range
func (r *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]model.Room, error) {
	var rooms []model.Room
	if time.Since(start) > 0 {
		return nil, errors.New("some error")
	}
	if start.Sub(end) == 0 {
		return rooms, nil
	}
	rooms = append(rooms, model.Room{})
	return rooms, nil
}

// GetRoomById returns a room by given id
func (r *testDBRepo) GetRoomById(id int) (model.Room, error) {
	var room model.Room
	if id > 2 {
		return room, errors.New("sql: no rows")
	}
	return room, nil
}

// GetUserById returns user struct by id
func (r *testDBRepo) GetUserById(id int) (model.User, error) {

	var user model.User

	return user, nil
}

// UpdateUser updates user information by id
func (r *testDBRepo) UpdateUser(user model.User) error {
	return nil
}

// Authenticate authenticates user
func (r *testDBRepo) Authenticate(email, password string) (int, error) {
	var (
		id int
	)
	return id, nil
}

// AllReservations - returns a slice of all reservations
func (r *testDBRepo) AllReservations() (reservations []model.Reservation, err error) {
	return
}

// NewReservations - returns a slice of new not processed reservations
func (r *testDBRepo) NewReservations() (reservations []model.Reservation, err error) {
	return
}

// GetReservationById - returns one reservation by given id
func (r *testDBRepo) GetReservationById(id int) (reservation model.Reservation, err error) {
	return
}

// UpdateReservation updates reservation information by id
func (r *testDBRepo) UpdateReservation(res model.Reservation) error {
	return nil
}

// DeleteReservation deletes reservation by id
func (r *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedReservation updates reservation processed field by id
func (r *testDBRepo) UpdateProcessedReservation(id int, processed bool) error {
	return nil
}

// AllRooms returns information about all rooms
func (r *testDBRepo) AllRooms() (rooms []model.Room, err error) {
	return
}

// GetRestrictionsForRoomByDate returns room restrictions by given date
func (r *testDBRepo) GetRestrictionsForRoomByDate(roomId int, startDate, endDate time.Time) (restrictions []model.RoomRestriction, err error) {
	return
}

// InsertBlockForRoom inserts block for given room by id
func (r *testDBRepo) InsertBlockForRoom(roomId int, startDate time.Time) error {
	return nil
}

// DeleteBlockById deletes block by id
func (r *testDBRepo) DeleteBlockById(id int) error {
	return nil
}
