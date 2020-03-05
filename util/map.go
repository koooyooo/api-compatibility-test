package util

import (
	"fmt"
	"strings"
)

// Map内からSkip要素を削除
func RemoveElmFromMap(m *map[string]interface{}, paths []string) {
	for _, path := range paths {
		currentMap := *m
		path = strings.TrimLeft(path, "/")
		pathElms := strings.Split(path, "/")
		for i, elm := range pathElms {
			if i == len(pathElms)-1 {
				fmt.Println(elm) // TODO
				delete(currentMap, elm)
				break
			}
			currentMap = currentMap[elm].(map[string]interface{})
		}
	}
}
