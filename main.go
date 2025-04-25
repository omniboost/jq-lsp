package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime/debug"

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
	omniboostSubdirRegex := regexp.MustCompile("/jqs/(.+/)?jqs/")
	jqsReadFile := func(name string) ([]byte, error) {
		name = omniboostSubdirRegex.ReplaceAllString(name, "/jqs/")
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
