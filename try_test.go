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
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackoff(t *testing.T) {
	minBackoff := DefaultMinRetryBackoff
	maxBackOff := DefaultMaxRetryBackoff

	at := assert.New(t)
	for i := -1; i < 10; i++ {
		backoff := Backoff(i, minBackoff, maxBackOff)
		v := int64(backoff)

		at.GreaterOrEqual(v, int64(minBackoff))
		at.LessOrEqual(v, int64(maxBackOff))
	}

	backoff := Backoff(math.MaxInt64, -1, 0)
	at.EqualValues(0, backoff)

	backoff = Backoff(math.MaxInt64, -1, -0)
	at.EqualValues(0, backoff)

	DefaultJitter = 10
	backoff = Backoff(math.MaxInt64, 1, 0)
	at.EqualValues(0, backoff)
}

func TestDo1(t *testing.T) {
	var counter int
	maxRetries := 2

	err := Do(context.TODO(), func(attempt int) (retry bool, err error) {
		counter++

		retry = attempt < maxRetries
		if attempt == maxRetries {
			return retry, nil
		}

		return retry, errors.New("something error")
	}, Options{
		MaxRetries:      maxRetries,
		MinRetryBackoff: DefaultMinRetryBackoff,
		MaxRetryBackoff: DefaultMaxRetryBackoff,
	})

	at := assert.New(t)
	at.NoError(err)
	at.EqualValues(maxRetries, counter)
}

func TestDo2(t *testing.T) {
	err := Do(context.TODO(), func(attempt int) (retry bool, err error) {
		return true, errors.New("something error")
	})

	assert.True(t, errors.Is(err, ErrMaxRetriesReached))
}

func TestDo3(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := Do(ctx, func(attempt int) (retry bool, err error) {
		return true, errors.New("something error")
	})

	assert.True(t, errors.Is(err, context.DeadlineExceeded))
}
