// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import "context"

const (
	ExampleLine   = exampleLine
	ExampleColumn = exampleColumn
	ExampleOffset = exampleOffset
)

// RunServe runs the serve verb of app directly and returns its error, rather
// than going through tool.Main, which exits the process on error.
func RunServe(ctx context.Context, app *Application, args ...string) error {
	app.Serve.app = app
	return app.Serve.Run(ctx, args...)
}
