package backoff

import (
	"time"
)

type BackOff interface {
	Next() time.Duration
}
