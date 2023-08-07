package main

import (
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type AdcRequest http.Request

func panic_if_error(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (req *AdcRequest) Prepare(rawurl string) *AdcRequest {
	u, err := url.Parse(rawurl)
	panic_if_error(err)

	r, err := http.NewRequest("POST", rawurl, nil)
	panic_if_error(err)

	r.Header.Add("Accept", `text/html`)
	r.Header.Add("Accept-Language", `it-IT,it;q=0.9`)
	r.Header.Add("Cache-Control", `max-age=0`)
	r.Header.Add("Connection", `keep-alive`)
	r.Header.Add("Content-Type", `application/x-www-form-urlencoded`)
	r.Header.Add("DNT", `1`)
	r.Header.Add("Origin", u.Hostname())
	r.Header.Add("Referer", rawurl)
	r.Header.Add("Upgrade-Insecure-Request", `1`)
	r.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36`)

	return req
}

func PrepareBody(s, bodyfile string) io.Reader {
	b, err := ioutil.ReadFile(bodyfile)
	panic_if_error(err)

	b = append(b, s...)
	ioutil.WriteFile("temp.txt", b, fs.ModeAppend)

	f, err := os.Open("temp.txt")
	panic_if_error(err)

	return f
}
