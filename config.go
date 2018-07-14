package config

import (
	"errors"
	"sync"
)

var errNotFound = errors.New("not found")

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

// String returns the value of the named prop as a string,
// or `ok=false` if the name is not found.
//
// Each `Source` is checked, in the order they are added via
// `ApppendSource` and `PrependSource`, until a value for the
// property is found.
func (config *Config) String(name string) (val string, ok bool) {
	config.mutexSrc.RLock()
	defer config.mutexSrc.RUnlock()

	for _, src := range config.srcs {
		if val, ok = src.GetProp(name); ok {
			return
		}
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
// Returns non-nil error if the listener was not present.
func (config *Config) RemoveChangedListener(l ChangedListener) error {
	config.mutexListeners.Lock()
	defer config.mutexListeners.Unlock()

	dest := make([]ChangedListener, 0, len(config.chgListeners))
	err := errNotFound

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
// Returns non-nil error if the listener was not present.
func (config *Config) RemoveChangedPropListener(l ChangedPropListener) error {
	config.mutexListeners.Lock()
	defer config.mutexListeners.Unlock()

	dest := make([]ChangedPropListener, 0, len(config.propListeners))
	err := errNotFound

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
