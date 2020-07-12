package main

import (
	"github.com/imonke/monkebase"

	"context"
	"net/http"
	"strings"
)

const (
	BEARER_PREFIX     = "Bearer "
	MULTIPART_MEM_MAX = 4 << 20
)

func rejectBadAuth(request *http.Request) (modified *http.Request, ok bool, code int, r_map map[string]interface{}, err error) {
	code = 401
	var bearer string = strings.TrimPrefix(request.Header.Get("Authorization"), BEARER_PREFIX)

	var owner string
	if owner, ok, err = monkebase.ReadTokenStat(bearer); err != nil || !ok {
		r_map = bad_auth
		return
	}

	modified = request.WithContext(context.WithValue(
		request.Context(),
		"requester",
		owner,
	))

	return
}

func parseMultipart(request *http.Request) (_ *http.Request, ok bool, code int, r_map map[string]interface{}, err error) {
	if err = request.ParseMultipartForm(MULTIPART_MEM_MAX); err != nil {
		err = nil
		code = 400
		return
	}

	ok = true
	return
}
