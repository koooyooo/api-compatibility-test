package executor

import (
	"bytes"
	"net/http"

	"github.com/koooyooo/api-compatibility-test/model"
)

func CallAPIs(r1, r2 *model.Request) (*model.Responses, error) {
	res1, err := callAPI(r1.Method, r1.Url, r1.Header, r1.Body)
	if err != nil {
		return nil, err
	}
	res2, err := callAPI(r1.Method, r1.Url, r1.Header, r1.Body)
	if err != nil {
		return nil, err
	}
	return &model.Responses{res1, res2}, nil
}

func callAPI(method, url string, header *http.Header, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header = *header
	if err != nil {
		return nil, err
	}
	cli := &http.Client{}
	resp, err := cli.Do(req)
	return resp, err
}
