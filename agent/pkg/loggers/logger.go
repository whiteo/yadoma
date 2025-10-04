// Package loggers initializes and configures the global structured logger used across the project.
//
// It sets up zerolog for JSON logging — suitable for production and containerized environments.
// The logger is initialized automatically on package import and writes logs to stdout.
//
// Example log output:
//
//	{"level":"info","time":"2025-10-04T10:12:33Z","message":"agent started","module":"core"}
//
// Author: Leo Tanas (https://github.com/whiteo)
package loggers

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type safeWriter struct {
	w io.Writer
}

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = false
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.TimestampFieldName = "time"

	writer := &safeWriter{w: os.Stdout}

	log.Logger = zerolog.New(writer).With().Timestamp().Logger()
}

// Write implements io.Writer for safeWriter.
// It writes the given bytes to the underlying writer.
// On error, it logs a diagnostic message to stderr.
func (s *safeWriter) Write(p []byte) (int, error) {
	n, err := s.w.Write(p)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Logger write error: %v, data: %s\n", err, string(p))
	}
	return n, err
}
