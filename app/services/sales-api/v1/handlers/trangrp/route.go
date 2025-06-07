package trangrp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/owezzy/service-5/business/core/event"
	"github.com/owezzy/service-5/business/core/product"
	"github.com/owezzy/service-5/business/core/product/stores/productdb"
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
	Log  *logger.Logger
	Auth *auth.Auth
	DB   *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	envCore := event.NewCore(cfg.Log)
	usrCore := user.NewCore(cfg.Log, envCore, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))
	prdCore := product.NewCore(cfg.Log, envCore, usrCore, productdb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)
	tran := mid.ExecuteInTransation(cfg.Log, db.NewBeginner(cfg.DB))

	hdl := New(usrCore, prdCore)
	app.Handle(http.MethodPost, version, "/tranexample", hdl.Create, authen, tran)
}
