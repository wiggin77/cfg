package config

import (
	"math"
	"strconv"
	"strings"
	"testing"
	"time"

	timeutil "github.com/wiggin77/config/time"
)

func TestConfig_PrependSource(t *testing.T) {
	map1 := map[string]string{"prop1": "1"}
	src1 := NewSrcMapFromMap(map1)
	map2 := map[string]string{"prop2": "2", "prop1": "2"}
	src2 := NewSrcMapFromMap(map2)
	map3 := map[string]string{"prop3": "3", "prop2": "3", "prop1": "3"}
	src3 := NewSrcMapFromMap(map3)

	cfg := &Config{}
	defer cfg.Shutdown()

	// Prepend one to empty config.
	cfg.PrependSource(src1)
	if len(cfg.srcs) != 1 {
		t.Errorf("Prepend src to empty config; expected len=1, got len=%d", len(cfg.srcs))
	}
	if val, err := cfg.String("prop1", ""); err != nil || val != "1" {
		t.Errorf("Prepend src to empty config; expected err=nil, val=1; got err=%v, val=%s", err, val)
	}
	if _, err := cfg.String("blap", ""); err == nil {
		t.Errorf("Prepend src to empty config; expected err=false for missing prop, got err=%v", err)
	}

	// Prepend second
	cfg.PrependSource(src2)
	if len(cfg.srcs) != 2 {
		t.Errorf("Prepend second src; expected len=2, got len=%d", len(cfg.srcs))
	}
	if val, err := cfg.String("prop1", ""); err != nil || val != "2" {
		t.Errorf("Prepend second src; for prop1 expected err=nil, val=2; got err=%v, val=%s", err, val)
	}
	if val, err := cfg.String("prop2", ""); err != nil || val != "2" {
		t.Errorf("Prepend second src; for prop2 expected err=nil, val=2; got err=%v, val=%s", err, val)
	}

	// Prepend third
	cfg.PrependSource(src3)
	if len(cfg.srcs) != 3 {
		t.Errorf("Prepend third src; expected len=3, got len=%d", len(cfg.srcs))
	}
	if val, err := cfg.String("prop1", ""); err != nil || val != "3" {
		t.Errorf("Prepend third src; for prop1 expected err=nil, val=3; got err=%v, val=%s", err, val)
	}
	if val, err := cfg.String("prop2", ""); err != nil || val != "3" {
		t.Errorf("Prepend third src; for prop2 expected err=nil, val=3; got err=%v, val=%s", err, val)
	}
	if val, err := cfg.String("prop3", ""); err != nil || val != "3" {
		t.Errorf("Prepend third src; for prop3 expected err=nil, val=3; got err=%v, val=%s", err, val)
	}
	if _, err := cfg.String("blap", ""); err == nil {
		t.Errorf("Prepend third src; expected err=false for missing prop, got err=%v", err)
	}
}

func TestConfig_AppendSource(t *testing.T) {
	map1 := map[string]string{"prop1": "1"}
	src1 := NewSrcMapFromMap(map1)
	map2 := map[string]string{"prop2": "2", "prop1": "2"}
	src2 := NewSrcMapFromMap(map2)
	map3 := map[string]string{"prop3": "3", "prop2": "3", "prop1": "3"}
	src3 := NewSrcMapFromMap(map3)

	cfg := &Config{}
	defer cfg.Shutdown()

	// Append to empty config.
	cfg.AppendSource(src1)
	if len(cfg.srcs) != 1 {
		t.Errorf("Append src to empty config; expected len=1, got len=%d", len(cfg.srcs))
	}
	if val, err := cfg.String("prop1", ""); err != nil || val != "1" {
		t.Errorf("Append src to empty config; expected err=nil, val=1; got err=%v, val=%s", err, val)
	}
	if _, err := cfg.String("blap", ""); err == nil {
		t.Errorf("Append src to empty config; expected err=false for missing prop, got err=%v", err)
	}

	// Append second
	cfg.AppendSource(src2)
	if len(cfg.srcs) != 2 {
		t.Errorf("Append second src; expected len=2, got len=%d", len(cfg.srcs))
	}
	if val, err := cfg.String("prop1", ""); err != nil || val != "1" {
		t.Errorf("Append second src; for prop1 expected err=nil, val=1; got err=%v, val=%s", err, val)
	}
	if val, err := cfg.String("prop2", ""); err != nil || val != "2" {
		t.Errorf("Append second src; for prop2 expected err=nil, val=2; got err=%v, val=%s", err, val)
	}

	// Append third
	cfg.AppendSource(src3)
	if len(cfg.srcs) != 3 {
		t.Errorf("Append third src; expected len=3, got len=%d", len(cfg.srcs))
	}
	if val, err := cfg.String("prop1", ""); err != nil || val != "1" {
		t.Errorf("Append third src; for prop1 expected err=nil, val=1; got err=%v, val=%s", err, val)
	}
	if val, err := cfg.String("prop2", ""); err != nil || val != "2" {
		t.Errorf("Append third src; for prop2 expected err=nil, val=2; got err=%v, val=%s", err, val)
	}
	if val, err := cfg.String("prop3", ""); err != nil || val != "3" {
		t.Errorf("Append third src; for prop3 expected err=nil, val=3; got err=%v, val=%s", err, val)
	}
	if _, err := cfg.String("blap", ""); err == nil {
		t.Errorf("Append third src; expected err=false for missing prop, got err=%v", err)
	}
}

func TestConfig_AddRemoveChangedListener(t *testing.T) {
	map1 := map[string]string{"prop1": "1"}
	src1 := NewSrcMapFromMap(map1)
	cfg := &Config{}
	defer cfg.Shutdown()
	cfg.AppendSource(src1)

	tl1 := &TestListener{}
	tl2 := &TestListener{}
	tl3 := &TestListener{}

	// Test Add
	cfg.AddChangedListener(tl1)
	if len(cfg.chgListeners) != 1 {
		t.Errorf("AddChangeListener; expected len=1, got len=%d", len(cfg.chgListeners))
	}
	cfg.AddChangedListener(tl2)
	if len(cfg.chgListeners) != 2 {
		t.Errorf("AddChangeListener; expected len=2, got len=%d", len(cfg.chgListeners))
	}
	cfg.AddChangedListener(tl3)
	cfg.AddChangedListener(tl1)
	if len(cfg.chgListeners) != 4 {
		t.Errorf("AddChangeListener; expected len=4, got len=%d", len(cfg.chgListeners))
	}

	// Test remove
	cfg.RemoveChangedListener(tl1)
	if len(cfg.chgListeners) != 2 {
		t.Errorf("RemoveChangeListener; expected len=2, got len=%d", len(cfg.chgListeners))
	}
	cfg.RemoveChangedListener(tl2)
	if len(cfg.chgListeners) != 1 {
		t.Errorf("RemoveChangeListener; expected len=1, got len=%d", len(cfg.chgListeners))
	}
	cfg.RemoveChangedListener(tl3)
	if len(cfg.chgListeners) != 0 {
		t.Errorf("RemoveChangeListener; expected len=0, got len=%d", len(cfg.chgListeners))
	}
}

// ChangeListener and ChangedPropListener
type TestListener struct {
	_ [1]byte // cannot be empty struct
}

func (l *TestListener) ConfigChanged(cfg *Config, src SourceMonitored) {
}

const noerror = ""

// tests contains parameters for one test case.
type test struct {
	srcLevel        int // 0 means don't add prop to Source
	propName        string
	propVal         string
	defVal          interface{}
	expectedVal     interface{}
	expectedErrText string
}

// numberOfSources determines the highest source number in the array.
func numberOfSources(tests []test) (num int) {
	for _, p := range tests {
		if p.srcLevel > num {
			num = p.srcLevel
		}
	}
	return
}

// makeTestConfig creates a `Config` populated with test data.
func makeTestConfig(tests []test) *Config {
	num := numberOfSources(tests)
	srcs := make([]Source, num)
	for i := range srcs {
		srcs[i] = NewSrcMap()
	}

	for _, p := range tests {
		level := p.srcLevel
		if level > 0 {
			src := srcs[level-1].(*SrcMap)
			src.Put(p.propName, p.propVal)
		}
	}

	config := &Config{}
	config.AppendSource(srcs...)
	return config
}

func runTest(tests []test, t *testing.T) {
	config := makeTestConfig(tests)
	defer config.Shutdown()

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			var gotVal interface{}
			var gotErr error
			var testName string

			switch def := tt.defVal.(type) {
			case string:
				testName = "Config.String()"
				gotVal, gotErr = config.String(tt.propName, def)
			case int:
				testName = "Config.Int()"
				gotVal, gotErr = config.Int(tt.propName, def)
			case int64:
				testName = "Config.Int64()"
				gotVal, gotErr = config.Int64(tt.propName, def)
			case float64:
				testName = "Config.Float64()"
				gotVal, gotErr = config.Float64(tt.propName, def)
			case bool:
				testName = "Config.Bool()"
				gotVal, gotErr = config.Bool(tt.propName, def)
			case time.Duration:
				testName = "Config.Duration()"
				gotVal, gotErr = config.Duration(tt.propName, time.Duration(def))
			default:
				t.Errorf("%s no test defined for type %T", testName, def)
				return
			}

			if gotErr == nil && tt.expectedErrText != noerror {
				t.Errorf("%s for prop = %s; gotErr = nil, expected %v", testName, tt.propName, tt.expectedErrText)
				return
			}

			if gotErr != nil && tt.expectedErrText == noerror {
				t.Errorf("%s for prop = %s; gotErr = %v, expected no error", testName, tt.propName, gotErr)
				return
			}

			if gotErr != nil && !strings.Contains(gotErr.Error(), tt.expectedErrText) {
				t.Errorf("%s for prop = %s; gotErr = %v, expected %v", testName, tt.propName, gotErr, tt.expectedErrText)
				return
			}

			if gotVal != tt.expectedVal {
				t.Errorf("%s for prop = %s; [gotType=%T, expType=%T]", testName, tt.propName, gotVal, tt.expectedVal)
				t.Errorf("%s for prop = %s; gotVal = %v, expected %v", testName, tt.propName, gotVal, tt.expectedVal)
			}
		})
	}
}

func TestConfig_String(t *testing.T) {
	tests := []test{
		// srcLevel, propName, propVal, defVal, expectedVal, expectedErrText
		{1, "prop1", "1", "x", "1", noerror},
		{2, "prop2", "2", "x", "2", noerror},
		{3, "prop3", "3", "x", "3", noerror},
		{0, "blap", "x", "", "", "not found"},
		{0, "", "", "x", "x", "not found"},
	}
	runTest(tests, t)
}

func TestConfig_Int(t *testing.T) {
	tests := []test{
		// srcLevel, propName, propVal, defVal, expectedVal, expectedErrText
		{0, "missing", "1", -1, -1, "not found"},
		{1, "prop1", "1", -1, 1, noerror},
		{1, "zero", "0", -1, 0, noerror},
		{1, "neg", "-5", -1, -5, noerror},
		{1, "neg_zero", "-0", -1, 0, noerror},
		{1, "pos_zero", "+0", -1, 0, noerror},
		{1, "big", "2147483647", -1, 2147483647, noerror},
		{3, "min", strconv.Itoa(math.MinInt32), -1, math.MinInt32, noerror},
		{3, "max", strconv.Itoa(math.MaxInt32), -1, math.MaxInt32, noerror},
		{3, "overflow", strconv.Itoa(math.MaxInt32 * 2), -1, -1, "out of range"},
		{3, "bad", "00x55", -1, -1, "invalid syntax"},
		{3, "bad2", "1.025", -1, -1, "invalid syntax"},
		{3, "bad3", "0x11", -1, -1, "invalid syntax"},
	}
	runTest(tests, t)
}

func TestConfig_Int64(t *testing.T) {
	tests := []test{
		// srcLevel, propName, propVal, defVal, expectedVal, expectedErrText
		{0, "missing", "1", int64(-1), int64(-1), "not found"},
		{1, "prop1", "1", int64(-1), int64(1), noerror},
		{1, "zero", "0", int64(-1), int64(0), noerror},
		{1, "neg", "-5", int64(-1), int64(-5), noerror},
		{1, "neg_zero", "-0", int64(-1), int64(0), noerror},
		{1, "pos_zero", "+0", int64(-1), int64(0), noerror},
		{1, "big", "2147483647000", int64(-1), int64(2147483647000), noerror},
		{3, "min", strconv.Itoa(math.MinInt64), int64(-1), int64(math.MinInt64), noerror},
		{3, "max", strconv.Itoa(math.MaxInt64), int64(-1), int64(math.MaxInt64), noerror},
		{3, "overflow", strconv.Itoa(math.MaxInt64) + "0", int64(-1), int64(-1), "out of range"},
		{3, "bad", "00x55", int64(-1), int64(-1), "invalid syntax"},
		{3, "bad2", "1.025", int64(-1), int64(-1), "invalid syntax"},
		{3, "bad3", "0x11", int64(-1), int64(-1), "invalid syntax"},
	}
	runTest(tests, t)
}

func TestConfig_Float64(t *testing.T) {
	tests := []test{
		// srcLevel, propName, propVal, defVal, expectedVal, expectedErrText
		{0, "missing", "1", float64(-1), float64(-1), "not found"},
		{1, "prop1", "1", float64(-1), float64(1), noerror},
		{1, "prop2", "1.025", float64(-1), float64(1.025), noerror},
		{1, "zero", "0", float64(-1), float64(0), noerror},
		{1, "neg", "-5.25", float64(-1), float64(-5.25), noerror},
		{1, "neg_zero", "-0.0", float64(-1), float64(0), noerror},
		{1, "pos_zero", "+0.0", float64(-1), float64(0), noerror},
		{1, "big", "2147483647000", float64(-1), float64(2147483647000), noerror},
		{3, "min", strconv.FormatFloat(math.SmallestNonzeroFloat64, 'f', -1, 64), float64(-1), float64(math.SmallestNonzeroFloat64), noerror},
		{3, "max", strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64), float64(-1), float64(math.MaxFloat64), noerror},
		{3, "overflow", strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64) + "0", float64(-1), float64(-1), "out of range"},
		{3, "bad", "00x55", float64(-1), float64(-1), "invalid syntax"},
		{3, "bad2", "0x11", float64(-1), float64(-1), "invalid syntax"},
		{3, "bad3", "0..314", float64(-1), float64(-1), "invalid syntax"},
		{3, "bad3", "bad", float64(-1), float64(-1), "invalid syntax"},
	}
	runTest(tests, t)
}

func TestConfig_Bool(t *testing.T) {
	tests := []test{
		// srcLevel, propName, propVal, defVal, expectedVal, expectedErrText
		{0, "missing", "1", false, false, "not found"},
		{1, "prop1", " 1 ", false, true, noerror},
		{1, "prop2", " T ", false, true, noerror},
		{1, "prop3", " TRUE ", false, true, noerror},
		{1, "prop4", " Y ", false, true, noerror},
		{1, "prop5", " YES ", false, true, noerror},
		{1, "prop10", "\t0", true, false, noerror},
		{1, "prop11", "F", true, false, noerror},
		{1, "prop12", "false", true, false, noerror},
		{1, "prop13", "N", true, false, noerror},
		{1, "prop14", "NO", true, false, noerror},
		{3, "bad", "00x55", false, false, "invalid syntax"},
		{3, "bad2", "-F", false, false, "invalid syntax"},
		{3, "bad3", "", true, true, "invalid syntax"},
		{3, "bad3", "TRUE DAT", false, false, "invalid syntax"},
	}
	runTest(tests, t)
}

func TestConfig_Duration(t *testing.T) {
	// convert milliseconds to Duration
	ms2dur := func(ms int64) time.Duration {
		return time.Duration(ms) * time.Millisecond
	}
	tests := []test{
		// srcLevel, propName, propVal, defVal, expectedVal, expectedErrText
		{0, "missing", "1", time.Duration(-1), time.Duration(-1), "not found"},

		// All supported units of measure tested in "github.com/wiggin77/config/time"
		{1, "none", "1", ms2dur(-1), ms2dur(1), noerror},
		{1, "ms", "1ms", ms2dur(-1), ms2dur(1), noerror},
		{1, "sec", "1sec", ms2dur(-1), ms2dur(timeutil.MillisPerSecond), noerror},
		{1, "min", "1min", ms2dur(-1), ms2dur(timeutil.MillisPerMinute), noerror},
		{1, "hour", "1hour", ms2dur(-1), ms2dur(timeutil.MillisPerHour), noerror},
		{1, "day", "1day", ms2dur(-1), ms2dur(timeutil.MillisPerDay), noerror},
		{1, "week", "1week", ms2dur(-1), ms2dur(timeutil.MillisPerWeek), noerror},
		{1, "year", "1year", ms2dur(-1), ms2dur(timeutil.MillisPerYear), noerror},

		{3, "fraction1", "1.025", ms2dur(-1), ms2dur(1), noerror},
		{3, "fraction2", "1.5 minutes", ms2dur(-1), ms2dur(90000), noerror},

		{1, "zero", "0", ms2dur(-1), ms2dur(0), noerror},
		{1, "neg", "-5", ms2dur(-1), ms2dur(-5), noerror},
		{1, "neg_zero", "-0", ms2dur(-1), ms2dur(0), noerror},
		{1, "pos_zero", "+0", ms2dur(-1), ms2dur(0), noerror},
		{1, "big", "400 years", ms2dur(-1), ms2dur(timeutil.MillisPerYear * 400), noerror},
		{3, "overflow", "400000000 years", ms2dur(-1), ms2dur(-1), "out of range"},
		{3, "bad1", "00x55", ms2dur(-1), ms2dur(-1), "invalid syntax"},
		{3, "bad2", "1..025 days", ms2dur(-1), ms2dur(-1), "invalid syntax"},
		{3, "bad3", "0x11 week", ms2dur(-1), ms2dur(-1), "invalid syntax"},
		{3, "bad4", "garbage", ms2dur(-1), ms2dur(-1), "invalid syntax"},
	}
	runTest(tests, t)
}
