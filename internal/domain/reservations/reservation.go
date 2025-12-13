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

func (r Reservations) ForUser(userId int) Reservations {
	var filtered Reservations

	for _, reservation := range r {
		if reservation.UserId == userId {
			filtered = append(filtered, reservation)
		}
	}

	return filtered
}

