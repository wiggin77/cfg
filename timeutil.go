package config

import "errors"

// MillisPerSecond is the number of millseconds per second.
const MillisPerSecond int64 = 1000

// MillisPerMinute is the number of millseconds per minute.
const MillisPerMinute int64 = MillisPerSecond * 60

// MillisPerHour is the number of millseconds per hour.
const MillisPerHour int64 = MillisPerMinute * 60

// MillisPerDay is the number of millseconds per day.
const MillisPerDay int64 = MillisPerHour * 24

// MillisPerWeek is the number of millseconds per week.
const MillisPerWeek int64 = MillisPerDay * 7

// MillisPerYear is the approximate number of millseconds per year.
const MillisPerYear int64 = MillisPerDay*365 + int64((float64(MillisPerDay) * 0.25))

// ParseMilliseconds parses a string containing a number plus
// a unit of measure for time and returns the number of milliseconds
// it represents.
//
// Example:
// * "1 second" returns 1000
// * "1 minute" returns 60000
// * "1 hour" returns 3600000
//
// See config.unitsToMillis for a list of supported units of measure.
func ParseMilliseconds(str string) (int64, error) {
	// TODO
	return 0, errors.New("not implemented")
}

// UnitsToMillis returns the number of milliseconds represented by the specified unit of measure.
//
// Example:
// * "second" returns 1000	<br/>
// * "minute" returns 60000	<br/>
// * "hour" returns 3600000	<br/>
//
// Supported units of measure:
// * "milliseconds", "millis", "ms", "millisecond"
// * "seconds", "sec", "s", "second"
// * "minutes", "mins", "min", "m", "minute"
// * "hours", "h", "hour"
// * "days", "d", "day"
// * "weeks", "w", "week"
func UnitsToMillis(unit string) (int64, error) {
	// TODO
	return 0, errors.New("not implemented")
}
