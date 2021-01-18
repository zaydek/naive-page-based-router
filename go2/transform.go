package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var filenames = []string{
	"./retro-app/pages/index.js",
	"./retro-app/pages/nested/index.js",
}

func camelCase(filename string) string {
	byteIsLetter := func(b byte) bool {
		ok := ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z')
		return ok
	}

	trim := filename
	trim = strings.TrimPrefix(trim, "./")           // Remove ./etc
	trim = trim[:len(trim)-len(filepath.Ext(trim))] // Remove etc.*

	var ret string
	for x := 0; x < len(trim); x++ {
		switch trim[x] {
		case '/':
			ret += "Slash"
			x++
			if x < len(trim) {
				ret += strings.ToUpper(string(trim[x]))
			}
		case '-':
			x++
			if x < len(trim) && byteIsLetter(trim[x]) {
				ret += strings.ToUpper(string(trim[x]))
			}
		default:
			ret += string(trim[x])
		}
	}

	return ret
}

func pageCase(filename string) string {
	camelCase := camelCase(filename)
	return "Page" + strings.Split(camelCase, "PagesSlash")[1]
}

func main() {
	var requires string
	for x, filename := range filenames {
		var sep string
		if x > 0 {
			sep = "\n"
		}
		requires += sep + fmt.Sprintf("const %s = require(%q)", pageCase(filename), filename)
	}

	var importsAsArr string
	for x, filename := range filenames {
		var sep string
		if x > 0 {
			sep = ", "
		}
		importsAsArr += sep + fmt.Sprintf("{ name: %[1]q, imports: %[1]s }", pageCase(filename))
	}
	importsAsArr = "[" + strings.Join(strings.Split(importsAsArr, "{ "), "\n\t\t{ ") + ",\n\t]"

	js := `import React from "react"
import ReactDOMServer from "react-dom/server"

// Synthetic requires
` + requires + `

async function asyncRun(imports) {
	const chain = []
	for (const each of imports) {
		const p = new Promise(async resolve => {
			const { load, head: Head } = each.imports
			const loadProps = await load()
			const head = ReactDOMServer.renderToStaticMarkup(<Head {...loadProps} />)
			resolve({ name: each.name, loadProps, head })
		})
		chain.push(p)
	}
	const resolvedAsArr = await Promise.all(chain)
	const resolvedAsMap = resolvedAsArr.reduce((acc, each) => {
		acc[each.name] = { ...each, name: undefined }
		return acc
	}, {})
	console.log(JSON.stringify(resolvedAsMap))
}

;(async () => {
	// Synthetic imports array
	await asyncRun(` + importsAsArr + `)
})()
`

	err := ioutil.WriteFile("app2.js", []byte(js), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write file: %w", err))
	}

	// result := api.Build(api.BuildOptions{
	// 	Bundle:      true,
	// 	Define:      map[string]string{"process.env.NODE_ENV": "\"production\""},
	// 	EntryPoints: []string{"app2.js"},
	// 	Loader:      map[string]api.Loader{".js": api.LoaderJSX},
	// })
	// if len(result.Errors) > 0 {
	// 	bstr, _ := json.MarshalIndent(result.Errors, "", "\t")
	// 	fmt.Println(string(bstr))
	// 	os.Exit(1)
	// }
	// fmt.Print(string(result.OutputFiles[0].Contents))
}
