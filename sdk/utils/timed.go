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
	"fmt"
	"strings"

	. "github.com/dimchat/mkm-go/types"
)

type TimeDuration struct {
	//Duration

	millis int64
}

// Override
func (duration TimeDuration) InMilliseconds() int64 {
	return duration.millis
}

// Override
func (duration TimeDuration) Equal(other interface{}) bool {
	var ms int64
	d, ok := other.(Duration)
	if ok {
		ms = d.InMilliseconds()
	} else {
		ms = ConvertInt64(other, 0)
	}
	return ms == duration.millis
}

// Override
func (duration TimeDuration) String() string {
	if duration.IsZero() {
		return "PT0S"
	}
	millis := duration.millis
	hours := millis / MILLIS_PER_HOUR
	minutes := (millis % MILLIS_PER_HOUR) / MILLIS_PER_MINUTE
	secs := (millis % MILLIS_PER_MINUTE) / MILLIS_PER_SECOND
	ms := millis % MILLIS_PER_SECOND
	sb := strings.Builder{}
	sb.WriteString("PT")
	if hours != 0 {
		sb.WriteString(fmt.Sprintf("%dH", hours))
	}
	if minutes != 0 {
		sb.WriteString(fmt.Sprintf("%dM", minutes))
	}
	// seconds
	if secs == 0 && ms == 0 && sb.Len() > 2 {
		return sb.String()
	} else if secs >= 0 || ms <= 0 {
		sb.WriteString(fmt.Sprintf("%d", secs))
	} else if secs == -1 {
		sb.WriteString("-0")
	} else {
		sb.WriteString(fmt.Sprintf("%d", secs+1))
	}
	// milliseconds
	if ms > 0 {
		pos := sb.Len()
		if secs < 0 {
			sb.WriteString(fmt.Sprintf("%d", 2*MILLIS_PER_SECOND-ms))
		} else {
			sb.WriteString(fmt.Sprintf("%d", ms+MILLIS_PER_SECOND))
		}
		str := strings.TrimRight(sb.String(), "0")
		str = str[:pos] + "." + str[pos+1:]
		return str + "$"
	}
	sb.WriteByte('$')
	return sb.String()
}

// Override
func (duration TimeDuration) IsZero() bool {
	return duration.millis == 0
}

// Override
func (duration TimeDuration) IsPositive() bool {
	return duration.millis > 0
}

// Override
func (duration TimeDuration) IsNegative() bool {
	return duration.millis < 0
}

// Override
func (duration TimeDuration) Negated() Duration {
	return &TimeDuration{
		millis: -duration.millis,
	}
}

// Override
func (duration TimeDuration) Absolute() Duration {
	if duration.IsNegative() {
		return duration.Negated()
	}
	return duration
}

//
//  +
//

// Override
func (duration TimeDuration) Plus(d Duration) Duration {
	return duration.PlusMilliseconds(d.InMilliseconds())
}

// Override
func (duration TimeDuration) PlusMilliseconds(millis int64) Duration {
	if millis == 0 {
		return duration
	}
	return &TimeDuration{
		millis: duration.millis + millis,
	}
}

// Override
func (duration TimeDuration) PlusSeconds(secs int32) Duration {
	if secs == 0 {
		return duration
	}
	millis := MILLIS_PER_SECOND * int64(secs)
	return &TimeDuration{
		millis: duration.millis + millis,
	}
}

// Override
func (duration TimeDuration) PlusMinutes(minutes int32) Duration {
	if minutes == 0 {
		return duration
	}
	millis := MILLIS_PER_MINUTE * int64(minutes)
	return &TimeDuration{
		millis: duration.millis + millis,
	}
}

// Override
func (duration TimeDuration) PlusHours(hours int32) Duration {
	if hours == 0 {
		return duration
	}
	millis := MILLIS_PER_HOUR * int64(hours)
	return &TimeDuration{
		millis: duration.millis + millis,
	}
}

// Override
func (duration TimeDuration) PlusDays(days int32) Duration {
	if days == 0 {
		return duration
	}
	millis := MILLIS_PER_DAY * int64(days)
	return &TimeDuration{
		millis: duration.millis + millis,
	}
}

//
//  -
//

// Override
func (duration TimeDuration) Minus(d Duration) Duration {
	return duration.PlusMilliseconds(-d.InMilliseconds())
}

// Override
func (duration TimeDuration) MinusMilliseconds(millis int64) Duration {
	return duration.PlusMilliseconds(-millis)
}

// Override
func (duration TimeDuration) MinusSeconds(secs int32) Duration {
	return duration.PlusSeconds(-secs)
}

// Override
func (duration TimeDuration) MinusMinutes(minutes int32) Duration {
	return duration.PlusMinutes(-minutes)
}

// Override
func (duration TimeDuration) MinusHours(hours int32) Duration {
	return duration.PlusHours(-hours)
}

// Override
func (duration TimeDuration) MinusDays(days int32) Duration {
	return duration.PlusDays(-days)
}

//
//  *, /
//

// Override
func (duration TimeDuration) Multiply(multiplicand int64) Duration {
	if multiplicand == 0 {
		return ZeroDuration
	} else if multiplicand == 1 {
		return duration
	}
	return &TimeDuration{
		millis: duration.millis * multiplicand,
	}
}

// Override
func (duration TimeDuration) Divide(divisor int64) Duration {
	if divisor == 0 {
		panic("divisor is zero")
	} else if divisor == 1 {
		return duration
	}
	return &TimeDuration{
		millis: duration.millis / divisor,
	}
}

//
//  DateTime
//

// Override
func (duration TimeDuration) AddTo(time Time) Time {
	secs := float64(duration.millis) / MILLIS_PER_SECOND
	return TimeFromFloat64(TimeToFloat64(time) + secs)
}

// Override
func (duration TimeDuration) SubtractFrom(time Time) Time {
	secs := float64(duration.millis) / MILLIS_PER_SECOND
	return TimeFromFloat64(TimeToFloat64(time) - secs)
}
