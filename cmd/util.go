package cmd

import (
	"crypto/tls"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
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

func retrieve(from string, insecure bool) (string, error) {
	client := http.DefaultClient
	if insecure {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	res, err := client.Get(from)
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

func send(to, body string, insecure bool) (string, error) {
	client := http.DefaultClient
	if insecure {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	res, err := client.Post(to, http.DetectContentType([]byte(body)), strings.NewReader(body))
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

func del(from string, insecure bool) (string, error) {
	req, err := http.NewRequest(http.MethodDelete, from, nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to create DELETE http request")
	}

	client := http.DefaultClient
	if insecure {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	res, err := client.Do(req)
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

func osUser() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve current OS user")
	}

	return u.Username, nil
}
