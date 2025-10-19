package util

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func SetupZeroLogger(isDebugMode bool) {
	currentDay := time.Now().Format("2006-01-02")
	logFileName := "/tmp/" + fmt.Sprintf("app-%s.log", currentDay)

	// Open the file in append mode, create if it doesn't exist
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}

	var logLevel zerolog.Level
	if isDebugMode {
		logLevel = zerolog.DebugLevel
	} else {
		logLevel = zerolog.WarnLevel
	}
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(logLevel)

	// MultiWriter for console + file if in debug, otherwise just file
	var writer io.Writer
	if isDebugMode {
		writer = io.MultiWriter(consoleWriter, logFile)
	} else {
		writer = logFile
	}

	// Configure global logger
	log.Logger = zerolog.New(writer).With().Timestamp().Logger()

	if logLevel == zerolog.DebugLevel {
		LogDebug("## Debug mode enabled! ## ")
	}
}

func LogErrorWithStackTrace(err error) {
	log.Error().Stack().Err(err).Msg("")
}

func LogError(err error) {
	log.Error().Err(err).Msg("")
}

func LogWarning(message string) {
	log.Warn().Msg(message)
}

func LogInfo(message string) {
	log.Info().Msg(message)
}

func LogDebug(message string) {
	log.Debug().Msg(message)
}
