package inmemory

import (
	"fmt"
	"time"
	"errors"
	"github.com/SneedusSnake/Reservations/domain/reservations"
)

type ReservationsStore struct
{
	counter int
	reservations reservations.Reservations
}

func NewReservationStore() *ReservationsStore {
	return &ReservationsStore{}
}

func (r *ReservationsStore) NextIdentity() int {
	r.counter++

	return r.counter
}

func (r *ReservationsStore) ReservedAt(t time.Time) reservations.Reservations {
	var result []reservations.Reservation

	for _, reservation := range(r.reservations) {
		if reservation.ActiveAt(t) {
			result = append(result, reservation)
		}
	}

	return result
}

func (r *ReservationsStore) Add(reservation reservations.Reservation) error {
	r.reservations = append(r.reservations, reservation)

	return nil
}

func (r *ReservationsStore) Get(id int) (reservations.Reservation, error) {
	for _, reservation := range r.reservations {
		if (reservation.Id == id) {
			return reservation, nil
		}
	}

	return reservations.Reservation{}, errors.New(fmt.Sprintf("Reservation with id %d was not found", id))
}

func (r *ReservationsStore) Remove(id int) error {

	for index, reservation := range r.reservations {
		if (reservation.Id == id) {
			r.reservations = append(r.reservations[:index], r.reservations[index+1:]...)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Reservation with id %d was not found", id));
}

