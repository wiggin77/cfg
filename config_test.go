package config

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestConfig_PrependSource(t *testing.T) {
	map1 := map[string]string{"prop1": "1"}
	src1 := NewSrcMapFromMap(map1)
	map2 := map[string]string{"prop2": "2", "prop1": "2"}
	src2 := NewSrcMapFromMap(map2)
	map3 := map[string]string{"prop3": "3", "prop2": "3", "prop1": "3"}
	src3 := NewSrcMapFromMap(map3)

	cfg := &Config{}

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

func TestConfig_AddRemoveChangedPropListener(t *testing.T) {
	map1 := map[string]string{"prop1": "1"}
	src1 := NewSrcMapFromMap(map1)
	cfg := &Config{}
	cfg.AppendSource(src1)

	tl1 := &TestListener{}
	tl2 := &TestListener{}
	tl3 := &TestListener{}

	// Test Add
	cfg.AddChangedPropListener(tl1)
	if len(cfg.propListeners) != 1 {
		t.Errorf("AddChangedPropListener; expected len=1, got len=%d", len(cfg.propListeners))
	}
	cfg.AddChangedPropListener(tl2)
	if len(cfg.propListeners) != 2 {
		t.Errorf("AddChangedPropListener; expected len=2, got len=%d", len(cfg.propListeners))
	}
	cfg.AddChangedPropListener(tl3)
	cfg.AddChangedPropListener(tl1)
	if len(cfg.propListeners) != 4 {
		t.Errorf("AddChangedPropListener; expected len=4, got len=%d", len(cfg.propListeners))
	}

	// Test remove
	cfg.RemoveChangedPropListener(tl1)
	if len(cfg.propListeners) != 2 {
		t.Errorf("RemoveChangedPropListener; expected len=2, got len=%d", len(cfg.propListeners))
	}
	cfg.RemoveChangedPropListener(tl2)
	if len(cfg.propListeners) != 1 {
		t.Errorf("RemoveChangedPropListener; expected len=1, got len=%d", len(cfg.propListeners))
	}
	cfg.RemoveChangedPropListener(tl3)
	if len(cfg.propListeners) != 0 {
		t.Errorf("RemoveChangedPropListener; expected len=0, got len=%d", len(cfg.propListeners))
	}
}

// ChangeListener and ChangedPropListener
type TestListener struct {
	_ [1]byte // cannot be empty struct
}

func (l *TestListener) Changed(cfg *Config, src *Source) {
}

func (l *TestListener) ChangedProp(cfg *Config, src *Source, name string) {
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
	src1 := NewSrcMapFromMap(map[string]string{"prop1": "-1", "big": "21474836470000", "pi": "3.14", "neg": "-3.14"})
	src2 := NewSrcMapFromMap(map[string]string{"prop2": "2", "prop1": "2", "zero": "0", "fl_zero": "0.0"})
	src3 := NewSrcMapFromMap(map[string]string{"prop3": "3", "prop2": "3", "prop1": "3", "neg_zero": "-0.0"})
	config := &Config{}
	config.AppendSource(src1, src2, src3)
	type args struct {
		prop string
		def  float64
	}
	tests := []struct {
		name    string
		args    args
		wantVal float64
		wantErr error
	}{
		{"Missing prop", args{"blap", 77.7}, 77.7, ErrNotFound},
		{"Prop1", args{"prop1", 77.7}, -1, nil},
		{"Prop2", args{"prop2", 77.7}, 2, nil},
		{"Pi", args{"pi", 77.7}, 3.14, nil},
		{"Neg", args{"neg", 77.7}, -3.14, nil},
		{"Float Zero", args{"fl_zero", 77.7}, 0, nil},
		{"Neg Zero", args{"neg_zero", 77.7}, 0, nil},
		{"zero", args{"zero", 77.7}, 0, nil},
		{"big", args{"big", 77.7}, 2147483647 * 10000, nil},
		{"Blank prop", args{"", 77.7}, 77.7, ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotErr := config.Float64(tt.args.prop, tt.args.def)
			if gotVal != tt.wantVal {
				t.Errorf("Config.Float64() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Config.Float64() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
