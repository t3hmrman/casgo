package cas

import (
	"log"
	"fmt"
)

const (
	_ = iota
	DEBUG
	INFO
	WARN
)

// Utility functions for logging messages
func logMessage(actualLogLevel, msgLogLevel, msg string, ) {
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
