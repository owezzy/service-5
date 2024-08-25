package handlers

import (
	"github.com/owezzy/service-5/app/services/sales-api/v1/handlers/hackgrp"
	v1 "github.com/owezzy/service-5/business/web/v1"
	"github.com/owezzy/service-5/foundation/web"
)

type Routes struct {
}

// Add implements the routerAdder interface

func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	hackgrp.Routes(app)

}
