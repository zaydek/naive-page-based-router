package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type PageBasedRoute struct {
	Path      string `json:"path"`      // path/to/component.tsx
	Page      string `json:"page"`      // /component
	Component string `json:"component"` // Component
}

type PageBasedRouter []PageBasedRoute

var routeFileTypes = []string{
	".js",  // JavaScript
	".jsx", // React JavaScript
	".ts",  // TypeScript
	".tsx", // React TypeScript
	".md",  // Markdown
	".mdx", // MDX (React Markdown)
}

func isRouteFileType(path string) bool {
	ext := filepath.Ext(path)
	for _, fileType := range routeFileTypes {
		if ext == fileType {
			return true
		}
	}
	return false
}

// Initializes a new page-based route. A page-based route provides a layer of
// indirection so that a path name can be queried as a page name or as a React-
// constructable component name.
func newPageBasedRoute(config Configuration, path string) PageBasedRoute {
	// TODO: Sanitize `path`; should be limited to set of cross-platform ASCII
	// characters. In the future, this can be broadened to support Unicode
	// characters more generally.
	// TODO: For now, let’s lazily qualify the path name against a regex. Later,
	// we should qualify more carefully using `parseParts` or equivalent.

	x1 := len(config.PagesDir + "/")
	x2 := len(path) - len(filepath.Ext(path))

	route := PageBasedRoute{
		Path:      path,
		Page:      path[x1:x2],
		Component: "TODO",
	}
	return route
}

// TODO: We still need to remove `config.PagesDir` from `route.Page`.
func InitPageBasedRouter(config Configuration) (PageBasedRouter, error) {
	var router PageBasedRouter
	err := filepath.Walk(config.PagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// TODO: Extract `noOpFilesAndFolders`?
		if info.IsDir() && info.Name() == "internal" {
			return filepath.SkipDir
		}
		if isRouteFileType(path) {
			router = append(router, newPageBasedRoute(config, path))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot get page-based router; %w", err)
	}
	return router, nil
}