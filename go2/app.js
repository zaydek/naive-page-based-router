import React from "react"
import ReactDOMServer from "react-dom/server"

async function awaitRun() {
	const { load, head: Head } = require("%RETRO_STATIC_PATH%")
	const loadProps = await load()
	const head = ReactDOMServer.renderToStaticMarkup(<Head {...loadProps} />)
	console.log(JSON.stringify({ loadProps, head }))
}

;(async () => {
	await awaitRun()
})()
