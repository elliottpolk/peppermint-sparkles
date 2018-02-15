// Created by Elliott Polk on 29/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/backend/unit_test.go
//

package backend

import "testing"

func TestKey(t *testing.T) {
	wants := map[string][][]string{
		"97df3588b5a3f24babc3851b372f0ba71a9dcdded43b14b9d06961bfc1707d9d": {{"foo", "bar", "baz"}},
		"ebf3b019bb7e36bdc0fbc4159345c04af54193ccf43ae4572922f6d4aa94bd5b": {{"foo", "bar", "baz", "biz"}},
		"2c60dbf3773104dce76dfbda9b82a729e98a42a7a0b3f9bae5095c7bed752b90": {
			{"foo", "bar", "bazz"},
			{"foo", "bar", "baz", "z"},
		},
		"796362b8b4289fca4d666ab486487d6699e828f9c098fc1c91566c291ef682f6": {
			{"foo", "bar", "baz z"},
			{"foo", "bar", "baz", " z"},
		},
	}

	for want, vals := range wants {
		for _, val := range vals {
			if got := Key(val...); got != want {
				t.Errorf("\nwant %s\ngot %s\nwith %+v", want, got, val)
			}
		}
	}
}
