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

type IRecentTimeChecker[K comparable] interface {
	SetLastTime(key K, now Time) bool
	IsExpired(key K, now Time) bool
}

func NewRecentTimeChecker[K comparable]() IRecentTimeChecker[K] {
	checker := &RecentTimeChecker[K]{}
	return checker.Init()
}

type RecentTimeChecker[K comparable] struct {
	//IRecentTimeChecker

	times map[K]Time
	lock  sync.Mutex
}

func (checker *RecentTimeChecker[K]) Init() IRecentTimeChecker[K] {
	checker.times = make(map[K]Time, 1024)
	return checker
}

// Override
func (checker *RecentTimeChecker[K]) SetLastTime(key K, now Time) bool {
	if now == nil {
		//panic("recent time empty")
		return false
	}
	// TODO: calibration clock

	checker.lock.Lock()
	defer checker.lock.Unlock()
	last := checker.times[key]
	if last == nil || TimeIsBefore(now, last) {
		checker.times[key] = now
		return true
	}
	return false
}

// Override
func (checker *RecentTimeChecker[K]) IsExpired(key K, now Time) bool {
	if now == nil {
		//panic("recent time empty")
		return true
	}
	last := checker.times[key]
	return last != nil && TimeIsAfter(now, last)
}
