import parseParts from "./parseParts"
import { RouteInfo } from "./types"
import { toTitleCase } from "../utils"

function index(pageStr: string) {
	if (pageStr.endsWith("/index")) {
		return pageStr.replace(/\/index$/, "/")
	}
	return pageStr
}

// Ex:
//
// parseRoutes("/[hello]/[world]")
//
// -> {
// ->   page: "/[hello]/[world]",
// ->   component: "PageDynamicHelloSlashDynamicWorld",
// -> }
//
export default function parseRoutes(partsStr: string) {
	const parts = parseParts(partsStr)
	if (!parts) {
		return null
	}
	const componentStr = parts
		.map(each => {
			let str = ""
			if (each.dynamic) {
				str += "Dynamic"
			}
			str += toTitleCase(!each.dynamic ? each.part : each.part.slice(1, -1))
			if (each.nests) {
				str += "Slash"
			}
			return str
		})
		.join("")

	const info: RouteInfo = {
		page: index(partsStr),
		component: "Page" + componentStr,
	}
	return info
}
