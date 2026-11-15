package ini_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/wiggin77/cfg/ini"
)

type entry struct {
	section string
	key     string
	val     string
}

var sample1 = "[sec1]\nkey1=val1\nkey2 = val2\n[sec2]\n\n  key1=val1  \nkey2 = val2"
var sample2 = "[ sec1 ]\nkey1=val1\nkey2 = val2\n[  sec2]\n\n  key1=val1  \nkey2 = val2"

func TestLoadFromReader(t *testing.T) {
	ini := ini.Ini{}
	r := strings.NewReader(sample1)
	err := ini.LoadFromReader(r)
	if err != nil {
		t.Error(err)
	}

	v, ok := ini.GetProp("sec2", "key2")
	if !ok {
		t.Errorf("sec2.key2 missing")
	}
	if v != "val2" {
		t.Errorf("sec2.key2 should equal val2")
	}
}

func TestLoadFromReaderError(t *testing.T) {
	ini := ini.Ini{}
	r := strings.NewReader("123blap")
	err := ini.LoadFromReader(r)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestGetProp(t *testing.T) {
	ini := ini.Ini{}

	// test empty section (Linux style config)
	s := "key1=val1 \n\n key2=val2 \n"
	err := ini.LoadFromString(s)
	if err != nil {
		t.Error(err)
	}
	v, ok := ini.GetProp("", "key1")
	if !ok {
		t.Errorf("key1 missing")
	}
	if v != "val1" {
		t.Errorf("key1 should equal val1")
	}

	// test empty section, plus non-empty sections
	s = "key1=val1 \n\n key2=val2 \n[sec1]\n\nkey1=sec1val1\nkey2=sec1val2"
	err = ini.LoadFromString(s)
	if err != nil {
		t.Error(err)
	}
	data := []entry{
		{section: "", key: "key1", val: "val1"},
		{section: "", key: "key2", val: "val2"},
		{section: "sec1", key: "key1", val: "sec1val1"},
		{section: "sec1", key: "key2", val: "sec1val2"},
	}
	for _, d := range data {
		v, ok := ini.GetProp(d.section, d.key)
		if !ok {
			t.Errorf("key %s missing from section %s", d.key, d.section)
		}
		if v != d.val {
			t.Errorf("%s.%s should equal %s (actual %s)", d.section, d.key, d.val, v)
		}
	}

	// test missing section
	v, ok = ini.GetProp("blap", "key1")
	if ok {
		t.Errorf("%s.%s should be missing", "blap", "key1")
	}
	// test missing prop
	v, ok = ini.GetProp("sec1", "blap")
	if ok {
		t.Errorf("%s.%s should be missing", "sec1", "blap")
	}
}

func TestToMap(t *testing.T) {
	s := "key1=val1 \n\n key2=val2 \n[sec1]\n\nkey1=sec1val1\nkey2=sec1val2"
	m := map[string]string{"key1": "val1", "key2": "val2", "sec1.key1": "sec1val1", "sec1.key2": "sec1val2"}
	ini := ini.Ini{}
	err := ini.LoadFromString(s)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(m, ini.ToMap()) {
		t.Errorf("maps not equal")
	}
}

func TestGetFlattenedKeys(t *testing.T) {
	s := "key1=val1 \n\n key2=val2 \n[sec1]\n\nkey1=sec1val1\nkey2=sec1val2"
	arr := []string{"key1", "key2", "sec1.key1", "sec1.key2"}
	ini := ini.Ini{}
	err := ini.LoadFromString(s)
	if err != nil {
		t.Error(err)
	}
	keys := ini.GetFlattenedKeys() // order is undefined

	sort.Strings(arr)
	sort.Strings(keys)

	if !reflect.DeepEqual(arr, keys) {
		t.Errorf("arrays not equal -- expected:%v, got %v", arr, keys)
	}
}

func TestGetSectionNames(t *testing.T) {
	s := "key1=val1 \n\n key2=val2 \n[sec1]\n\nkey1=sec1val1\nkey2=sec1val2 \n \n [sec2]\n [sec3]\n key1=sec3val1"
	arr := []string{"", "sec1", "sec2", "sec3"}
	ini := ini.Ini{}
	err := ini.LoadFromString(s)
	if err != nil {
		t.Error(err)
	}
	names := ini.GetSectionNames()

	sort.Strings(arr)
	sort.Strings(names)

	if !reflect.DeepEqual(arr, names) {
		t.Errorf("arrays not equal -- expected:%v, got %v", arr, names)
	}
}

func TestGetKeys(t *testing.T) {
	s := "key1=val1 \n\n key2=val2 \n[sec1]\n\nkey1=sec1val1\nkey2=sec1val2 \n \n [sec2]\n [sec3]\n key1=sec3val1"
	arr := []string{"key1", "key2"}
	ini := ini.Ini{}
	err := ini.LoadFromString(s)
	if err != nil {
		t.Error(err)
	}
	keys, err := ini.GetKeys("sec1")
	if err != nil {
		t.Error(err)
	}

	sort.Strings(arr)
	sort.Strings(keys)

	if !reflect.DeepEqual(arr, keys) {
		t.Errorf("arrays not equal -- expected:%v, got %v", arr, keys)
	}

	// test missing section
	keys, err = ini.GetKeys("blap")
	if err == nil {
		t.Errorf("expected error")
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestLoadFromReaderReadError(t *testing.T) {
	ini := ini.Ini{}
	if err := ini.LoadFromReader(errReader{}); err == nil {
		t.Errorf("expected error")
	}
}

func TestLoadFromFilespec(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.ini")
	if err := os.WriteFile(file, []byte(sample1), 0600); err != nil {
		t.Fatal(err)
	}
	mtime := time.Now().Add(-time.Hour).Truncate(time.Second)
	if err := os.Chtimes(file, mtime, mtime); err != nil {
		t.Fatal(err)
	}

	ini := ini.Ini{}
	if err := ini.LoadFromFilespec(file); err != nil {
		t.Fatal(err)
	}
	if v, ok := ini.GetProp("sec1", "key1"); !ok || v != "val1" {
		t.Errorf("sec1.key1 = %q, %v; want val1", v, ok)
	}
	if lm := ini.GetLastModified(); !lm.Equal(mtime) {
		t.Errorf("GetLastModified = %v, want file mtime %v", lm, mtime)
	}
}

func TestLoadFromFilespecErrors(t *testing.T) {
	dir := t.TempDir()
	ini := ini.Ini{}
	if err := ini.LoadFromFilespec(filepath.Join(dir, "missing.ini")); err == nil {
		t.Error("expected error for missing file")
	}

	bad := filepath.Join(dir, "bad.ini")
	if err := os.WriteFile(bad, []byte("123blap"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := ini.LoadFromFilespec(bad); err == nil {
		t.Error("expected error for invalid content")
	}
}

func TestLoadFromFileClosed(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.ini")
	if err := os.WriteFile(file, []byte(sample1), 0600); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	ini := ini.Ini{}
	if err := ini.LoadFromFile(f); err == nil {
		t.Error("expected error for closed file")
	}
}
