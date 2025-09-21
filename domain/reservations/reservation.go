package reservations

import (
	"time"
)

type Reservation struct {
	Id  	  int
	UserId    int
	SubjectId int
	Start     time.Time
	End       time.Time
}

func (r Reservation) ActiveAt(t time.Time) bool {
	return r.End.After(t) && (r.Start.Before(t) || r.Start.Equal(t))
}

type Reservations []Reservation

func (r Reservations) ForSubject(subjectId int) Reservations {
	var filtered Reservations

	for _, reservation := range r {
		if reservation.SubjectId == subjectId {
			filtered = append(filtered, reservation)
		}
	}

	return filtered
}

type ReservationsRegistry interface {
	NextIdentity() int
	ReservedAt(t time.Time) Reservations
	Add(reservation Reservation) error
	Get(id int) (Reservation, error)
	Remove(id int) error
}


