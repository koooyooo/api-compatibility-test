package main

import (
	"net/http"
	"testing"

	"github.com/koooyooo/api-compatibility-test/executor"

	"github.com/stretchr/testify/assert"

	"github.com/koooyooo/api-compatibility-test/model"
)

func TestHttpBin(t *testing.T) {
	// リクエスト生成
	req1 := &model.Request{"GET", "https://httpbin.org/get", &http.Header{"HELLO": []string{"WORLD"}}, []byte{}}
	req2 := &model.Request{"GET", "https://httpbin.org/get", nil, []byte{}}

	// クライアントコール
	resps, err := executor.CallAPIs(req1, req2)
	assert.NoError(t, err, "fails in client call")

	// 結果確認
	resps.AssertStatus(t)
	resps.AssertHeader(t, []string{"Date", "Content-Length"})
	resps.AssertBody(t, []string{"headers/X-Amzn-Trace-Id", "headers/Hello", "origin"})
}
