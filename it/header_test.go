package it

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koooyooo/api-compatibility-test/executor"
	"github.com/koooyooo/api-compatibility-test/model"
	"github.com/stretchr/testify/assert"
)

// Header値が別の場合 Headerの比較で失敗する
func TestDifferentHeader(t *testing.T) {
	// Headerで別の結果を返すテストサーバを作成
	var sampleHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/path1":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Num", "One")
			w.WriteHeader(200)
			w.Write([]byte(`{"headers": {"hello": "world"}, "body": {"foo": "bar"}}`))
		case "/path2":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Num", "Two")
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

	// Headerの比較で Equalにならない
	headerPair, err := resps.HeaderPair(nil)
	assert.NoError(t, err)
	assert.NotEqual(t, headerPair.First, headerPair.Second)

	bodyPair, err := resps.BodyPair(nil)
	assert.NoError(t, err)
	assert.Equal(t, bodyPair.First, bodyPair.Second)
}

// Header値が別の場合でも 当該属性がSkip指定していればテストは成功する
func TestDifferentHeaderWithSkip(t *testing.T) {
	// Headerで別の結果を返すテストサーバを作成
	var sampleHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/path1":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Num", "One")
			w.WriteHeader(200)
			w.Write([]byte(`{"headers": {"hello": "world"}, "body": {"foo": "bar"}}`))
		case "/path2":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Num", "Two")
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

	// 差分をSkipすればHeaderの比較で Equalになる
	headerPair, err := resps.HeaderPair([]string{"Num"})
	assert.NoError(t, err)
	assert.Equal(t, headerPair.First, headerPair.Second)

	bodyPair, err := resps.BodyPair(nil)
	assert.NoError(t, err)
	assert.Equal(t, bodyPair.First, bodyPair.Second)
}
