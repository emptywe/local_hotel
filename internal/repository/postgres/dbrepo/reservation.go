package dbrepo

import (
	"context"
	"time"

	"github.com/emptywe/local_hotel/model"
)

// InsertReservation inserts reservation into database
func (r *postgresDBRepo) InsertReservation(res model.Reservation) (int, error) {

	var rId int

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	sql := `INSERT INTO reservations (first_name, last_name, email, phone, 
                          start_date,end_date, room_id,created_at,updated_at) 
						  values ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`

	row := r.DB.QueryRow(ctx, sql,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	)

	if err := row.Scan(&rId); err != nil {
		return 0, err
	}

	return rId, nil
}

// InsertRoomRestriction inserts a room restriction into database
func (r *postgresDBRepo) InsertRoomRestriction(res model.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	sql := `INSERT INTO room_restrictions ( start_date,end_date, room_id,reservation_id,created_at,updated_at,restriction_id) 
						  values ($1,$2,$3,$4,$5,$6,$7) RETURNING id`

	_, err := r.DB.Exec(ctx, sql,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		res.ReservationId,
		time.Now(),
		time.Now(),
		res.RestrictionId,
	)

	return err
}

// SearchAvailabilityByDatesByRoomId returns true if passed dates are available for specific room
func (r *postgresDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {
	var (
		reservationsNum int
		availability    bool
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	sql := `SELECT COUNT(id) FROM reservations WHERE $1 < end_date AND $2 > start_date AND room_id = $3 `

	row := r.DB.QueryRow(ctx, sql, start, end, roomId)

	if err := row.Scan(&reservationsNum); err != nil {
		return availability, err
	}
	if reservationsNum == 0 {
		availability = true
	}
	return availability, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms for given date range
func (r *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]model.Room, error) {

	var rooms []model.Room
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `SELECT id, room_name FROM rooms WHERE id NOT IN (SELECT room_id FROM room_restrictions WHERE $1 < end_date AND $2 > start_date)`
	rows, err := r.DB.Query(ctx, sql, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room model.Room
		if err = rows.Scan(&room.Id, &room.RoomName); err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// GetRoomById returns a room by given id
func (r *postgresDBRepo) GetRoomById(id int) (model.Room, error) {
	var room model.Room
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	sql := `SELECT * from rooms WHERE id=$1`

	row := r.DB.QueryRow(ctx, sql, id)
	if err := row.Scan(&room.Id, &room.RoomName, &room.DayPrice, &room.CreatedAt, &room.UpdatedAt); err != nil {
		return room, err
	}
	return room, nil
}
