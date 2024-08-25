package hackgrp

import (
	"context"
	"github.com/owezzy/service-5/foundation/web"
	"net/http"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK MF AGN",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
