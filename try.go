// Package try provides retry functionality.
//     err := try.Do(context.TODO(), func(attempt int) (retry bool, err error) {
//       retry = attempt < 3 // try 3 times
//       err = doSomeThing()
//       return retry, err
//     })
//     if err != nil {
//       log.Fatalln("error:", err)
//     }
//
// Copyright 2020 lewisay. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package try

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

var (
	// DefaultMultiplier is the factor with which to multiply backoffs after a
	// failed retry. Should ideally be greater than 1.
	DefaultMultiplier = 1.6

	// DefaultJitter is the factor with which backoffs are randomized.
	DefaultJitter = 0.2

	// DefaultMaxRetries default max retries
	DefaultMaxRetries = 10

	// DefaultMinRetryBackoff default minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	DefaultMinRetryBackoff time.Duration = 8 * time.Millisecond

	// DefaultMaxRetryBackoff maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	DefaultMaxRetryBackoff time.Duration = 512 * time.Millisecond

	// ErrMaxRetriesReached error of ErrMaxRetriesReached
	ErrMaxRetriesReached = errors.New("exceeded retry limit")
)

type (
	// Func function that can be retried.
	Func func(attempt int) (retry bool, err error)

	// Options for retry function
	Options struct {
		// MaxRetries retries
		MaxRetries int

		// MinRetryBackoff minimum backoff between each retry.
		MinRetryBackoff time.Duration

		// MinRetryBackoff maximum backoff between each retry.
		MaxRetryBackoff time.Duration
	}
)

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func Do(ctx context.Context, fn Func, opts ...Options) error {
	opt := Options{DefaultMaxRetries, DefaultMinRetryBackoff, DefaultMaxRetryBackoff}
	for _, v := range opts {
		if v.MaxRetries != 0 {
			opt.MaxRetries = v.MaxRetries
		}

		if v.MinRetryBackoff != 0 {
			opt.MinRetryBackoff = v.MinRetryBackoff
		}

		if v.MaxRetryBackoff != 0 {
			opt.MaxRetryBackoff = v.MaxRetryBackoff
		}
	}

	attempt := 1
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			retry, err := fn(attempt)
			if !retry || err == nil {
				return err
			}
			attempt++
			if attempt > opt.MaxRetries {
				return ErrMaxRetriesReached
			}
			time.Sleep(Backoff(attempt, opt.MinRetryBackoff, opt.MaxRetryBackoff))
		}
	}
}

// Backoff backoff with jitter sleep to prevent overloaded conditions during intervals
// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md
func Backoff(retry int, minBackoff, maxBackoff time.Duration) time.Duration {
	if minBackoff == -1 || maxBackoff == -1 {
		return 0
	}

	if retry == 0 {
		return minBackoff
	}

	backoff, max := float64(minBackoff), float64(maxBackoff)
	for backoff < max && retry > 0 {
		backoff *= DefaultMultiplier
		retry--
	}

	if backoff > max {
		backoff = max
	}

	backoff *= 1 + DefaultJitter*(rand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}

	return time.Duration(backoff)
}
