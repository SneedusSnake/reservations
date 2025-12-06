package ports

import "time"

type Clock interface {
	Current() time.Time
}
