package timeconv

import "testing"

func TestParseMilliseconds(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"no_units", args{"1000"}, 1000, false},
		{"no_units_neg", args{"-360000"}, -360000, false},

		{"milliseconds", args{"07 milliseconds"}, 7, false},
		{"millisecond", args{"17 millisecond"}, 17, false},
		{"millis", args{"27 millis"}, 27, false},
		{"ms", args{"37 ms"}, 37, false},

		{"seconds", args{"07seconds"}, 7 * MillisPerSecond, false},
		{"second", args{"17second"}, 17 * MillisPerSecond, false},
		{"sec", args{"27sec"}, 27 * MillisPerSecond, false},
		{"s", args{"37s"}, 37 * MillisPerSecond, false},

		{"minutes", args{"07  Minutes"}, 7 * MillisPerMinute, false},
		{"minute", args{"17  Minute"}, 17 * MillisPerMinute, false},
		{"min", args{"27  Min"}, 27 * MillisPerMinute, false},
		{"m", args{"37 M"}, 37 * MillisPerMinute, false},

		{"hours", args{"-07 hours"}, -7 * MillisPerHour, false},
		{"hour", args{"-17 hour"}, -17 * MillisPerHour, false},
		{"h", args{"-27 h"}, -27 * MillisPerHour, false},

		{"days", args{"007 days"}, 7 * MillisPerDay, false},
		{"day", args{"17 day"}, 17 * MillisPerDay, false},
		{"d", args{"27 d"}, 27 * MillisPerDay, false},

		{"weeks", args{"+0 Weeks"}, 0 * MillisPerWeek, false},
		{"week", args{"+17 Week"}, 17 * MillisPerWeek, false},
		{"w", args{"+27 W"}, 27 * MillisPerWeek, false},

		{"years", args{"   7    years   "}, 7 * MillisPerYear, false},
		{"year", args{"   17    year   "}, 17 * MillisPerYear, false},
		{"y", args{"   27    y   "}, 27 * MillisPerYear, false},

		{"fractions1", args{"-0.5 minutes"}, -30000, false},
		{"fractions2", args{"+1.50 minutes"}, 90000, false},

		{"bad1", args{"x years"}, 0, true},
		{"bad2", args{"17 leap years"}, 0, true},
		{"bad3", args{"27px"}, 0, true},
		{"bad4", args{""}, 0, true},
		{"bad5", args{"27..1 years"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMilliseconds(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMilliseconds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseMilliseconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
