package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	timeutil "github.com/wiggin77/config/time"
)

// ErrNotFound returned when an operation is attempted on a
// resource that doesn't exist, such as fetching a non-existing
// property name.
var ErrNotFound = errors.New("not found")

// Config provides methods for retrieving property values from one or more
// configuration sources.
type Config struct {
	mutexSrc       sync.RWMutex
	mutexListeners sync.RWMutex
	srcs           []Source
	chgListeners   []ChangedListener
	shutdown       chan interface{}
}

// PrependSource inserts one or more `Sources` at the beginning of
// the list of sources such that the first source will be the
// source checked first when resolving a property value.
func (config *Config) PrependSource(srcs ...Source) {
	config.mutexSrc.Lock()
	defer config.mutexSrc.Unlock()

	if config.shutdown == nil {
		config.shutdown = make(chan interface{})
	}
	config.srcs = append(srcs, config.srcs...)

	for _, src := range srcs {
		if sm, ok := src.(SourceMonitored); ok {
			config.monitor(sm)
		}
	}
}

// AppendSource appends one or more `Sources` at the end of
// the list of sources such that the last source will be the
// source checked last when resolving a property value.
func (config *Config) AppendSource(srcs ...Source) {
	config.mutexSrc.Lock()
	defer config.mutexSrc.Unlock()

	if config.shutdown == nil {
		config.shutdown = make(chan interface{})
	}
	config.srcs = append(config.srcs, srcs...)

	for _, src := range srcs {
		if sm, ok := src.(SourceMonitored); ok {
			config.monitor(sm)
		}
	}
}

// String returns the value of the named prop as a string.
// If the property is not found then the supplied default `def`
// and `ErrNotFound` are returned.
//
// Each `Source` is checked, in the order they are added via
// `ApppendSource` and `PrependSource`, until a value for the
// property is found.
func (config *Config) String(name string, def string) (val string, err error) {
	config.mutexSrc.RLock()
	defer config.mutexSrc.RUnlock()

	var ok bool
	var s string
	for _, src := range config.srcs {
		if s, ok = src.GetProp(name); ok {
			val = strings.TrimSpace(s)
			err = nil
			return
		}
	}
	err = ErrNotFound
	val = def
	return
}

// Int returns the value of the named prop as an `int`.
// If the property is not found then the supplied default `def`
// and `ErrNotFound` are returned.
//
// See config.String
func (config *Config) Int(name string, def int) (val int, err error) {
	var s string
	if s, err = config.String(name, ""); err == nil {
		var i int64
		if i, err = strconv.ParseInt(s, 10, 32); err == nil {
			val = int(i)
		}
	}
	if err != nil {
		val = def
	}
	return
}

// Int64 returns the value of the named prop as an `int64`.
// If the property is not found then the supplied default `def`
// and `ErrNotFound` are returned.
//
// See config.String
func (config *Config) Int64(name string, def int64) (val int64, err error) {
	var s string
	if s, err = config.String(name, ""); err == nil {
		val, err = strconv.ParseInt(s, 10, 64)
	}
	if err != nil {
		val = def
	}
	return
}

// Float64 returns the value of the named prop as a `float64`.
// If the property is not found then the supplied default `def`
// and `ErrNotFound` are returned.
//
// See config.String
func (config *Config) Float64(name string, def float64) (val float64, err error) {
	var s string
	if s, err = config.String(name, ""); err == nil {
		val, err = strconv.ParseFloat(s, 64)
	}
	if err != nil {
		val = def
	}
	return
}

// Bool returns the value of the named prop as a `bool`.
// If the property is not found then the supplied default `def`
// and `ErrNotFound` are returned.
//
// Supports (t, true, 1, y, yes) for true, and (f, false, 0, n, no) for false,
// all case-insensitive.
//
// See config.String
func (config *Config) Bool(name string, def bool) (val bool, err error) {
	var s string
	if s, err = config.String(name, ""); err == nil {
		switch strings.ToLower(s) {
		case "t", "true", "1", "y", "yes":
			val = true
		case "f", "false", "0", "n", "no":
			val = false
		default:
			err = errors.New("invalid syntax")
		}
	}
	if err != nil {
		val = def
	}
	return
}

// Duration returns the value of the named prop as a `time.Duration`, representing
// a span of time.
//
// Units of measure are supported: ms, sec, min, hour, day, week, year.
// See config.UnitsToMillis for a complete list of units supported.
//
// If the property is not found then the supplied default `def`
// and `ErrNotFound` are returned.
//
// See config.String
func (config *Config) Duration(name string, def time.Duration) (val time.Duration, err error) {
	var s string
	if s, err = config.String(name, ""); err == nil {
		var ms int64
		ms, err = timeutil.ParseMilliseconds(s)
		val = time.Duration(ms) * time.Millisecond
	}
	if err != nil {
		val = def
	}
	return
}

// AddChangedListener adds a listener that will receive notifications
// whenever one or more property values change within the config.
func (config *Config) AddChangedListener(l ChangedListener) {
	config.mutexListeners.Lock()
	defer config.mutexListeners.Unlock()

	config.chgListeners = append(config.chgListeners, l)
}

// RemoveChangedListener removes all instances of a ChangedListener.
// Returns `ErrNotFound` if the listener was not present.
func (config *Config) RemoveChangedListener(l ChangedListener) error {
	config.mutexListeners.Lock()
	defer config.mutexListeners.Unlock()

	dest := make([]ChangedListener, 0, len(config.chgListeners))
	err := ErrNotFound

	for _, s := range config.chgListeners {
		if s != l {
			dest = append(dest, s)
		} else {
			err = nil
		}
	}
	config.chgListeners = dest
	return err
}

// Shutdown can be called to stop monitoring of all config sources.
func (config *Config) Shutdown() {
	config.mutexSrc.RLock()
	defer config.mutexSrc.RUnlock()
	if config.shutdown != nil {
		close(config.shutdown)
	}
}

// onSourceChanged is called whenever a config source properties
// have changed.
func (config *Config) onSourceChanged(src SourceMonitored) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p)
		}
	}()
	config.mutexListeners.RLock()
	defer config.mutexListeners.RUnlock()
	for _, l := range config.chgListeners {
		l.Changed(config, src)
	}
}

// monitor periodically checks a config source for changes.
func (config *Config) monitor(src SourceMonitored) {
	go func(sm SourceMonitored, shutdown <-chan interface{}) {
		paused := false
		last := time.Time{}
		freq := src.GetMonitorFreq()
		if freq <= 0 {
			paused = true
			freq = 10
			last = src.GetLastModified()
		}
		timer := time.NewTimer(freq)
		for {
			select {
			case <-timer.C:
				if !paused {
					latest := src.GetLastModified()
					if last.Before(latest) {
						last = latest
						config.onSourceChanged(sm)
					}
				}
				freq = src.GetMonitorFreq()
				if freq <= 0 {
					paused = true
					freq = 10
				} else {
					paused = false
				}
				timer.Reset(freq)
			case <-shutdown:
				// stop the timer and exit
				if !timer.Stop() {
					<-timer.C
				}
				return
			}
		}
	}(src, config.shutdown)
}
