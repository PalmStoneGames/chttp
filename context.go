// Copyright 2015 Palm Stone Games, Inc.

package chttp

import (
	"net/http"

	"golang.org/x/net/context"
)

const (
	keyRequest        = loaderKey("request")
	keyWriter         = loaderKey("writer")
	keyErrFunc        = loaderKey("errFunc")
	keyRequestCreator = loaderKey("requestCreatorFunc")
	keyDefaultLoaders = loaderKey("defaultLoaders")
	keyDefaultChain   = loaderKey("defaultChain")
)

type Context struct {
	context.Context
}

// RequestCreatorFunc can be used to manipulate the context before passing it off to other handlers
// This is useful when you need a special context everywhere (for example, on appengine)
type RequestCreatorFunc func(ctx context.Context) context.Context

// NewContext creates a new empty context
func NewContext() Context {
	return Context{context.Background()}
}

func createContext(parent context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	ctx := parent
	ctx = context.WithValue(ctx, keyRequest, r)
	ctx = context.WithValue(ctx, keyWriter, w)

	creator := getRequestCreatorFunc(ctx)
	if creator != nil {
		ctx = creator(ctx)
	}

	return ctx
}

func getErrFunc(ctx context.Context) LoadingErrorFunc {
	f := ctx.Value(keyErrFunc)
	if f == nil {
		return defaultErrFunc
	}

	return f.(LoadingErrorFunc)
}

func getDefaultLoaders(ctx context.Context) []LoaderFunc {
	l := ctx.Value(keyDefaultLoaders)
	if l == nil {
		return nil
	}

	return l.([]LoaderFunc)
}

func getDefaultChain(ctx context.Context) []ChainFunc {
	c := ctx.Value(keyDefaultChain)
	if c == nil {
		return nil
	}

	return c.([]ChainFunc)
}

func getRequestCreatorFunc(ctx context.Context) RequestCreatorFunc {
	c := ctx.Value(keyRequestCreator)
	if c == nil {
		return nil
	}

	return c.(RequestCreatorFunc)
}

// GetRequest will return the *http.Request given a context
func GetRequest(ctx context.Context) *http.Request {
	r := ctx.Value(keyRequest)
	if r == nil {
		return nil
	}

	return r.(*http.Request)
}

// GetWriter will return the http.ResponseWriter given a context
func GetWriter(ctx context.Context) http.ResponseWriter {
	w := ctx.Value(keyWriter)
	if w == nil {
		panic("Attempted write on read only context")
	}

	return w.(http.ResponseWriter)
}

// AsReadOnly returns a read only version of the context
// The read only context lacks a http.ResponseWriter and GetWriter will panic if called on it
func AsReadOnly(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyWriter, nil)
}

// WithRequestCreatorFunc returns a context with the passed request creator applied to it
// the request creator is called right after creating the request specific context so it can get the chance to modify it
// this is useful when other parts of your application also use contexts
func (ctx Context) WithRequestCreatorFunc(creator RequestCreatorFunc) Context {
	return Context{context.WithValue(ctx, keyRequestCreator, creator)}
}

// WithLoadingErrorFunc returns a context with the passed loading error handler applied to it
// the LoadingErrorFunc will be called when an error happens in one of the loader functions passed to NewLoader
// a Loader is considered in error when it closes its channel before sending on it
func (ctx Context) WithLoadingErrorFunc(errFunc LoadingErrorFunc) Context {
	return Context{context.WithValue(ctx.Context, keyErrFunc, errFunc)}
}

// WithDefaultLoaders returns a context on which all subsequent NewLoader calls will have these additional loaders prepended to it
func (ctx Context) WithDefaultLoaders(defaultLoaders ...LoaderFunc) Context {
	return Context{context.WithValue(ctx.Context, keyDefaultLoaders, append(getDefaultLoaders(ctx.Context), defaultLoaders...))}
}

// WithDefaultChain returns a context on which all subsequent NewChain calls with have the passed middleware prepended to it
func (ctx Context) WithDefaultChain(defaultChain ...ChainFunc) Context {
	return Context{context.WithValue(ctx.Context, keyDefaultChain, append(getDefaultChain(ctx.Context), defaultChain...))}
}
