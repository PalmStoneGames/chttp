// Copyright 2015 Palm Stone Games, Inc. All rights reserved.

package chttp

import "golang.org/x/net/context"

// ChainFunc is a function that's part of a chain
type ChainFunc func(Handler) Handler

// Chain represents a chain of functions to be called in succession
type Chain []ChainFunc

// NewChain creates a new chain of ChainFuncs
func NewChain(ctx context.Context, handlers ...ChainFunc) Chain {
	return Chain(append(getDefaultChain(ctx), handlers...))
}

// Then assembles a chain and sets finalHandler as the last handler in the chain
// It is safe to call Then many times on the same chain
func (chain Chain) Then(finalHandler Handler) Handler {
	curr := finalHandler
	for i := len(chain) - 1; i >= 0; i-- {
		curr = chain[i](curr)
	}

	return curr
}

// ThenFunc works similarly to Then, but accepts a HandlerFunc instead of a Handler
func (chain Chain) ThenFunc(finalHandler HandlerFunc) Handler {
	return chain.Then(finalHandler)
}
