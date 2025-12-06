package reservations

import (
	"slices"
	"testing"
	"time"
	"github.com/SneedusSnake/Reservations/internal/utils"
)

type ReservationsRegistryContract struct {
	NewRegistry func() ReservationsRegistry
}

func (r ReservationsRegistryContract) Test (t *testing.T) {
	t.Run("it returns error when the reservation was not found", func (t *testing.T) {
		registry := r.NewRegistry()
		_, err := registry.Get(1234567)

		if err == nil {
			t.Error("expected to see error, got nil")
		}
	})

	t.Run("it adds a new reservation into the registry", func (t *testing.T) {
		registry := r.NewRegistry()
		reservation := Reservation{1,2,3, time.Now(), time.Now()}

		err := registry.Add(reservation)

		if err != nil {
			t.Fatal(err)
		}

		foundReservation, err := registry.Get(reservation.Id)

		if err != nil {
			t.Fatal(err)
		}

		if reservation != foundReservation {
			t.Errorf("expected to see %#v got %#v", reservation, foundReservation)
		}
	})

	t.Run("it removes reservation from the registry", func (t *testing.T) {
		registry := r.NewRegistry()
		reservation := Reservation{1,2,3, time.Now(), time.Now()}
		err := registry.Add(reservation)

		if err != nil {
			t.Fatal(err)
		}

		err = registry.Remove(reservation.Id)

		if err != nil {
			t.Fatal(err)
		}

		_, err = registry.Get(reservation.Id)

		if err == nil {
			t.Error("Expected reservation to be removed from the registry")
		}
	})

	t.Run("it fetches reservations active during given period", func (t *testing.T) {
		registry := r.NewRegistry()
		from := time.Now()
		to := time.Now().Add(time.Hour*1)
		
		expiredReservations := Reservations{
			Reservation{1,1,1, from.Add(-time.Hour*2), from.Add(-time.Second)},
			Reservation{2,2,2, from.Add(-time.Hour*4), from.Add(-time.Hour)},
		}
		activeReservations := Reservations{
			Reservation{3,3,3, from.Add(-time.Hour*2), from.Add(time.Second)},
			Reservation{4,4,4, from, to},
			Reservation{5,5,5, from.Add(time.Second), to.Add(-time.Second)},
			Reservation{6,6,6, to.Add(-time.Second), to.Add(time.Minute*20)},
		}
		futureReservations := Reservations{
			Reservation{7,7,7, to.Add(time.Second), to.Add(time.Hour)},
		}
		all := Reservations{}
		all = append(all, expiredReservations...)
		all = append(all, activeReservations...)
		all = append(all, futureReservations...)

		for _, reservation := range all {
			err := registry.Add(reservation)
			if err != nil {
				t.Fatal(err)
			}
		}

		reservations := registry.ForPeriod(from, to)
		
		if len(reservations) != len(activeReservations) {
			t.Errorf("Expected %d active reservations, got %d", len(activeReservations), len(reservations))
		}
	})

	t.Run("it generates next ID", func(t *testing.T) {
		registry := r.NewRegistry()
		ch := make(chan int, 5)
		var ids []int

		for range 5 {
			go (func (c chan int) {
				c <- registry.NextIdentity()
			})(ch)
		}
		
		for range 5 {
			ids = append(ids, <- ch)
		}

		if !slices.IsSorted(ids) {
			t.Errorf("Generated identities %v are not in ascending order", ids)
		}

		if len(utils.Unique(ids)) != len(ids) {
			t.Errorf("Generated identities %v contain duplicate values", ids)
		}
	})
}
