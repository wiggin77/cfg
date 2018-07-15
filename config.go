package config

import (
	"errors"
	"strconv"
	"sync"
)

var ErrNotFound = errors.New("not found")

// Config provides methods for retrieving property values from one or more
// configuration sources.
type Config struct {
	mutexSrc       sync.RWMutex
	mutexListeners sync.RWMutex
	srcs           []Source
	chgListeners   []ChangedListener
	propListeners  []ChangedPropListener
}

// PrependSource inserts one or more `Sources` at the beginning of
// the list of sources such that the first source will be the
// source first checked when resolving a property value.
func (config *Config) PrependSource(srcs ...Source) {
	config.mutexSrc.Lock()
	defer config.mutexSrc.Unlock()

	config.srcs = append(srcs, config.srcs...)
}

// AppendSource appends one or more `Sources` at the end of
// the list of sources such that the last source will be the
// source last checked when resolving a property value.
func (config *Config) AppendSource(srcs ...Source) {
	config.mutexSrc.Lock()
	defer config.mutexSrc.Unlock()

	config.srcs = append(config.srcs, srcs...)
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
	for _, src := range config.srcs {
		if val, ok = src.GetProp(name); ok {
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

//
// Types:  String, Int, Int64, Float64, Bool, milliseconds
//

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

// AddChangedPropListener adds a listener that will receive notifications
// whenever a property value changes.
func (config *Config) AddChangedPropListener(l ChangedPropListener) {
	config.mutexListeners.Lock()
	defer config.mutexListeners.Unlock()

	config.propListeners = append(config.propListeners, l)
}

// RemoveChangedPropListener removes all instances of a ChangedPropListener.
// Returns `ErrNotFound` if the listener was not present.
func (config *Config) RemoveChangedPropListener(l ChangedPropListener) error {
	config.mutexListeners.Lock()
	defer config.mutexListeners.Unlock()

	dest := make([]ChangedPropListener, 0, len(config.propListeners))
	err := ErrNotFound

	for _, s := range config.propListeners {
		if s != l {
			dest = append(dest, s)
		} else {
			err = nil
		}
	}
	config.propListeners = dest
	return err
}
