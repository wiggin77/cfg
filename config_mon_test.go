package config

import (
	"math/rand"
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

func (n *Notify) Changed(cfg *Config, src SourceMonitored) {
	atomic.AddInt32(&n.count, 1)
}

func TestConfig_Monitor(t *testing.T) {
	const loops int = 100
	mapSrc := makeSrc(time.Millisecond * 20)
	cfg := &Config{}
	cfg.AppendSource(mapSrc)
	var wg sync.WaitGroup

	notify := &Notify{}
	cfg.AddChangedListener(notify)

	r := rand.New(rand.NewSource(77))
	for i := 0; i < loops; i++ {
		wg.Add(1)
		go func(idx int, rnd int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rnd))
			mapSrc.Put("prop0", strconv.Itoa(idx))
		}(i, r.Intn(50)+50)
	}

	wg.Wait()
	count := atomic.LoadInt32(&notify.count)
	if count <= 0 {
		t.Error("ChangedListener was not called")
	}

	atomic.StoreInt32(&notify.count, 0)
	cfg.RemoveChangedListener(notify)

	for i := 0; i < loops; i++ {
		wg.Add(1)
		go func(idx int, rnd int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rnd))
			mapSrc.Put("prop1", strconv.Itoa(idx))
		}(i, r.Intn(50)+50)
	}

	wg.Wait()
	count = atomic.LoadInt32(&notify.count)
	if count != 0 {
		t.Errorf("ChangedListener was called %d times after being removed.", count)
	}
}
