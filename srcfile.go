package cfg

import (
	"io"
	"os"
	"time"

	"github.com/wiggin77/cfg/ini"
)

// SrcFile is a configuration `Source` backed by a file containing
// name/value pairs or INI format.
type SrcFile struct {
	AbstractSourceMonitor
	ini      ini.Ini
	file     *os.File
	filespec string
}

// NewSrcFileFromFilespec creates a new SrcFile with the specified filespec.
// The file is opened only while it is being read, so it is not held open
// between reloads.
func NewSrcFileFromFilespec(filespec string) (*SrcFile, error) {
	sf := &SrcFile{filespec: filespec}
	sf.freq = time.Minute
	if err := sf.ini.LoadFromFilespec(filespec); err != nil {
		return nil, err
	}
	return sf, nil
}

// NewSrcFile creates a new SrcFile with the specified os.File.
// The caller retains ownership of the file and must keep it open.
func NewSrcFile(file *os.File) (*SrcFile, error) {
	sf := &SrcFile{}
	sf.freq = time.Minute
	sf.file = file
	if err := sf.ini.LoadFromFile(file); err != nil {
		return nil, err
	}
	return sf, nil
}

// GetProps fetches all the properties from a source and returns
// them as a map.
func (sf *SrcFile) GetProps() (map[string]string, error) {
	lm, err := sf.GetLastModified()
	if err != nil {
		return nil, err
	}

	// Check if we need to reload.
	if sf.ini.GetLastModified() != lm {
		if err := sf.reload(); err != nil {
			return nil, err
		}
	}
	return sf.ini.ToMap(), nil
}

func (sf *SrcFile) reload() error {
	if sf.filespec != "" {
		return sf.ini.LoadFromFilespec(sf.filespec)
	}
	// The previous load left the offset at EOF.
	if _, err := sf.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	return sf.ini.LoadFromFile(sf.file)
}

// GetLastModified returns the time of the latest modification to any
// property value within the source.
func (sf *SrcFile) GetLastModified() (time.Time, error) {
	var fi os.FileInfo
	var err error
	if sf.filespec != "" {
		// Stat by path so a file replaced by rename (as many editors save) is seen.
		fi, err = os.Stat(sf.filespec)
	} else {
		fi, err = sf.file.Stat()
	}
	if err != nil {
		return time.Now(), err
	}
	return fi.ModTime(), nil
}
