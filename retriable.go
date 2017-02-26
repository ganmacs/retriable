package retriable

import (
	"github.com/ganmacs/retriable/backoff"
	"time"
)

const (
	defaultRetries = 3
)

type Operation func() error

type Options struct {
	operation Operation
	retries   int
}

func Retry(op Operation) error {
	opt := &Options{
		operation: op,
		retries:   defaultRetries,
	}

	return doRetry(backoff.NewExponentialBackOff(), opt)
}

func doRetry(bo backoff.BackOff, opt *Options) error {
	var retries = opt.retries

	var next time.Duration
	var err error

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
