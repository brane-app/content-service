package main

import (
	"github.com/gastrodon/groudon/v2"

	"os"
)

const (
	uuid_regex = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
)

var (
	prefix = os.Getenv("PATH_PREFIX")

	routeContent = "^" + prefix + "/" + uuid_regex + "/?$"
)

func register_handlers() {
	groudon.AddHandler("GET", routeContent, getContent)
}
