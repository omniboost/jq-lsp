package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/wader/jq-lsp/lsp"
	"github.com/wader/jq-lsp/profile"
)

var version string

func main() {
	defer profile.MaybeProfile()()

	if version == "" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			version = bi.Main.Version
		}
	}

	// Overwrite read file of JQ's, so relative paths work in our (Omniboost) integrations
	jqsReadFile := func(name string) ([]byte, error) {
		name = strings.Replace(name, `/jqs/jqs/`, `/jqs/`, -1)
		return os.ReadFile(name)
	}

	if err := lsp.Run(lsp.Env{
		Version:  version,
		ReadFile: jqsReadFile,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Args:     os.Args,
		Environ:  os.Environ(),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
