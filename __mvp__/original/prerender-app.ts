// import seedHash from "./seedHash"

import conf from "./conf"
import fs from "fs"
import path from "path"
import { buildSync } from "esbuild"
import { detab } from "../utils"
import { getPageSrcs } from "./utils"
import { parseRoutes } from "../Router/parts"

const App = require("../" + conf.PAGES_DIR + "/internal/app.tsx").default // FIXME: Change `/` for COMPAT

const srcs = getPageSrcs()

// prettier-ignore
const routes = srcs
	.map(each => path.parse(each).name)      // Pages
	.map(each => parseRoutes("/" + each)) // RouteInfo

function run() {
	// prettier-ignore
	const app = `
// THIS FILE IS AUTOGENERATED.
// THESE AREN’T THE FILES YOU’RE LOOKING FOR. MOVE ALONG.

import React from "react"
import ReactDOM from "react-dom"
import { Route, Router } from "../Router"

// App
${!App ? "/* No-op */" : `import App from ${JSON.stringify("../" + conf.PAGES_DIR + "/internal/app") /* FIXME: Change `/` for COMPAT */}` }

// Pages
${routes.map(each =>
	`import ${each!.component} from ${JSON.stringify("../" + conf.PAGES_DIR + each!.page) /* FIXME: Change `/` for COMPAT */}`
).join("\n")}

// Page props
import pageProps from "./pageProps"

export default function RoutedApp() {
	return (
		<Router>
			${routes.map(each =>
				!App
					? `
			<Route page=${JSON.stringify(each!.page)}>
				<${each!.component} {...pageProps[${JSON.stringify(each!.page)}]} />
			</Route>`
					: `
			<Route page=${JSON.stringify(each!.page)}>
				<App {...pageProps[${JSON.stringify(each!.page)}]}>
					<${each!.component} {...pageProps[${JSON.stringify(each!.page)}]} />
				</App>
			</Route>`
			).join("\n")}

		</Router>
	)
}

${
	!conf.STRICT_MODE
		? detab(`
			ReactDOM.hydrate(
				<RoutedApp />,
				document.getElementById("root"),
			)`)
		: detab(`
			ReactDOM.hydrate(
				<React.StrictMode>
					<RoutedApp />
				</React.StrictMode>,
				document.getElementById("root"),
			)`)
}`.trimStart()

	fs.writeFileSync(conf.CACHE_DIR + "/app.js", app + "\n") // FIXME: Change `/` for COMPAT

	buildSync({
		bundle: true,
		define: {
			__DEV__: JSON.stringify(conf.__DEV__),
			"process.env.NODE_ENV": JSON.stringify(conf.NODE_ENV),
		},
		entryPoints: [conf.CACHE_DIR + "/app.js"],
		loader: { ".js": "jsx" },
		minify: !conf.__DEV__,
		outfile: conf.BUILD_DIR + "/app.js",
	})
}

;(() => {
	run()
})()
