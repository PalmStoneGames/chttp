// +build !appengine

// Copyright 2015 Palm Stone Games, Inc. All rights reserved.

package chttp

import (
	"net/http"

	"golang.org/x/net/context"
)

func contextFactory(w http.ResponseWriter, r *http.Request) context.Context {
	return context.Background()
}
