package dbrepo

import (
	"errors"
	"github.com/samuelowad/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//GetUserByID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := getCtx()
	defer cancel()

	query := `SELECT id,first_name,last_name,email,access_token,password,created_at,updated_at FROM users WHERE id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.AccessLevel, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}
	return u, nil
}

//UpdateUser updates user data
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := getCtx()
	defer cancel()

	query := `UPDATE set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())

	if err != nil {
		return err
	}
	return nil
}

//Authenticate authenticates the user with the provided data
func (m *postgresDBRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := getCtx()
	defer cancel()
	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id,password from users where email=$1", email)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

//AllRepositories returns all reservation
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := getCtx()
	defer cancel()
	var reservation []models.Reservation

	query := `SELECT r.id,r.first_name,r.last_name,r.email,r.phone,r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at, rm.id,rm.room_name,r.processed  FROM reservations r left join rooms rm on (r.room_id=rm.id) ORDER BY r.start_date ASC`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservation, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName, &i.Email, &i.Phone, &i.StartDate, &i.EndDate, &i.RoomID, &i.CreatedAt, &i.UpdatedAt, &i.Room.ID, &i.Room.RoomName, &i.Processed)
		if err != nil {
			return reservation, err
		}
		reservation = append(reservation, i)

	}
	if err = rows.Err(); err != nil {
		return reservation, err
	}
	return reservation, nil

}

//AllNewRepositories returns all reservation
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := getCtx()
	defer cancel()
	var reservation []models.Reservation

	query := `SELECT r.id,r.first_name,r.last_name,r.email,r.phone,r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at, rm.id,rm.room_name  FROM reservations r left join rooms rm on (r.room_id=rm.id) where processed=0 ORDER BY r.start_date ASC`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservation, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName, &i.Email, &i.Phone, &i.StartDate, &i.EndDate, &i.RoomID, &i.CreatedAt, &i.UpdatedAt, &i.Room.ID, &i.Room.RoomName)
		if err != nil {
			return reservation, err
		}
		reservation = append(reservation, i)

	}
	if err = rows.Err(); err != nil {
		return reservation, err
	}
	return reservation, nil

}

//GetReservationByID returns one reservation by Id
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := getCtx()
	defer cancel()

	var res models.Reservation

	query := `SELECT r.id,r.first_name,r.last_name,r.email,r.phone,r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at, rm.id,rm.room_name  FROM reservations r left join rooms rm on (r.room_id=rm.id) where r.id=$1 `

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&res.ID, &res.FirstName, &res.LastName, &res.Email, &res.Phone, &res.StartDate, &res.EndDate, &res.RoomID, &res.CreatedAt, &res.UpdatedAt, &res.Room.ID, &res.Room.RoomName)
	if err != nil {
		return res, err
	}

	return res, nil

}
