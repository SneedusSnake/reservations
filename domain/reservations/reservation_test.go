package reservations_test

import (
	"testing"
	"time"
	"github.com/SneedusSnake/Reservations/domain/reservations"
)

func TestReservation(t *testing.T) {
	start, err := time.Parse("2006-01-02 15:04:05", "2021-01-02 14:00:00")
	if err != nil {
		t.Fatal(err)
	}
	end, err := time.Parse("2006-01-02 15:04:05", "2021-01-02 14:30:00")
	if err != nil {
		t.Fatal(err)
	}


	reservation := reservations.Reservation{1, 2, 3, start, end}

	cases := []struct{
		Description string
		Reservation reservations.Reservation
		Date string
		Expected bool
	} {
		{"Is not active given past date", reservation, "2021-01-01 14:02:00", false},
		{"Is not active given a second before start", reservation, "2021-01-02 13:59:59", false},
		{"Is active given same date as start", reservation, "2021-01-02 14:00:00", true},
		{"Is active given a second before end", reservation, "2021-01-02 14:29:59", true},
		{"Is not active given end date", reservation, "2021-01-02 14:30:00", false},
		{"Is not active given future date", reservation, "2021-01-03 14:00:00", false},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			date, err := time.Parse("2006-01-02 15:04:05", test.Date)
			if err != nil {
				t.Fatal(err)
			}
			result := test.Reservation.ActiveAt(date)

			if result != test.Expected {
				t.Errorf("expected %t got %t", test.Expected, result)
			}
		})
	}
}

func TestReservations(t *testing.T) {
	t.Run("it returns reservations filtered by subject id", func(t *testing.T) {
		reservations := reservations.Reservations{
			reservations.Reservation{1, 1, 1, time.Now(), time.Now()},
			reservations.Reservation{2, 2, 2, time.Now(), time.Now()},
			reservations.Reservation{3, 3, 3, time.Now(), time.Now()},
			reservations.Reservation{4, 4, 2, time.Now(), time.Now()},
		}

		result := reservations.ForSubject(reservations[1].SubjectId)


		if len(result) != 2 {
			t.Errorf("expected to get 2 reservations, got %d", len(result))
		}

		if result[0] != reservations[1] {
			t.Errorf("Expected %v, got %v", reservations[1], result[0])
		}
		if result[1] != reservations[3] {
			t.Errorf("Expected %v, got %v", reservations[3], result[1])
		}
	})
}
