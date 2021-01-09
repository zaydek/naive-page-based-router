import fs from "fs"
import { getPageSrcs, guards } from "./utils"

// prettier-ignore
interface PageProps {
	basename: string
	props:    any
}

type PagePropsMap = { [key: string]: PageProps }

async function run() {
	guards()

	const ps = []

	// Asynchronously prerender page props.
	const srcs = getPageSrcs()
	for (const each of srcs) {
		const p = new Promise<PageProps>(async resolve => {
			const basename = each.replace(/\.tsx$/, "")

			const { load } = require("../pages/" + each)
			let props = null
			if (load) {
				props = await load()
			}
			resolve({ basename, props })
		})
		ps.push(p)
	}

	// Convert response from an array of to a map:
	const propsArr = await Promise.all(ps)
	const propsMap = propsArr.reduce((acc, each) => {
		acc[each.basename] = each.props
		return acc
	}, {} as PagePropsMap)

	fs.writeFileSync("cache/__props.json", JSON.stringify(propsMap, null, "\t"))
}

;(async () => {
	await run()
})()