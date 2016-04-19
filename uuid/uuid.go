// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package uuid

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/golang/glog"
)

//  GetV4 returns a version 4 UUID based on RFC 4122 section 4.4 spec
func GetV4() string {
	defer glog.Flush()

	buf := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, buf)
	if n != len(buf) || err != nil {
		glog.Errorf("unable to retrieve random data: %v\n", err)

		return ""
	}

	// required to identify version 4 (pseudo-random) of section 4.1.3
	buf[6] = buf[6]&^0xf0 | 0x40

	//  variant bits via section 4.1.1 of RFC 4122
	buf[8] = buf[8]&^0xc0 | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}
