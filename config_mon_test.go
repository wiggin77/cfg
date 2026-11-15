package cfg

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func makeSrc(freq time.Duration) *SrcMap {
	m := make(map[string]string)
	m["prop0"] = "0"
	m["prop1"] = "1"
	m["prop2"] = "2"
	m["prop3"] = "3"
	m["prop4"] = "4"
	mapSrc := NewSrcMapFromMap(m)
	mapSrc.SetMonitorFreq(freq)
	return mapSrc
}

type Notify struct {
	count int32
}

func (n *Notify) ConfigChanged(cfg *Config, src SourceMonitored) {
	atomic.AddInt32(&n.count, 1)
}

func TestConfig_Monitor(t *testing.T) {
	config := &Config{}
	defer config.Shutdown()

	mapSrc := makeSrc(time.Millisecond * 20)
	config.AppendSource(mapSrc)

	notify := &Notify{}
	config.AddChangedListener(notify)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	done := time.After(2 * time.Second)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	isDone := false
	actual := int32(0)
	for !isDone {
		select {
		case <-ticker.C:
			mapSrc.Put("prop1", strconv.Itoa(rnd.Intn(30)))
			actual++
		case <-done:
			isDone = true
		}
	}

	deadline := time.Now().Add(time.Second)
	for atomic.LoadInt32(&notify.count) < actual && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}

	count := atomic.LoadInt32(&notify.count)
	if count != actual {
		t.Errorf("ChangedListener was called %d times; expected %d", count, actual)
	}

	// make sure listener not called after removal.
	atomic.StoreInt32(&notify.count, 0)
	config.RemoveChangedListener(notify)

	mapSrc.Put("prop1", "n/a")
	time.Sleep(50 * time.Millisecond)
	mapSrc.Put("prop2", "n/a")
	time.Sleep(50 * time.Millisecond)

	count = atomic.LoadInt32(&notify.count)
	if count != 0 {
		t.Errorf("ChangedListener was called %d times after being removed.", count)
	}
}

type faultySrc struct {
	AbstractSourceMonitor
	mu       sync.Mutex
	props    map[string]string
	lm       time.Time
	propsErr error
	lmErr    error
}

func newFaultySrc(freq time.Duration, props map[string]string) *faultySrc {
	fs := &faultySrc{props: props, lm: time.Now()}
	fs.SetMonitorFreq(freq)
	return fs
}

func (fs *faultySrc) GetProps() (map[string]string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.propsErr != nil {
		return nil, fs.propsErr
	}
	return fs.props, nil
}

func (fs *faultySrc) GetLastModified() (time.Time, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.lmErr != nil {
		return time.Time{}, fs.lmErr
	}
	return fs.lm, nil
}

func (fs *faultySrc) put(key, val string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	props := make(map[string]string, len(fs.props)+1)
	for k, v := range fs.props {
		props[k] = v
	}
	props[key] = val
	fs.props = props
	fs.lm = time.Now()
}

func (fs *faultySrc) setErrs(propsErr, lmErr error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.propsErr = propsErr
	fs.lmErr = lmErr
}

func waitForCount(t *testing.T, n *Notify, want int32) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for atomic.LoadInt32(&n.count) < want && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if got := atomic.LoadInt32(&n.count); got != want {
		t.Fatalf("ChangedListener was called %d times; expected %d", got, want)
	}
}

func assertCountStays(t *testing.T, n *Notify, want int32) {
	t.Helper()
	time.Sleep(100 * time.Millisecond)
	if got := atomic.LoadInt32(&n.count); got != want {
		t.Fatalf("ChangedListener was called %d times; expected %d", got, want)
	}
}

func assertConfigProp(t *testing.T, config *Config, key, want string) {
	t.Helper()
	if got, err := config.String(key, ""); err != nil || got != want {
		t.Fatalf("%s = %q (err %v), want %q", key, got, err, want)
	}
}

func TestConfig_MonitorNoChangeNoNotify(t *testing.T) {
	config := &Config{}
	defer config.Shutdown()
	config.AppendSource(makeSrc(10 * time.Millisecond))

	notify := &Notify{}
	config.AddChangedListener(notify)
	assertCountStays(t, notify, 0)
}

func TestConfig_MonitorSrcFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.conf")
	start := time.Now().Add(-time.Hour)
	writeConf(t, file, "name = one\n", start)

	sf, err := NewSrcFileFromFilespec(file)
	if err != nil {
		t.Fatal(err)
	}
	sf.SetMonitorFreq(10 * time.Millisecond)

	config := &Config{}
	defer config.Shutdown()
	config.AppendSource(sf)
	notify := &Notify{}
	config.AddChangedListener(notify)
	assertConfigProp(t, config, "name", "one")

	writeConf(t, file, "name = two\n", start.Add(time.Minute))
	waitForCount(t, notify, 1)
	assertConfigProp(t, config, "name", "two")

	// A missing file keeps the last loaded values and is picked up again once restored.
	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	assertCountStays(t, notify, 1)
	assertConfigProp(t, config, "name", "two")

	writeConf(t, file, "name = three\n", start.Add(2*time.Minute))
	waitForCount(t, notify, 2)
	assertConfigProp(t, config, "name", "three")
}

func TestConfig_MonitorPauseResume(t *testing.T) {
	config := &Config{pausedPollFreq: 10 * time.Millisecond}
	defer config.Shutdown()

	mapSrc := makeSrc(0)
	config.AppendSource(mapSrc)
	notify := &Notify{}
	config.AddChangedListener(notify)

	mapSrc.Put("prop1", "paused")
	assertCountStays(t, notify, 0)

	mapSrc.SetMonitorFreq(10 * time.Millisecond)
	waitForCount(t, notify, 1)
	assertConfigProp(t, config, "prop1", "paused")

	mapSrc.SetMonitorFreq(0)
	time.Sleep(50 * time.Millisecond)
	mapSrc.Put("prop1", "paused again")
	assertCountStays(t, notify, 1)

	mapSrc.SetMonitorFreq(10 * time.Millisecond)
	waitForCount(t, notify, 2)
	assertConfigProp(t, config, "prop1", "paused again")
}

func TestConfig_MonitorRetriesAfterGetPropsError(t *testing.T) {
	src := newFaultySrc(10*time.Millisecond, map[string]string{"prop": "1"})
	config := &Config{}
	defer config.Shutdown()
	config.AppendSource(src)
	notify := &Notify{}
	config.AddChangedListener(notify)

	src.put("prop", "2")
	src.setErrs(errors.New("props unavailable"), nil)
	assertCountStays(t, notify, 0)
	assertConfigProp(t, config, "prop", "1")

	src.setErrs(nil, nil)
	waitForCount(t, notify, 1)
	assertConfigProp(t, config, "prop", "2")
}

func TestConfig_MonitorSurvivesGetLastModifiedError(t *testing.T) {
	src := newFaultySrc(10*time.Millisecond, map[string]string{"prop": "1"})
	config := &Config{}
	defer config.Shutdown()
	config.AppendSource(src)
	notify := &Notify{}
	config.AddChangedListener(notify)

	src.setErrs(nil, errors.New("stat failed"))
	assertCountStays(t, notify, 0)

	src.setErrs(nil, nil)
	src.put("prop", "2")
	waitForCount(t, notify, 1)
	assertConfigProp(t, config, "prop", "2")
}

func TestConfig_GetPropsErrorOnAdd(t *testing.T) {
	config := &Config{}
	defer config.Shutdown()
	if config.ShouldPanicOnError() {
		t.Fatal("ShouldPanicOnError should default to false")
	}

	src := newFaultySrc(time.Minute, map[string]string{"prop": "1"})
	src.setErrs(errors.New("props unavailable"), nil)
	config.AppendSource(src)
	if _, err := config.String("prop", ""); err != ErrNotFound {
		t.Errorf("expected ErrNotFound for a source that failed to load, got %v", err)
	}

	config.SetWantPanicOnError(true)
	if !config.ShouldPanicOnError() {
		t.Fatal("ShouldPanicOnError should be true after SetWantPanicOnError(true)")
	}
	defer func() {
		if recover() == nil {
			t.Error("expected panic when GetProps fails with panic on error enabled")
		}
	}()
	config.AppendSource(src)
}

type panicListener struct{}

func (panicListener) ConfigChanged(cfg *Config, src SourceMonitored) {
	panic("listener panic (expected by test)")
}

func TestConfig_ListenerPanicDoesNotBlockOthers(t *testing.T) {
	config := &Config{}
	defer config.Shutdown()
	mapSrc := makeSrc(10 * time.Millisecond)
	config.AppendSource(mapSrc)

	notify := &Notify{}
	config.AddChangedListener(panicListener{})
	config.AddChangedListener(notify)

	mapSrc.Put("prop1", "a")
	waitForCount(t, notify, 1)
	mapSrc.Put("prop1", "b")
	waitForCount(t, notify, 2)
}

func TestConfig_ShutdownStopsMonitor(t *testing.T) {
	config := &Config{}
	mapSrc := makeSrc(10 * time.Millisecond)
	config.AppendSource(mapSrc)
	notify := &Notify{}
	config.AddChangedListener(notify)

	config.Shutdown()
	mapSrc.Put("prop1", "after shutdown")
	assertCountStays(t, notify, 0)
}
