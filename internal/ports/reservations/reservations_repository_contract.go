package reservations

import (
	"slices"
	"testing"
	"time"
	domain "github.com/SneedusSnake/Reservations/internal/domain/reservations"
	"github.com/SneedusSnake/Reservations/internal/utils"
)

type ReservationsRepositoryContract struct {
	NewRepository func() ReservationsRepository
}

func (r ReservationsRepositoryContract) Test (t *testing.T) {
	t.Run("it returns error when the reservation was not found", func (t *testing.T) {
		store := r.NewRepository()
		_, err := store.Get(1234567)

		if err == nil {
			t.Error("expected to see error, got nil")
		}
	})

	t.Run("it adds a new reservation into the store", func (t *testing.T) {
		store := r.NewRepository()
		reservation := domain.Reservation{Id: 1,UserId: 2,SubjectId: 3, Start: time.Now(), End: time.Now()}

		err := store.Add(reservation)

		if err != nil {
			t.Fatal(err)
		}

		foundReservation, err := store.Get(reservation.Id)

		if err != nil {
			t.Fatal(err)
		}

		if reservation != foundReservation {
			t.Errorf("expected to see %#v got %#v", reservation, foundReservation)
		}
	})

	t.Run("it removes reservation from the store", func (t *testing.T) {
		store := r.NewRepository()
		reservation := domain.Reservation{Id: 1,UserId: 2,SubjectId: 3, Start: time.Now(), End: time.Now()}
		err := store.Add(reservation)

		if err != nil {
			t.Fatal(err)
		}

		err = store.Remove(reservation.Id)

		if err != nil {
			t.Fatal(err)
		}

		_, err = store.Get(reservation.Id)

		if err == nil {
			t.Error("Expected reservation to be removed from the store")
		}
	})

	t.Run("it fetches reservations active during given period", func (t *testing.T) {
		store := r.NewRepository()
		from := time.Now()
		to := time.Now().Add(time.Hour*1)
		
		expiredReservations := domain.Reservations{
			domain.Reservation{Id: 1,UserId: 1,SubjectId: 1, Start: from.Add(-time.Hour*2), End: from.Add(-time.Second)},
			domain.Reservation{Id: 2,UserId: 2,SubjectId: 2, Start: from.Add(-time.Hour*4), End: from.Add(-time.Hour)},
		}
		activeReservations := domain.Reservations{
			domain.Reservation{Id: 3,UserId: 3,SubjectId: 3, Start: from.Add(-time.Hour*2), End: from.Add(time.Second)},
			domain.Reservation{Id: 4,UserId: 4,SubjectId: 4, Start: from, End: to},
			domain.Reservation{Id: 5,UserId: 5,SubjectId: 5, Start: from.Add(time.Second), End: to.Add(-time.Second)},
			domain.Reservation{Id: 6,UserId: 6,SubjectId: 6, Start: to.Add(-time.Second), End: to.Add(time.Minute*20)},
		}
		futureReservations := domain.Reservations{
			domain.Reservation{Id: 7,UserId: 7,SubjectId: 7, Start: to.Add(time.Second), End: to.Add(time.Hour)},
		}
		all := domain.Reservations{}
		all = append(all, expiredReservations...)
		all = append(all, activeReservations...)
		all = append(all, futureReservations...)

		for _, reservation := range all {
			err := store.Add(reservation)
			if err != nil {
				t.Fatal(err)
			}
		}

		reservations := store.ForPeriod(from, to)
		
		if len(reservations) != len(activeReservations) {
			t.Errorf("Expected %d active reservations, got %d", len(activeReservations), len(reservations))
		}
	})

	t.Run("it generates next ID", func(t *testing.T) {
		store := r.NewRepository()
		ch := make(chan int, 5)
		var ids []int

		for range 5 {
			go (func (c chan int) {
				id, _ := store.NextIdentity()
				c <- id
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
