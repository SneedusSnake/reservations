package reservations

import (
	"slices"
	"testing"

	"github.com/SneedusSnake/Reservations/utils"
)

type SubjectStoreContract struct {
	NewStore func() SubjectsStore
}

func (s SubjectStoreContract) Test (t *testing.T) {
	store := s.NewStore()
	cleanUp := func (subject Subject) {
		store.Remove(subject.Id)
	}

	t.Run("it returns error when subject is not found", func(t *testing.T) {
		_, err := store.Get(1234)

		if err == nil {
			t.Error("Expected to see error, got nil")
		}
	})

	t.Run("it adds a new subject to the store", func(t *testing.T) {
		subject := Subject{1, "Test subject"}

		err := store.Add(subject)

		if err != nil {
			t.Fatal(err)
		}

		foundSubject, err := store.Get(subject.Id)

		if err != nil {
			t.Fatal(err)
		}

		if foundSubject != subject {
			t.Errorf("Expected to find %v, got %v", subject, foundSubject)
		}

		cleanUp(subject)
	})

	t.Run("it removes a subject from the store", func(t *testing.T) {
		subject := Subject{1, "Test subject"}
		err := store.Add(subject)

		if err != nil {
			t.Fatal(err)
		}

		store.Remove(subject.Id)

		if err != nil {
			t.Fatal(err)
		}

		_, err = store.Get(subject.Id)

		if err == nil {
			t.Errorf("Expected subject with id %d to be removed from the store", subject.Id)
		}
	})

	t.Run("it generates next ID", func(t *testing.T) {
		ch := make(chan int, 5)
		var ids []int

		for i := 0; i < 5; i++ {
			go (func (c chan int) {
				c <- store.NextIdentity()
			})(ch)
		}
		
		for i := 0; i < 5; i++ {
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

