package it

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koooyooo/api-compatibility-test/executor"
	"github.com/koooyooo/api-compatibility-test/model"
	"github.com/stretchr/testify/assert"
)

// 全く同じレスポンスの場合テストが成功する
func TestBasic(t *testing.T) {
	// 同じ結果を返すテストサーバを作成
	var sampleHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"headers": {"hello": "world"}, "body": {"foo": "bar"}}`))
		}
	}
	server := httptest.NewServer(sampleHandler)
	defer server.Close()

	req1 := &model.Request{"GET", server.URL + "/path1", nil, []byte{}}
	req2 := &model.Request{"GET", server.URL + "/path2", nil, []byte{}}

	// クライアントコール
	resps, err := executor.CallAPIs(req1, req2)
	assert.NoError(t, err, "fails in client call")

	// 結果確認
	statusPair := resps.StatusPair()
	assert.Equal(t, statusPair.First, statusPair.Second)

	headerPair, err := resps.HeaderPair(nil)
	assert.NoError(t, err)
	assert.Equal(t, headerPair.First, headerPair.Second)

	bodyPair, err := resps.BodyPair(nil)
	assert.NoError(t, err)
	assert.Equal(t, bodyPair.First, bodyPair.Second)
}

// レスポンスのBodyJSON内が異なっていても、同階層の順序入れ替えのみであればテストが成功する
func TestBasicWithDifferentJSONOrder(t *testing.T) {
	// JSONの同列要素を入れ替えたBodyを用意
	var sampleHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/path1":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"headers": {"hello": "world"}, "body": {"foo": "bar"}}`))
		case "/path2":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"body": {"foo": "bar"}, "headers": {"hello": "world"}}`))
		}
	}
	server := httptest.NewServer(sampleHandler)
	defer server.Close()

	req1 := &model.Request{"GET", server.URL + "/path1", nil, []byte{}}
	req2 := &model.Request{"GET", server.URL + "/path2", nil, []byte{}}

	// クライアントコール
	resps, err := executor.CallAPIs(req1, req2)
	assert.NoError(t, err, "fails in client call")

	// 結果確認
	statusPair := resps.StatusPair()
	assert.Equal(t, statusPair.First, statusPair.Second)

	headerPair, err := resps.HeaderPair(nil)
	assert.NoError(t, err)
	assert.Equal(t, headerPair.First, headerPair.Second)

	bodyPair, err := resps.BodyPair(nil)
	assert.NoError(t, err)
	assert.Equal(t, bodyPair.First, bodyPair.Second)
}
