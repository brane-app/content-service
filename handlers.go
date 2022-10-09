package main

import (
	"github.com/brane-app/content-service/handlers"
	"github.com/brane-app/content-service/middleware"

	tools_middleware "github.com/brane-app/librane/tools/middleware"
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

	groudon.AddMiddleware("POST", rootRoute, tools_middleware.MustAuth)
	groudon.AddMiddleware("POST", rootRoute, tools_middleware.RejectBanned)
	groudon.AddMiddleware("POST", rootRoute, tools_middleware.ParseMultipart)
	groudon.AddMiddleware("POST", rootRoute, middleware.TransformBase64)

	groudon.AddHandler("POST", rootRoute, handlers.PostContent)
}
