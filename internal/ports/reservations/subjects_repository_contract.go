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
	store := subjectsRepositoryHelper{SubjectsRepository: s.NewStore(), t: t}
	cleanUp := store.CleanUp

	t.Run("it returns error when the subject was not found", func(t *testing.T) {
		_, err := store.Get(1234)

		assert.Error(t, err)
	})

	t.Run("it adds a new subject to the store", func(t *testing.T) {
		cleanUp(t)
		subject := reservations.Subject{Id: 1, Name: "Test subject"}

		err := store.Add(subject)
		assert.NoError(t, err)

		foundSubject, err := store.Get(subject.Id)
		assert.NoError(t, err)

		assert.Equal(t, subject, foundSubject)
	})

	t.Run("it gets the subject by name", func(t *testing.T) {
		cleanUp(t)
		subjects := store.SubjectsExist("first", "second", "third")

		s, err := store.GetByName("second")
		assert.NoError(t, err)
		assert.Equal(t, subjects[1].Id, s.Id)

		s, err = store.GetByName("does not exist")
		assert.Error(t, err)
	})

	t.Run("it removes the subject from the store", func(t *testing.T) {
		cleanUp(t)
		subject := store.SubjectExists("Test Subject")

		err := store.Remove(subject.Id)
		assert.NoError(t, err)

		_, err = store.Get(subject.Id)
		assert.Error(t, err)
	})

	t.Run("it returns list of all subjects", func(t *testing.T) {
		cleanUp(t)
		store.SubjectsExist("Subject 1", "Subject 2", "Subject 3")

		subjects, err := store.List()
		assert.NoError(t, err)
		assert.Equal(t, 3, len(subjects))
	})

	t.Run("it can add tags and filter by them", func(t *testing.T) {
		cleanUp(t)
		subjects := store.SubjectsExist(
			"Conference room #1",
			"Conference room #2",
			"Conference room #3",
			"Conference room #4",
		)
		expectedSpacious := subjects[1]
		expectedSpaciousAndSoundProof := subjects[3]

		store.AddTag(expectedSpacious.Id, "spacious")
		store.AddTag(expectedSpaciousAndSoundProof.Id, "spacious")
		store.AddTag(expectedSpaciousAndSoundProof.Id, "soundproof")

		spaciousRooms, err := store.GetByTags([]string{"spacious"})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(spaciousRooms))
		assert.Equal(t, expectedSpacious.Id, spaciousRooms[0].Id)
		assert.Equal(t, expectedSpaciousAndSoundProof.Id, spaciousRooms[1].Id)

		spaciousAndSoundProof, err := store.GetByTags([]string{"spacious", "soundproof"})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(spaciousAndSoundProof))
		assert.Equal(t, expectedSpaciousAndSoundProof.Id, spaciousAndSoundProof[0].Id)

	})

	t.Run("it cannot add same tag to the same subject twice", func(t *testing.T) {
		cleanUp(t)
		s := store.SubjectExists("Test Subject")
		err := store.AddTag(s.Id, "Test")
		assert.NoError(t, err)

		err = store.AddTag(s.Id, "Test")
		assert.Error(t, err)
	})

	t.Run("it returns all tags of a subject", func(t *testing.T) {
		cleanUp(t)
		subject := store.SubjectExists("Test")
		expectedTags := []string{"tag 1", "tag 2", "tag 3"}
		for _, tag := range expectedTags {
			store.AddTag(subject.Id, tag)
		}

		tags, err := store.GetTags(subject.Id)
		assert.NoError(t, err)

		assert.SliceContains(t, tags, expectedTags[0])
		assert.SliceContains(t, tags, expectedTags[1])
		assert.SliceContains(t, tags, expectedTags[2])
	})

	t.Run("it generates next ID", func(t *testing.T) {
		cleanUp(t)
		ch := make(chan int, 5)
		var ids []int

		for range 5 {
			go (func (c chan int, t *testing.T) {
				id, err := store.NextIdentity()
				assert.NoError(t, err)
				c <- id
			})(ch, t)
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

type subjectsRepositoryHelper struct {
	SubjectsRepository
	t testing.TB
}

func (h *subjectsRepositoryHelper) SubjectExists(name string) reservations.Subject {
	id, err := h.NextIdentity()
	assert.NoError(h.t, err)
	s := reservations.Subject{Id: id, Name: name}
	err = h.Add(s)
	assert.NoError(h.t, err)

	return s
}

func (h *subjectsRepositoryHelper) SubjectsExist(names ...string) reservations.Subjects {
	var subjects reservations.Subjects

	for _, name := range names {
		subjects = append(subjects, h.SubjectExists(name))
	}

	return subjects
}

func (h *subjectsRepositoryHelper) CleanUp(t testing.TB) {
	t.Cleanup(func() {
		subjects, err := h.List()
		assert.NoError(t, err)

		for _, s := range subjects {
			h.Remove(s.Id)
		}
	})
}
