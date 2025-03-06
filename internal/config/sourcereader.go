package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultConfigFilename = ".tgnotifier.yml"
)

// SourceReaderFactory is a config reader source factory.
type SourceReaderFactory func() (io.ReadCloser, error)

// FromDefaultFileDiscovery returns FromFileDiscovery
// factory with possible default paths of the config file:
//
//   - $HOME/.config/.tgnotifier.yml
//   - $HOME/.tgnotifier.yml
//   - path/to/tgnotifier/dir/.tgnotifier.yml
//
// The configFilePath argument points to the explicit config file path if not empty
// and must exist if defined.
func FromDefaultFileDiscovery(configFilePath string) SourceReaderFactory {
	if configFilePath != "" {
		return FromFile(configFilePath)
	}

	paths := make([]string, 0, 3)

	dirGetters := []func() (string, error){
		os.UserConfigDir,
		os.UserHomeDir,
		func() (string, error) {
			path, err := os.Executable()
			if err != nil {
				return "", err
			}

			return filepath.Dir(path), nil
		},
	}

	for _, dirGetter := range dirGetters {
		dir, err := dirGetter()
		if err != nil {
			continue
		}

		paths = append(paths, filepath.Join(dir, defaultConfigFilename))
	}

	return FromFileDiscovery(paths...)
}

// FromFileDiscovery uses first existing file from list as a source.
// Returns FromNil factory if no file found.
func FromFileDiscovery(paths ...string) SourceReaderFactory {
	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil || stat.IsDir() {
			continue
		}

		return FromFile(path)
	}

	return FromNil()
}

// FromFile is a file config source factory.
func FromFile(path string) SourceReaderFactory {
	if path == "" {
		return FromNil()
	}

	return func() (io.ReadCloser, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("read config file %s: %w", path, err)
		}

		return f, nil
	}
}

// FromReader is a io.Reader source factory.
func FromReader(r io.Reader) SourceReaderFactory {
	return func() (io.ReadCloser, error) {
		return io.NopCloser(r), nil
	}
}

// FromString is a string config source factory.
func FromString(s string) SourceReaderFactory {
	return FromReader(strings.NewReader(s))
}

// FromNil is a nil config source factory.
func FromNil() SourceReaderFactory {
	return func() (io.ReadCloser, error) {
		return nil, nil
	}
}
