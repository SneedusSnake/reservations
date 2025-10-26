package specifications

import (
	"testing"
	"github.com/alecthomas/assert/v2"
)

type SubjectsViewer interface{
	List() ([]string, error)
}

func ListSpecification(t testing.TB, viewer SubjectsViewer) {
	subjects, err := viewer.List()
	assert.Equal(t, err, nil)
	assert.Equal(t, subjects, []string{"Test Subject #1", "Test Subject #2", "Test Subject #3"})
}
