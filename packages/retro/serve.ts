import * as constants from "constants"
import * as fs from "fs"
import * as http from "http"
import * as p from "path"

const PORT = 3000

// convertToFilesystemPath converts a browser path to a filesystem path.
export function convertToFilesystemPath(path: string) {
	// "/" -> "/index"
	let path2 = path
	if (path2.endsWith("/")) {
		path2 += "index"
	}
	// "/index" -> "/index.html"
	if (p.extname(path2) === "") {
		path2 += ".html"
	}
	return path2
}

// This implementation is loosely based on https://stackoverflow.com/a/44188852.
function serve() {
	const server = http.createServer(async (req, res) => {
		if (req.url! === "/favicon.ico") {
			res.writeHead(204)
			return
		}

		// Convert the browser path to a filesystem path:
		req.url = convertToFilesystemPath(req.url!)

		let bytes: Buffer
		try {
			const path = p.join(process.cwd(), req.url!)
			bytes = await fs.promises.readFile(path)
		} catch (err) {
			if (err.code === constants.ENOENT) {
				res.writeHead(404)
				res.end(http.STATUS_CODES[404])
				return
			} else {
				res.writeHead(500)
				res.end(http.STATUS_CODES[500])
				return
			}
		}
		// Done:
		res.writeHead(200)
		res.end(bytes!)
	})

	setTimeout(() => {
		console.log(`📡 http://localhost:${PORT}`)
	}, 100)
	server.listen(PORT)
}

serve()
