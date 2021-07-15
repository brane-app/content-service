package main

import (
	"github.com/brane-app/database-library"
	"github.com/brane-app/types-library"

	"net/http"
	"strings"
)

func pathSplit(it rune) (ok bool) {
	ok = it == '/'
	return
}

func getContent(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var parts []string = strings.FieldsFunc(request.URL.Path, pathSplit)
	var id string = parts[len(parts)-1]

	var fetched types.Content
	var exists bool
	if fetched, exists, err = database.ReadSingleContent(id); err != nil {
		return
	}

	if !exists {
		code = 404
		r_map = map[string]interface{}{"error": "no_such_content"}
		return
	}

	code = 200
	r_map = map[string]interface{}{
		"content": fetched.Map(),
	}
	return
}
