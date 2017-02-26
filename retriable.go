package retriable

import (
	"errors"
	"github.com/ganmacs/retriable/backoff"
	"time"
)

const (
	defaultRetries        = 3
	defaultMaxElapsedTime = 15 * time.Minute
	defaultTimeout        = 0
)

type Operation func() error

type Options struct {
	operation      Operation
	retries        int
	maxElapsedTime time.Duration
	timeout        time.Duration
}

type clock struct {
	startTime time.Time
}

func newClock() *clock {
	return &clock{
		startTime: time.Now(),
	}
}

func (c *clock) getElapsedTime() time.Duration {
	return time.Now().Sub(c.startTime)
}

func Retry(op Operation) error {
	opt := &Options{
		operation:      op,
		retries:        defaultRetries,
		maxElapsedTime: defaultMaxElapsedTime,
		timeout:        defaultTimeout,
	}

	return doRetry(backoff.NewExponentialBackOff(), opt)
}

func RetryWithOptions(op Operation, opt *Options) error {
	opt.operation = op

	if opt != nil {
		if opt.retries == 0 {
			opt.retries = defaultRetries
		}
	}

	return doRetry(backoff.NewExponentialBackOff(), opt)
}

func timeout(t time.Duration, op Operation) error {
	c := make(chan error)

	go func() {
		c <- op()
	}()

	select {
	case err := <-c:
		return err
	case <-time.After(t):
		return errors.New("Timeout")
	}
}

func doRetry(bo backoff.BackOff, opt *Options) error {
	retry := func() error {
		var retries = opt.retries
		var next time.Duration
		var err error
		var clock = newClock()

		if retries < 1 {
			return errors.New("retires should be 1 or more")
		}

		for i := 0; i < retries; i++ {
			if err = opt.operation(); err == nil {
				return nil
			}

			if opt.maxElapsedTime < clock.getElapsedTime() {
				return errors.New("Runngin too long")
			}

			// TODO fix -1
			if next = bo.Next(); next == -1 {
				return err
			}

			time.Sleep(next)
		}
		return err
	}

	if opt.timeout > 0 {
		return timeout(opt.timeout, retry)
	}

	return retry()
}
