package dev

import (
	"text/template"
)

// PageBasedRoute describes a page-based route.
type PageBasedRoute struct {
	SrcPath   string `json:"srcPath"`   // pages/path/to/component.js
	DstPath   string `json:"dstPath"`   // build/path/to/component.html
	Path      string `json:"path"`      // path/to/component
	Component string `json:"component"` // Component
}

// DirectoryConfiguration describes persistent directory configuration.
type DirectoryConfiguration struct {
	AssetDirectory string
	PagesDirectory string
	CacheDirectory string
	BuildDirectory string
}

type Runtime struct {
	// TODO: Remove hashes?
	EpochUUID         string
	IndexHTMLTemplate *template.Template
	Command           interface{}
	DirConfiguration  DirectoryConfiguration
	PageBasedRouter   []PageBasedRoute
}

// ExperimentalReactSuspenseEnabled   bool // Wrap <React.Suspense>
// ExperimentalReactStrictModeEnabled bool // Wrap <React.StrictMode>
