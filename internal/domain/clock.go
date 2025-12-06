package domain

import "time"

type Clock interface {
	Current() time.Time
}
