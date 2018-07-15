package config

import (
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

func TestConfig_String(t *testing.T) {
	src1 := NewSrcMapFromMap(map[string]string{"prop1": "1"})
	src2 := NewSrcMapFromMap(map[string]string{"prop2": "2", "prop1": "2"})
	src3 := NewSrcMapFromMap(map[string]string{"prop3": "3", "prop2": "3", "prop1": "3"})
	config := &Config{}
	config.AppendSource(src1, src2, src3)
	type args struct {
		prop string
		def  string
	}
	tests := []struct {
		name    string
		args    args
		wantVal string
		wantErr error
	}{
		{"Missing prop", args{"blap", "x"}, "x", ErrNotFound},
		{"Prop1", args{"prop1", "x"}, "1", nil},
		{"Prop2", args{"prop2", "x"}, "2", nil},
		{"Prop3", args{"prop3", "x"}, "3", nil},
		{"Blank prop", args{"", "x"}, "x", ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotErr := config.String(tt.args.prop, tt.args.def)
			if gotVal != tt.wantVal {
				t.Errorf("Config.String() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Config.String() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestConfig_Int(t *testing.T) {
	src1 := NewSrcMapFromMap(map[string]string{"prop1": "-1", "big": "2147483647"})
	src2 := NewSrcMapFromMap(map[string]string{"prop2": "2", "prop1": "2", "zero": "0"})
	src3 := NewSrcMapFromMap(map[string]string{"prop3": "3", "prop2": "3", "prop1": "3"})
	config := &Config{}
	config.AppendSource(src1, src2, src3)
	type args struct {
		prop string
		def  int
	}
	tests := []struct {
		name    string
		args    args
		wantVal int
		wantErr error
	}{
		{"Missing prop", args{"blap", 77}, 77, ErrNotFound},
		{"Prop1", args{"prop1", 77}, -1, nil},
		{"Prop2", args{"prop2", 77}, 2, nil},
		{"Prop3", args{"prop3", 77}, 3, nil},
		{"zero", args{"zero", 77}, 0, nil},
		{"big", args{"big", 77}, 2147483647, nil},
		{"Blank prop", args{"", 77}, 77, ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotErr := config.Int(tt.args.prop, tt.args.def)
			if gotVal != tt.wantVal {
				t.Errorf("Config.Int() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Config.Int() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
