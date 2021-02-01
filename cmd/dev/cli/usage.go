package cli

import (
	"strings"

	"github.com/zaydek/retro/pkg/term"
)

var usageOnly = strings.TrimSpace(`
retro start [flags]  Starts the dev server
retro build [flags]  Builds the production-ready build
retro serve [flags]  Serves the production-ready build
`)

var usage = `
  ` + term.Bold("Usage:") + `

    retro start [flags]  Starts the dev server
    retro build [flags]  Builds the production-ready build
    retro serve [flags]  Serves the production-ready build

  ` + term.Bold("retro watch [flags] paths") + `

    Starts a dev server and watches paths for changes

    Flags:
      --cached           Use cached resources (default false)
      --port=<number>    Port number (default 8000)
      --source-map       Add source maps (default false)

  ` + term.Bold("retro build [flags]") + `

    Builds the production-ready build

    Flags:
      --cached           Use cached resources (default false)
      --source-map       Add source maps (default false)

  ` + term.Bold("retro serve [flags]") + `

    Serves the production-ready build

    Flags:
      --port=<number>    Port number (default 8000)

  ` + term.Bold("Repository:") + `

    ` + term.Underline("https://github.com/zaydek/retro") + `
`
