// Copyright 2015 Palm Stone Games, Inc. All rights reserved.

package chttp

import "golang.org/x/net/context"

// ChainFunc is a function that's part of a chain
type ChainFunc func(Handler) Handler

// Chain represents a chain of functions to be called in succession
type Chain struct {
	ctx   context.Context
	funcs []ChainFunc
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

	return chain.contextCreatorWrapper(curr)
}

func (chain Chain) contextCreatorWrapper(h Handler) Handler {
	wrapper := getRequestCreatorFunc(chain.ctx)

	if wrapper == nil {
		return h
	}

	return HandlerFunc(func(ctx context.Context) {
		h.ServeHTTPContext(wrapper(ctx))
	})
}

// ThenFunc works similarly to Then, but accepts a HandlerFunc instead of a Handler
func (chain Chain) ThenFunc(finalHandler HandlerFunc) Handler {
	return chain.Then(finalHandler)
}
