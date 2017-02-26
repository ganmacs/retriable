package backoff

import "time"

type ConstantBackOff struct {
	Interval time.Duration
}

func (bo *ConstantBackOff) Next() time.Duration {
	return bo.Interval
}

func (bo *ConstantBackOff) Reset() {}

func NewConstantBackOff(t time.Duration) *ConstantBackOff {
	return &ConstantBackOff{Interval: t}
}
