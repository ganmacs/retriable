package backoff

import (
	"testing"
	"time"
)

func TestExopnetialNext(t *testing.T) {
	var delta = 0.0
	var defaultIntervalFloat = float64(defaultInterval)

	// 1.5 is default Multiplier
	var intervals = []float64{
		1.5 * defaultIntervalFloat,
		1.5 * 1.5 * defaultIntervalFloat,
		1.5 * 1.5 * 1.5 * defaultIntervalFloat,
		1.5 * 1.5 * 1.5 * 1.5 * defaultIntervalFloat,
		1.5 * 1.5 * 1.5 * 1.5 * 1.5 * defaultIntervalFloat,
		1.5 * 1.5 * 1.5 * 1.5 * 1.5 * 1.5 * defaultIntervalFloat,
		1.5 * 1.5 * 1.5 * 1.5 * 1.5 * 1.5 * 1.5 * defaultIntervalFloat,
	}

	backoff := NewExponentialBackOff()
	for _, interval := range intervals {
		delta = interval * defaultRandomizationFactor
		assertBetween(t, float64(backoff.Next()), interval-delta, interval+delta)
	}
}

func TestMaxInterval(t *testing.T) {
	var delta = 0.0
	var defaultIntervalFloat = float64(defaultInterval)

	backoff := NewExponentialBackOff()
	backoff.MaxInterval = 1 * time.Second

	// 1.5 is default Multiplier
	var intervals = []float64{
		1.5 * defaultIntervalFloat,
		float64(backoff.MaxInterval),
		float64(backoff.MaxInterval),
	}

	for _, interval := range intervals {
		delta = interval * defaultRandomizationFactor
		assertBetween(t, float64(backoff.Next()), interval-delta, interval+delta)
	}

}

func assertBetween(t *testing.T, v float64, left float64, right float64) {
	if !(v >= left && v <= right) {
		t.Errorf("got: %v, expected between %v to %v", time.Duration(v), time.Duration(left), time.Duration(right))
	}
}
