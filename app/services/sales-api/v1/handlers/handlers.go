package handlers

import (
	"github.com/owezzy/service-5/app/services/sales-api/v1/handlers/checkgrp"
	"github.com/owezzy/service-5/app/services/sales-api/v1/handlers/hackgrp"
	v1 "github.com/owezzy/service-5/business/web/v1"
	"github.com/owezzy/service-5/foundation/web"
)

type Routes struct {
}

// All implements the routerAdder interface

func (Routes) Add(app *web.App, apiCfg v1.APIMuxConfig) {

	cfg := hackgrp.Config{Auth: apiCfg.Auth}
	hackgrp.Routes(app, cfg)

	checkgrp.Routes(app, checkgrp.Config{
		Build: cfg.Build,
	})
}
