package mocks

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

func AppendPathParamToChiContext(req *http.Request, key string, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	return req
}
