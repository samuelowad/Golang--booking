package dbrepo

import (
	"errors"
	"github.com/samuelowad/bookings/internal/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool { return true }

//InsertReservation inserts a new reservation
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//if room id is 2 fail
	if res.RoomID == 2 {
		return 0, errors.New("error inserting reservation")
	}
	return 1, nil
}

//InsertRoomRestriction insert room restriction in database
func (m *testDBRepo) InsertRoomRestriction(rr models.RoomRestriction) error {
	if rr.RoomID == 1000 {
		return errors.New("error inserting room restriction")
	}
	return nil
}

// retruns if Availability exists

func (m *testDBRepo) SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil

}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("Error: invalid room ID")
	}

	return room, nil
}

//GetUserByID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {

	var u models.User

	return u, nil
}
