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
	. "github.com/dimchat/mkm-go/types"
)

func TimeIsBefore(otherTime, thisTime Time) bool {
	return TimeToFloat64(thisTime) < TimeToFloat64(otherTime)
}

func TimeIsAfter(otherTime, thisTime Time) bool {
	return TimeToFloat64(thisTime) > TimeToFloat64(otherTime)
}

//goland:noinspection GoSnakeCaseUsage
const (
	/**
	 * Hours per day.
	 */
	HOURS_PER_DAY = 24
	/**
	 * Minutes per hour.
	 */
	MINUTES_PER_HOUR = 60
	/**
	 * Seconds per minute.
	 */
	SECONDS_PER_MINUTE = 60
	/**
	 * Milliseconds per second.
	 */
	MILLIS_PER_SECOND = 1000
	/**
	 * Milliseconds per minute.
	 */
	MILLIS_PER_MINUTE = MILLIS_PER_SECOND * SECONDS_PER_MINUTE
	/**
	 * Milliseconds per hour.
	 */
	MILLIS_PER_HOUR = MILLIS_PER_MINUTE * MINUTES_PER_HOUR
	/**
	 * Milliseconds per day.
	 */
	MILLIS_PER_DAY = MILLIS_PER_HOUR * HOURS_PER_DAY
)

var ZeroDuration = DurationOfMilliseconds(0)

/**
 *  A span of time
 *  ~~~~~~~~~~~~~~
 *  such as 27 days, 4 hours, 12 minutes, and 3 seconds.
 */
type Duration interface {
	Equal(other interface{}) bool
	String() string

	InMilliseconds() int64

	IsZero() bool
	IsPositive() bool
	IsNegative() bool

	Negated() Duration
	Absolute() Duration

	//
	//  +
	//
	Plus(d Duration) Duration
	PlusMilliseconds(millis int64) Duration
	PlusSeconds(secs int32) Duration
	PlusMinutes(minutes int32) Duration
	PlusHours(hours int32) Duration
	PlusDays(days int32) Duration

	//
	//  -
	//
	Minus(d Duration) Duration
	MinusMilliseconds(millis int64) Duration
	MinusSeconds(secs int32) Duration
	MinusMinutes(minutes int32) Duration
	MinusHours(hours int32) Duration
	MinusDays(days int32) Duration

	//
	//  *, /
	//
	Multiply(multiplicand int64) Duration

	Divide(divisor int64) Duration

	//
	//  DateTime
	//
	AddTo(time Time) Time
	SubtractFrom(time Time) Time
}

//
//  Factories
//

func DurationOfMilliseconds(millis int64) Duration {
	return &TimeDuration{
		millis: millis,
	}
}

func DurationOfSeconds(secs int32) Duration {
	millis := int64(MILLIS_PER_SECOND * secs)
	return &TimeDuration{
		millis: millis,
	}
}

func DurationOfMinutes(minutes int32) Duration {
	millis := int64(MILLIS_PER_MINUTE * minutes)
	return &TimeDuration{
		millis: millis,
	}
}

func DurationOfHours(hours int32) Duration {
	millis := int64(MILLIS_PER_HOUR * hours)
	return &TimeDuration{
		millis: millis,
	}
}

func DurationOfDays(days int32) Duration {
	millis := int64(MILLIS_PER_DAY * days)
	return &TimeDuration{
		millis: millis,
	}
}

/**
 *  The result of this method can be a negative period if the end is before the start.
 */
func DurationBetween(startTime, endTime Time) Duration {
	seconds := TimeToFloat64(endTime) - TimeToFloat64(startTime)
	millis := int64(seconds * MILLIS_PER_SECOND)
	return &TimeDuration{
		millis: millis,
	}
}
