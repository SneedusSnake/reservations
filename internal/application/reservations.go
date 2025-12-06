package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/SneedusSnake/Reservations/internal/ports"
	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	usersPort "github.com/SneedusSnake/Reservations/internal/ports/users"
	reservationsPort "github.com/SneedusSnake/Reservations/internal/ports/reservations"
)

type CreateReservation struct {
	SubjectId int
	UserId int
	From time.Time
	To time.Time
}

type AlreadyReservedError struct {
	ReservationIds []int
}

func (e AlreadyReservedError) Error() string {
	return string(fmt.Sprintf("Unable to create reservation: conflict with reservations {IDs: %v}", e.ReservationIds))
}

type ReservationService struct {
	subjectsStore reservationsPort.SubjectsRepository
	reservationsRegistry reservationsPort.ReservationsRepository
	usersStore usersPort.UsersRepository
	clock ports.Clock
}

func NewReservationService(
	subjStore reservationsPort.SubjectsRepository,
	registry reservationsPort.ReservationsRepository,
	usersStore usersPort.UsersRepository,
	clock ports.Clock,
) *ReservationService {
	return &ReservationService{
		subjectsStore: subjStore,
		reservationsRegistry: registry,
		usersStore: usersStore,
		clock: clock,
	}
}

func (s *ReservationService) Create(cmd CreateReservation) (reservations.Reservation, error) {
	_, err := s.usersStore.Get(cmd.UserId)

	if err != nil {
		return reservations.Reservation{}, err
	}

	_, err = s.subjectsStore.Get(cmd.SubjectId)

	if err != nil {
		return reservations.Reservation{}, err
	}

	if s.clock.Current().After(cmd.From.Add(time.Minute)) {
		return reservations.Reservation{}, errors.New("Attempt to make a reservation in the past")
	}

	activeReservations := s.reservationsRegistry.ForPeriod(cmd.From, cmd.To).ForSubject(cmd.SubjectId)

	if len(activeReservations) > 0 {
		var ids []int
		for _, r := range activeReservations {
			ids = append(ids, r.Id)
		}
		return reservations.Reservation{}, AlreadyReservedError{ReservationIds: ids}
	}

	reservation := reservations.Reservation{
		Id: s.reservationsRegistry.NextIdentity(),
		UserId: cmd.UserId,
		SubjectId: cmd.SubjectId,
		Start: cmd.From,
		End: cmd.To,
	}
	s.reservationsRegistry.Add(reservation)

	return reservation, nil
}

func (s *ReservationService) Get(id int) (reservations.Reservation, error) {
	return s.reservationsRegistry.Get(id)
}
