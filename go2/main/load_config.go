package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/zaydek/retro/static"
)

// Configuration describes global configuration.
type Configuration struct {
	// The asset directory.
	AssetDir string `JSON:"ASSET_DIR"`

	// The pages directory.
	PagesDir string `json:"PAGES_DIR"`

	// The cache directory.
	CacheDir string `json:"CACHE_DIR"`

	// The build directory.
	BuildDir string `json:"BUILD_DIR"`

	// Wrap the app with <React.StrictMode>.
	ReactStrictMode bool
}

// Validates the presence of a directory.
func validateDirectory(field, dirname string) error {
	if info, err := os.Stat(dirname); os.IsNotExist(err) {
		if err := os.MkdirAll(dirname, 0755); err != nil {
			return fmt.Errorf("server guard: cannot create a directory for configuration field %s=%s; %w", field, dirname, err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("server guard: configuration field %s=%s must be a directory", field, dirname)
	}
	return nil
}

// Validates the presence of public/index.html.
func validatePublicIndexHTML(rc Configuration) error {
	if _, err := os.Stat(rc.PagesDir); os.IsNotExist(err) {
		src, err := static.StaticFS.Open("static/public/index.html")
		if err != nil {
			return fmt.Errorf("an unexpected error occurred; %w", err)
		}
		dst, err := os.Create(path.Join(rc.AssetDir, "index.html"))
		if err != nil {
			return fmt.Errorf("an unexpected error occurred; %w", err)
		}
		if _, err := io.Copy(dst, src); err != nil {
			return fmt.Errorf("an unexpected error occurred; %w", err)
		}
		src.Close()
		dst.Close()
	}
	return nil
}

func validateServerGuards(rc Configuration) error {
	dirs := []struct {
		field   string
		dirname string
	}{
		{field: "ASSET_DIR", dirname: rc.AssetDir},
		{field: "PAGES_DIR", dirname: rc.PagesDir},
		{field: "CACHE_DIR", dirname: rc.CacheDir},
		{field: "PAGES_DIR", dirname: rc.PagesDir},
	}
	// Passthrough error:
	for _, each := range dirs {
		if err := validateDirectory(each.field, each.dirname); err != nil {
			return err
		}
	}
	// Passthrough error:
	if err := validatePublicIndexHTML(rc); err != nil {
		return err
	}
	return nil
}

// Matches comments:
//
// - [ \t]*       spaces
// - \/\/         //
// - .*           any (to EOL)
//
// - \/\*         /*
// - (?:\n?.*?)+  one or more paragraphs (lazy)
// - \*\/         */
//
var commentRe = regexp.MustCompile(`[ \t]*\/\/.*` + `|` + `\/\*(?:\n?.*?)+\*\/`)

// loadConfiguration loads or creates retro.config.jsonc.
//
// TODO: Add support for user overriding configuration fields using
// environmental variables; environmental variables take precedence over
// configuration fields.
// TODO: Add support for JavaScript (or TypeScript-based?) configuration files.
func loadConfiguration() (Configuration, error) {
	config := Configuration{
		AssetDir:        "public",
		PagesDir:        "pages",
		CacheDir:        "cache",
		BuildDir:        "build",
		ReactStrictMode: false,
	}
	// Stat or create retro.config.jsonc:
	if _, err := os.Stat("retro.config.jsonc"); os.IsNotExist(err) {
		bstr, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return Configuration{}, fmt.Errorf("cannot write retro.config.jsonc; %w", err)
		}
		bstr = append(bstr, '\n') // EOF
		if err := ioutil.WriteFile("retro.config.jsonc", bstr, 0644); err != nil {
			return Configuration{}, fmt.Errorf("cannot write retro.config.jsonc; %w", err)
		}
	}
	// Read retro.config.jsonc:
	bstr, err := ioutil.ReadFile("retro.config.jsonc")
	if err != nil {
		return Configuration{}, fmt.Errorf("cannot read retro.config.jsonc; %w", err)
	}
	// Remove comments:
	bstr = commentRe.ReplaceAll(bstr, []byte(""))
	if err = json.Unmarshal(bstr, &config); err != nil {
		return Configuration{}, fmt.Errorf("cannot read retro.config.jsonc; %w", err)
	}
	// Passthrough error:
	if err := validateServerGuards(config); err != nil {
		return Configuration{}, err
	}
	return config, nil
}