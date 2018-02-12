// Created by Elliott Polk on 17/02/2017
// Copyright © 2017 Manulife AM. All rights reserved.
// go-common/log/log.go
//
package log

import (
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
		prevFmt, prevLevel := "", ""
		for {
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
	logger.Out = os.Stdout
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Out = os.Stdout
	logger.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	logger.Out = os.Stdout
	logger.Println(args...)
}

func Debug(args ...interface{}) {
	logger.Out = os.Stdout
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Out = os.Stdout
	logger.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	logger.Out = os.Stdout
	logger.Debugln(args...)
}

func NewError(format string, args ...interface{}) {
	logger.Out = os.Stderr
	logger.Errorf(format, args...)
}

func Error(err error, message string) {
	logger.Out = os.Stderr
	logger.Error(errors.Wrap(err, message))
}

func Errorf(format string, args ...interface{}) {
	logger.Out = os.Stderr
	logger.Errorf(format, args...)
}

func Errorln(err error, message string) {
	logger.Out = os.Stderr
	logger.Errorln(errors.Wrap(err, message))
}

func Fatal(args ...interface{}) {
	logger.Out = os.Stderr
	logger.Panic(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Out = os.Stderr
	logger.Panicf(format, args...)
}

func Fatalln(args ...interface{}) {
	logger.Out = os.Stderr
	logger.Panicln(args...)
}
