// Copyright 2015 Palm Stone Games, Inc.

package chttp

import (
	"net/http"

	"golang.org/x/net/context"
)

// Error replies to the request with the specified error message and HTTP code.
// The error message should be plain text.
func Error(ctx context.Context, err string, code int) {
	http.Error(GetWriter(ctx), err, code)
}

// Redirect replies to the request with a redirect to url,
// which may be a path relative to the request path.
func Redirect(ctx context.Context, urlStr string, code int) {
	http.Redirect(GetWriter(ctx), GetRequest(ctx), urlStr, code)
}
