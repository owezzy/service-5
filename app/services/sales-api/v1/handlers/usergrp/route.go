package usergrp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/owezzy/service-5/business/core/event"
	"github.com/owezzy/service-5/business/core/user"
	"github.com/owezzy/service-5/business/core/user/stores/usercache"
	"github.com/owezzy/service-5/business/core/user/stores/userdb"
	db "github.com/owezzy/service-5/business/data/dbsql/pgx"
	"github.com/owezzy/service-5/business/web/v1/auth"
	"github.com/owezzy/service-5/business/web/v1/mid"
	"github.com/owezzy/service-5/foundation/logger"
	"github.com/owezzy/service-5/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	Auth  *auth.Auth
	DB    *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	envCore := event.NewCore(cfg.Log)
	usrCore := user.NewCore(cfg.Log, envCore, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))

	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)
	ruleAdminOrSubject := mid.Authorize(cfg.Auth, auth.RuleAdminOrSubject)
	tran := mid.ExecuteInTransation(cfg.Log, db.NewBeginner(cfg.DB))

	hdl := New(usrCore, cfg.Auth)
	app.Handle(http.MethodGet, version, "/users/token/:kid", hdl.Token)
	app.Handle(http.MethodGet, version, "/users", hdl.Query, authen, ruleAdmin)
	app.Handle(http.MethodGet, version, "/users/:user_id", hdl.QueryByID, authen, ruleAdminOrSubject)
	app.Handle(http.MethodPost, version, "/users", hdl.Create, authen, ruleAdmin)
	app.Handle(http.MethodPut, version, "/users/:user_id", hdl.Update, authen, ruleAdminOrSubject, tran)
	app.Handle(http.MethodDelete, version, "/users/:user_id", hdl.Delete, authen, ruleAdminOrSubject, tran)
}
