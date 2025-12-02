# Learn Go with Tests - Claude Code Skill

A Claude Code skill for practicing Test-Driven Development (TDD) in Go, based on the principles from [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/).

## What This Skill Provides

When installed, Claude Code will automatically apply TDD best practices when helping you write Go code:

- **TDD Cycle Guidance** - Red-Green-Refactor workflow
- **Go Testing Patterns** - Table tests, subtests, benchmarks, examples
- **Test Design Principles** - Avoiding anti-patterns, listening to your tests
- **Mocking & Dependency Injection** - Stubs, spies, fakes, and when to use each
- **Advanced Testing** - Property-based testing, acceptance tests, contract tests
- **Refactoring Practices** - Safe refactoring within the TDD cycle

## Installation

Copy the skill file to your Claude Code skills directory:

```bash
# For personal use (available in all projects)
mkdir -p ~/.claude/skills/learn-go-with-tests
curl -o ~/.claude/skills/learn-go-with-tests/SKILL.md \
  https://raw.githubusercontent.com/quii/learn-go-with-tests/main/claude-code/skills/learn-go-with-tests/SKILL.md

# Or for a specific project (available only in that project)
mkdir -p .claude/skills/learn-go-with-tests
curl -o .claude/skills/learn-go-with-tests/SKILL.md \
  https://raw.githubusercontent.com/quii/learn-go-with-tests/main/claude-code/skills/learn-go-with-tests/SKILL.md
```

## Usage

Once installed, the skill is automatically available. Claude Code will apply these TDD practices when you:

- Ask to write tests for Go code
- Request help with TDD workflow
- Need guidance on testing patterns
- Want to improve test design

Example prompts that benefit from this skill:

- "Help me write tests for this Go function"
- "What's the TDD approach for implementing a cache?"
- "How should I structure my table-driven tests?"
- "Review my tests and suggest improvements"

## Future: Marketplace Distribution

This skill is structured as a Claude Code plugin for potential future distribution via a [plugin marketplace](https://docs.anthropic.com/en/docs/claude-code/plugins). This would enable installation via `/plugin install` commands.

## Contributing

Contributions to improve the skill content are welcome via pull request.

## License

Same license as the Learn Go with Tests repository.
