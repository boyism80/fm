package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boyism80/fm/services/game/wz"
)

func main() {
	wzPath := flag.String("wz", "resources/wz", "path to WZ XML root")
	outDir := flag.String("out", "wz_lookup", "output directory for YAML lookup files")
	only := flag.String("only", "all", "all|quests|maps|mobs|reactors|shops|strings|drops")
	flag.Parse()

	resolvedWZ := wz.FindWzPath(*wzPath)
	absOut, err := filepath.Abs(*outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve output path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("loading WZ from %s\n", resolvedWZ)
	resources := wz.NewResources(resolvedWZ)
	if resources == nil {
		fmt.Fprintln(os.Stderr, "failed to load WZ resources")
		os.Exit(1)
	}

	fmt.Printf("writing lookup YAML to %s\n", absOut)
	if err := exportAll(resources, resolvedWZ, absOut, *only); err != nil {
		fmt.Fprintf(os.Stderr, "export failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("done")
}
