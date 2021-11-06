#!/usr/bin/env bash

set -e

#docker run -v `pwd`:/source jagregory/pandoc -o learn-go-with-tests.pdf -H meta.tex --latex-engine=xelatex --variable urlcolor=blue --toc --toc-depth=1 pdf-cover.md \
#    gb-readme.md \
#    why.md \
#    hello-world.md \
#    integers.md \
#    arrays-and-slices.md \
#    structs-methods-and-interfaces.md \
#    pointers-and-errors.md \
#    maps.md \
#    dependency-injection.md \
#    mocking.md \
#    concurrency.md \
#    select.md \
#    reflection.md \
#    sync.md \
#    context.md \
#    roman-numerals.md \
#    math.md \
#    app-intro.md \
#    http-server.md \
#    json.md \
#    io.md \
#    command-line.md \
#    time.md \
#    websockets.md \
#    os-exec.md \
#    error-types.md \

docker run --rm -v `pwd`:/data pandoc/latex --from=gfm --to=epub --file-scope title.txt -o learn-go-with-tests.epub --pdf-engine=xelatex --toc --toc-depth=1  \
    gb-readme.md \
    why.md \
    hello-world.md \
    integers.md \
    iteration.md \
    arrays-and-slices.md \
    structs-methods-and-interfaces.md \
    pointers-and-errors.md \
    maps.md \
    dependency-injection.md \
    mocking.md \
    concurrency.md \
    select.md \
    reflection.md \
    sync.md \
    context.md \
    roman-numerals.md \
    math.md \
    reading-files.md \
    intro-to-generics.md \
    app-intro.md \
    http-server.md \
    json.md \
    io.md \
    command-line.md \
    time.md \
    websockets.md \
    os-exec.md \
    error-types.md \
    context-aware-reader.md \
    http-handlers-revisited.md \
    anti-patterns.md
