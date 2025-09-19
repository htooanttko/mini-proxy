package models

import "net/http"

type ProxiedRequest struct {
	Req  *http.Request
	Resp *http.Response
}
