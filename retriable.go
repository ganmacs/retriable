package retriable

import (
	"context"
	"errors"
	"time"

	"github.com/ganmacs/retriable/backoff"
)

const (
	defaultRetries        = 3
	defaultMaxElapsedTime = 15 * time.Minute
	defaultTimeout        = 0
)

type Operation func() error

type Options struct {
	retries        int
	maxElapsedTime time.Duration
	timeout        time.Duration
	interval       time.Duration
	backoff        backoff.BackOff
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
		retries:        defaultRetries,
		maxElapsedTime: defaultMaxElapsedTime,
		timeout:        defaultTimeout,
		backoff:        backoff.NewExponentialBackOff(),
	}

	return doRetry(op, opt)
}

func RetryWithOptions(op Operation, opt *Options) error {
	if opt != nil {
		if opt.retries == 0 {
			opt.retries = defaultRetries
		}

		if opt.maxElapsedTime == 0 {
			opt.maxElapsedTime = defaultMaxElapsedTime
		}

		if opt.backoff == nil && opt.interval != 0 {
			opt.backoff = backoff.NewExponentialBackOffWithInterval(opt.interval)
		} else if opt.backoff == nil {
			opt.backoff = backoff.NewExponentialBackOff()
		}
	}

	return doRetry(op, opt)
}

func timeout(t time.Duration, op func(context.Context) error) error {
	c := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c <- op(ctx)
	}()

	select {
	case err := <-c:
		return err
	case <-time.After(t):
		return errors.New("Timeout")
	}
}

func doRetry(op Operation, opt *Options) error {
	retry := func(ctx context.Context) error {
		var (
			next time.Duration
			err  error
		)

		clock := newClock()
		for i := 0; i < opt.retries; i++ {
			select {
			case <-ctx.Done():
				return errors.New("Timeout")
			default:
				if err = op(); err == nil {
					return nil
				}

				if opt.maxElapsedTime < clock.getElapsedTime() {
					return errors.New("Runngin too long")
				}

				next = opt.backoff.Next()
				time.Sleep(next)
			}
		}
		return err
	}

	if opt.timeout > 0 {
		return timeout(opt.timeout, retry)
	}

	return retry(context.Background())
}
