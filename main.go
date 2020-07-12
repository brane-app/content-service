package main

import (
	"github.com/gastrodon/groudon"
	"github.com/imonke/monkebase"

	"log"
	"net/http"
	"os"
)

var (
	bad_auth     = map[string]interface{}{"error": "bad_auth"}
	expired_auth = map[string]interface{}{"error": "expired_auth"}
)

func main() {
	monkebase.Connect(os.Getenv("MONKEBASE_CONNECTION"))
	groudon.RegisterCatch(401, bad_auth)
	groudon.RegisterMiddleware(rejectBadAuth)
	groudon.RegisterMiddleware(parseMultipart)
	groudon.RegisterHandler("POST", "^/$", postContent)
	http.Handle("/", http.HandlerFunc(groudon.Route))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
