package main

import (
	"bytes"
	"context"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/sagernet/sing-box/docgen"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/schema"

	E "github.com/sagernet/sing/common/exceptions"
)

const schemaPath = "docs/schema.json"

func main() {
	flagCheck := flag.Bool("check", false, "verify generated documentation is up to date without writing")
	flagBase := flag.String("base", docgen.DefaultBaseSha, "baseline git ref containing the legacy documentation")
	flagOnly := flag.String("only", "", "restrict processing to the listed page files (comma-separated glob)")
	flag.Parse()
	var only []string
	if *flagOnly != "" {
		only = strings.Split(*flagOnly, ",")
		for index, pattern := range only {
			only[index] = filepath.ToSlash(strings.TrimSpace(pattern))
		}
	}
	err := docgen.Run(*flagCheck, *flagBase, only)
	if err != nil {
		log.Fatal(err)
	}
	err = runSchema(*flagCheck)
	if err != nil {
		log.Fatal(err)
	}
}

// runSchema regenerates docs/schema.json with descriptions from the option
// package doc comments. The sing-box schema subcommand itself stays
// dependency-free: it emits a plain schema and works outside the repository.
func runSchema(check bool) error {
	repoRoot, err := docgen.FindRepoRoot()
	if err != nil {
		return err
	}
	schemaFile := filepath.Join(repoRoot, filepath.FromSlash(schemaPath))
	commentSource, err := docgen.NewSchemaCommentSource()
	if err != nil {
		return err
	}
	ctx := schema.ContextWithCommentSource(include.Context(context.Background()), commentSource)
	content, err := schema.Generate(ctx, reflect.TypeFor[option.Options]())
	if err != nil {
		return err
	}
	if check {
		current, err := os.ReadFile(schemaFile)
		if err != nil {
			if os.IsNotExist(err) {
				return E.New(schemaPath, " does not exist, run `go run ./cmd/internal/generate_docs`")
			}
			return err
		}
		if !bytes.Equal(current, content) {
			return E.New("docs/schema.json is not up to date, run `go run ./cmd/internal/generate_docs`")
		}
		return nil
	}
	return os.WriteFile(schemaFile, content, 0o644)
}
