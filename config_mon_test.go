package config

import (
	"math/rand"
	"strconv"
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
	cfg := &Config{}
	defer cfg.Shutdown()

	mapSrc := makeSrc(time.Millisecond * 20)
	cfg.AppendSource(mapSrc)

	notify := &Notify{}
	cfg.AddChangedListener(notify)

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

	count := atomic.LoadInt32(&notify.count)
	if count != actual {
		t.Errorf("ChangedListener was called %d times; expected %d", count, actual)
	}

	// make sure listener not called after removal.
	atomic.StoreInt32(&notify.count, 0)
	cfg.RemoveChangedListener(notify)

	mapSrc.Put("prop1", "n/a")
	time.Sleep(50 * time.Millisecond)
	mapSrc.Put("prop2", "n/a")
	time.Sleep(50 * time.Millisecond)

	count = atomic.LoadInt32(&notify.count)
	if count != 0 {
		t.Errorf("ChangedListener was called %d times after being removed.", count)
	}
}
