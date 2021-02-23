import chalk from "chalk"

// format converts tabs to spaces.
function format(...args: unknown[]): string {
	if (args.length === 1 && args[0] instanceof Error) {
		return format(args[0].message)
	}

	return args
		.join(" ")
		.split("\n")
		.map((each, x) => {
			if (x === 0) return each
			if (each === "") return each
			return " ".repeat(2) + each.replace("\t", "  ")
		})
		.join("\n")
}

// "> OK: ..."
export function info(...args: unknown[]): void {
	const message = format(...args)
	console.log(`${" ".repeat(2)}${chalk.bold(">")} ${chalk.bold.green("OK:")} ${chalk.bold(message)}`)
	console.log() // "\n"
}

// > Error: ..."
export function error(...args: unknown[]): void {
	const message = format(...args)
	const traceEnabled = process.env["STACK_TRACE"] === "true"
	if (!traceEnabled) {
		console.error(`${" ".repeat(2)}${chalk.bold(">")} ${chalk.bold.red("Error:")} ${chalk.bold(message)}`)
		console.error()
	} else {
		console.error(`${" ".repeat(2)}${chalk.bold(">")} ${chalk.bold.red("Error:")} ${chalk.bold(message)}`)
		console.error()
		console.error({ error })
	}
	process.exit(0)
}
