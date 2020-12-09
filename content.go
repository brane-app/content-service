package main

import (
	"github.com/imonke/monkebase"
	"github.com/imonke/monketype"

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

	var fetched monketype.Content
	var exists bool
	if fetched, exists, err = monkebase.ReadSingleContent(id); err != nil {
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
