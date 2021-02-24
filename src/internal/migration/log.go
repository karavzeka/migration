package migration

import "fmt"

// Log represents the logger for migration lib
type Log struct {
	verbose bool
}

// Printf prints out formatted string into a log
func (l *Log) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// Verbose shows if verbose print enabled
func (l *Log) Verbose() bool {
	return l.verbose
}
