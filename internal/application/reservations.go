package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	"github.com/SneedusSnake/Reservations/internal/ports"
	reservationsPort "github.com/SneedusSnake/Reservations/internal/ports/reservations"
	usersPort "github.com/SneedusSnake/Reservations/internal/ports/users"
	readmodel "github.com/SneedusSnake/Reservations/internal/read_model"
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
	reservationsStore reservationsPort.ReservationsRepository
	reservationsReadStore reservationsPort.ReservationsReadRepository
	usersStore usersPort.UsersRepository
	clock ports.Clock
}

func NewReservationService(
	subjStore reservationsPort.SubjectsRepository,
	registry reservationsPort.ReservationsRepository,
	reservationsReadStore reservationsPort.ReservationsReadRepository,
	usersStore usersPort.UsersRepository,
	clock ports.Clock,
) *ReservationService {
	return &ReservationService{
		subjectsStore: subjStore,
		reservationsStore: registry,
		reservationsReadStore: reservationsReadStore,
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

	activeReservations := s.reservationsStore.ForPeriod(cmd.From, cmd.To).ForSubject(cmd.SubjectId)

	if len(activeReservations) > 0 {
		var ids []int
		for _, r := range activeReservations {
			ids = append(ids, r.Id)
		}
		return reservations.Reservation{}, AlreadyReservedError{ReservationIds: ids}
	}
	id, err := s.reservationsStore.NextIdentity()
	if err != nil {
		return reservations.Reservation{}, err
	}

	reservation := reservations.Reservation{
		Id: id,
		UserId: cmd.UserId,
		SubjectId: cmd.SubjectId,
		Start: cmd.From,
		End: cmd.To,
	}
	s.reservationsStore.Add(reservation)

	return reservation, nil
}

func (s *ReservationService) Get(id int) (reservations.Reservation, error) {
	return s.reservationsStore.Get(id)
}

type RemoveReservations struct {
	UserId int
	SubjectId int
}

func (s *ReservationService) Remove(cmd RemoveReservations) error {
	//checking reservations for a year in advance will suffice for now
	subjReservations := s.reservationsStore.ForPeriod(s.clock.Current(), s.clock.Current().Add(time.Hour*8760)).ForUser(cmd.UserId).ForSubject(cmd.SubjectId)

	if len(subjReservations) == 0 {
		return errors.New("No active reservations found")
	}

	for _, r := range subjReservations {
		s.reservationsStore.Remove(r.Id)
	}

	return nil
}

func (s *ReservationService) ActiveReservations(t time.Time, tags ...string) ([]readmodel.Reservation, error) {
	return s.reservationsReadStore.Active(t, tags...)
}
