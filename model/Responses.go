package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/koooyooo/api-compatibility-test/util"

	"github.com/stretchr/testify/assert"
)

type Responses struct {
	Res1 *http.Response
	Res2 *http.Response
}

func (r Responses) AssertAll(t *testing.T, skipHeaders, skipBodyPath []string) {
	r.AssertStatus(t)
	r.AssertHeader(t, skipHeaders)
	r.AssertBody(t, skipBodyPath)
}

func (r Responses) AssertStatus(t *testing.T) {
	assert.Equal(t, r.Res1.Status, r.Res2.Status)
}

func (r Responses) AssertHeader(t *testing.T, skipKeys []string) {
	h1 := r.Res1.Header
	h2 := r.Res2.Header
	// ヘッダ数を比較
	assert.Equal(t, len(h1), len(h2))

HEADERS:
	for hk1, _ := range h1 {
		hv1 := h1[hk1]
		hv2 := h2[hk1]
		// ヘッダ毎のValue数を比較
		assert.Equal(t, len(hv1), len(hv2))
		for _, sk := range skipKeys {
			if hk1 == sk {
				continue HEADERS
			}
		}
		for i, _ := range hv1 {
			v1 := hv1[i]
			v2 := hv2[i]
			if hk1 == "X-Request-Id" {
				assert.NotEmpty(t, v1)
				assert.NotEmpty(t, v2)
				continue
			}
			// Value値を比較
			assert.Equal(t, v1, v2, "Header Key:"+hk1)
		}
	}
}

func (r Responses) AssertBody(t *testing.T, skipPaths []string) error {
	switch r.Res1.Header["Content-Type"][0] {
	case "application/json":
		b1, err := ioutil.ReadAll(r.Res1.Body)
		if err != nil {
			return err
		}
		bMap1 := make(map[string]interface{})
		if err := json.Unmarshal(b1, &bMap1); err != nil {
			return err
		}

		b2, err := ioutil.ReadAll(r.Res2.Body)
		assert.Nil(t, err)
		bMap2 := make(map[string]interface{})
		if err := json.Unmarshal(b2, &bMap2); err != nil {
			return err
		}

		util.RemoveElmFromMap(&bMap1, skipPaths)
		util.RemoveElmFromMap(&bMap2, skipPaths)

		eq := reflect.DeepEqual(bMap1, bMap2)
		assert.True(t, eq, string(b1)+"\n Not Equal \n"+string(b2))

	default:
		b1, err := ioutil.ReadAll(r.Res1.Body)
		if err != nil {
			return err
		}
		b2, err := ioutil.ReadAll(r.Res2.Body)
		if err != nil {
			return err
		}
		assert.Equal(t, string(b1), string(b2))
	}
	return nil
}
