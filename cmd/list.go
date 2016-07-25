// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const List = "list"

func ListCfgs() error {
	addr := GetConfgrAddr()

	res, err := http.Get(fmt.Sprintf("%s/list", addr))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
