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
	if val, ok := cfg.String("prop1"); ok != true || val != "1" {
		t.Errorf("Prepend src to empty config; expected ok=true, val=1; got ok=%v, val=%s", ok, val)
	}
	if _, ok := cfg.String("blap"); ok == true {
		t.Errorf("Prepend src to empty config; expected ok=false for missing prop, got ok=%v", ok)
	}

	// Prepend second
	cfg.PrependSource(src2)
	if len(cfg.srcs) != 2 {
		t.Errorf("Prepend second src; expected len=2, got len=%d", len(cfg.srcs))
	}
	if val, ok := cfg.String("prop1"); ok != true || val != "2" {
		t.Errorf("Prepend second src; for prop1 expected ok=true, val=2; got ok=%v, val=%s", ok, val)
	}
	if val, ok := cfg.String("prop2"); ok != true || val != "2" {
		t.Errorf("Prepend second src; for prop2 expected ok=true, val=2; got ok=%v, val=%s", ok, val)
	}

	// Prepend third
	cfg.PrependSource(src3)
	if len(cfg.srcs) != 3 {
		t.Errorf("Prepend third src; expected len=3, got len=%d", len(cfg.srcs))
	}
	if val, ok := cfg.String("prop1"); ok != true || val != "3" {
		t.Errorf("Prepend third src; for prop1 expected ok=true, val=3; got ok=%v, val=%s", ok, val)
	}
	if val, ok := cfg.String("prop2"); ok != true || val != "3" {
		t.Errorf("Prepend third src; for prop2 expected ok=true, val=3; got ok=%v, val=%s", ok, val)
	}
	if val, ok := cfg.String("prop3"); ok != true || val != "3" {
		t.Errorf("Prepend third src; for prop3 expected ok=true, val=3; got ok=%v, val=%s", ok, val)
	}
	if _, ok := cfg.String("blap"); ok == true {
		t.Errorf("Prepend third src; expected ok=false for missing prop, got ok=%v", ok)
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
	if val, ok := cfg.String("prop1"); ok != true || val != "1" {
		t.Errorf("Append src to empty config; expected ok=true, val=1; got ok=%v, val=%s", ok, val)
	}
	if _, ok := cfg.String("blap"); ok == true {
		t.Errorf("Append src to empty config; expected ok=false for missing prop, got ok=%v", ok)
	}

	// Append second
	cfg.AppendSource(src2)
	if len(cfg.srcs) != 2 {
		t.Errorf("Append second src; expected len=2, got len=%d", len(cfg.srcs))
	}
	if val, ok := cfg.String("prop1"); ok != true || val != "1" {
		t.Errorf("Append second src; for prop1 expected ok=true, val=1; got ok=%v, val=%s", ok, val)
	}
	if val, ok := cfg.String("prop2"); ok != true || val != "2" {
		t.Errorf("Append second src; for prop2 expected ok=true, val=2; got ok=%v, val=%s", ok, val)
	}

	// Append third
	cfg.AppendSource(src3)
	if len(cfg.srcs) != 3 {
		t.Errorf("Append third src; expected len=3, got len=%d", len(cfg.srcs))
	}
	if val, ok := cfg.String("prop1"); ok != true || val != "1" {
		t.Errorf("Append third src; for prop1 expected ok=true, val=1; got ok=%v, val=%s", ok, val)
	}
	if val, ok := cfg.String("prop2"); ok != true || val != "2" {
		t.Errorf("Append third src; for prop2 expected ok=true, val=2; got ok=%v, val=%s", ok, val)
	}
	if val, ok := cfg.String("prop3"); ok != true || val != "3" {
		t.Errorf("Append third src; for prop3 expected ok=true, val=3; got ok=%v, val=%s", ok, val)
	}
	if _, ok := cfg.String("blap"); ok == true {
		t.Errorf("Append third src; expected ok=false for missing prop, got ok=%v", ok)
	}
}

func TestConfig_AddChangedListener(t *testing.T) {
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

type TestListener struct {
	_ [1]byte // cannot be empty struct
}

func (l *TestListener) Changed(cfg *Config, src *Source) {
}

func (l *TestListener) ChangedProp(cfg *Config, src *Source, name string) {
}
