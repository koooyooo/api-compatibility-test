package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/koooyooo/api-compatibility-test/util"

	"github.com/stretchr/testify/assert"
)

type Responses struct {
	Res1 *http.Response
	Res2 *http.Response
}

type IntPair struct {
	First  int
	Second int
}

type StringPair struct {
	First  string
	Second string
}

func (r Responses) StatusPair() *IntPair {
	return &IntPair{r.Res1.StatusCode, r.Res2.StatusCode}
}

func (r Responses) AssertStatus(t *testing.T) {
	p := r.StatusPair()
	assert.Equal(t, p.First, p.Second)
}

func (r Responses) HeaderPair(skipKeys []string) (*StringPair, error) {
	stringfy := func(h http.Header, skipKeys []string) (string, error) {
		for _, key := range skipKeys {
			h.Del(key)
		}
		headerB, err := json.Marshal(h)
		if err != nil {
			return "", err
		}
		return string(headerB), nil
	}
	h1Str, err := stringfy(r.Res1.Header, skipKeys)
	if err != nil {
		return nil, err
	}
	h2Str, err := stringfy(r.Res2.Header, skipKeys)
	if err != nil {
		return nil, err
	}
	return &StringPair{h1Str, h2Str}, nil
}

func (r Responses) AssertHeader(t *testing.T, skipKeys []string) {
	p, err := r.HeaderPair(skipKeys)
	assert.NoError(t, err)
	assert.Equal(t, p.First, p.Second)
}

func (r *Responses) BodyPair(skipBodyPaths []string) (*StringPair, error) {
	switch r.Res1.Header["Content-Type"][0] {
	case "application/json":
		// レスポンスよりMap型でJSONを取得
		_, bMap1, err := util.CreateBytesAndMapFromJSONBody(r.Res1)
		if err != nil {
			return nil, err
		}
		_, bMap2, err := util.CreateBytesAndMapFromJSONBody(r.Res2)
		if err != nil {
			return nil, err
		}

		// スキップ属性を除去し[]byte型に復元
		util.RemoveElmFromMap(&bMap1, skipBodyPaths)
		util.RemoveElmFromMap(&bMap2, skipBodyPaths)

		removedB1, err := json.Marshal(bMap1)
		removedB2, err := json.Marshal(bMap2)

		return &StringPair{
			string(removedB1),
			string(removedB2),
		}, nil

	default:
		b1, err := ioutil.ReadAll(r.Res1.Body)
		if err != nil {
			return nil, err
		}
		b2, err := ioutil.ReadAll(r.Res2.Body)
		if err != nil {
			return nil, err
		}
		return &StringPair{
			string(b1),
			string(b2),
		}, nil
	}
}

func (r *Responses) AssertBody(t *testing.T, skipBodyPaths []string) {
	p, err := r.BodyPair(skipBodyPaths)
	assert.NoError(t, err)
	assert.Equal(t, p.First, p.Second)

	//switch r.Res1.Header["Content-Type"][0] {
	//case "application/json":
	//	// レスポンスを読み取り生JSONと Map形式JSONを取得
	//	bJSON1, bMap1, err := util.CreateBytesAndMapFromJSONBody(r.Res1)
	//	assert.NoError(t, err)
	//	bJSON2, bMap2, err := util.CreateBytesAndMapFromJSONBody(r.Res2)
	//	assert.NoError(t, err)
	//
	//	// スキップ属性を除去
	//	util.RemoveElmFromMap(&bMap1, skipBodyPaths)
	//	util.RemoveElmFromMap(&bMap2, skipBodyPaths)
	//
	//	// 階層的に確認
	//	eq := reflect.DeepEqual(bMap1, bMap2)
	//	assert.True(t, eq, "### Body1:\n"+string(bJSON1)+"\n### Body2:\n"+string(bJSON2)+"")
	//default:
	//	b1, err := ioutil.ReadAll(r.Res1.Body)
	//	assert.NoError(t, err, "fails in reading response body 1")
	//	b2, err := ioutil.ReadAll(r.Res2.Body)
	//	assert.NoError(t, err, "fails in reading response body 2")
	//	assert.Equal(t, string(b1), string(b2))
	//}
}
