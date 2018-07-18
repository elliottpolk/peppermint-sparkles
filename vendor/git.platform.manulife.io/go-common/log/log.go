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

var (
	logger  *logrus.Logger
	version string = "0.0.0"
)

func Init(ver string) {
	version = ver
	logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: formatter(""),
		Level:     level(""),
	}

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
	case "text":
		return &logrus.TextFormatter{}

	default:
		return &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		}
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

func Info(tag string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Info(args...)
}

func Infof(tag string, format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Infof(format, args...)
}

func Infoln(tag string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Println(args...)
}

func Debug(tag string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Debug(args...)
}

func Debugf(tag string, format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Debugf(format, args...)
}

func Debugln(tag string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Debugln(args...)
}

func NewError(tag string, format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Errorf(format, args...)
}

func Error(tag string, err error, message string) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Error(errors.Wrap(err, message))
}

func Errorf(tag string, err error, format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Errorf("%v: "+format, []interface{}{err, args}...)
}

func Errorln(tag string, err error, message string) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Errorln(errors.Wrap(err, message))
}

func Fatal(tag string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Panic(args...)
}

func Fatalf(tag string, format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Panicf(format, args...)
}

func Fatalln(tag string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"@version":    version,
		"logger_name": tag,
	}).Panicln(args...)
}
