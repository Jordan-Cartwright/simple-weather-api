package util

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// InitializeLogging takes an io writer to output (defaulting to os.Stderr
// when nil), a log level string and a format string that when set to json
// will output logs in a json format. Otherwise a text formatter will be used.
func InitializeLogging(wr io.Writer, level string, format string) {
	if wr == nil {
		wr = os.Stderr
	}

	log.SetOutput(wr)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	if format == "json" {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	setLoggingLevel(level)
}

// setLoggingLevel parses a string and attempts to set the logging level
// using that string. If the string is not valid the info log level will
// be used as a fall back.
func setLoggingLevel(level string) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		lvl = log.InfoLevel
		log.Warnf("failed to parse log-level '%s', defaulting to 'info'", level)
	}
	log.SetLevel(lvl)
}
