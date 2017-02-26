package retriable

import (
	"github.com/ganmacs/retriable/backoff"
	"time"
)

const (
	defaultInterval = 0
	defaultRetries  = 3
)

type Operation func() error

type Options struct {
	operation Operation
	backoff   backoff.BackOff
	retries   int
}

func Retry(op Operation) error {
	opt := &Options{
		operation: op,
		retries:   defaultRetries,
		backoff:   backoff.NewConstantBackOff(defaultInterval),
	}

	return doRetry(opt)
}

func doRetry(opt *Options) error {
	var bo = opt.backoff
	var retries = opt.retries

	var next time.Duration
	var err error

	bo.Reset()

	for i := 0; i < retries; i++ {
		if err = opt.operation(); err == nil {
			return nil
		}

		// TODO fix -1
		if next = bo.Next(); next == -1 {
			return err
		}

		time.Sleep(next)
	}

	return err
}
