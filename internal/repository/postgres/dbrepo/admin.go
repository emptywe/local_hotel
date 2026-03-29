package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/emptywe/local_hotel/model"
	"golang.org/x/crypto/bcrypt"
)

// GetUserById returns user struct by id
func (r *postgresDBRepo) GetUserById(id int) (model.User, error) {

	var user model.User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `SELECT id, first_name, last_name, email, access_level, created_at, updated_at FROM users WHERE id=$1`
	row := r.DB.QueryRow(ctx, sql, id)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser updates user information by id
func (r *postgresDBRepo) UpdateUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `UPDATE users SET first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5 WHERE id=$6 `

	_, err := r.DB.Exec(ctx, sql, user.FirstName, user.LastName, user.Email, user.AccessLevel, time.Now(), user.Id)

	return err
}

// Authenticate authenticates user
func (r *postgresDBRepo) Authenticate(email, password string) (int, error) {

	var (
		id             int
		hashedPassword string
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `SELECT id, password FROM users WHERE email=$1`

	row := r.DB.QueryRow(ctx, sql, email)
	if err := row.Scan(&id, &hashedPassword); err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

// AllReservations - returns a slice of all reservations
func (r *postgresDBRepo) AllReservations() (reservations []model.Reservation, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, 
			r.start_date, r.end_date, r.created_at,r.updated_at,r.processed,rm.room_name
			FROM reservations as r
			LEFT JOIN rooms as rm ON r.room_id = rm.id
			ORDER BY r.start_date ASC
			`
	rows, err := r.DB.Query(ctx, sql)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var rs model.Reservation
		if err = rows.Scan(
			&rs.Id,
			&rs.FirstName,
			&rs.LastName,
			&rs.Email,
			&rs.Phone,
			&rs.StartDate,
			&rs.EndDate,
			&rs.CreatedAt,
			&rs.UpdatedAt,
			&rs.Processed,
			&rs.Room.RoomName,
		); err != nil {
			return reservations, err
		}
		reservations = append(reservations, rs)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return
}

// NewReservations - returns a slice of new not processed reservations
func (r *postgresDBRepo) NewReservations() (reservations []model.Reservation, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, 
			r.start_date, r.end_date, r.created_at,r.updated_at,r.processed,rm.room_name
			FROM reservations as r
			LEFT JOIN rooms as rm ON r.room_id = rm.id
			WHERE processed = false
			ORDER BY r.start_date ASC
			`
	rows, err := r.DB.Query(ctx, sql)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var rs model.Reservation
		if err = rows.Scan(
			&rs.Id,
			&rs.FirstName,
			&rs.LastName,
			&rs.Email,
			&rs.Phone,
			&rs.StartDate,
			&rs.EndDate,
			&rs.CreatedAt,
			&rs.UpdatedAt,
			&rs.Processed,
			&rs.Room.RoomName,
		); err != nil {
			return reservations, err
		}
		reservations = append(reservations, rs)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return
}

// GetReservationById - returns one reservation by given id
func (r *postgresDBRepo) GetReservationById(id int) (reservation model.Reservation, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, 
			r.start_date, r.end_date, r.created_at,r.updated_at,r.processed, rm.id, rm.room_name
			FROM reservations as r
			LEFT JOIN rooms as rm ON r.room_id = rm.id
			WHERE r.id = $1
			`

	row := r.DB.QueryRow(ctx, sql, id)
	if err = row.Scan(
		&reservation.Id,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.Id,
		&reservation.Room.RoomName,
	); err != nil {
		return reservation, err
	}

	return
}

// UpdateReservation updates reservation information
func (r *postgresDBRepo) UpdateReservation(res model.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `UPDATE reservations SET first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5 WHERE id=$6 `

	_, err := r.DB.Exec(ctx, sql, res.FirstName, res.LastName, res.Email, res.Phone, time.Now(), res.Id)

	return err
}

// DeleteReservation deletes reservation by id
func (r *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `DELETE FROM  reservations WHERE id=$1 `

	_, err := r.DB.Exec(ctx, sql, id)

	return err
}

// UpdateProcessedReservation updates reservation processed field by id
func (r *postgresDBRepo) UpdateProcessedReservation(id int, processed bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `UPDATE reservations SET processed=$1 WHERE id=$2 `

	_, err := r.DB.Exec(ctx, sql, processed, id)

	return err
}

// AllRooms returns information about all rooms
func (r *postgresDBRepo) AllRooms() (rooms []model.Room, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `SELECT * FROM rooms ORDER BY room_name`

	rows, err := r.DB.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var room model.Room
		if err = rows.Scan(
			&room.Id,
			&room.RoomName,
			&room.DayPrice,
			&room.CreatedAt,
			&room.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return
}

// GetRestrictionsForRoomByDate returns room restrictions by given date range
func (r *postgresDBRepo) GetRestrictionsForRoomByDate(roomId int, startDate, endDate time.Time) (restrictions []model.RoomRestriction, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `SELECT id, COALESCE(reservation_id, 0), restriction_id, room_id, start_date, end_date FROM room_restrictions WHERE $1 < end_date AND $2 >= start_date AND room_id=$3`

	rows, err := r.DB.Query(ctx, sql, startDate, endDate, roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var restriction model.RoomRestriction
		if err = rows.Scan(
			&restriction.Id,
			&restriction.ReservationId,
			&restriction.RestrictionId,
			&restriction.RoomId,
			&restriction.StartDate,
			&restriction.EndDate,
		); err != nil {
			return nil, err
		}
		restrictions = append(restrictions, restriction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return
}

// InsertBlockForRoom inserts block for given room by id
func (r *postgresDBRepo) InsertBlockForRoom(roomId int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `INSERT INTO room_restrictions (start_date, end_date, room_id,restriction_id,created_at,updated_at) values ($1,$2,$3,$4,$5,$6)`

	_, err := r.DB.Exec(ctx, sql, startDate, startDate.AddDate(0, 0, 1), roomId, 2, time.Now(), time.Now())
	return err
}

// DeleteBlockById deletes block by id
func (r *postgresDBRepo) DeleteBlockById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	sql := `DELETE FROM room_restrictions WHERE id=$1`

	_, err := r.DB.Exec(ctx, sql, id)
	return err
}
