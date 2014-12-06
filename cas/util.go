package cas

import (
	"log"
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
func logMessagef(actualLogLevel, msgLogLevel, format string, msgArgs...interface{} ) {
	if actualLogLevel <= msgLogLevel {
		log.Printf("[%s] " + format, msgLogLevel, msgArgs)
	}
}

