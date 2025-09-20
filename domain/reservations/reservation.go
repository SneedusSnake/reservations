package reservations

import (
	"fmt"
	"testing"
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

type ReservationsRegistry interface {
	NextIdentity() int
	ReservedAt(t time.Time) Reservations
	Add(reservation Reservation) error
	Get(id int) (Reservation, error)
	Remove(id int) error
}


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

		s:= []int{1,2,3,4,5,6}
		s = append(s[:3], s[4:6]...)
		fmt.Print(s)
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

	t.Run("it fetches reservations active at given time", func (t *testing.T) {
		registry := r.NewRegistry()
		now := time.Now()
		pastReservations := Reservations{
			Reservation{1,1,1, now.Add(-time.Hour*2), now},
			Reservation{2,2,2, now.Add(-time.Hour*4), now.Add(-time.Second)},
		}
		currentReservations := Reservations{
			Reservation{3,3,3, now.Add(-time.Hour*2), now.Add(time.Hour)},
			Reservation{4,4,4, now.Add(-time.Hour*4), now.Add(time.Second)},
		}
		futureReservations := Reservations{
			Reservation{5,5,5, now.Add(time.Second), now.Add(time.Hour)},
		}
		all := Reservations{}
		all = append(all, pastReservations...)
		all = append(all, currentReservations...)
		all = append(all, futureReservations...)

		for _, reservation := range all {
			err := registry.Add(reservation)
			if err != nil {
				t.Fatal(err)
			}
		}

		reservations := registry.ReservedAt(now)
		
		if len(reservations) != len(currentReservations) {
			t.Errorf("Expected %d active reservations, got %d", len(currentReservations), len(reservations))
		}
	})
}
