package reservations

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	domain "github.com/SneedusSnake/Reservations/internal/domain/reservations"
	"github.com/SneedusSnake/Reservations/internal/utils"
	"github.com/alecthomas/assert/v2"
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
		reservation := builder(t, store).Persist()

		foundReservation, err := store.Get(reservation.Id)
		assert.NoError(t, err)
		assert.Equal(t, reservation, foundReservation)
	})

	t.Run("it removes reservation from the store", func (t *testing.T) {
		store := r.NewRepository()
		reservation := builder(t, store).Persist()

		err := store.Remove(reservation.Id)
		assert.NoError(t, err)

		_, err = store.Get(reservation.Id)
		assert.Error(t, err)
	})

	t.Run("it fetches reservations active during given period", func (t *testing.T) {
		store := r.NewRepository()
		from, err := time.Parse(time.DateTime, "2025-09-20 14:00:00")
		assert.NoError(t, err)
		to := from.Add(time.Hour*1)
		blueprint := builder(t, store)
		
		//expired reservations relative to given period
		blueprint.StartsAt(from.Add(-time.Hour*2)).EndsAt(from.Add(-time.Second)).Persist()
		blueprint.StartsAt(from.Add(-time.Hour*4)).EndsAt(from.Add(-time.Hour)).Persist()

		expectedReservations := domain.Reservations{
			blueprint.StartsAt(from.Add(-time.Hour*2)).EndsAt(from.Add(time.Second)).Persist(),
			blueprint.StartsAt(from).EndsAt(to).Persist(),
			blueprint.StartsAt(from.Add(time.Second)).EndsAt(to.Add(-time.Second)).Persist(),
			blueprint.StartsAt(to.Add(-time.Second)).EndsAt(to.Add(time.Minute*20)).Persist(),
		}
		//future reservations relative to given period
		blueprint.StartsAt(to.Add(time.Second)).EndsAt(to.Add(time.Hour)).Persist()

		reservations, err := store.ForPeriod(from, to)
		assert.NoError(t, err)
		assert.Equal(t, len(reservations), len(expectedReservations))

		for _, r := range reservations {
			assert.SliceContains(t, expectedReservations, r)
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

type reservationBuilder struct {
	t testing.TB
	store ReservationsRepository
	blueprint domain.Reservation
}

func (builder reservationBuilder) Make() domain.Reservation {
	id, err := builder.store.NextIdentity()
	assert.NoError(builder.t, err)
	userId := builder.blueprint.UserId
	subjectId := builder.blueprint.SubjectId
	start := builder.blueprint.Start
	end := builder.blueprint.End

	if userId == 0 {
		userId = rand.N(999999) + 1
	}
	
	if subjectId == 0 {
		subjectId = rand.N(9999) + 1
	}

	if start.IsZero() {
		start = time.Now().UTC().Truncate(time.Second)
	} 

	if end.IsZero() {
		end = time.Now().Add(time.Hour*2).UTC().Truncate(time.Second)
	} 

	return domain.Reservation{
		Id: id,
		UserId: userId,
		SubjectId: subjectId,
		Start: start,
		End: end,
	}
}

func (builder reservationBuilder) Persist() domain.Reservation {
	result := builder.Make()
	err := builder.store.Add(result)
	assert.NoError(builder.t, err)
	builder.t.Cleanup(func() {
		builder.store.Remove(result.Id)
	})

	return result
}

func (builder reservationBuilder) StartsAt(t time.Time) reservationBuilder {
	blueprint := builder.blueprint
	blueprint.Start = t
	builder.blueprint = blueprint

	return builder
}

func (builder reservationBuilder) EndsAt(t time.Time) reservationBuilder {
	blueprint := builder.blueprint
	blueprint.End = t
	builder.blueprint = blueprint

	return builder
}

func (builder reservationBuilder) UserId(id int) reservationBuilder {
	blueprint := builder.blueprint
	blueprint.UserId = id
	builder.blueprint = blueprint

	return builder
}

func (builder reservationBuilder) SubjectId(id int) reservationBuilder {
	blueprint := builder.blueprint
	blueprint.SubjectId = id
	builder.blueprint = blueprint

	return builder
}

func (builder reservationBuilder) CleanUp(t testing.TB) {
	t.Cleanup(func() {
		rs, err := builder.store.List()
		assert.NoError(t, err)

		for _, r := range rs {
			builder.store.Remove(r.Id)
		}
	})
}

func builder(t testing.TB, store ReservationsRepository) reservationBuilder {
	return reservationBuilder{
		t: t,
		store: store,
		blueprint: domain.Reservation{},
	}
}
