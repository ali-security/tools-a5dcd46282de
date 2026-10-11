// Command create_mod writes the proxy.golang.org-style module zip for
// golang.org/x/tools from a source directory, using golang.org/x/mod/zip.
//
// Usage: go run . <source-dir> <output.zip>
//
// It is built from a temporary directory with its own go.mod (see
// .github/workflows/seal-build.yml); it is not part of the x/tools module.
package main

import (
	"log"
	"os"

	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

const (
	modulePath = "golang.org/x/tools"
	version    = "v0.0.0-20190425150028-36563e24a262"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s <source-dir> <output.zip>", os.Args[0])
	}
	srcDir, outPath := os.Args[1], os.Args[2]

	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	mv := module.Version{Path: modulePath, Version: version}
	if err := zip.CreateFromDir(out, mv, srcDir); err != nil {
		out.Close()
		log.Fatal(err)
	}
	if err := out.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s for %s@%s from %s", outPath, modulePath, version, srcDir)
}
