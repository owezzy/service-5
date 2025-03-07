package checkgrp

import (
	"github.com/owezzy/service-5/foundation/web"
	"net/http"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	hdl := New(cfg.Build)
	app.HandleNoMiddleware(http.MethodGet, version, "/readiness", hdl.Readiness)
	app.HandleNoMiddleware(http.MethodGet, version, "/liveness", hdl.Liveness)

}
