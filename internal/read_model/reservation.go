package readmodel

import "time"

type Reservation struct {
	Id int
	Subject string
	User string
	Start time.Time
	End time.Time
}

