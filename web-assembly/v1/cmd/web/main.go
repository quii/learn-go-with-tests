package main

import (
	"log"
	"net/http"
)

func wasmFileServer() http.HandlerFunc {
	fs := http.FileServer(http.Dir("./html"))

	return func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/test.wasm" {
			res.Header().Set("content-type", "application/wasm")
		}

		fs.ServeHTTP(res, req)
	}
}

func main() {
	log.Print("Serving on http://localhost:8080")
	http.ListenAndServe(":8080", wasmFileServer())
}
