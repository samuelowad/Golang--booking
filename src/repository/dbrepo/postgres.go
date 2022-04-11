package dbrepo

import (
	"context"
	"github.com/samuelowad/bookings/src/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool { return true }

//InsertReservation inserts a new reservation
func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `insert into reservations (first_name, last_name,email,phone,start_date,end_date,room_id,created_at,updated_at) values($1,$2,$3,$4,$5,$6,$8,$9)`

	_, err := m.DB.ExecContext(ctx, statement, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}
