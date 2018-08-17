# Web assembly

A common criticism of TDD is it is too hard when you dont know exactly how something works. As of writing I know almost nothing about Web Assembly. In these instances it is totally acceptable to do a "spike"

## Spikes

A spike is simply where you do short experiments to figure out how to implement something. You can write the code as awful as you like, dont worry about coming up with something production ready

The idea is to learn what you don't know. Once you've done this it is important to **throw the code away** and then write the code again "properly" (i.e, we will TDD it). 

## Getting started

To 

- `go get golang.org/dl/go1.11rc1`
- `go1.11rc1 download`

## Hello, world

In order for this to work, we need two applications

1. The web assembly artifact. This is where we will write our web assembly code and we will use the go compiler with some special flags to output a `.wasm` file. This will live in the root of our project as it is our "library" code.
2. A simple web server that serves our html, js and wasm assets. We need a JS file that is included in Go to bootstrap the wasm into the browser for it to run. This will live inside `/cmd/web`

### The web server

Inside `cmd/web/main.go` add the following

```go
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

```

This is just a normal Go web server, no special build instructions required. Just kick it off and leave it running. 

In order for browsers to interpret our wasm file as wasm we need to set a `content-type` header in the response.

You will need to add `index.html` and `wasm_exec.js` which you can fetch with

`cp $$(go1.11rc1 env GOROOT)/misc/wasm/wasm_exec.html cmd/web/html/index.html`
`cp $$(go1.11rc1 env GOROOT)/misc/wasm/wasm_exec.js cmd/web/html/wasm_exec.js`

### Building WASM

In main.go add the following

```go
// +build js,wasm

package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}

```

The comment at the top is a _build tag_

> They mean your file will only be built when you use `GOOS=js GOARCH=wasm`

- Johan Brandhorst

Run `GOOS=js GOARCH=wasm go1.11rc1 build -o cmd/web/html/test.wasm main.go`

This will build `test.wasm` from this code copying it into our web server. `fmt.Println` just maps to `console.log`.

Assuming you're running the web server if you go to `http://localhost:8080` you should see a `run` button. Click it and inspect the console. 

## DOM manipulation

This is cool, but in order to do anything interesting we'll want to interact with the DOM to manipulate elements in an interesting way.

In the HTML file, add 

In our wasm main, change the code to the following

## Links

- [https://www.youtube.com/watch?v=iTrx0BbUXI4](Get going with WebAssembly)
- [https://grpcweb.jbrandhorst.com/](Johan Brandhorst's GopherJS gRPC-Web Client Examples)
- [https://brianketelsen.com/web-assembly-and-go-a-look-to-the-future](Brian Ketelsen's Web Assembly and Go: A look to the future)

## Notes

- Web standard, cross browser support. 
- Intended for "fast numerical calculations"
