package v1

import (
	"github.com/dimfeld/httptreemux/v5"
	"github.com/owezzy/service-5/foundation/logger"
	"os"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *httptreemux.ContextMux, cfg APIMuxConfig)
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig, routerAdder RouteAdder) *httptreemux.ContextMux {
	mux := httptreemux.NewContextMux()

	routerAdder.Add(mux, cfg)

	return mux
}
