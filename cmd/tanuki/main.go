package main

import (
	"fmt"
	"os"

	"github.com/example/tanuki/internal/catalog"
	"github.com/example/tanuki/internal/formatter"
	"github.com/example/tanuki/internal/schema"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "list":
		runList()
	case "status":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: tanuki status <service-name>")
			os.Exit(1)
		}
		runStatus(args[0])
	case "owners":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: tanuki owners <service-name>")
			os.Exit(1)
		}
		runOwners(args[0])
	case "search":
		team := ""
		for i := 0; i < len(args); i++ {
			if args[i] == "--team" && i+1 < len(args) {
				team = args[i+1]
				break
			}
		}
		if team == "" {
			fmt.Fprintln(os.Stderr, "Usage: tanuki search --team <team-name>")
			os.Exit(1)
		}
		runSearch(team)
	case "validate":
		runValidate()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprint(os.Stderr, `tanuki - Service catalog CLI

Usage:
  tanuki list                      List all registered services
  tanuki status <name>             Show health, version, owner, last deploy
  tanuki owners <name>            Show owner and on-call info
  tanuki search --team <team>     Filter services by team
  tanuki validate                 Validate local registry against schema

Catalog source: TANUKI_CATALOG_URL or ./catalog.json or ./dist/catalog.json
`)
}

func runList() {
	svcs, err := catalog.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	formatter.List(svcs)
}

func runStatus(name string) {
	svcs, err := catalog.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	s := catalog.FindByName(svcs, name)
	formatter.Status(s)
	if s == nil {
		os.Exit(1)
	}
}

func runOwners(name string) {
	svcs, err := catalog.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	s := catalog.FindByName(svcs, name)
	formatter.Owners(s)
	if s == nil {
		os.Exit(1)
	}
}

func runSearch(team string) {
	svcs, err := catalog.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	filtered := catalog.FilterByTeam(svcs, team)
	formatter.List(filtered)
}

func runValidate() {
	repoRoot, err := catalog.RepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := schema.Validate(repoRoot); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
