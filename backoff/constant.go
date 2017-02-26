package backoff

import "time"

type ConstantBackOff struct {
	Interval time.Duration
}

func NewConstantBackOff(t time.Duration) *ConstantBackOff {
	return &ConstantBackOff{Interval: t}
}

func (bo *ConstantBackOff) Next() time.Duration {
	return bo.Interval
}
