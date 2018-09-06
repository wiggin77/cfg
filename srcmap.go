package config

import (
	"time"
)

// SrcMap is a configuration `Source` backed by a simple map.
type SrcMap struct {
	AbstractSourceMonitor
	m  map[string]string
	lm time.Time
}

// NewSrcMap creates an empty `SrcMap`.
func NewSrcMap() *SrcMap {
	sm := &SrcMap{}
	sm.m = make(map[string]string)
	sm.lm = time.Now()
	sm.freq = time.Minute
	return sm
}

// NewSrcMapFromMap creates a `SrcMap` containing a copy of the
// specified map.
func NewSrcMapFromMap(mapIn map[string]string) *SrcMap {
	sm := NewSrcMap()
	sm.PutAll(mapIn)
	return sm
}

// Put inserts or updates a value in the `SrcMap`.
func (sm *SrcMap) Put(key string, val string) {
	sm.mutex.Lock()
	sm.m[key] = val
	sm.lm = time.Now()
	sm.mutex.Unlock()
}

// PutAll inserts a copy of `mapIn` into the `SrcMap`
func (sm *SrcMap) PutAll(mapIn map[string]string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for k, v := range mapIn {
		sm.m[k] = v
	}
}

// GetProp fetches the value of a named property. The value is
// returned as a string unless the name is not found, in which
// case `ok=false` is returned.
func (sm *SrcMap) GetProp(name string) (val string, ok bool) {
	sm.mutex.Lock()
	val, ok = sm.m[name]
	sm.mutex.Unlock()
	return
}

// GetLastModified returns the time of the latest modification to any
// property value within the source. If a source does not support
// modifying properties at runtime then the zero value for `Time`
// should be returned to ensure reload events are not generated.
func (sm *SrcMap) GetLastModified() (last time.Time) {
	sm.mutex.Lock()
	last = sm.lm
	sm.mutex.Unlock()
	return
}
