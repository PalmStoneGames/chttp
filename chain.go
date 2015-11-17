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
	"net/http"

	"golang.org/x/net/context"
)

// ChainFunc is a function that's part of a chain
type ChainFunc func(Handler) Handler

// Chain represents a chain of functions to be called in succession
type Chain struct {
	ctx   context.Context
	funcs []ChainFunc
}

type chainHandler struct {
	ctx     context.Context
	handler Handler
}

// NewChain creates a new chain of ChainFuncs
func NewChain(ctx context.Context, handlers ...ChainFunc) Chain {
	var funcs []ChainFunc
	funcs = append(funcs, getDefaultChain(ctx)...)
	funcs = append(funcs, handlers...)
	return Chain{
		ctx:   ctx,
		funcs: funcs,
	}
}

// Then assembles a chain and sets finalHandler as the last handler in the chain
// It is safe to call Then many times on the same chain
func (chain Chain) Then(finalHandler Handler) Handler {
	curr := finalHandler
	for i := len(chain.funcs) - 1; i >= 0; i-- {
		curr = chain.funcs[i](curr)
	}

	return chainHandler{
		ctx:     chain.ctx,
		handler: curr,
	}
}

func (chain chainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chain.ServeHTTPContext(CreateContext(chain.ctx, w, r))
}

func (chain chainHandler) ServeHTTPContext(ctx context.Context) {
	wrapper := getRequestCreatorFunc(chain.ctx)

	if wrapper != nil {
		ctx = wrapper(ctx)
	}

	chain.handler.ServeHTTPContext(ctx)
}

// ThenFunc works similarly to Then, but accepts a HandlerFunc instead of a Handler
func (chain Chain) ThenFunc(finalHandler HandlerFunc) Handler {
	return chain.Then(finalHandler)
}
