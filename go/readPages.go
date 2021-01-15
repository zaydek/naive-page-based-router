package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

// This service is responsible for resolving bytes for `build/page.html`.
func ReadPages(config Configuration, router PageBasedRouter) ([]byte, error) {
	// var buf bytes.Buffer

	dot := struct {
		Config Configuration   `json:"config"`
		Router PageBasedRouter `json:"router"`
	}{Config: config, Router: router}

	dotBytes, err := json.MarshalIndent(dot, "", "\t")
	if err != nil {
		return nil, err
	}

	stdout, stderr, err := execcmd("yarn", "-s", "ts-node", "-T", "go/services/pages.tsx", string(dotBytes))
	if stderr != nil { // Takes precedence
		return nil, errors.New(string(stderr))
	} else if err != nil {
		return nil, err
	}

	fmt.Println(string(stdout))

	return nil, nil

	// 	buf.Write([]byte(`
	// // THIS FILE IS AUTO-GENERATED.
	// // THESE AREN’T THE FILES YOU’RE LOOKING FOR.
	// // MOVE ALONG.
	//
	// module.exports = `))
	// 	buf.Write(stdout)
	//
	// 	contents := buf.Bytes()
	// 	contents = bytes.TrimLeft(contents, "\n") // Remove BOF
	//
	// return contents, nil

	//	var buf bytes.Buffer
	//
	//	dot := struct {
	//		Config Configuration   `json:"config"`
	//		Router PageBasedRouter `json:"router"`
	//	}{Config: config, Router: router}
	//
	//	const data = `
	//<!DOCTYPE html>
	//<html lang="en">
	//	<head>
	//		<meta charset="utf-8">
	//		<meta name="viewport" content="width=device-width, initial-scale=1">
	//		{{if .Head}}{{.Head}}{{end}}
	//	</head>
	//	<body>
	//		<noscript>You need to enable JavaScript to run this app.</noscript>
	//		<div id="root">{{if .Root}}{{.Root}}{{end}}</div>
	//		<script src="/app.js"></script>
	//	</body>
	//</html>
	//`
	//	tmpl := template.Must(template.New("").Parse(data))
	//	err := tmpl.Execute(&buf, dot)
	//	if err != nil {
	//		return nil, err
	//	}
	//	contents := bytes.TrimLeft(buf.Bytes(), "\n") // Remove BOF
	//	return contents, nil
}
