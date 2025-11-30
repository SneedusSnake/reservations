package application

import (
	"fmt"
	"time"

	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/SneedusSnake/Reservations/domain/users"
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

type CreateReservationHandler struct {
	subjectsStore reservations.SubjectsStore
	reservationsRegistry reservations.ReservationsRegistry
	usersStore users.UsersStore
}

func NewCreateReservationHandler(subjStore reservations.SubjectsStore, registry reservations.ReservationsRegistry, usersStore users.UsersStore) *CreateReservationHandler {
	return &CreateReservationHandler{
		subjectsStore: subjStore,
		reservationsRegistry: registry,
		usersStore: usersStore,
	}
}

func (h *CreateReservationHandler) Handle(cmd CreateReservation) (reservations.Reservation, error) {
	_, err := h.usersStore.Get(cmd.UserId)

	if err != nil {
		return reservations.Reservation{}, err
	}

	_, err = h.subjectsStore.Get(cmd.SubjectId)

	if err != nil {
		return reservations.Reservation{}, err
	}
	activeReservations := h.reservationsRegistry.ForPeriod(cmd.From, cmd.To).ForSubject(cmd.SubjectId)

	if len(activeReservations) > 0 {
		var ids []int
		for _, r := range activeReservations {
			ids = append(ids, r.Id)
		}
		fmt.Print(ids)
		return reservations.Reservation{}, AlreadyReservedError{ReservationIds: ids}
	}

	reservation := reservations.Reservation{
		Id: h.reservationsRegistry.NextIdentity(),
		UserId: cmd.UserId,
		SubjectId: cmd.SubjectId,
		Start: cmd.From,
		End: cmd.To,
	}
	h.reservationsRegistry.Add(reservation)

	return reservation, nil
}
