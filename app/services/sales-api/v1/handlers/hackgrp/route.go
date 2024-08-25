package hackgrp

import (
	"github.com/owezzy/service-5/foundation/web"
	"net/http"
)

func Routes(app *web.App) {
	app.Handle(http.MethodGet, "/hack", Hack)
}
