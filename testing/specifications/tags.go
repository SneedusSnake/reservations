package specifications

import (
	"testing"
	"github.com/SneedusSnake/Reservations/testing/drivers"
)

func SubjectTagsSpecification(t testing.TB, driver drivers.Reservations) {
	driver.UserRequestsSubjectTags("Subject#1")
	driver.UserSeesSubjectTags("Subject#1", "This_is_a_first_subject", "test")

	driver.UserRequestsSubjectTags("Subject#2")
	driver.UserSeesSubjectTags("Subject#2", "This_is_a_second_subject", "test")
}
