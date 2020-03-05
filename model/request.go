package model

import (
	"net/http"
)

type Request struct {
	Method string
	Url    string
	Header *http.Header
	Body   []byte
}
