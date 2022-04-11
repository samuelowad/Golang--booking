package models

import "time"

//DB models
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

type RestrictionName struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	Room          Room
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Reservation   Reservation
	ReservationID int
	RestrictionID int
	Restriction   Restriction
}