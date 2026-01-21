package reservations

import (
	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	"github.com/SneedusSnake/Reservations/internal/read_model"
	"time"
)

type ReservationsRepository interface {
	NextIdentity() (int, error)
	List() reservations.Reservations
	Add(reservation reservations.Reservation) error
	Get(id int) (reservations.Reservation, error)
	Remove(id int) error
	ForPeriod(from time.Time, to time.Time) reservations.Reservations
}

type ReservationsReadRepository interface {
	Get(id int) (readmodel.Reservation, error)
	Active(t time.Time, tags ...string) ([]readmodel.Reservation, error)
}


