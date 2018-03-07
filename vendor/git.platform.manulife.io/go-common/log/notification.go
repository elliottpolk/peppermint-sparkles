// Created by Elliott Polk on 17/02/2017
// Copyright © 2017 Manulife AM. All rights reserved.
// go-common/log/notification.go
//
package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	ErrorNotification string = "error"
	InfoNotification  string = "info"
)

const (
	DefaultAddr    string = ":8080"
	DefaultService string = "log"

	EnvAddr    string = "NOTIFY_API"
	EnvService string = "SERVICE_NAME"

	PathNotify string = "/v1/api/notify/add"

	//  TODO ... convert to using https when certs are available
	protocol string = "http"
)

type Notification struct {
	Message string `json:"message"`
	Service string `json:"service"`
	Type    string `json:"type"`
	Created int64  `json:"created"`
}

func NofityAddr() string {
	if v := os.Getenv(EnvAddr); v != "" {
		v = strings.TrimPrefix(v, "http://")
		v = strings.TrimPrefix(v, "https://")
		v = strings.TrimSuffix(v, "/")

		return fmt.Sprintf("%s://%s%s", protocol, v, PathNotify)
	}

	return fmt.Sprintf("%s://%s%s", protocol, DefaultAddr, PathNotify)
}

func ServiceName() string {
	if v := os.Getenv(EnvService); v != "" {
		return v
	}

	return DefaultService
}

func (n *Notification) Post(where string) error {
	out, err := json.Marshal(n)
	if err != nil {
		return errors.Wrap(err, "unable to marshal notification for post")
	}

	res, err := http.Post(where, "application/json", bytes.NewReader(out))
	if err != nil {
		return errors.Wrap(err, "unable to post notification to notify service")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "unable to read notify service response")
		}

		return errors.Errorf("notify service responded with status code '%d' and message '%s'", res.StatusCode, string(b))
	}

	return nil
}

func NotifyInfo(args ...interface{}) {
	Info(args...)
	notify(fmt.Sprint(args...), InfoNotification)
}

func NotifyInfof(format string, args ...interface{}) {
	Infof(format, args...)
	notify(fmt.Sprintf(format, args...), InfoNotification)
}

func NotifyInfoln(args ...interface{}) {
	Infoln(args...)
	notify(fmt.Sprint(args...), InfoNotification)
}

func NotifyError(err error, args ...interface{}) {
	Error(err, "")
	notify(fmt.Sprint(args...), ErrorNotification)
}

func NotifyErrorf(err error, format string, args ...interface{}) {
	Error(err, "")
	notify(fmt.Sprintf(format, args...), ErrorNotification)
}

func NotifyErrorln(err error, args ...interface{}) {
	Error(err, "")
	notify(fmt.Sprint(args...), ErrorNotification)
}

func notify(msg, t string) {
	service := ServiceName()

	n := &Notification{msg, service, t, time.Now().Unix()}
	if err := n.Post(NofityAddr()); err != nil {
		Errorf("%v: unable to post notification for service %s", err, service)
	}
}
