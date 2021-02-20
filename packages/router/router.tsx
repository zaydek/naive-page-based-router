import * as React from "react"
import * as store from "../store"
import * as types from "./types"
import * as utils from "./utils"

type RouteEventType = "PUSH" | "REPLACE"

interface RouterState {
	type: RouteEventType
	path: string
	scrollTo?: types.ScrollTo
}

// routerStore describes the router as a store. The router can change the
// current path (window.location.pathname) as a pushState or replaceState event,
// and scrollTo to a number or [number, number].
//
// TODO: Implement synthetic pushState and replaceState functions?
const routerStore = store.createStore<RouterState>({
	// TODO: Add support for key.
	path: utils.getCurrentPath(),
	type: "PUSH",
	scrollTo: [0, 0],
})

export const Link: types.Link = ({ path, scrollTo, children, ...props }) => {
	const setRouter = store.useStoreSetState(routerStore)

	function handleClick(e: React.MouseEvent) {
		e.preventDefault()
		setRouter({ type: "PUSH", path, scrollTo })
	}

	const scoped = path.startsWith("/")
	return (
		// prettier-ignore
		<a href={path} target={scoped ? undefined : "_blank"} rel={scoped ? undefined : "noreferrer noopener"}
				onClick={scoped ? handleClick : undefined} {...props}>
			{children}
		</a>
	)
}

export const Route: types.Route = ({ children }) => {
	return <>{children}</>
}

export const Router: types.Router = ({ children }) => {
	const [router, setRouter] = store.useStore(routerStore)

	React.useEffect(() => {
		function handlePopState() {
			setRouter({
				type: "REPLACE",
				path: utils.getCurrentPath(),
				scrollTo: [0, 0],
			})
		}
		window.addEventListener("popstate", handlePopState)
		return () => window.removeEventListener("popstate", handlePopState)
	}, [])

	let onceRef = React.useRef(false)
	React.useEffect(() => {
		if (!onceRef.current) {
			onceRef.current = true
			return
		}
		// window.pushState / window.replaceState:
		let { path, scrollTo } = router
		if (path !== utils.getCurrentPath()) {
			let emitHistoryEvent: Function
			if (router.type === "PUSH") {
				emitHistoryEvent = () => window.history.pushState({}, "", path)
			} else if (router.type === "REPLACE") {
				emitHistoryEvent = () => window.history.replaceState({}, "", path)
			}
			emitHistoryEvent!()
		}
		// window.scrollTo:
		if (scrollTo !== "no-op") {
			scrollTo = scrollTo || 0
			window.scrollTo(...(scrollTo as any))
		}
	}, [router])

	// cachedRoutes caches routes so rerenders are O(1).
	const cachedRoutes = React.useMemo(() => {
		type RouteMap = { [key: string]: React.ReactElement<types.Route> }

		// routeMap maps paths to components.
		const routeMap: RouteMap = {}
		React.Children.forEach(children, child => {
			if (!React.isValidElement(child)) return
			if (child?.type === Route && !!child?.props?.path) {
				routeMap[child.props.path] = child
			}
		})
		return routeMap
	}, [children])

	// Match the current route or the "/404" route or undefined:
	const route = cachedRoutes[router.path] || cachedRoutes["/404"]
	return <>{route}</>
}
