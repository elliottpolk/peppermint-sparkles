package cmd

import (
	"bufio"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const tag string = "peppermint-sparkles.cmd"

var (
	schemeExp = regexp.MustCompile(`^(?P<scheme>http(s)?):\/\/`)

	ErrNoPipe       = errors.New("no piped input")
	ErrDataTooLarge = errors.New("data to large")

	MaxData = (int(math.Pow10(7)) * 3)
)

func asURL(addr, path, params string) string {
	scheme := "https"
	if schemeExp.MatchString(addr) {
		matches, res := schemeExp.FindStringSubmatch(addr), make(map[string]string)
		for i, n := range schemeExp.SubexpNames() {
			if i > 0 && i <= len(matches) {
				res[n] = matches[i]
			}
		}
		scheme, addr = res["scheme"], schemeExp.ReplaceAllString(addr, "")
	}

	return (&url.URL{
		Scheme:   scheme,
		Host:     addr,
		Path:     path,
		RawQuery: params,
	}).String()
}

func retrieve(from string) (string, error) {
	res, err := http.Get(from)
	if err != nil {
		return "", errors.Wrap(err, "unable to call secrets service")
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read secrets service response body")
	}

	if code := res.StatusCode; code != http.StatusOK {
		switch code {
		case http.StatusNotFound:
			return "", errors.New("no valid secret")

		default:
			return "", errors.Errorf("secrets service responded with status code %d and message %s", code, string(b))
		}
	}

	return string(b), nil
}

func send(to, body string) (string, error) {
	res, err := http.Post(to, http.DetectContentType([]byte(body)), strings.NewReader(body))
	if err != nil {
		return "", errors.Wrap(err, "unable to post secret to secrets service")
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable read secrets service response")
	}

	if code := res.StatusCode; code < 200 || code > 299 {
		return "", errors.Errorf("secrets service responded with status code %d and message %s", code, string(b))
	}

	return string(b), nil
}

func del(from string) (string, error) {
	req, err := http.NewRequest(http.MethodDelete, from, nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to create DELETE http request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "unable to perform DELETE request")
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read secrets service response")
	}

	if code := res.StatusCode; code < 200 || code > 299 {
		return "", errors.Errorf("secrets service responded with status code %d and message %s", code, string(b))
	}

	return string(b), nil
}

func pipe() (string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", errors.Wrap(err, "unable to stat stdin")
	}

	if fi.Mode()&os.ModeCharDevice != 0 || fi.Size() < 1 {
		return "", ErrNoPipe
	}

	buf, res := bufio.NewReader(os.Stdin), make([]byte, 0)
	for {
		in, _, err := buf.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		res = append(res, in...)

		if len(res) > MaxData {
			return "", ErrDataTooLarge
		}
	}

	return string(res), nil
}

func osUser() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve current OS user")
	}

	return u.Username, nil
}
