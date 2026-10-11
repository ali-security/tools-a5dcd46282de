// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"golang.org/x/tools/go/packages/packagestest"
	"golang.org/x/tools/internal/lsp/cmd"
	"golang.org/x/tools/internal/lsp/tests"
)

var isRace = false

type runner struct {
	data *tests.Data
	app  *cmd.Application
}

func TestCommandLine(t *testing.T) {
	packagestest.TestAll(t, testCommandLine)
}

func testCommandLine(t *testing.T, exporter packagestest.Exporter) {
	data := tests.Load(t, exporter, "../testdata")
	defer data.Exported.Cleanup()

	r := &runner{
		data: data,
		app: &cmd.Application{
			Config: *data.Exported.Config,
		},
	}
	tests.Run(t, r, data)
}

// TestServeNoPortFlag checks that the serve verb no longer has a -port flag,
// which listened on all network interfaces.
func TestServeNoPortFlag(t *testing.T) {
	typ := reflect.TypeOf(cmd.Serve{})
	for i := 0; i < typ.NumField(); i++ {
		if field := typ.Field(i); field.Tag.Get("flag") == "port" {
			t.Errorf("Serve.%s provides a -port flag, which binds all network interfaces", field.Name)
		}
	}
}

// runServe runs the serve verb with -listen=addr, failing the test if it does
// not return in time (i.e. it is listening for connections).
func runServe(t *testing.T, addr string) error {
	app := &cmd.Application{}
	app.Serve.Address = addr
	errc := make(chan error, 1)
	go func() { errc <- cmd.RunServe(context.Background(), app) }()
	select {
	case err := <-errc:
		return err
	case <-time.After(10 * time.Second):
		t.Fatalf("-listen=%s: server is listening, want it to return an error", addr)
		return nil
	}
}

// TestServeListenImplicitAllInterfaces checks that -listen rejects addresses
// with an empty host, which would silently bind all network interfaces.
func TestServeListenImplicitAllInterfaces(t *testing.T) {
	const want = "implicitly binds all network interfaces"
	for _, addr := range []string{":0", ":12345"} {
		err := runServe(t, addr)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("-listen=%s: got error %v, want %q", addr, err, want)
		}
	}
	// An explicit host is not rejected; an invalid port makes the
	// listener fail without accepting connections.
	if err := runServe(t, "127.0.0.1:99999"); err == nil || strings.Contains(err.Error(), want) {
		t.Errorf("-listen=127.0.0.1:99999: got error %v, want a listen error", err)
	}
}

func (r *runner) Completion(t *testing.T, data tests.Completions, items tests.CompletionItems) {
	//TODO: add command line completions tests when it works
}

func (r *runner) Format(t *testing.T, data tests.Formats) {
	//TODO: add command line formatting tests when it works
}

func (r *runner) Highlight(t *testing.T, data tests.Highlights) {
	//TODO: add command line highlight tests when it works
}
func (r *runner) Symbol(t *testing.T, data tests.Symbols) {
	//TODO: add command line symbol tests when it works
}

func (r *runner) Signature(t *testing.T, data tests.Signatures) {
	//TODO: add command line signature tests when it works
}

func captureStdOut(t testing.TB, f func()) string {
	r, out, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	defer func() {
		os.Stdout = old
		out.Close()
		r.Close()
	}()
	os.Stdout = out
	f()
	out.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(data))
}
