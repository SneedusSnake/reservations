package users

import "testing"

type UsersStoreContract struct {
	NewStore func() UsersStore
}

func (s UsersStoreContract) Test(t *testing.T) {
	store := s.NewStore()

	t.Run("it returns error when user was not found", func(t *testing.T) {
		_, err := store.Get(1234)
	
		if err == nil {
			t.Error("Expected to see error, got nil")
		}
	})

	t.Run("it adds a new user to the store", func(t *testing.T) {
		user := User{Id: store.NextIdentity(), Name: "Adam"}
		t.Cleanup(func() {
			store.Remove(user.Id)
		})

		err := store.Add(user)

		if err != nil {
			t.Fatal(err)
		}

		foundUser, err := store.Get(user.Id)
		if err != nil {
			t.Fatal(err)
		}

		if foundUser != user {
			t.Errorf("Expected to find %v, got %v", user, foundUser)
		}
	})
	
	t.Run("It cannot add user with same id twice", func(t *testing.T) {
		user := User{Id: store.NextIdentity(), Name: "Eve"}
		t.Cleanup(func() {
			store.Remove(user.Id)
		})

		err := store.Add(user)

		if err != nil {
			t.Fatal(err)
		}

		err = store.Add(user)

		if err == nil {
			t.Errorf("Expected a duplicate error")
		}
	})
}
