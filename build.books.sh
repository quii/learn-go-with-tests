#!/usr/bin/env bash

set -e

# safer separator for sed
sep=$'\001'

if [ -v GITHUB_REF_NAME ]; then
    sed "s${sep}%%FOOTER_VERSION%%${sep}${GITHUB_REF_NAME}${sep}" meta.tmpl.tex > meta.tex
else
    sed "s${sep}%%FOOTER_VERSION%%${sep}UNDEFINED VERSION${sep}" meta.tmpl.tex > meta.tex
fi

docker run --rm -v `pwd`:/data uppalabharath/pandoc-latex-cjk:latest --from=gfm+rebase_relative_paths -o learn-go-with-tests.pdf \
    -H meta.tex --pdf-engine=xelatex --variable urlcolor=blue --toc --toc-depth=1 \
    -B pdf-cover.tex --lua-filter=fix-internal-links.lua \
    gb-readme.md \
    meta/why.md \
    go-fundamentals/hello-world.md \
    go-fundamentals/integers.md \
    go-fundamentals/iteration.md \
    go-fundamentals/arrays-and-slices.md \
    go-fundamentals/structs-methods-and-interfaces.md \
    go-fundamentals/pointers-and-errors.md \
    go-fundamentals/maps.md \
    go-fundamentals/dependency-injection.md \
    go-fundamentals/mocking.md \
    go-fundamentals/concurrency.md \
    go-fundamentals/select.md \
    go-fundamentals/reflection.md \
    go-fundamentals/sync.md \
    go-fundamentals/context.md \
    go-fundamentals/roman-numerals.md \
    go-fundamentals/math.md \
    go-fundamentals/reading-files.md \
    go-fundamentals/html-templates.md \
    go-fundamentals/generics.md \
    go-fundamentals/revisiting-arrays-and-slices-with-generics.md \
    testing-fundamentals/intro-to-acceptance-tests.md \
    testing-fundamentals/scaling-acceptance-tests.md \
    testing-fundamentals/working-without-mocks.md \
    testing-fundamentals/refactoring-checklist.md \
    build-an-application/app-intro.md \
    build-an-application/http-server.md \
    build-an-application/json.md \
    build-an-application/io.md \
    build-an-application/command-line.md \
    build-an-application/time.md \
    build-an-application/websockets.md \
    questions-and-answers/os-exec.md \
    questions-and-answers/error-types.md \
    questions-and-answers/context-aware-reader.md \
    questions-and-answers/http-handlers-revisited.md \
    meta/anti-patterns.md

docker run --rm -v `pwd`:/data pandoc/latex:latest --from=gfm+rebase_relative_paths --to=epub --file-scope title.txt -o learn-go-with-tests.epub --pdf-engine=xelatex --toc --toc-depth=1 --lua-filter=fix-internal-links.lua \
    gb-readme.md \
    meta/why.md \
    go-fundamentals/hello-world.md \
    go-fundamentals/integers.md \
    go-fundamentals/iteration.md \
    go-fundamentals/arrays-and-slices.md \
    go-fundamentals/structs-methods-and-interfaces.md \
    go-fundamentals/pointers-and-errors.md \
    go-fundamentals/maps.md \
    go-fundamentals/dependency-injection.md \
    go-fundamentals/mocking.md \
    go-fundamentals/concurrency.md \
    go-fundamentals/select.md \
    go-fundamentals/reflection.md \
    go-fundamentals/sync.md \
    go-fundamentals/context.md \
    go-fundamentals/roman-numerals.md \
    go-fundamentals/math.md \
    go-fundamentals/reading-files.md \
    go-fundamentals/html-templates.md \
    go-fundamentals/generics.md \
    go-fundamentals/revisiting-arrays-and-slices-with-generics.md \
    testing-fundamentals/intro-to-acceptance-tests.md \
    testing-fundamentals/scaling-acceptance-tests.md \
    testing-fundamentals/working-without-mocks.md \
    testing-fundamentals/refactoring-checklist.md \
    build-an-application/app-intro.md \
    build-an-application/http-server.md \
    build-an-application/json.md \
    build-an-application/io.md \
    build-an-application/command-line.md \
    build-an-application/time.md \
    build-an-application/websockets.md \
    questions-and-answers/os-exec.md \
    questions-and-answers/error-types.md \
    questions-and-answers/context-aware-reader.md \
    questions-and-answers/http-handlers-revisited.md \
    meta/anti-patterns.md
