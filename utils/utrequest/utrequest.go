package utrequest

import (
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"logging-example-go/serror"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type RequestStruct struct {
	Client    *http.Client
	CookieJar *cookiejar.Jar
}

type RequestOption struct {
	Req     *http.Request
	Headers []string
	Proxy   map[string]string
}

func Construct() (*RequestStruct, serror.SError) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, serror.NewFromError(err)
	}

	request := &RequestStruct{
		CookieJar: cookieJar,
	}

	request.Client = &http.Client{
		Jar:     request.CookieJar,
		Timeout: time.Duration(1 * time.Hour),
	}

	return request, nil
}

func (request *RequestStruct) BasicRequest(input RequestOption) (*http.Response, []byte, serror.SError) {
	if len(input.Headers) > 0 {
		for _, v := range input.Headers {
			input.ApplyHeader(v)
		}
	}

	if len(input.Proxy) > 0 {
		proxyURL, err := url.Parse(input.Proxy["url"])
		if err != nil {
			return nil, nil, serror.NewFromError(err)
		}

		request.Client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

	} else {
		request.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}
	}

	resp, err := request.Client.Do(input.Req)
	if err != nil {
		return nil, nil, serror.NewFromError(err)
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, nil, serror.NewFromErrorc(err, "Failed to create reader")
		}

		defer reader.Close()

	default:
		reader = resp.Body
	}

	bbyte, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, nil, serror.NewFromError(err)
	}

	return resp, bbyte, nil
}

func (ox *RequestOption) ApplyHeader(s string) {
	spt := strings.Split(s, ":")
	if len(spt) >= 2 {
		f := spt[0]
		spt = spt[1:]
		ox.Req.Header.Set(f, strings.Trim(strings.Join(spt[:], ":"), " "))
	}
}

func ApplyHeader(c *http.Request, s string) {
	spt := strings.Split(s, ":")
	if len(spt) >= 2 {
		f := spt[0]
		spt = spt[1:]
		c.Header.Set(f, strings.Trim(strings.Join(spt[:], ":"), " "))
	}
}
