package hackgrp

import (
	"context"
	"errors"
	"github.com/owezzy/service-5/business/web/v1/response"
	"github.com/owezzy/service-5/foundation/web"
	"math/rand"
	"net/http"
)

type Handlers struct {
	build string
}

// New constructs a Handlers api for the check group.
func New(build string) *Handlers {
	return &Handlers{
		build: build,
	}
}
func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100) % 2; n == 0 {
		return response.NewError(errors.New("CUSTOM ERROR"), http.StatusBadRequest)
	}
	status := struct {
		Status string
	}{
		Status: "OK MF AGN",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
