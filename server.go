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
