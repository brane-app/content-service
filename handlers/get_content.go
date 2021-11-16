package handlers

import (
	"github.com/brane-app/database-library"
	"github.com/brane-app/tools-library"
	"github.com/brane-app/types-library"

	"net/http"
)

var (
	err_no_such_content = map[string]interface{}{"error": "no_such_content"}
)

func GetContent(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var parts []string = tools.SplitPath(request.URL.Path)
	var id string = parts[len(parts)-1]

	var fetched types.Content
	var exists bool
	if fetched, exists, err = database.ReadSingleContent(id); err != nil {
		return
	}

	if !exists {
		code = 404
		r_map = err_no_such_content
		return
	}

	code = 200
	r_map = map[string]interface{}{
		"content": fetched.Map(),
	}
	return
}
