// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package log

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

const (
	EnvOutput string = "LOGGER_OUTPUT"
	EnvFmt    string = "LOGGER_FMT"
	EnvLevel  string = "LOGGER_LEVEL"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()

	//	Run a process to adjust the logging parameters if the environment variables o
	//	are updated. This allows for immediate logging changes without restarting
	// 	the affected service.
	go func() {
		prevOut, prevFmt, prevLevel := "", "", ""
		for {
			if out := strings.ToLower(os.Getenv(EnvOutput)); len(out) > 1 {
				if prevOut != out {
					prevOut = out
					logger.Out = output(out)
				}
			}

			if fmt := strings.ToLower(os.Getenv(EnvFmt)); len(fmt) > 1 {
				if prevFmt != fmt {
					prevFmt = fmt
					logger.Formatter = formatter(fmt)
				}
			}

			if l := strings.ToLower(os.Getenv(EnvLevel)); len(l) > 1 {
				if prevLevel != l {
					prevLevel = l
					logger.Level = level(l)
				}
			}

			//	no need to run constantly
			time.Sleep(800 * time.Millisecond)
		}
	}()
}

func output(key string) io.Writer {
	switch key {
	case "stdout":
		return os.Stdout

	default:
		return os.Stderr
	}
}

func formatter(key string) logrus.Formatter {
	switch key {
	case "json":
		return &logrus.JSONFormatter{}

	default:
		return &logrus.TextFormatter{}
	}
}

func level(key string) logrus.Level {
	switch key {
	case "debug":
		return logrus.DebugLevel

	case "warn":
		return logrus.WarnLevel

	case "fatal":
		return logrus.FatalLevel

	case "panic":
		return logrus.PanicLevel

	default:
		return logrus.InfoLevel
	}
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	logger.Println(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

func NewError(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Error(err error, message string) {
	logger.Error(errors.Wrap(err, message))
}

func Errorf(err error, format string, args ...interface{}) {
	logger.Error(errors.Wrapf(err, format, args...))
}

func Errorln(err error, message string) {
	logger.Errorln(errors.Wrap(err, message))
}

func Fatal(args ...interface{}) {
	logger.Panic(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Fatalln(args ...interface{}) {
	logger.Panicln(args...)
}
