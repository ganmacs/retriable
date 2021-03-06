package retriable

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	var i int
	fn := func() error {
		i++
		return nil
	}

	err := Retry(fn)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if i != 1 {
		t.Errorf("should be 1 but, %v", i)
	}
}

func TestFailRetry(t *testing.T) {
	var i int
	fn := func() error {
		i++
		return errors.New("error")
	}

	err := Retry(fn)

	if err == nil {
		t.Errorf("should have error")
	}

	if err.Error() != "error" {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if i != 3 {
		t.Errorf("should be 3 but, %v", i)
	}
}

func TestRetriesAndInterval(t *testing.T) {
	var i int
	opt := &Options{
		retries:  4,
		interval: 10 * time.Millisecond,
	}

	err := RetryWithOptions(func() error {
		i++
		return errors.New("error")
	}, opt)

	if err == nil {
		t.Errorf("should have error")
	}

	if err.Error() != "error" {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if i != 4 {
		t.Errorf("should be 4 but, %v", i)
	}
}

func TestMaxElapsedTimeRetry(t *testing.T) {
	fn := func() error {
		time.Sleep(60 * time.Millisecond)
		return errors.New("error")
	}

	err := RetryWithOptions(fn, &Options{
		maxElapsedTime: 50 * time.Millisecond,
	})

	if err == nil {
		t.Errorf("should have error")
	}

	if err.Error() != "Exceeds ElapsedTime" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestTimeoutRetry(t *testing.T) {
	fn := func() error {
		time.Sleep(1 * time.Second)
		return nil
	}

	err := RetryWithOptions(fn, &Options{
		timeout: 300 * time.Millisecond,
	})

	if err == nil {
		t.Errorf("should have error")
	}

	if err.Error() != "Timeout" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestTimeoutRetry2(t *testing.T) {
	i := 0
	fn := func() error {
		i++
		time.Sleep(200 * time.Millisecond)
		return errors.New("error")
	}

	err := RetryWithOptions(fn, &Options{
		interval: 1 * time.Millisecond,
		timeout:  100 * time.Millisecond,
	})

	// wait next retry if exist
	time.Sleep(200 * time.Millisecond)

	if err == nil {
		t.Errorf("should have error")
	}

	if err.Error() != "Timeout" {
		t.Errorf("unexpected error: %s", err.Error())
	}

	// goroutine should be killed
	if i != 1 {
		t.Errorf("should be 1 but, %v", i)
	}

}
