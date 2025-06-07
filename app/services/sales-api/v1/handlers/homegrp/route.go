package homegrp

import (
	"github.com/owezzy/service-5/business/core/event"
	"github.com/owezzy/service-5/business/core/home"
	"github.com/owezzy/service-5/business/core/home/stores/homedb"
	"github.com/owezzy/service-5/business/core/user"
	"github.com/owezzy/service-5/business/core/user/stores/usercache"
	"github.com/owezzy/service-5/business/core/user/stores/userdb"
	db "github.com/owezzy/service-5/business/data/dbsql/pgx"
	"github.com/owezzy/service-5/business/web/v1/auth"
	"github.com/owezzy/service-5/business/web/v1/mid"
	"github.com/owezzy/service-5/foundation/logger"
	"github.com/owezzy/service-5/foundation/web"
	"net/http"

	"github.com/jmoiron/sqlx"
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
	hmeCore := home.NewCore(cfg.Log, envCore, usrCore, homedb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)
	tran := mid.ExecuteInTransation(cfg.Log, db.NewBeginner(cfg.DB))

	hdl := New(hmeCore)
	app.Handle(http.MethodGet, version, "/homes", hdl.Query, authen)
	app.Handle(http.MethodGet, version, "/homes/:home_id", hdl.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/homes", hdl.Create, authen)
	app.Handle(http.MethodPut, version, "/homes/:home_id", hdl.Update, authen, tran)
	app.Handle(http.MethodDelete, version, "/homes/:home_id", hdl.Delete, authen, tran)
}
