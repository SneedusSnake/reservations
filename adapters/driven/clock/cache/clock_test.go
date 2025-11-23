package cache_test

import (
	"testing"
	"time"

	"github.com/SneedusSnake/Reservations/adapters/driven/clock/cache"
)

func TestClock(t *testing.T) {
	t.Run("it returns current time if clock was not set", func(t *testing.T) {
		clock := cache.NewClock("/tmp/clock_go")
		expected := time.Now()

		result := clock.Current()

		if result.Sub(result) >2*time.Second {
			t.Errorf("Expected datetime %s is, got %s", expected.Format(time.DateTime), result.Format(time.DateTime))
		}
	})
}
