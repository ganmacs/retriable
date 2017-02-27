# Retriable

[![Build Status](https://travis-ci.org/ganmacs/retriable.svg?branch=master)](https://travis-ci.org/ganmacs/retriable)

Retriable is a simple Library to retry failed code with randomized [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff) algorithm.
This Library is inspired by [retriable](https://github.com/kamui/retriable) which is implemented in Ruby.

# Examples

Here is a simple example of what retriable provides

```go
retriable.Retry(func() error {
	// your code
	return nil
})
```

You can specify retry time and interval time if you want. This example will retry 5 times and retries interval is 1 second.


```go
opt := retriable.Options{
  retries:  5,
  interval: 1 * time.Second,
}

retriable.RetryWithOptions(func() error {
  // your code
  return nil
}, opt)
```

You can also specify a timeout as follow.

```go
opt := retriable.Options{
	timeout: 5 * time.Second,
}

retriable.RetryWithOptions(func() error {
	// your code
	return nil
}, opt)
```

You can retry your code every second instead of using exponential backoff.

```go
opt := retriable.Options{
	backoff: backoff.NewConstantBackOff(1 * time.Second),
}

retriable.RetryWithOptions(func() error {
	// your code
	return nil
}, opt)
```

# Options

`retries` (default: 3) - Number of attempts to run your code

`interval` (default: 500 msec) - The interval between tries

`maxElapsedTime` (default: 15 min) - The maximum amount of total time

`timeout` (default: no) - Number of Second
