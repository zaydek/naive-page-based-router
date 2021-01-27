package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	pathpkg "path"

	"github.com/zaydek/retro/color"
	"github.com/zaydek/retro/embedded"
	"github.com/zaydek/retro/errs"
	"github.com/zaydek/retro/loggers"
)

var (
	reactVersion        = "latest"
	reactDOMVersion     = "latest"
	retroVersion        = "latest"
	retroScriptsVersion = "latest"
)

// TODO: npx create-retro-app is functionally equivalent to retro create [dir].
func (r Runtime) Create() {
	fsys := embedded.JavaScriptFS
	if r.CreateCommand.Language == "ts" {
		fsys = embedded.TypeScriptFS
	}

	if r.CreateCommand.Directory != "." {
		if info, err := os.Stat(r.CreateCommand.Directory); !os.IsNotExist(err) {
			var typ string
			if !info.IsDir() {
				typ = "file"
			} else {
				typ = "directory"
			}
			loggers.Stderr.Println("Aborted. A " + typ + " named " + color.Boldf("'%s'", r.CreateCommand.Directory) + " already exists.\n\n" +
				"- " + color.Bold("retro create [dir]") + "\n\n" +
				"Or\n\n" +
				"- " + color.Boldf("rm -r %[1]s && retro create %[1]s", r.CreateCommand.Directory))
			os.Exit(1)
		}

		if err := os.MkdirAll(r.CreateCommand.Directory, 0755); err != nil {
			loggers.Stderr.Println(errs.MkdirAll(r.CreateCommand.Directory, err))
			os.Exit(1)
		}
	}

	var paths []copyPath
	if err := fs.WalkDir(fsys, ".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dirEntry.IsDir() {
			src := path
			dst := pathpkg.Join(r.CreateCommand.Directory, path)
			paths = append(paths, copyPath{src: src, dst: dst})
		}
		return nil
	}); err != nil {
		entry := fmt.Sprintf("<embedded:%s>", r.CreateCommand.Language)
		loggers.Stderr.Println(errs.Walk(entry, err))
		os.Exit(1)
	}

	for _, each := range paths {
		if dir := pathpkg.Dir(each.dst); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				loggers.Stderr.Println(errs.MkdirAll(dir, err))
				os.Exit(1)
			}
		}
	}

	for _, each := range paths {
		bstr, err := fs.ReadFile(fsys, each.src)
		if err != nil {
			entry := fmt.Sprintf("<embedded:%s>", each)
			loggers.Stderr.Println(errs.ReadFile(entry, err))
			os.Exit(1)
		}
		if err := ioutil.WriteFile(each.dst, bstr, 0644); err != nil {
			loggers.Stderr.Println(errs.WriteFile(each.dst, err))
			os.Exit(1)
		}
	}

	repoName := r.CreateCommand.Directory
	if repoName == "." {
		repoName = "retro-app"
	}

	dot := embedded.PkgStruct{
		RepoName:            repoName,
		ReactVersion:        reactVersion,
		ReactDOMVersion:     reactDOMVersion,
		RetroVersion:        retroVersion,
		RetroScriptsVersion: retroScriptsVersion,
	}

	var buf bytes.Buffer
	if err := embedded.PkgTemplate.Execute(&buf, dot); err != nil {
		loggers.Stderr.Println(errs.ExecuteTemplate("package.json", err))
		os.Exit(1)
	}

	if err := ioutil.WriteFile("package.json", buf.Bytes(), 0644); err != nil {
		loggers.Stderr.Println(errs.WriteFile("package.json", err))
		os.Exit(1)
	}

	if r.CreateCommand.Directory == "." {
		loggers.Stdout.Println(`Successfully created a new Retro app.

` + color.Bold("# npm") + `

	1. npm
	2. npm run watch

` + color.Bold("# yarn") + `

	1. yarn
	2. yarn watch

Happy hacking!`)
	} else {
		loggers.Stdout.Println(`Successfully created a new Retro app.

` + color.Bold("# npm") + `

	1. cd ` + r.CreateCommand.Directory + `
	2. npm
	3. npm run watch

` + color.Bold("# yarn") + `

	1. cd ` + r.CreateCommand.Directory + `
	2. yarn
	3. yarn watch

Happy hacking!`)
	}
}
