/* For license and copyright information please see the LEGAL file in the code repository */

package timer

import (
	"memar/time/duration"
	"memar/time/monotonic"
)

// when is a helper function for setting the 'when' field of a Timer.
// It returns what the time will be, in nanoseconds, Duration d in the future.
// If d is negative, it is ignored. If the returned value would be less than
// zero because of an overflow, MaxInt64 is returned.
func when(d duration.NanoSecond) (t monotonic.Time) {
	t.Now()
	if d <= 0 {
		return
	}
	t.Add(d)
	// check for overflow.
	if t < 0 {
		// monotonic.Now() and d are always positive, so addition
		// (including overflow) will never result in t == 0.
		t = maxWhen
	}
	return
}
