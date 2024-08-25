package hackgrp

import (
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
)

func Routes(mux *httptreemux.ContextMux) {
	mux.Handle(http.MethodGet, "/hack", Hack)
}
