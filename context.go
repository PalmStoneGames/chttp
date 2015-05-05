// Copyright 2015 Palm Stone Games, Inc. All rights reserved.

package chttp // import "code.delta-mmo.com/chttp"

import (
	"net/http"

	"golang.org/x/net/context"
)

func createContext(w http.ResponseWriter, r *http.Request) context.Context {
	ctx := contextFactory(w, r)
	ctx = context.WithValue(ctx, loaderKey("request"), r)
	ctx = context.WithValue(ctx, loaderKey("writer"), w)

	return ctx
}

// GetRequest will return the *http.Request given a context
func GetRequest(ctx context.Context) *http.Request {
	return ctx.Value(loaderKey("request")).(*http.Request)
}

// GetWriter will return the http.ResponseWriter given a context
func GetWriter(ctx context.Context) http.ResponseWriter {
	return ctx.Value(loaderKey("writer")).(http.ResponseWriter)
}
