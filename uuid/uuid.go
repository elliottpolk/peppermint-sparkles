// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/main.go
//
package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
)

//  GetV4 returns a version 4 UUID based on RFC 4122 section 4.4 spec
func GetV4() string {
	buf := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, buf)
	if n != len(buf) || err != nil {
		fmt.Printf("unable to retrieve random data: %v\n", err)
		return ""
	}

	// required to identify version 4 (pseudo-random) of section 4.1.3
	buf[6] = buf[6]&^0xf0 | 0x40

	//  variant bits via section 4.1.1 of RFC 4122
	buf[8] = buf[8]&^0xc0 | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}
