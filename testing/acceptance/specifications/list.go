package specifications

import (
	"testing"
	"github.com/SneedusSnake/Reservations/testing/acceptance/drivers"
)

func ListSpecification(t testing.TB, driver drivers.Reservations) {
	driver.UserRequestsSubjectsList()

	driver.UserSeesSubjects("Subject#1", "Subject#2", "Subject#3")
}
