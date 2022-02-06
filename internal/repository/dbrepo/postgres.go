package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/tsawler/bookings-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation insert a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
            end_date, room_id, created_at, updated_at)
          values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
					created_at, updated_at, restriction_id)
					values
					($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDates returns true if availability exists for roomID and false if no availability exists
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	var numRows int
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id)
						from room_restrictions
						where
						room_id = $1
						and $2 < end_date
						and $3 > start_date;`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)

	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `select r.id, r.room_name
						from rooms r
						where r.id not in (
								select rr.room_id
								from room_restrictions rr
								where $1 < rr.end_date and $2 > rr.start_date
						)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at
						from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil
}

func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
						from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	} else {
		return u, nil
	}
}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUserByID(u models.User) error {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
						from users where id = $1`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil

}

// Authenticate authenticates the users
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	query := `select id, password from users where email = $1`
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return 0, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return id, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllReservations returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.phone,
						r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
						rm.id, rm.room_name
						from reservations as r
						left join rooms as rm on (r.room_id = rm.id)
						order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		defer rows.Close()
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.phone,
						r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
						rm.id, rm.room_name
						from reservations as r
						left join rooms as rm on (r.room_id = rm.id)
						where processed = 0
						order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		defer rows.Close()
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

// GetReservationByID returns one reservation by id
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	// Set default context to three seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.phone,
						r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
						rm.id, rm.room_name
						from reservations as r
						left join rooms as rm on (r.room_id = rm.id)
						where r.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomID,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
	)
	if err != nil {
		return reservation, err
	}
	return reservation, nil
}
