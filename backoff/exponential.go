package backoff

import (
	"math"
	"math/rand"
	"time"
)

/*
RandomizedInterval = CurrentInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])
*/

const (
	DefaultInterval            = 500 * time.Millisecond
	DefaultMaxInterval         = 60 * time.Second
	DefaultRandomizationFactor = 0.25
	DefaultMultiplier          = 1.5
)

type ExponentialBackOff struct {
	CurrentInterval     time.Duration
	MaxInterval         time.Duration
	RandomizationFactor float64
	Multiplier          float64
}

func NewExponentialBackOff() *ExponentialBackOff {
	return &ExponentialBackOff{
		CurrentInterval:     DefaultInterval,
		MaxInterval:         DefaultMaxInterval,
		RandomizationFactor: DefaultRandomizationFactor,
		Multiplier:          DefaultMultiplier,
	}
}

func (bo *ExponentialBackOff) getRandomizedInterval(interval float64) float64 {
	delta := bo.RandomizationFactor * interval
	min := interval - delta
	max := interval + delta
	dt := rand.Float64() * (max - min)

	return min + dt
}

func (bo *ExponentialBackOff) Next() time.Duration {
	interval := math.Min(float64(bo.MaxInterval), bo.Multiplier*float64(bo.CurrentInterval))
	bo.CurrentInterval = time.Duration(interval)

	return time.Duration(bo.getRandomizedInterval(interval))
}
