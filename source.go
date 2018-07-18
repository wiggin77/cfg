package config

import (
	"time"
)

// Source is the interface required for any source of name/value pairs.
type Source interface {

	// GetProp fetches the value of a named property. The value is
	// returned as a string unless the name is not found, in which
	// case `ok=false` is returned.
	GetProp(name string) (val string, ok bool)
}

// SourceMonitored is the interface required for any config source that is
// monitored for changes.
type SourceMonitored interface {

	// GetLastModified returns the time of the latest modification to any
	// property value within the source. If a source does not support
	// modifying properties at runtime then the zero value for `Time`
	// should be returned to ensure reload events are not generated.
	GetLastModified() time.Time

	// GetMonitorFreq returns the frequency in milliseconds between
	// checks for changes to this config source.
	GetMonitorFreq() int64
}
