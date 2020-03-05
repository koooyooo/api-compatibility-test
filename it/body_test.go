package it

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koooyooo/api-compatibility-test/executor"
	"github.com/koooyooo/api-compatibility-test/model"
	"github.com/stretchr/testify/assert"
)

// Bodyが別実装の場合 Bodyの比較で失敗する
func TestDifferentBody(t *testing.T) {
	// Body部分で別の結果を返すテストサーバを作成
	var sampleHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/path1":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"headers": {"hello": "world1"}, "body": {"foo": "bar"}}`))
		case "/path2":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"headers": {"hello": "world2"}, "body": {"foo": "bar"}}`))
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

	// Bodyの比較で Equalにならない
	bodyPair, err := resps.BodyPair(nil)
	assert.NoError(t, err)
	assert.NotEqual(t, bodyPair.First, bodyPair.Second)
}

// Bodyが別実装の場合 でも該当属性を Skipしていれば引っかからない
func TestDifferentBodyWithSkip(t *testing.T) {
	// Body部分で別の結果を返すテストサーバを作成
	var sampleHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/path1":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"headers": {"hello": "world1"}, "body": {"foo": "bar"}}`))
		case "/path2":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"headers": {"hello": "world2"}, "body": {"foo": "bar"}}`))
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

	// 該当属性の Skipをしていれば Bodyの比較で Equalにならない
	bodyPair, err := resps.BodyPair([]string{"headers/hello"})
	assert.NoError(t, err)
	assert.Equal(t, bodyPair.First, bodyPair.Second)
}
