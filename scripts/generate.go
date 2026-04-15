//go:build ignore

// Generate typed models from the AllToken OpenAPI specs.
//
// Reads from a sibling ../megaopenrouter/openapi/{chat,anthropic}.yml and
// writes to internal/gen/{chat,anthropic}/types.go.
//
// Run from the module root:
//
//	go run scripts/generate.go
//
// This uses `go run <pkg>@version` to invoke oapi-codegen without requiring a
// pre-installed binary on PATH. First run downloads to the module cache;
// subsequent runs are fast.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const oapiCodegenPkg = "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.6.0"

var specs = []string{"chat", "anthropic"}

func main() {
	root, err := os.Getwd()
	if err != nil {
		die(err)
	}

	specDir := filepath.Join(root, "..", "megaopenrouter", "openapi")
	if _, err := os.Stat(specDir); err != nil {
		die(fmt.Errorf("spec dir not found at %s\n  Clone megaopenrouter as a sibling:\n    git clone git@gitlab.53site.com:ai-innovation-lab/megaopenrouter.git ../megaopenrouter", specDir))
	}

	for _, name := range specs {
		spec := filepath.Join(specDir, name+".yml")
		if _, err := os.Stat(spec); err != nil {
			die(fmt.Errorf("spec not found at %s", spec))
		}

		out := filepath.Join(root, "internal", "gen", name, "types.go")
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			die(err)
		}

		cmd := exec.Command("go", "run", oapiCodegenPkg,
			"-generate", "types",
			"-package", name,
			"-o", out,
			spec,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			die(fmt.Errorf("codegen %s: %w", name, err))
		}
		fmt.Printf("[generate] %s.yml -> internal/gen/%s/types.go\n", name, name)
	}

	fmt.Println("[generate] Done.")
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "[generate] ERROR: %v\n", err)
	os.Exit(1)
}
