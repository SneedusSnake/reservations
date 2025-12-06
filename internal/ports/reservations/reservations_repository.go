package reservations

import (
	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	"time"
)

type ReservationsRepository interface {
	NextIdentity() int
	Add(reservation reservations.Reservation) error
	Get(id int) (reservations.Reservation, error)
	Remove(id int) error
	ForPeriod(from time.Time, to time.Time) reservations.Reservations
}


