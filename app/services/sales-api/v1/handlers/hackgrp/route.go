package hackgrp

import (
	"github.com/owezzy/service-5/business/web/v1/auth"
	"github.com/owezzy/service-5/business/web/v1/mid"
	"github.com/owezzy/service-5/foundation/web"
	"net/http"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Auth *auth.Auth
	Build string
}

func Routes(app *web.App, cfg Config) {

	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)

	app.Handle(http.MethodGet, "v1", "/hack", Hack)
	app.Handle(http.MethodGet, "v1", "/hack-auth", Hack, authen, ruleAdmin)
}
