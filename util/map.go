package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func CreateBytesAndMapFromJSONBody(r *http.Response) ([]byte, map[string]interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	// 読みきったBody部分を復元
	r.Body = ioutil.NopCloser(bytes.NewReader(b))
	bMap := make(map[string]interface{})
	if err := json.Unmarshal(b, &bMap); err != nil {
		return nil, nil, err
	}
	return b, bMap, nil
}

// Map内からSkip要素を削除
func RemoveElmFromMap(m *map[string]interface{}, paths []string) {
	for _, path := range paths {
		currentMap := *m
		path = strings.TrimLeft(path, "/")
		pathElms := strings.Split(path, "/")
		for i, elm := range pathElms {
			if i == len(pathElms)-1 {
				delete(currentMap, elm)
				break
			}
			currentMap = currentMap[elm].(map[string]interface{})
		}
	}
}
