// Copyright 2015 Palm Stone Games, Inc. All rights reserved.

package chttp

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

// LoadingErrorHandler is the handler to be called when an error happens in one of the loaders
// The default implementation will simply print to http.Error, but may be replaced
var LoadingErrorHandler = func(ctx context.Context, err string) {
	Error(ctx, err, http.StatusInternalServerError)
}

// LoaderFunc describes a function used to load in data
// it should return a channel on which the data is sent and the key in which to store it
type LoaderFunc func(context.Context) (<-chan interface{}, interface{})

// HandlerFunc describes a function to handle a http request
type HandlerFunc func(context.Context)

// loaderKey is the type we use for our context keys so they can't conflict with any other module
type loaderKey string

// Handler is the interface implemented by HandlerFunc and LoaderFunc
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	ServeHTTPContext(context.Context)
}

// ServeHTTPContext will serve a http request given a context
func (h HandlerFunc) ServeHTTPContext(ctx context.Context) {
	h(ctx)
}

// ServeHTTP will serve a http request given a responseWriter and the request
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(createContext(w, r))
}

// Loader is a structure that can be used to schedule many parallel LoaderFuncs to be ran when the request is served
// Each LoaderFunc is ran in parallel, if all run successfully, HandlerFunc will be called, otherwise, LoadingErrorHandler will be called
type Loader struct {
	handler Handler
	loaders []LoaderFunc
}

// NewLoader creates a new loader
func NewLoader(handler Handler, loaders ...LoaderFunc) Loader {
	return Loader{
		handler: handler,
		loaders: loaders,
	}
}

func NewLoaderFunc(handler HandlerFunc, loaders ...LoaderFunc) Loader {
	return NewLoader(handler, loaders...)
}

// ServeHTTP will serve a http request given a responseWriter and the request
func (l Loader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.ServeHTTPContext(createContext(w, r))
}

// ServeHTTPContext will serve a http request given a context
func (l Loader) ServeHTTPContext(ctx context.Context) {
	var chans []<-chan interface{}
	var keys []interface{}

	// Startup all the loaders
	for _, loader := range l.loaders {
		loaderCh, key := loader(ctx)
		chans = append(chans, loaderCh)
		keys = append(keys, key)
	}

	// Wait for all the data
	failed := make(map[interface{}]struct{})
	for i, ch := range chans {
		val, ok := <-ch
		if !ok {
			failed[keys[i]] = struct{}{}
		} else {
			ctx = context.WithValue(ctx, keys[i], val)
		}
	}

	// Error reporting
	if len(failed) != 0 {
		var errText string
		for _, k := range keys {
			_, isFailed := failed[k]
			errText += fmt.Sprintf("%v: %v\n", k, !isFailed)
		}

		LoadingErrorHandler(ctx, errText)
		return
	}

	// Call the passed handler with all the data
	l.handler.ServeHTTPContext(ctx)
}