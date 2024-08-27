package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func New(level string) *zerolog.Logger {
	var l zerolog.Level
	level = strings.ToLower(level)

	levelMap := map[string]zerolog.Level{
		"error": zerolog.ErrorLevel,
		"warn":  zerolog.WarnLevel,
		"info":  zerolog.InfoLevel,
		"debug": zerolog.DebugLevel,
	}

	l = zerolog.InfoLevel
	if level, exists := levelMap[level]; exists {
		l = level
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &logger
}
