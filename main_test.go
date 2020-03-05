package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koooyooo/api-compatibility-test/executor"

	"github.com/stretchr/testify/assert"

	"github.com/koooyooo/api-compatibility-test/model"
)

// HTTPBinを利用した利用法のサンプル
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

// 同じ実装の場合テストが成功する
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

// 同改装のBodyJSONの順序が異なってもテストが成功する
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

// Headerが別実装の場合 Headerの比較で失敗する
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

// Headerが別実装の場合 Headerの比較で失敗する
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
