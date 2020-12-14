package main

import (
	"github.com/gastrodon/groudon"
	"github.com/google/uuid"

	"net/http"
	"os"
	"strings"
)

var (
	auth = os.Getenv("FERROTHORN_SECRET")
)

func splitIgnoreEmpty(it rune) (ok bool) {
	ok = it == '/'
	return
}

func mustAuth(request *http.Request) (_ *http.Request, ok bool, code int, r_map map[string]interface{}, err error) {
	ok = auth == "" || request.Header.Get("Authorization") == auth
	code = 401
	return
}

func upload(name string) (code int, r_map map[string]interface{}, _ error) {
	code = 200
	r_map = map[string]interface{}{"id": name}
	return
}

func postID(_ *http.Request) (code int, r_map map[string]interface{}, _ error) {
	code, r_map, _ = upload(uuid.New().String())
	return
}

func postName(request *http.Request) (code int, r_map map[string]interface{}, _ error) {
	code, r_map, _ = upload(strings.FieldsFunc(request.URL.Path, splitIgnoreEmpty)[0])
	return
}

func get(_ *http.Request) (code int, r_map map[string]interface{}, _ error) {
	code = 200
	r_map = map[string]interface{}{"file": "pretend this is a file"}
	return
}

func delete(_ *http.Request) (_ int, _ map[string]interface{}, _ error) {
	return
}

func main() {
	groudon.RegisterMiddlewareRoute([]string{"POST", "DELETE"}, ".*", mustAuth)
	groudon.RegisterHandler("POST", `^/$`, postID)
	groudon.RegisterHandler("POST", `^/[a-zA-Z0-9\-\.]+/?$`, postName)
	groudon.RegisterHandler("GET", `^/[a-zA-Z0-9\-\.]+/?$`, get)
	groudon.RegisterHandler("DELETE", `^/[a-zA-Z0-9\-\.]+/?$`, delete)
	http.Handle("/", http.HandlerFunc(groudon.Route))
	http.ListenAndServe(":4000", nil)
}
