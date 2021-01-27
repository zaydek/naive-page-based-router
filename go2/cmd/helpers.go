package main

import (
	"io/ioutil"
	"os"
	p "path"
	"path/filepath"

	"github.com/zaydek/retro/errs"
)

// getCmd gets the current command.
func (r Runtime) getCmd() string {
	if r.CreateCommand != nil {
		return "create"
	} else if r.WatchCommand != nil {
		return "watch"
	} else if r.BuildCommand != nil {
		return "build"
	} else if r.ServeCommand != nil {
		return "serve"
	}
	return ""
}

// getPort gets the current port.
func (r Runtime) getPort() int {
	if cmd := r.getCmd(); cmd == "watch" {
		return r.WatchCommand.Port
	} else if cmd == "serve" {
		return r.ServeCommand.Port
	}
	return 0
}

type copyPath struct {
	src string
	dst string
}

// copyAssetDirectoryToBuildDirectory recursively copies the asset directory to
// the build directory.
func copyAssetDirectoryToBuildDirectory(config DirConfiguration) error {

	var paths []copyPath
	if err := filepath.Walk(config.AssetDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if info.Name() == "index.html" {
				return nil
			}
			src := path
			dst := p.Join(config.BuildDirectory, path)
			paths = append(paths, copyPath{src: src, dst: dst})
		}
		return nil
	}); err != nil {
		errs.Walk(config.AssetDirectory, err)
	}

	for _, each := range paths {
		if dir := p.Dir(each.dst); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return errs.MkdirAll(dir, err)
			}
		}
	}

	for _, each := range paths {
		bstr, err := os.ReadFile(each.src)
		if err != nil {
			return errs.ReadFile(each.src, err)
		}
		if err := ioutil.WriteFile(each.dst, bstr, 0644); err != nil {
			return errs.ReadFile(each.dst, err)
		}
	}
	return nil
}
