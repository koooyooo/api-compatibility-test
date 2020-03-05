package executor

import (
	"bytes"
	"net/http"

	"github.com/koooyooo/api-compatibility-test/model"
)

func CallAPIs(r1, r2 *model.Request) (*model.Responses, error) {
	call := func(r *model.Request) (*http.Response, error) {
		req, err := http.NewRequest(r.Method, r.Url, bytes.NewBuffer(r.Body))
		if err != nil {
			return nil, err
		}
		if r.Header != nil {
			req.Header = *r.Header
		}
		return callRawAPI(req)
	}
	res1, err := call(r1)
	if err != nil {
		return nil, err
	}
	res2, err := call(r2)
	if err != nil {
		return nil, err
	}
	return &model.Responses{res1, res2}, nil
}

func callRawAPI(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
