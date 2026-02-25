/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package utils

import (
	"sync"

	. "github.com/dimchat/mkm-go/types"
)

type IFrequencyChecker[K comparable] interface {
	IsExpired(key K, now Time, force bool) bool
}

// FrequencyChecker provides thread-safe frequency control for duplicate query prevention
//
// Generic type K: Comparable key type (e.g., ID, string) to track query identifiers
//
// Core functionality: Tracks when each key was last used to prevent duplicate/too-frequent queries
type FrequencyChecker[K comparable] struct {
	//IFrequencyChecker

	expires Duration
	records map[K]Time
	lock    sync.Mutex
}

func NewFrequencyChecker[K comparable](lifeSpan Duration) IFrequencyChecker[K] {
	return &FrequencyChecker[K]{
		expires: lifeSpan,
		records: make(map[K]Time, 1024),
		lock:    sync.Mutex{},
	}
}

func (checker *FrequencyChecker[K]) checkExpired(key K, now Time) bool {
	expired := checker.records[key]
	if expired != nil && TimeIsAfter(now, expired) {
		// record exists and not expired yet
		return false
	}
	checker.records[key] = checker.expires.AddTo(now)
	return true
}

func (checker *FrequencyChecker[K]) forceExpired(key K, now Time) bool {
	checker.records[key] = checker.expires.AddTo(now)
	return true
}

// Override
func (checker *FrequencyChecker[K]) IsExpired(key K, now Time, force bool) bool {
	if now == nil {
		now = TimeNow()
	}
	checker.lock.Lock()
	defer checker.lock.Unlock()
	// if force == true:
	//     ignore last updated time, force to update now
	// else:
	//     check last update time
	if force {
		return checker.forceExpired(key, now)
	}
	return checker.checkExpired(key, now)
}
