package main

import (
	"github.com/brane-app/tools-library/library"
	"github.com/gastrodon/groudon/v2"

	"os"
)

var (
	prefix = os.Getenv("PATH_PREFIX")

	rootRoute = "^" + prefix + "/?$"

	bad_auth = map[string]interface{}{"error": "bad_auth"}
)

func register_handlers() {
	groudon.AddCodeResponse(401, bad_auth)

	groudon.AddMiddleware("POST", rootRoute, middleware.MustAuth)
	groudon.AddMiddleware("POST", rootRoute, middleware.RejectBanned)
	groudon.AddMiddleware("POST", rootRoute, middleware.ParseMultipart)
	groudon.AddMiddleware("POST", rootRoute, transformBase64)

	groudon.AddHandler("POST", rootRoute, postContent)
}
