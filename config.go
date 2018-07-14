package config

import (
	"sync"
)

type Config struct {
	mutex         sync.Mutex
	src           []Source
	chglisteners  []ChangedListener
	proplisteners []ChangedPropListener
}

// PrependSource inserts one or more `Sources` at the beginning of
// the list of sources such that the first source will be the
// source first checked when resolving a property value.
func (config *Config) PrependSource(src ...Source) {

}

// AppendSource appends one or more `Sources` at the end of
// the list of sources such that the last source will be the
// source last checked when resolving a property value.
func (config *Config) AppendSource(src ...Source) {

}

func (config *Config) AddChangedListener(l ChangedListener) {

}

func (config *Config) RemoveChangedListener(l ChangedListener) {

}

func (config *Config) AddChangedPropListener(l ChangedPropListener) {

}

func (config *Config) RemoveChangedPropListener(l ChangedPropListener) {

}
