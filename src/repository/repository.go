package repository

import "github.com/samuelowad/bookings/src/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) error
}
