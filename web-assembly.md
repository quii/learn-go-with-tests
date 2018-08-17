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

1. The web assembly artifact. This is where we will write our web assembly code and we will use the go compiler with some special flags to output a `.wasm` file.
2. A simple web server that serves our html, js and wasm assets. We need a JS file that is included in Go to bootstrap the wasm into the browser for it to run.

## Links

- [https://www.youtube.com/watch?v=iTrx0BbUXI4](Get going with WebAssembly)
- [https://grpcweb.jbrandhorst.com/](Johan Brandhorst's GopherJS gRPC-Web Client Examples)
- [https://brianketelsen.com/web-assembly-and-go-a-look-to-the-future/](Brian Ketelsen's Web Assembly and Go: A look to the future)

## Notes

- Web standard, cross browser support. 
- Intended for "fast numerical calculations"
- Need to set content type if its wasm
