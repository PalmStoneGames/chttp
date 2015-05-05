// Copyright 2015 Palm Stone Games, Inc. All rights reserved.

package chttp // import "code.delta-mmo.com/chttp"

import (
	"net/http"

	"golang.org/x/net/context"
)

const (
	requestKey        = loaderKey("request")
	writerKey         = loaderKey("writer")
	errFuncKey        = loaderKey("errFunc")
	defaultLoadersKey = loaderKey("defaultLoaders")
	defaultChainKey   = loaderKey("defaultChain")
)

func NewContext() context.Context {
	return context.Background()
}

func createContext(w http.ResponseWriter, r *http.Request) context.Context {
	ctx := contextFactory(w, r)
	ctx = context.WithValue(ctx, requestKey, r)
	ctx = context.WithValue(ctx, writerKey, w)

	return ctx
}

func getErrFunc(ctx context.Context) LoadingErrorFunc {
	f := ctx.Value(errFuncKey)
	if f == nil {
		return defaultErrFunc
	}

	return f.(LoadingErrorFunc)
}

func getDefaultLoaders(ctx context.Context) []LoaderFunc {
	l := ctx.Value(defaultLoadersKey)
	if l == nil {
		return nil
	}

	return l.([]LoaderFunc)
}

func getDefaultChain(ctx context.Context) []ChainFunc {
	c := ctx.Value(defaultChainKey)
	if c == nil {
		return nil
	}

	return c.([]ChainFunc)
}

// GetRequest will return the *http.Request given a context
func GetRequest(ctx context.Context) *http.Request {
	r := ctx.Value(requestKey)
	if r == nil {
		return nil
	}

	return r.(*http.Request)
}

// GetWriter will return the http.ResponseWriter given a context
func GetWriter(ctx context.Context) http.ResponseWriter {
	w := ctx.Value(writerKey)
	if w == nil {
		panic("Attempted write on read only context")
	}

	return w.(http.ResponseWriter)
}

// WithReadOnly returns a read only version of the context
// The read only context lacks a http.ResponseWriter and GetWriter will panic if called on it
func WithReadOnly(ctx context.Context) context.Context {
	return context.WithValue(ctx, writerKey, nil)
}

func WithLoadingErrorFunc(parent context.Context, errFunc LoadingErrorFunc) context.Context {
	return context.WithValue(parent, errFuncKey, errFunc)
}

func WithDefaultLoaders(parent context.Context, defaultLoaders ...LoaderFunc) context.Context {
	return context.WithValue(parent, defaultLoadersKey, append(getDefaultLoaders(parent), defaultLoaders...))
}

func WithDefaultChain(parent context.Context, defaultChain ...ChainFunc) context.Context {
	return context.WithValue(parent, defaultChainKey, append(getDefaultChain(parent), defaultChain...))
}
