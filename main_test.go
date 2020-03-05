package main

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/koooyooo/api-compatibility-test/model"
)

func TestXxxxx(t *testing.T) {
	g1Res, _ := CallAPI("GET", "https://httpbin.org/get", &http.Header{
		"HELLO": []string{"WORLD"},
	}, "")
	g2Res, _ := CallAPI("GET", "https://httpbin.org/get", &http.Header{}, "")
	model.Responses{
		g1Res,
		g2Res,
	}.AssertAll(t, []string{"Date", "Content-Length"}, []string{"headers/X-Amzn-Trace-Id", "headers/Hello"}) // TODO HEaderは Camelで指定
}

func CallAPI(method, url string, header *http.Header, body string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	req.Header = *header
	if err != nil {
		return nil, err
	}
	cli := &http.Client{}
	resp, err := cli.Do(req)
	return resp, err
}
