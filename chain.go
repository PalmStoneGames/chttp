// Copyright 2015 Palm Stone Games, Inc.

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
	return Chain{
		ctx:   ctx,
		funcs: append(getDefaultChain(ctx), handlers...)}
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
	chain.ServeHTTPContext(createContext(chain.ctx, w, r))
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
