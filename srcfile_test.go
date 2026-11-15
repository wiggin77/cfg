package cfg

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeConf(t *testing.T, file string, content string, mtime time.Time) {
	t.Helper()
	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(file, mtime, mtime); err != nil {
		t.Fatal(err)
	}
}

func assertProp(t *testing.T, sf *SrcFile, key, want string) {
	t.Helper()
	props, err := sf.GetProps()
	if err != nil {
		t.Fatal(err)
	}
	if got := props[key]; got != want {
		t.Errorf("%s = %q, want %q (props %v)", key, got, want, props)
	}
}

func TestSrcFileFromFilespecReloadsAfterEdit(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.conf")
	start := time.Now().Add(-time.Hour)
	writeConf(t, file, "name = one\n", start)

	sf, err := NewSrcFileFromFilespec(file)
	if err != nil {
		t.Fatal(err)
	}
	assertProp(t, sf, "name", "one")

	writeConf(t, file, "name = two\n", start.Add(time.Minute))
	assertProp(t, sf, "name", "two")

	writeConf(t, file, "name = three\n", start.Add(2*time.Minute))
	assertProp(t, sf, "name", "three")
}

func TestSrcFileFromFilespecReloadsAfterReplace(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.conf")
	start := time.Now().Add(-time.Hour)
	writeConf(t, file, "name = one\n", start)

	sf, err := NewSrcFileFromFilespec(file)
	if err != nil {
		t.Fatal(err)
	}

	tmp := filepath.Join(dir, "test.conf.new")
	writeConf(t, tmp, "name = two\n", start.Add(time.Minute))
	if err := os.Rename(tmp, file); err != nil {
		t.Fatal(err)
	}
	assertProp(t, sf, "name", "two")
}

func TestSrcFileFromFilespecDoesNotHoldFileOpen(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.conf")
	writeConf(t, file, "name = one\n", time.Now())

	if _, err := NewSrcFileFromFilespec(file); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(file); err != nil {
		t.Fatalf("file still in use after load: %v", err)
	}
}

func TestSrcFileFromHandleReloadsAfterEdit(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.conf")
	start := time.Now().Add(-time.Hour)
	writeConf(t, file, "name = one\n", start)

	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	sf, err := NewSrcFile(f)
	if err != nil {
		t.Fatal(err)
	}
	assertProp(t, sf, "name", "one")

	writeConf(t, file, "name = two\n", start.Add(time.Minute))
	assertProp(t, sf, "name", "two")
}
