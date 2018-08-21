


//

package backend

import "testing"

func TestKey(t *testing.T) {
	wants := map[string][][]string{
		"97df3588b5a3f24babc3851b372f0ba71a9dcdded43b14b9d06961bfc1707d9d34b379d1190316e4d5f8492140a8781d7ee37caa3cc7521bb58d49cf3b659e0e": {{"foo", "bar", "baz"}},
		"ebf3b019bb7e36bdc0fbc4159345c04af54193ccf43ae4572922f6d4aa94bd5b34b379d1190316e4d5f8492140a8781d7ee37caa3cc7521bb58d49cf3b659e0e": {{"foo", "bar", "baz", "biz"}},
		"2c60dbf3773104dce76dfbda9b82a729e98a42a7a0b3f9bae5095c7bed752b9034b379d1190316e4d5f8492140a8781d7ee37caa3cc7521bb58d49cf3b659e0e": {
			{"foo", "bar", "bazz"},
			{"foo", "bar", "baz", "z"},
		},
		"796362b8b4289fca4d666ab486487d6699e828f9c098fc1c91566c291ef682f634b379d1190316e4d5f8492140a8781d7ee37caa3cc7521bb58d49cf3b659e0e": {
			{"foo", "bar", "baz z"},
			{"foo", "bar", "baz", " z"},
		},
	}

	app, env := "test", "testing"
	for want, vals := range wants {
		for _, val := range vals {
			if got := Key(app, env, val...); got != want {
				t.Errorf("\nwant %s\ngot %s\nwith %+v", want, got, val)
			}
		}
	}
}
