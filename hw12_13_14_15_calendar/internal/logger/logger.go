package logger

import (
	"errors"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

var (
	ErrUndefinedLevel = errors.New("level is not found")
	logger            = &Logger{&logrus.Logger{}}
)

func init() {
	logger.log.SetFormatter(&logrus.JSONFormatter{})
}

type Fields = logrus.Fields

type Logger struct {
	log *logrus.Logger
}

func GetLogger() *logrus.Logger {
	return logger.log
}

func SetLevel(levelStr string) {
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		fmt.Printf("failed to parse level '%s', err: %s \n'InfoLevel' will be used\n", level, err)
		level = logrus.InfoLevel
	}
	logger.log.SetLevel(level)
}

func SetOutput(output io.Writer) {
	logger.log.SetOutput(output)
}

func Info(msg string, f Fields) {
	if f != nil {
		logger.log.WithFields(f).Info(msg)
	} else {
		logger.log.Info(msg)
	}
}

func Error(msg string, f Fields) {
	if f != nil {
		logger.log.WithFields(f).Error(msg)
	} else {
		logger.log.Error(msg)
	}
}

func Warning(msg string, f Fields) {
	if f != nil {
		logger.log.WithFields(f).Warning(msg)
	} else {
		logger.log.Warning(msg)
	}
}

func Fatal(msg string, f Fields) {
	if f != nil {
		logger.log.WithFields(f).Fatal(msg)
	} else {
		logger.log.Fatal(msg)
	}
}
