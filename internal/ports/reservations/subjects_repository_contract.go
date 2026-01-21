package reservations

import (
	"slices"
	"testing"

	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	"github.com/SneedusSnake/Reservations/internal/utils"
	"github.com/alecthomas/assert/v2"
)

type SubjectsRepositoryContract struct {
	NewStore func() SubjectsRepository
}

func (s SubjectsRepositoryContract) Test (t *testing.T) {
	store := s.NewStore()
	cleanUp := func (subjects ...reservations.Subject) {
		for _, subject := range subjects {
			store.Remove(subject.Id)
		}
	}

	t.Run("it returns error when subject was not found", func(t *testing.T) {
		_, err := store.Get(1234)

		if err == nil {
			t.Error("Expected to see error, got nil")
		}
	})

	t.Run("it adds a new subject to the store", func(t *testing.T) {
		subject := reservations.Subject{Id: 1, Name: "Test subject"}

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
		subject := reservations.Subject{Id: 1, Name: "Test subject"}
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

	t.Run("it returns list of all subjects", func(t *testing.T) {
		store.Add(reservations.Subject{Id: 1, Name: "Subject 1"})
		store.Add(reservations.Subject{Id: 2, Name: "Subject 2"})
		store.Add(reservations.Subject{Id: 3, Name: "Subject 3"})

		subjects := store.List()

		if len(subjects) != 3 {
			t.Errorf("Expected 3 subjects, got %d", len(subjects))
		}

		cleanUp(subjects...)
	})

	t.Run("it can add tags and filter by them", func(t *testing.T) {
		subjects := []reservations.Subject{
			{Id: 1, Name: "Conference room #1"},
			{Id: 2, Name: "Conference room #2"},
			{Id: 3, Name: "Conference room #3"},
			{Id: 4, Name: "Conference room #4"},
		}

		for _, subject := range subjects {
			err := store.Add(subject)
			if err != nil {
				t.Fatal(err)
			}
		}
		store.AddTag(2, "spacious")
		store.AddTag(4, "spacious")
		store.AddTag(4, "soundproof")

		spaciousRooms := store.GetByTags([]string{"spacious"})
		spaciousAndSoundProof := store.GetByTags([]string{"spacious", "soundproof"})

		if len(spaciousRooms) != 2 {
			t.Errorf("expected to find 2 spacious rooms, found %d instead", len(spaciousRooms))
		}

		if spaciousRooms[0].Id != 2 || spaciousRooms[1].Id != 4 {
			t.Errorf("Unexpected item retrieved by 'spacious' tag, expected subjects with ids 2 and 4, got %v", spaciousRooms)
		}

		if len(spaciousAndSoundProof) != 1 {
			t.Errorf("expected to find 1 spacious soundproof room, found %d instead", len(spaciousAndSoundProof))
		}
		cleanUp(subjects...)
	})

	t.Run("it cannot add same tag to the same subject twice", func(t *testing.T) {
		err := store.Add(reservations.Subject{Id: 1, Name: "Test"})

			if err != nil {
				t.Fatal(err)
			}
	})

	t.Run("it returns all tags of a subject", func(t *testing.T) {
		subject := reservations.Subject{Id: 1, Name: "Test"}
		err := store.Add(subject)
		if err != nil {
			t.Fatal(err)
		}
		expectedTags := []string{"tag 1", "tag 2", "tag 3"}
		for _, tag := range expectedTags {
			store.AddTag(subject.Id, tag)
		}

		tags, err := store.GetTags(subject.Id)
		if err != nil {
			t.Fatal(err)
		}

		assert.SliceContains(t, tags, expectedTags[0])
		assert.SliceContains(t, tags, expectedTags[1])
		assert.SliceContains(t, tags, expectedTags[2])

		cleanUp(subject)
	})

	t.Run("it generates next ID", func(t *testing.T) {
		ch := make(chan int, 5)
		var ids []int

		for range 5 {
			go (func (c chan int) {
				id, _ := store.NextIdentity()
				c <- id
			})(ch)
		}
		
		for range 5 {
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

