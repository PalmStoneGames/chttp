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
	r := ctx.Value(loaderKey("request"))
	if r == nil {
		return nil
	}

	return r.(*http.Request)
}

// GetWriter will return the http.ResponseWriter given a context
func GetWriter(ctx context.Context) http.ResponseWriter {
	w := ctx.Value(loaderKey("writer"))
	if w == nil {
		panic("Attempted write on read only context")
	}

	return w.(http.ResponseWriter)
}

// ReadOnlyContext returns a read only version of the context
// The read only context lacks a http.ResponseWriter and GetWriter will return nil on it
func ReadOnlyContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loaderKey("writer"), nil)
}
