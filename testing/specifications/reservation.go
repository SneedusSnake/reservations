package specifications

import (
	"testing"
	"github.com/SneedusSnake/Reservations/testing/drivers"
)

func ReserveSubjectSpecification(t testing.TB, driver drivers.Reservations) {
	driver.ClockSet("12:00")
	driver.UserRequestsReservationForSubject("Alice", "Subject#1", 30)
	driver.UserAcquiredReservationForSubject("Alice", "Subject#1", "12:30")

	driver.ClockSet("12:10")
	driver.UserRequestsReservationForSubject("Bob", "Subject#1", 30)
	driver.SubjectHasAlreadyBeenReservedBy("Alice", "12:30")
}

func RemoveReservationSpecification(t testing.TB, driver drivers.Reservations) {
	driver.ClockSet("13:00")
	driver.UserRequestsReservationForSubject("Alice", "Subject#2", 5)
	driver.UserAcquiredReservationForSubject("Alice", "Subject#2", "13:05")

	driver.UserRequestsReservationRemoval("Alice", "Subject#2")

	driver.UserRequestsReservationForSubject("Bob", "Subject#2", 30)
	driver.UserAcquiredReservationForSubject("Bob", "Subject#2", "13:30")
}

func ListReservedSubjects(t testing.TB, driver drivers.Reservations) {
	driver.ClockSet("14:00")
	driver.UserRequestsReservationForSubject("Alice", "Subject#1", 5)
	driver.UserRequestsReservationForSubject("Bob", "Subject#3", 10)

	driver.UserRequestsReservationsList()

	driver.UserSeesReservations("Alice Subject#1 14:05", "Bob Subject#3 14:10")
}
