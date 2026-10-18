-- Pandoc Lua filter: rewrite cross-chapter relative markdown links
-- (e.g. "pointers-and-errors.md" or "iteration.md#benchmarking") into the
-- internal anchor IDs pandoc generates when combining multiple --file-scope
-- markdown files into a single PDF/EPUB.
--
-- This only runs at build time (see build.books.sh) and never touches the
-- markdown source, so links continue to work as plain relative .md links on
-- GitBook, which resolves them independently.
--
-- Pandoc's --file-scope combines each source file's blocks into one document,
-- with header identifiers prefixed as "{filenameWithDotsStripped}__{slug}" to
-- avoid collisions, e.g. pointers-and-errors.md's first heading becomes
-- "pointers-and-errorsmd__pointers--errors". Pandoc's reader already resolves
-- links that include an explicit "#fragment" (e.g. "iteration.md#benchmarking")
-- correctly, including adding the right split-chapter filename. But links to
-- just a bare filename (e.g. "pointers-and-errors.md", no fragment) come out
-- two different ways depending on context:
--   1. left untouched as a literal ".md" string (pandoc has no target
--      fragment to resolve it to), or
--   2. half-resolved to "#{prefix}" (e.g. "#pointers-and-errorsmd") - a
--      dangling id, since only "{prefix}__{slug}" ids actually exist.
-- This filter fixes both by pointing them at that chapter's first (title)
-- heading id, and leaves anything it can't confidently resolve untouched
-- (logged as a warning) rather than guessing.
--
-- Written using the classic two-pass "list of filters" style (rather than a
-- single Pandoc(doc) filter using newer Blocks:walk() methods) since this
-- needs to run on both a modern pandoc (epub build) and an older pandoc
-- 2.16 (PDF build), which don't share the same Lua filter API.

-- The markdown files passed to pandoc for the PDF/EPUB build (build.books.sh).
-- Used to whitelist which "#{prefix}" / "{file}.md" links are genuinely
-- internal chapter references, so we don't misfire on something that merely
-- looks similar.
local knownFiles = {
  "gb-readme.md", "why.md", "hello-world.md", "integers.md", "iteration.md",
  "arrays-and-slices.md", "structs-methods-and-interfaces.md",
  "pointers-and-errors.md", "maps.md", "dependency-injection.md", "mocking.md",
  "concurrency.md", "select.md", "reflection.md", "sync.md", "context.md",
  "roman-numerals.md", "math.md", "reading-files.md", "html-templates.md",
  "generics.md", "revisiting-arrays-and-slices-with-generics.md",
  "intro-to-acceptance-tests.md", "scaling-acceptance-tests.md",
  "working-without-mocks.md", "refactoring-checklist.md", "app-intro.md",
  "http-server.md", "json.md", "io.md", "command-line.md", "time.md",
  "websockets.md", "os-exec.md", "error-types.md", "context-aware-reader.md",
  "http-handlers-revisited.md", "anti-patterns.md",
}

local function prefixFor(mdFilename)
  -- pandoc's prefix is the filename with dots removed, e.g.
  -- "pointers-and-errors.md" -> "pointers-and-errorsmd"
  return (mdFilename:gsub("%.", ""))
end

local knownPrefixes = {}
for _, f in ipairs(knownFiles) do
  knownPrefixes[prefixFor(f)] = true
end

local function fileBasename(path)
  return path:match("([^/]+)$") or path
end

-- Files that are genuinely linked from the book but aren't among the
-- chapters passed to pandoc, so there's no internal anchor to point at.
-- Redirect these to their GitHub-hosted copy instead of leaving a dangling
-- reference to a file that isn't in the EPUB/PDF at all.
local externalOverrides = {
  ["LICENSE.md"] = "https://github.com/quii/learn-go-with-tests/blob/main/LICENSE.md",
}

-- Case A: a literal relative link straight to one of our markdown files,
-- e.g. "pointers-and-errors.md" or "iteration.md#benchmarking", possibly
-- with a leading "./" or "../".
local function asPlainMdLink(target)
  if target:match("^https?://") or target:match("^mailto:") or target:match("^#") then
    return nil
  end
  local filePart, fragment = target:match("^([^#]+)#?(.*)$")
  if not filePart then
    return nil
  end
  local basename = fileBasename(filePart)
  if not knownPrefixes[prefixFor(basename)] then
    return nil
  end
  return prefixFor(basename), fragment
end

-- Case B: pandoc's own half-resolved dangling id, e.g. "#pointers-and-errorsmd".
local function asDanglingPrefix(target)
  local prefix = target:match("^#(.+)$")
  if prefix and knownPrefixes[prefix] and not prefix:match("__") then
    return prefix
  end
  return nil
end

local firstHeaderByPrefix = {}
local allIds = {}
local rewritten = 0
local unmatched = {}

local function collectHeader(el)
  if el.identifier and el.identifier ~= "" then
    allIds[el.identifier] = true
    local prefix = el.identifier:match("^(.-md)__.+$")
    if prefix and not firstHeaderByPrefix[prefix] then
      firstHeaderByPrefix[prefix] = el.identifier
    end
  end
  return nil
end

local function resolve(prefix, fragment)
  if fragment and fragment ~= "" then
    local candidate = prefix .. "__" .. fragment
    if allIds[candidate] then
      return candidate
    end
  end
  return firstHeaderByPrefix[prefix]
end

local function fixLink(el)
  local override = externalOverrides[fileBasename(el.target)]
  if override then
    el.target = override
    rewritten = rewritten + 1
    return el
  end

  local prefix, fragment = asPlainMdLink(el.target)
  if not prefix then
    prefix = asDanglingPrefix(el.target)
    fragment = nil
  end
  if not prefix then
    return nil
  end

  local newId = resolve(prefix, fragment)
  if newId then
    el.target = "#" .. newId
    rewritten = rewritten + 1
    return el
  end

  table.insert(unmatched, el.target)
  return nil
end

local function report(doc)
  io.stderr:write("fix-internal-links.lua: rewrote " .. rewritten .. " internal links\n")
  if #unmatched > 0 then
    io.stderr:write("fix-internal-links.lua: WARNING, could not resolve " .. #unmatched .. " link(s):\n")
    for _, t in ipairs(unmatched) do
      io.stderr:write("  " .. t .. "\n")
    end
  end
  return doc
end

-- Two full passes: collect every header id first, then fix links, since a
-- link earlier in the document can point at a header that appears later.
return {
  { Header = collectHeader },
  { Link = fixLink, Pandoc = report },
}
