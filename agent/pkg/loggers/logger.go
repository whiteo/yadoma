package loggers

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
	format      = "%s| %s|%s"
)

type safeWriter struct {
	w io.Writer
}

func init() {
	level := getEnv("LOG_LEVEL", "debug")
	switch strings.ToLower(level) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	out := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}

	out.FormatLevel = func(i interface{}) string {
		if i == nil {
			return ""
		}
		l := strings.ToUpper(fmt.Sprintf("%-6s", i))
		switch l {
		case "INFO  ":
			return fmt.Sprintf(format, colorGreen, l, colorReset)
		case "WARN  ":
			return fmt.Sprintf(format, colorYellow, l, colorReset)
		case "ERROR ":
			return fmt.Sprintf(format, colorRed, l, colorReset)
		case "DEBUG ":
			return fmt.Sprintf(format, colorBlue, l, colorReset)
		default:
			return fmt.Sprintf("| %s|", l)
		}
	}

	out.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s***", i)
	}

	out.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}

	out.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	//nolint:reassign //reassign is needed to set the logger
	log.Logger = zerolog.New(&safeWriter{out}).With().Timestamp().Logger()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (s *safeWriter) Write(p []byte) (int, error) {
	n, err := s.w.Write(p)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Logger write error: %v, data: %s\n", err, string(p))
	}
	return n, err
}
