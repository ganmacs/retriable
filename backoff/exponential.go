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
	defaultInterval            = 500 * time.Millisecond
	defaultMaxInterval         = 60 * time.Second
	defaultRandomizationFactor = 0.25
	defaultMultiplier          = 1.5
)

type ExponentialBackOff struct {
	CurrentInterval     time.Duration
	MaxInterval         time.Duration
	RandomizationFactor float64
	Multiplier          float64
}

func NewExponentialBackOff() *ExponentialBackOff {
	return &ExponentialBackOff{
		CurrentInterval:     defaultInterval,
		MaxInterval:         defaultMaxInterval,
		RandomizationFactor: defaultRandomizationFactor,
		Multiplier:          defaultMultiplier,
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
