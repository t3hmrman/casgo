package cas

import (
	"fmt"
	"github.com/t3hmrman/casgo/cas/Godeps/_workspace/src/github.com/GeertJohan/go.rice"
	"log"
	"os"
)

const (
	_ = iota
	DEBUG
	INFO
	WARN
)

// Utility functions for logging messages
func logMessage(actualLogLevel, msgLogLevel, msg string) {
	if actualLogLevel <= msgLogLevel {
		log.Printf("[%s] %s", msgLogLevel, msg)
	}
}

// Utility function for logging message
func logMessagef(actualLogLevel, msgLogLevel, format string, msgArgs ...interface{}) {
	if actualLogLevel <= msgLogLevel {
		log.Printf("[%s] %s", msgLogLevel, fmt.Sprintf(format, msgArgs...))
	}
}

// Small tuple implementation
func (t *StringTuple) First() string {
	return t[0]
}

// Small tuple implementation
func (t *StringTuple) Second() string {
	return t[1]
}

func ListFilesInBox(box *rice.Box, prefix string) ([]string, error) {
	var files []string

	// Walk files, make list of templates
	walkErr := box.Walk("", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, prefix+path)
		}

		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	return files, nil
}
