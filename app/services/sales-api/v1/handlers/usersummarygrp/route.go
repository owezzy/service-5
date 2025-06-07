package usersummarygrp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/owezzy/service-5/business/core/usersummary"
	"github.com/owezzy/service-5/business/core/usersummary/stores/usersummarydb"
	"github.com/owezzy/service-5/business/web/v1/auth"
	"github.com/owezzy/service-5/business/web/v1/mid"
	"github.com/owezzy/service-5/foundation/logger"
	"github.com/owezzy/service-5/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *logger.Logger
	Auth *auth.Auth
	DB   *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	usmCore := usersummary.NewCore(usersummarydb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)

	hdl := New(usmCore)
	app.Handle(http.MethodGet, version, "/usersummary", hdl.Query, authen, ruleAdmin)
}
