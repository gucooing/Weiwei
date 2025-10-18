// Copyright 2025 gucooing, gucooing@alsl.xyz
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backoff

import (
	"io"
	"math"
	"math/rand"
	"time"
)

type Backoff interface {
	Interval() time.Duration
	NextBackoff() (time.Duration, bool)
	Reset()
}

type ExponentialBackoff struct {
	BaseInterval time.Duration
	MaxRetries   int
	MaxInterval  time.Duration
	retryCount   int
}

func (eb *ExponentialBackoff) Interval() time.Duration {
	return eb.BaseInterval
}

func (eb *ExponentialBackoff) NextBackoff() (time.Duration, bool) {
	if eb.MaxRetries > 0 &&
		eb.retryCount >= eb.MaxRetries {
		return 0, false
	}

	if eb.MaxInterval <= 0 {
		return eb.BaseInterval, true
	}

	delay := eb.BaseInterval * time.Duration(math.Pow(2, float64(eb.retryCount)))
	jitter := time.Duration(rand.Float64() * 0.3 * float64(delay))
	delay += jitter
	if delay > eb.MaxInterval ||
		delay <= 0 {
		delay = eb.MaxInterval
	}
	eb.retryCount++
	return delay, true
}

func (eb *ExponentialBackoff) Reset() {
	eb.retryCount = 0
}

func BackoffStart(f func() error, doneChan <-chan struct{}, backoff Backoff) (lastErr error) {
	delay := backoff.Interval()
	var ok bool
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-doneChan:
			return io.EOF
		default:
		}

		if err := f(); err != nil {
			lastErr = err
		} else {
			return nil
		}

		delay, ok = backoff.NextBackoff()
		if !ok {
			return lastErr
		}
		ticker.Reset(delay)
		select {
		case <-doneChan:
			return io.EOF
		case <-ticker.C:
		}
	}
}
