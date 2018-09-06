package config

import (
	"os"
	"time"

	"github.com/wiggin77/config/ini"
)

// SrcFile is a configuration `Source` backed by a file containing
// name/value pairs or INI format.
type SrcFile struct {
	AbstractSourceMonitor
	ini  ini.Ini
	file os.File
}

// NewSrcFileFromFilespec creates a new SrcFile with the specified filespec.
func NewSrcFileFromFilespec(filespec string) (*SrcFile, error) {
	file, err := os.Open(filespec)
	if err != nil {
		return nil, err
	}
	return NewSrcFile(file)
}

// NewSrcFile creates a new SrcFile with the specified os.File.
func NewSrcFile(file *os.File) (*SrcFile, error) {
	sf := &SrcFile{}
	sf.freq = time.Minute
	if err := sf.ini.LoadFromFile(file); err != nil {
		return nil, err
	}
	return sf, nil
}

// GetLastModified returns the time of the latest modification to any
// property value within the source.
func (sf *SrcFile) GetLastModified() time.Time {

	// Since a source could change between calls to GetLastModified,
	// then it is currently possible to fetch changed properties
	// before the listeners are notified of changes.
	//
	// This needs a push ability for sources to notify of changes,
	// while using the GetLastModified calls as a convenient way to
	// know when to check for changes (if that makes sense for source type).

}
