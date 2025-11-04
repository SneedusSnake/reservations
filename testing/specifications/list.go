package specifications

import (
	"testing"
)

type ReservationsDriver interface{
	UserRequestsSubjectsList()
	UserSeesSubjects(subject ...string)
}

func ListSpecification(t testing.TB, driver ReservationsDriver) {
	driver.UserRequestsSubjectsList()

	driver.UserSeesSubjects("Subject #1", "Subject #2", "Subject #3")
}
