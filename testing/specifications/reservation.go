package specifications

import (
	"testing"
	"github.com/SneedusSnake/Reservations/testing/drivers"
)

func ReserveSubjectSpecification(t testing.TB, driver drivers.Reservations) {
	driver.ClockSet("12:00")
	driver.UserRequestsReservationForSubject("Alice", "Subject#1", 30)
	driver.UserAcquiredReservationForSubject("Alice", "Subject#1", "12:30")

	driver.UserRequestsReservationForSubject("Bob", "Subject#1", 30)
	driver.SubjectHasAlreadyBeenReservedBy("Alice", "12:30")
}
