/*
Copyright 2015 Palm Stone Games, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package chttp

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

// LoaderFunc describes a function used to load in data
// it should return a channel on which the data is sent and the key in which to store it
type LoaderFunc func(context.Context) (<-chan interface{}, interface{})

// HandlerFunc describes a function to handle a http request
type HandlerFunc func(context.Context)

// LoadingErrorFunc is the signature of the function called when a loading error occurs
type LoadingErrorFunc func(context.Context, error)

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
	h(CreateContext(nil, w, r))
}

// Loader is a structure that can be used to schedule many parallel LoaderFuncs to be ran when the request is served
// Each LoaderFunc is ran in parallel, if all run successfully, HandlerFunc will be called, otherwise, LoadingErrorHandler will be called
// A LoaderFunc is considered to be successful as long as it doesn't close its channel before sending on it once
type Loader struct {
	ctx     context.Context
	handler Handler
	loaders []LoaderFunc
}

// NewLoader creates a new loader
func NewLoader(ctx context.Context, handler Handler, loaders ...LoaderFunc) Loader {
	var funcs []LoaderFunc
	funcs = append(funcs, getDefaultLoaders(ctx)...)
	funcs = append(funcs, loaders...)
	return Loader{
		ctx:     ctx,
		handler: handler,
		loaders: funcs,
	}
}

// NewLoaderFunc is a wrapper around NewLoader for when the handler is a HandlerFunc
func NewLoaderFunc(ctx context.Context, handler HandlerFunc, loaders ...LoaderFunc) Loader {
	return NewLoader(ctx, handler, loaders...)
}

// ServeHTTP will serve a http request given a responseWriter and the request
func (l Loader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.ServeHTTPContext(CreateContext(l.ctx, w, r))
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

		getErrFunc(ctx)(ctx, errors.New(errText))
		return
	}

	// Call the passed handler with all the data
	l.handler.ServeHTTPContext(ctx)
}

func defaultErrFunc(ctx context.Context, err error) {
	Error(ctx, err.Error(), http.StatusInternalServerError)
}
