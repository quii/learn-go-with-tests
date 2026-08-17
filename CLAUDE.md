# Learn Go with Tests — writing guide

This is Chris James's ("quii") book. When drafting or editing chapter prose, match his voice, not generic technical-writing tone. This file is a distilled style guide; when in doubt, go read a chapter he actually wrote (see exclusions below) rather than guessing.

## Voice, in one paragraph

Casual, collaborative, "we're solving this together" — heavy on "we"/"let's", second person "you" for instructions. Contractions everywhere. Short-to-medium sentences; paragraphs before a code block are almost always 1–3 sentences. Dry, understated humour lands as single-word or single-line beats ("Yuck.", "Ouch.", "Perfect!", "As expected") rather than running jokes — it's seasoning, not the dish. Opinions are stated with real conviction but usually hedged with "I feel", "in my experience", or an immediate "but/however" qualifier; alternatives are acknowledged rather than dismissed. Blunt, unhedged claims are rare and therefore load-bearing — don't spend that currency casually.

## Structure

- Tutorial chapters open with a concrete scenario (a product owner, a colleague, a Slack/Reddit question) *before* any code, and the first line is almost always a bolded link: `**[You can find all the code for this chapter here](...)**`.
- New concepts are motivated by a failing test or a compiler error wherever possible, not front-loaded as "here is concept X" — except when introducing an entire new stdlib package's mental model (context, sync, select, synctest), which typically gets a short upfront "just enough information" section since there's no single failing test that can teach it.
- The TDD ritual headings recur near-verbatim: `## Write the test first`, `## Try to run the test`, `## Write the minimal amount of code for the test to run and check the failing test output`, `## Write enough code to make it pass`, `## Refactor`. Steps can be merged when one is trivial (e.g. skip the "minimal code" step if the real implementation is one small step away), but don't rename them.
- Chapters close with `## Wrapping up`, usually a bulleted "what we've covered", sometimes split into named sub-lists. An `### Additional material` (or "Further reading") section with links is common at the very end.
- Failures (compiler errors, panics, race detector output) are shown verbatim, then narrated calmly afterwards — the error is treated as a useful signal ("listen to the compiler"), not an obstacle. Reproduce real output by actually running the command; don't hand-write plausible-looking error text.
- Asides and caveats are inline (parentheticals, a one-line "Note:") or, if substantial, get their own end-of-chapter section. Blockquotes (`>`) are reserved for quoting external sources (Kent Beck, Go docs, etc.) — never for his own voice.

## Prose smells to avoid

- Long, hedged, multi-clause sentences stacking qualifiers ("it's worth being precise about what that means, because it's easy to over-generalise from..."). Cut it into two or three shorter sentences.
- Dense paragraphs (4+ sentences) explaining a mechanism with zero personality or paragraph breaks. If you find yourself writing an unbroken technical essay paragraph, look for the moment that deserves a one-line reaction beat instead.
- Third-person distance ("one might consider...", "it can be observed that..."). Use "we"/"you".
- Explaining a concept exhaustively before showing any code. Get to code fast, then explain what surprised you.

## Housekeeping specific to this repo

- The book is a single Go module (`go.mod` at the repo root) — chapter code usually lives in its own top-level directory (`time/v3`, `select/v1`, `synctest/v1`, etc.), each an independent package so old chapters keep compiling untouched even as later ones evolve the same domain differently.
- Before treating any chapter's code as finished: `go build`, `go test ./...`, `go vet ./...`, and where the chapter is about concurrency, run the relevant test with `-race` — repeatedly, since races don't always trigger on the first run. Also run `gofmt -l` on the code and `misspell -error` + `gofmtmd -r` (both already installed at `~/go/bin`) on the `.md` file itself, matching what `build.sh` does in CI.
- Don't hand-write compiler errors, panic output, or race detector traces in prose — reproduce them for real (a scratch directory or the actual chapter code) and paste the real output, then trim only if truly excessive. This book's existing race-detector example (`concurrency.md`) does exactly this.
- `maps.md` and `math.md` are substantially community-authored (by hackeryarn and Ruth Baker respectively, not Chris) — don't use them as a voice reference, and be aware their tone may not match the rest of the book.
- `concurrency.md` may also be partially co-authored (its example code references `github.com/gypsydave5`) — treat it as a slightly less certain voice reference than the rest.
- When Go's stdlib gains something new the book should cover (as with `testing/synctest` needing Go 1.25), check `.github/workflows/go.yml` and `go.mod`'s `go` directive — bumping them is a repo-wide change with real blast radius (every chapter's code must still build/test on the new version), so verify the full suite (`go test ./...`, `go vet ./...`) before and after, and flag the version bump explicitly rather than doing it silently.
