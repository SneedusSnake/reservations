package inmemory

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
)

type ReservationsStore struct
{
	counter int
	reservations reservations.Reservations
	mu sync.Mutex
}

func NewReservationStore() *ReservationsStore {
	return &ReservationsStore{}
}

func (r *ReservationsStore) NextIdentity() (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++

	return r.counter, nil
}

func (r *ReservationsStore) List() reservations.Reservations {
	return slices.Clone(r.reservations)
}

func (r *ReservationsStore) ForPeriod(from time.Time, to time.Time) reservations.Reservations {
	var result []reservations.Reservation

	if (from.After(to)) {
		return result
	}

	for _, reservation := range(r.reservations) {
		if inInterval(reservation.Start, from, to) || inInterval(reservation.End, from, to) {
			result = append(result, reservation)
		}
	}

	return result
}

func (r *ReservationsStore) Add(reservation reservations.Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reservations = append(r.reservations, reservation)

	return nil
}

func (r *ReservationsStore) Get(id int) (reservations.Reservation, error) {
	for _, reservation := range r.reservations {
		if (reservation.Id == id) {
			return reservation, nil
		}
	}

	return reservations.Reservation{}, fmt.Errorf("Reservation with id %d was not found", id)
}

func (r *ReservationsStore) Remove(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for index, reservation := range r.reservations {
		if (reservation.Id == id) {
			r.reservations = append(r.reservations[:index], r.reservations[index+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Reservation with id %d was not found", id);
}

func inInterval(t time.Time, from time.Time, to time.Time) bool {
	return t.Add(time.Second).After(from) && t.Add(-time.Second).Before(to)
}
