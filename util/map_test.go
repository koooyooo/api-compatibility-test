package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveElmFromMap(t *testing.T) {
	m := map[string]interface{}{
		"hello": map[string]interface{}{
			"world": 5,
			"champ": 4,
			"baby":  3,
		},
		"foo": "bar",
	}
	RemoveElmFromMap(&m, []string{"/hello/world", "/hello/champ", "/foo"})
	// "/hello/world", "hello/champ", "/foo" が消えたMapとなっていることを確認
	assert.Equal(t, map[string]interface{}{
		"hello": map[string]interface{}{
			"baby": 3,
		},
	}, m)
}
