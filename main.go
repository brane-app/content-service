package main

import (
	"github.com/gastrodon/groudon"
	"github.com/imonke/monkebase"

	"log"
	"net/http"
	"os"
)

const (
	uuid_regex = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
)

func main() {
	monkebase.Connect(os.Getenv("MONKEBASE_CONNECTION"))
	groudon.RegisterHandler("GET", "^/"+uuid_regex+"/?$", getContent)
	http.Handle("/", http.HandlerFunc(groudon.Route))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
