package reservations_test

import (
	"testing"
	"time"
	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
)

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
