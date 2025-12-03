---
name: learn-go-with-tests
description: TDD with Go skill for practicing Test-Driven Development. Use when writing tests for Go code, following the Red-Green-Refactor cycle, designing testable code, or improving testing practices.
---

# Test-Driven Development with Go

A comprehensive guide for practicing Test-Driven Development (TDD) in Go, based on the principles from "Learn Go with Tests".

## When to Use This Skill

Use this skill when:
- Practicing Test-Driven Development in Go
- Writing tests for Go code
- Improving your testing practices
- Setting up testing infrastructure for Go projects
- Building robust, well-tested Go systems

**Note**: This skill focuses on **testing techniques and TDD practices** in Go, not on teaching Go language features.

## Core TDD Principles

### The TDD Cycle: Red-Green-Refactor

Follow this strict cycle for all development:

1. **Red**: Write a small test for a small amount of desired behavior
   - See the test fail with a clear, meaningful error message
   - Verify the error message makes sense
   - This proves the test is actually testing something

2. **Green**: Write the minimal amount of code to make the test pass
   - Focus only on making the test pass, nothing more
   - Don't worry about perfect code yet
   - Get to working software quickly

3. **Refactor**: Improve the code while keeping tests green
   - Clean up implementation details
   - Extract functions, add types, improve names
   - **Rule**: Do not change behavior during refactoring
   - Run tests after each small refactor

### Why This Matters

- **Refactoring is NOT changing behavior** - it's restructuring code while maintaining the same behavior
- Tests give you confidence to refactor safely
- Working in small steps prevents going down rabbit holes
- Fast feedback loops (< 1 second) enable flow state
- Tests document the expected behavior for future developers

### Quick Start: Your First TDD Session

If you're new to TDD in Go, start here:

1. **Create test file first**: `mycode_test.go`
2. **Write a failing test**: Define what you want the code to do
3. **Run `go test`**: See it fail (compilation error or assertion failure)
4. **Write minimal code**: Just enough to make test pass
5. **Run `go test`**: See it pass (green!)
6. **Refactor**: Clean up while keeping tests green
7. **Repeat**: Small steps, constant feedback

**Golden rule**: Never write production code without a failing test first.

## Contents Overview

This skill is organized into the following major sections:

1. **Go Testing Fundamentals** - Test structure, helpers, subtests, benchmarks, coverage
2. **TDD Anti-Patterns** - Common pitfalls and how to avoid them
3. **Listening to Your Tests** - Using tests as design feedback
4. **Practical TDD Workflow** - Step-by-step examples and workflow
5. **Testing Best Practices** - Tools, coverage, conventions
6. **Advanced TDD Concepts** - Interfaces, DI, table tests, design principles
7. **Common Go Testing Patterns** - Slices, errors, helpers, HTTP, concurrency
8. **Property-Based Testing** - Testing properties vs examples
9. **Testing File Operations** - Using testing/fstest
10. **Advanced Testing Strategies** - Acceptance tests, test pyramid, fakes, contract tests
11. **Refactoring Principles** - Safe refactoring within TDD cycle

## Go Testing Fundamentals

### Test File Structure

```go
package mypackage

import "testing"

func TestSomething(t *testing.T) {
    got := FunctionToTest()
    want := "expected value"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

### Rules for Go Tests

- Test files must end with `_test.go`
- Test functions must start with `Test`
- Test functions take one argument: `t *testing.T`
- Use `t.Errorf()` for formatted failure messages
- Use meaningful format strings: `%q` for strings, `%d` for ints, `%v` for general values

### Helper Functions

```go
func assertCorrectMessage(t testing.TB, got, want string) {
    t.Helper()  // Makes failures point to the caller, not this helper
    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

- Accept `testing.TB` interface (works with `*testing.T` and `*testing.B`)
- Always call `t.Helper()` to get proper error line numbers

### Subtests

Use subtests to group related test cases:

```go
func TestHello(t *testing.T) {
    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        want := "Hello, Chris"
        assertCorrectMessage(t, got, want)
    })

    t.Run("empty string defaults to 'world'", func(t *testing.T) {
        got := Hello("")
        want := "Hello, World"
        assertCorrectMessage(t, got, want)
    })
}
```

### Table Tests

Table-driven tests are powerful for testing multiple inputs with the same logic:

```go
func TestAdd(t *testing.T) {
    cases := []struct {
        name string
        a, b int
        want int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 5, 5},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := Add(tc.a, tc.b)
            if got != tc.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}
```

**Benefits**:
- Test many scenarios with less code
- Easy to add new test cases
- Clear which case failed

**Warning**: Don't create complicated table tests with many optional fields and booleans. Break into separate tests when scenarios diverge.

**See "Table-Driven Tests for Different Scenarios" in Advanced TDD Concepts section for detailed best practices.**

### Testable Examples

```go
func ExampleAdd() {
    sum := Add(1, 5)
    fmt.Println(sum)
    // Output: 6
}
```

- Examples appear in documentation
- They're compiled and executed as tests
- Use `// Output:` comment to verify output

### Benchmarks

Benchmarking is a first-class feature in Go, similar to testing:

```go
func BenchmarkRepeat(b *testing.B) {
    for b.Loop() {
        Repeat("a")
    }
}
```

**Key concepts**:
- Benchmark functions start with `Benchmark` and take `*testing.B`
- `b.Loop()` returns true while the benchmark should continue running
- The framework automatically determines how many iterations to run
- Only code inside the loop is timed (setup/cleanup excluded)

**Running benchmarks**:
```bash
go test -bench=.              # Run all benchmarks
go test -bench=. -benchmem    # Include memory allocation stats
```

**Understanding output**:
```
BenchmarkRepeat-8    10000000    136 ns/op    8 B/op    1 allocs/op
```
- `10000000` - number of iterations
- `136 ns/op` - average time per operation
- `8 B/op` - bytes allocated per operation (with -benchmem)
- `1 allocs/op` - allocations per operation (with -benchmem)

**Typical benchmark structure**:
```go
func BenchmarkSomething(b *testing.B) {
    // ... setup code (not timed) ...
    for b.Loop() {
        // ... code to measure ...
    }
    // ... cleanup code (not timed) ...
}
```

**Use benchmarks to**:
- Verify performance improvements
- Compare different implementations
- Catch performance regressions
- Guide optimization decisions with data

## TDD Anti-Patterns to Avoid

### 1. Not Doing TDD At All

**Problem**: Writing tests after code leads to:
- Tests that can never fail (evergreen tests)
- Tests coupled to implementation details
- Poor test error messages

**Solution**: Always write the test first, see it fail, then implement.

### 2. Misunderstanding Refactoring Constraints

**Problem**: Thinking you need to write tests before refactoring working code.

**Truth**: When tests are green, you can refactor freely. You're only not allowed to add or change behavior.

### 3. Not Seeing the Test Fail First

**Problem**: Writing test and implementation together, never verifying the test can fail.

**Solution**: Always run the test and see it fail with a meaningful error before writing implementation.

### 4. Useless Assertions

**Problem**: Error messages like "false was not equal to true" that don't explain what broke.

**Solution**: Write descriptive assertions:
```go
// Bad
if !result {
    t.Error("expected true")
}

// Good
if got != want {
    t.Errorf("Hello(%q) = %q, want %q", name, got, want)
}
```

### 5. Testing Implementation Details

**Problem**: Tests coupled to internal structure break when refactoring.

**Example**: Testing that a square is made of triangles when you only care that it's a valid square.

**Solution**: Test behavior through public APIs. Use test packages (`package mypackage_test`) to enforce this.

### 6. Too Many Assertions in One Test

**Problem**: Multiple assertions make tests hard to read and debug.

**Solution**:
- Aim for one assertion per test
- Use subtests to separate different scenarios
- Exception: Integration/acceptance tests where setup is expensive

### 7. Asserting on Irrelevant Details

**Problem**: Asserting on entire complex objects when you only care about one field.

```go
// Bad - tightly coupled to whole object
if !cmp.Equal(complexObject, want) {
    t.Error("got %+v, want %+v", complexObject, want)
}

// Good - specific and loosely coupled
got := complexObject.fieldYouCareAbout
if got != want {
    t.Error("got %q, want %q", got, want)
}
```

### 8. Excessive Test Setup

**Problem**: 20-100 lines of setup with many mocks indicates design problems.

**Solution**: Simplify dependencies. Consider consolidating interfaces. Listen to your tests - complicated tests mean complicated code.

### 9. Violating Encapsulation

**Problem**: Making private functions public just to test them.

**Solution**:
- Test through public APIs only
- Use `package mypackage_test` to enforce this
- If private code needs testing, it might need to be its own package

### 10. Interface Pollution

**Problem**: Large interfaces force users to mock many unused methods.

**Solution**:
- Keep interfaces small (1-3 methods)
- "The bigger the interface, the weaker the abstraction"
- Let consumers define interfaces for their needs
- Export concrete types, not interfaces (like `redis.Client`)

## Listening to Your Tests

### Tests Are Design Feedback

If testing is hard, using the code will be hard. Tests are the first user of your code.

### Common Signals

**Signal**: Lots of mocks/test doubles
**Meaning**: Too many dependencies
**Action**: Consolidate or simplify dependencies

**Signal**: Complex test setup
**Meaning**: Complex code structure
**Action**: Break into smaller, focused units

**Signal**: Tests break when refactoring
**Meaning**: Testing implementation, not behavior
**Action**: Move tests to higher abstraction level

**Signal**: Hard to write test
**Meaning**: Code is hard to use
**Action**: Redesign the API

## Practical TDD Workflow in Go

### Starting a New Package

```bash
# Create directory
mkdir mypackage
cd mypackage

# Initialize module
go mod init example.com/mypackage

# Create test file
touch mypackage_test.go
```

### The Micro-Cycle

1. Write a failing test (< 2 minutes)
2. Run `go test` - see it fail with clear error
3. Write minimal code to pass (< 5 minutes)
4. Run `go test` - see it pass
5. Refactor in small steps, running tests after each change
6. Commit when green and happy

### Example Session: Building an Add Function

**Step 1: Write test first**

```go
// adder_test.go
package integers

import "testing"

func TestAdder(t *testing.T) {
    sum := Add(2, 2)
    expected := 4

    if sum != expected {
        t.Errorf("expected '%d' but got '%d'", expected, sum)
    }
}
```

**Step 2: Run test, see compilation error**

```
$ go test
./adder_test.go:6:9: undefined: Add
```

**Step 3: Write minimal code to compile**

```go
// adder.go
package integers

func Add(x, y int) int {
    return 0
}
```

**Step 4: Run test, see meaningful failure**

```
$ go test
adder_test.go:10: expected '4' but got '0'
```

**Step 5: Make it pass**

```go
func Add(x, y int) int {
    return x + y
}
```

**Step 6: Refactor (add documentation)**

```go
// Add takes two integers and returns the sum of them.
func Add(x, y int) int {
    return x + y
}
```

### Example Session: Building a Repeat Function with Iteration

**Step 1: Write test first**

```go
// repeat_test.go
package iteration

import "testing"

func TestRepeat(t *testing.T) {
    repeated := Repeat("a")
    expected := "aaaaa"

    if repeated != expected {
        t.Errorf("expected %q but got %q", expected, repeated)
    }
}
```

**Step 2: Write minimal code to compile**

```go
// repeat.go
package iteration

func Repeat(character string) string {
    return ""
}
```

**Step 3: Run test, see meaningful failure**

```
repeat_test.go:10: expected 'aaaaa' but got ''
```

**Step 4: Make it pass**

```go
func Repeat(character string) string {
    var repeated string
    for i := 0; i < 5; i++ {
        repeated = repeated + character
    }
    return repeated
}
```

**Step 5: Refactor**

```go
const repeatCount = 5

func Repeat(character string) string {
    var repeated string
    for i := 0; i < repeatCount; i++ {
        repeated += character
    }
    return repeated
}
```

**Step 6: Add benchmark to measure performance**

```go
func BenchmarkRepeat(b *testing.B) {
    for b.Loop() {
        Repeat("a")
    }
}
```

Run with `go test -bench=. -benchmem`

**Step 7: Use benchmark to guide refactoring**

The benchmark shows we can improve performance. Refactor to use `strings.Builder`:

```go
const repeatCount = 5

func Repeat(character string) string {
    var repeated strings.Builder
    for i := 0; i < repeatCount; i++ {
        repeated.WriteString(character)
    }
    return repeated.String()
}
```

**Run benchmark again** to verify the improvement:
- Before: 136 ns/op, 40 B/op, 5 allocs/op
- After: 25 ns/op, 8 B/op, 1 allocs/op

**Testing lesson**: Benchmarks let you refactor for performance with confidence. Without the benchmark, you wouldn't know if your optimization actually helped.

## Testing Best Practices in Go

### Listen to the Compiler During TDD

The compiler is your friend in the TDD cycle. When tests don't compile, the error tells you exactly what to do next in the Red phase.

### Use Go's Built-in Testing Tools

- `go test` - run tests
- `go test -v` - verbose output
- `go test -cover` - show test coverage percentage
- `go test -coverprofile=coverage.out` - generate detailed coverage report
- `go test -bench=.` - run benchmarks
- `go test -bench=. -benchmem` - run benchmarks with memory stats
- `go doc` - view documentation offline
- `pkgsite -open .` - view documentation in browser

### Test Coverage

Go has built-in coverage tooling:

```bash
go test -cover
```

Output:
```
PASS
coverage: 100.0% of statements
```

**Important principles**:
- **100% coverage is not the goal** - confidence is the goal
- If you follow strict TDD, you'll likely have near 100% coverage naturally
- Coverage helps identify untested code, not guarantee quality
- Use coverage to find gaps, not as a success metric

**Every test has a cost**:
- Maintenance overhead
- Execution time
- Cognitive load when reading/understanding

**Avoid redundant tests**. Example:
```go
// Both tests below are redundant - if Sum works for one size, it works for all
t.Run("collection of 5 numbers", func(t *testing.T) {
    got := Sum([]int{1, 2, 3, 4, 5})
    want := 15
    // ...
})

t.Run("collection of 3 numbers", func(t *testing.T) {
    got := Sum([]int{1, 2, 3})
    want := 6
    // ... same logic, no new behavior tested
})
```

Delete one test. Coverage will remain 100% because both tests exercise the same code path.

### Package Organization

- One package per directory
- Test files should use `package mypackage_test` (enforces public API testing)

### Naming Conventions

- Test files: `*_test.go`
- Test functions: `TestXxx(*testing.T)`
- Benchmark functions: `BenchmarkXxx(*testing.B)`
- Example functions: `ExampleXxx()`
- Helper functions: Accept `testing.TB`

### Format Strings

- `%q` - quoted string (great for seeing empty strings, whitespace)
- `%d` - integer
- `%v` - default format (works well for arrays, slices, and most types)
- `%+v` - struct with field names
- `%#v` - Go syntax representation
- `%T` - type

**Pro tip**: Include inputs in error messages to make debugging easier:
```go
if got != want {
    t.Errorf("got %d want %d given, %v", got, want, numbers)
}
```
This shows both the result AND the input that caused the failure.

## Advanced TDD Concepts

### Testing Interfaces

Define interfaces in the consumer, not the producer:

```go
// Consumer defines what it needs
type DataStore interface {
    Save(data Data) error
}

type Handler struct {
    store DataStore
}

// Tests use simple stub
type StubDataStore struct {
    savedData []Data
}

func (s *StubDataStore) Save(data Data) error {
    s.savedData = append(s.savedData, data)
    return nil
}
```

### Dependency Injection

Pass dependencies as parameters, not globals:

```go
// Good - testable
func ProcessData(store DataStore, data Data) error {
    return store.Save(data)
}

// Bad - not testable without global state
func ProcessData(data Data) error {
    return globalStore.Save(data)
}
```

### Table-Driven Tests for Different Scenarios

Table-driven tests let you build a list of test cases that share the same testing logic.

**Basic pattern with anonymous struct**:

```go
func TestArea(t *testing.T) {
    areaTests := []struct {
        name    string
        shape   Shape
        hasArea float64
    }{
        {name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
        {name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
        {name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
    }

    for _, tt := range areaTests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.shape.Area()
            if got != tt.hasArea {
                t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
            }
        })
    }
}
```

**Key testing practices**:

1. **Use named fields for clarity**:
   ```go
   // Hard to understand
   {Rectangle{12, 6}, 72.0}

   // Clear and self-documenting
   {shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0}
   ```

2. **Use t.Run for each case**:
   - Gives clear output: `--- FAIL: TestArea/Rectangle`
   - Can run specific cases: `go test -run TestArea/Rectangle`
   - Shows which case failed without debugging

3. **Use %#v in error messages**:
   ```go
   t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
   // Output: main.Rectangle{Width:12, Height:6} got 72.00 want 72.10
   ```
   Shows the exact struct values that caused the failure.

4. **Use descriptive field names**:
   - `hasArea` is better than `want` - it describes the assertion
   - `name` field documents what each case tests

5. **Tests as assertions of truth**:
   > "The test speaks to us more clearly, as if it were an assertion of truth, not a sequence of operations" - Kent Beck

**When to use table tests**:
- Testing various implementations of an interface
- Same logic with different inputs
- Adding new test cases is easy (just append to slice)

**When NOT to use**:
- Different test cases require different logic
- Would need many optional fields or booleans
- Makes tests harder to read (break into separate tests instead)

## Design Principles from TDD

### Units Should Be:

- **Self-contained**: Few dependencies, clear boundaries
- **Decoupled**: Changes don't ripple through system
- **Coherent**: Centered around one domain concept
- **Simple**: Easy to understand and use

### Good Design Properties:

- Public APIs don't leak implementation details
- Easy to compose units together (like Lego bricks)
- Minimal surface area (small interfaces)
- Clear, documented behavior

### When Design Needs Work:

- Many test doubles needed
- Complex test setup
- Tests fail when refactoring
- Hard to write tests
- Unclear what to test

## Iterative Development

### Small Steps Win

- Don't over-engineer for unknown future
- Build iteratively with tests
- Refactor as you learn about the domain
- Keep cycles tight (< 10 minutes per cycle)

### Managing Complexity (Lehman's Laws)

**Law of Continuous Change**: Software must evolve or become less useful.

**Law of Increasing Complexity**: Complexity increases unless you actively refactor.

**Solution**: TDD + continuous refactoring with test safety net.

## Common Go Testing Patterns

### Comparing Slices and Collections

**Problem**: Go doesn't allow `==` for slices:
```go
got := []int{1, 2, 3}
want := []int{1, 2, 3}
if got == want {  // Compile error: "slice can only be compared to nil"
    // ...
}
```

**Solution 1: reflect.DeepEqual** (any Go version)

```go
import "reflect"

if !reflect.DeepEqual(got, want) {
    t.Errorf("got %v want %v", got, want)
}
```

**Warning**: `reflect.DeepEqual` is **not type-safe**. This compiles but makes no sense:
```go
got := []int{1, 2, 3}
want := "bob"  // Wrong type!
if !reflect.DeepEqual(got, want) {  // Compiles, but logically wrong
    t.Errorf("got %v want %v", got, want)
}
```

**Solution 2: slices.Equal** (Go 1.21+, recommended)

```go
import "slices"

if !slices.Equal(got, want) {
    t.Errorf("got %v want %v", got, want)
}
```

**Benefits of slices.Equal**:
- Type-safe (won't compile if types don't match)
- Simpler, more explicit
- Works with any comparable element type

**Limitation**: Elements must be comparable. Won't work with slices of slices or slices of maps.

### Compile Errors vs Runtime Errors

**Key insight from TDD**:

> Compile time errors are our **friends** - they help us write software that works.
> Runtime errors are our **enemies** - they affect our users.

**Example**:
```go
func TestSumAllTails(t *testing.T) {
    got := SumAllTails([]int{}, []int{3, 4, 5})  // Empty slice!
    want := []int{0, 9}
    // ...
}
```

This test compiles fine but crashes at runtime:
```
panic: runtime error: slice bounds out of range
```

**TDD helps catch runtime errors**:
1. Write test that exercises edge case (empty slice)
2. See runtime panic
3. Fix code to handle edge case
4. Test passes

Now your code won't crash in production.

**Always test edge cases**:
- Empty collections
- Nil values
- Zero values
- Boundary conditions

### Helper Functions in Tests

**Pattern 1: Package-level helper** (traditional)

```go
func assertNoError(t testing.TB, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("got error: %v", err)
    }
}

func TestSomething(t *testing.T) {
    err := DoThing()
    assertNoError(t, err)
}
```

**Pattern 2: Local helper variable** (scoped to one test)

```go
func TestSumAllTails(t *testing.T) {
    // Helper function defined as local variable
    checkSums := func(t testing.TB, got, want []int) {
        t.Helper()
        if !slices.Equal(got, want) {
            t.Errorf("got %v want %v", got, want)
        }
    }

    t.Run("normal case", func(t *testing.T) {
        got := SumAllTails([]int{1, 2}, []int{0, 9})
        want := []int{2, 9}
        checkSums(t, got, want)  // Use the local helper
    })

    t.Run("empty slices", func(t *testing.T) {
        got := SumAllTails([]int{}, []int{3, 4, 5})
        want := []int{0, 9}
        checkSums(t, got, want)
    })
}
```

**Benefits of local helpers**:
- **Cannot be used outside the test** - reduces API surface area
- **Scoped to relevant context** - only available where needed
- **Type safety** - compiler catches misuse:
  ```go
  checkSums(t, got, "dave")  // Compile error!
  // cannot use "dave" (type string) as type []int
  ```
- **Binds to local variables** - can access test-specific state

**When to use each**:
- Package-level: Reusable across many tests
- Local variable: Specific to one test function, benefits from type safety and scoping

### Testing Errors

Go uses explicit error returns. Testing them properly is critical.

**Pattern 1: Test that error is returned**

```go
t.Run("withdraw insufficient funds", func(t *testing.T) {
    wallet := Wallet{balance: Bitcoin(20)}
    err := wallet.Withdraw(Bitcoin(100))

    if err == nil {
        t.Error("wanted an error but didn't get one")
    }
})
```

**Pattern 2: Test error message**

```go
t.Run("withdraw insufficient funds", func(t *testing.T) {
    wallet := Wallet{balance: Bitcoin(20)}
    err := wallet.Withdraw(Bitcoin(100))

    if err == nil {
        t.Fatal("didn't get an error but wanted one")
    }

    if err.Error() != "cannot withdraw, insufficient funds" {
        t.Errorf("got %q, want %q", err, "cannot withdraw, insufficient funds")
    }
})
```

**Problem**: Testing against string messages is brittle and couples tests to implementation.

**Pattern 3: Errors as values (recommended)**

```go
// In production code - define error as package variable
var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {
    if amount > w.balance {
        return ErrInsufficientFunds
    }
    w.balance -= amount
    return nil
}

// In test - compare error values directly
t.Run("withdraw insufficient funds", func(t *testing.T) {
    wallet := Wallet{balance: Bitcoin(20)}
    err := wallet.Withdraw(Bitcoin(100))

    if err != ErrInsufficientFunds {
        t.Errorf("got %q, want %q", err, ErrInsufficientFunds)
    }
})
```

**Benefits of error values**:
- Single source of truth
- Easy to change wording without breaking tests
- Users can check errors with simple equality: `if err == ErrInsufficientFunds`
- Tests focus on behavior, not exact wording

**Pattern 4: Test both success and error cases**

Always test that functions return NO error when they should succeed:

```go
t.Run("withdraw with funds", func(t *testing.T) {
    wallet := Wallet{balance: Bitcoin(20)}
    err := wallet.Withdraw(Bitcoin(10))

    if err != nil {
        t.Fatalf("got an error but didn't want one: %v", err)
    }

    // Continue with other assertions...
})
```

**Using t.Fatal vs t.Error**:

```go
// Use t.Fatal when you can't continue the test
if err == nil {
    t.Fatal("didn't get an error but wanted one")
    // Stops here - can't call err.Error() on nil
}

// Use t.Error when test can continue
if got != want {
    t.Error("values don't match")
    // Test continues
}
```

`t.Fatal` stops test execution immediately. Use it when:
- Checking for nil errors before calling methods on them
- Prerequisites for the rest of the test aren't met
- Continuing would cause a panic

**Helper functions for error testing**:

```go
func assertError(t testing.TB, got, want error) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but wanted one")
    }
    if got != want {
        t.Errorf("got %s, want %s", got, want)
    }
}

func assertNoError(t testing.TB, err error) {
    t.Helper()
    if err != nil {
        t.Fatal("got an error but didn't want one")
    }
}
```

**Unchecked errors linter**:

Install `errcheck` to find unchecked error returns:

```bash
go install github.com/kisielk/errcheck@latest
errcheck .
```

Example output:
```
wallet_test.go:17:18: wallet.Withdraw(Bitcoin(10))
```

This means you're not checking the error from `Withdraw`. Fix it:

```go
// Before (bad)
wallet.Withdraw(Bitcoin(10))

// After (good)
err := wallet.Withdraw(Bitcoin(10))
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
```

### Dependency Injection for Testing

Dependency Injection makes code testable by allowing you to inject controllable dependencies.

**Problem**: Hard to test functions that write to stdout or have hard-coded dependencies:

```go
func Greet(name string) {
    fmt.Printf("Hello, %s", name)  // Prints to stdout - hard to test!
}
```

**Solution**: Inject the Writer interface:

```go
func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}
```

**In tests**: Use `bytes.Buffer` which implements `io.Writer`:

```go
func TestGreet(t *testing.T) {
    buffer := bytes.Buffer{}
    Greet(&buffer, "Chris")

    got := buffer.String()
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

**In production**: Use `os.Stdout`:

```go
func main() {
    Greet(os.Stdout, "Elodie")
}
```

**Benefits**:
- Code becomes testable - can verify output
- Code becomes reusable - works with any `io.Writer`
- Separates concerns - "what to write" vs "where to write"
- No need for mocking frameworks - standard library has what you need

**Key testing insight**: If you can't test a function easily, it usually has dependencies hard-wired into it. DI lets you inject test doubles.

### Mocking and Test Doubles

Mocking lets you replace real dependencies with test versions you can control and inspect.

**When to mock**:
- Testing code that has slow dependencies (time.Sleep, HTTP calls, database)
- Testing error conditions that are hard to trigger
- Verifying interactions between components
- Avoiding fragile tests that depend on external services

**Types of test doubles**:

1. **Stub**: Returns predetermined values
   ```go
   type StubWebsiteChecker struct{}
   func (s *StubWebsiteChecker) Check(url string) bool {
       return true  // Always returns true
   }
   ```

2. **Spy**: Records how it was called
   ```go
   type SpySleeper struct {
       Calls int
   }
   func (s *SpySleeper) Sleep() {
       s.Calls++
   }

   // In test
   spy := &SpySleeper{}
   Countdown(buffer, spy)
   if spy.Calls != 3 {
       t.Errorf("expected 3 calls, got %d", spy.Calls)
   }
   ```

3. **Spy that records call order**:
   ```go
   type SpyCountdownOperations struct {
       Calls []string
   }

   func (s *SpyCountdownOperations) Sleep() {
       s.Calls = append(s.Calls, "sleep")
   }

   func (s *SpyCountdownOperations) Write(p []byte) (n int, err error) {
       s.Calls = append(s.Calls, "write")
       return
   }

   // Test verifies order of operations
   want := []string{"write", "sleep", "write", "sleep", "write"}
   if !reflect.DeepEqual(want, spy.Calls) {
       t.Errorf("wrong order of operations")
   }
   ```

**Mocking anti-patterns**:

1. **Too many mocks (>3)** - Redesign needed
2. **Mocking implementation details** - Test behavior, not implementation
3. **Tests break on refactoring** - Coupled to implementation, not behavior
4. **Complicated mock setup** - Code has too many dependencies

**When mocking becomes painful**:
- Your code has too many dependencies → Break it apart
- Your dependencies are too fine-grained → Consolidate them
- Your test is too concerned with implementation → Test behavior instead

**Key principle**:
> "The definition of refactoring is that the code changes but the behaviour stays the same. If you have decided to do some refactoring, in theory you should be able to make the commit without any test changes."

**Mocking without frameworks**:

You don't need a mocking framework in Go. Use interfaces and simple structs:

```go
// 1. Define interface for dependency
type Sleeper interface {
    Sleep()
}

// 2. Create spy for testing
type SpySleeper struct {
    Calls int
}
func (s *SpySleeper) Sleep() {
    s.Calls++
}

// 3. Create real implementation for production
type ConfigurableSleeper struct {
    duration time.Duration
    sleep    func(time.Duration)
}
func (c *ConfigurableSleeper) Sleep() {
    c.sleep(c.duration)
}

// 4. Inject the dependency
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := 3; i > 0; i-- {
        fmt.Fprintln(out, i)
        sleeper.Sleep()
    }
}
```

**Testing with mocks**:
- Tests run fast (no waiting for real time.Sleep)
- Can verify exact behavior (how many times, what order)
- Tests are deterministic (no flaky tests from timing)

**Value of tests with mocks**:
- Without mocking, important areas remain untested
- Without mocks, tests become slow (spinning up databases, web services)
- Without mocks, tests become fragile (unreliable external services)

**Listen to your tests**:
- If mocking is complicated, your design needs work
- Over-testing every implementation detail is a code smell
- Always ask: "Am I testing behavior or implementation?"

### Testing HTTP Handlers and Servers

**Using httptest.NewRecorder and NewRequest**:

```go
func TestPlayerServer(t *testing.T) {
    request := httptest.NewRequest(http.MethodGet, "/players/Pepper", nil)
    response := httptest.NewRecorder()

    PlayerServer(response, request)

    got := response.Body.String()
    want := "20"

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

**Using httptest.NewServer for integration tests**:

```go
func TestRacer(t *testing.T) {
    slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(20 * time.Millisecond)
        w.WriteHeader(http.StatusOK)
    }))
    defer slowServer.Close()

    fastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))
    defer fastServer.Close()

    got, _ := Racer(slowServer.URL, fastServer.URL)
    want := fastServer.URL

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

**Benefits of httptest.NewServer**:
- Creates real HTTP server for tests
- Automatically finds available port
- Returns URL you can use in tests
- Avoids flaky tests from calling real external APIs
- Can control response timing, errors, edge cases
- Use `defer server.Close()` for cleanup

### Testing Concurrent Code

**Race Detector** - Use `go test -race` to catch data races:

```bash
go test -race
```

Output shows concurrent access problems:
```
WARNING: DATA RACE
Write at 0x00c420084d20 by goroutine 8:
  runtime.mapassign_faststr()
```

**Benchmarking concurrent improvements**:

```go
func BenchmarkCheckWebsites(b *testing.B) {
    urls := make([]string, 100)
    for i := 0; i < len(urls); i++ {
        urls[i] = "a url"
    }

    for b.Loop() {
        CheckWebsites(slowStubWebsiteChecker, urls)
    }
}
```

Benchmark before concurrency: `2249228637 ns/op` (2.2 seconds)
Benchmark after concurrency: `23406615 ns/op` (0.023 seconds) - 100x faster!

**Testing with sync.WaitGroup**:

```go
func TestCounter(t *testing.T) {
    t.Run("it runs safely concurrently", func(t *testing.T) {
        wantedCount := 1000
        counter := NewCounter()

        var wg sync.WaitGroup
        wg.Add(wantedCount)

        for i := 0; i < wantedCount; i++ {
            go func() {
                counter.Inc()
                wg.Done()
            }()
        }
        wg.Wait()

        if counter.Value() != wantedCount {
            t.Errorf("got %d, want %d", counter.Value(), wantedCount)
        }
    })
}
```

**sync.WaitGroup** - Coordinate goroutines in tests:
- `wg.Add(n)` - Set number of goroutines to wait for
- `wg.Done()` - Call when goroutine finishes
- `wg.Wait()` - Block until all goroutines finish

**Avoiding mutex copies** - Use `go vet`:

```bash
go vet
```

Output:
```
sync_test.go:16: call of assertCounter copies lock value
```

Fix: Pass pointers to structs containing mutexes:
```go
func assertCounter(t testing.TB, got *Counter, want int) {
    // ...
}
```

### Testing with Channels and Select

**Testing concurrent operations with select**:

```go
func TestRacer(t *testing.T) {
    slowServer := makeDelayedServer(20 * time.Millisecond)
    fastServer := makeDelayedServer(0 * time.Millisecond)

    defer slowServer.Close()
    defer fastServer.Close()

    got, _ := Racer(slowServer.URL, fastServer.URL)
    want := fastServer.URL

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

**Testing timeouts with time.After**:

```go
func TestRacer(t *testing.T) {
    t.Run("returns an error if server doesn't respond within timeout", func(t *testing.T) {
        server := makeDelayedServer(25 * time.Millisecond)
        defer server.Close()

        _, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

        if err == nil {
            t.Error("expected an error but didn't get one")
        }
    })
}
```

**Making timeouts configurable for fast tests**:

```go
// Production code with sensible default
func Racer(a, b string) (string, error) {
    return ConfigurableRacer(a, b, 10*time.Second)
}

// Test-friendly version with configurable timeout
func ConfigurableRacer(a, b string, timeout time.Duration) (string, error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("timed out")
    }
}
```

**Key principle**: Don't make tests wait for real timeouts. Make timing configurable so tests run fast.

### Testing with Context

**Testing context cancellation**:

```go
func TestServer(t *testing.T) {
    t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
        store := &SpyStore{response: "data", t: t}
        svr := Server(store)

        request := httptest.NewRequest(http.MethodGet, "/", nil)

        cancellingCtx, cancel := context.WithCancel(request.Context())
        time.AfterFunc(5*time.Millisecond, cancel)
        request = request.WithContext(cancellingCtx)

        response := httptest.NewRecorder()
        svr.ServeHTTP(response, request)

        if !store.cancelled {
            t.Error("store was not told to cancel")
        }
    })
}
```

**Testing both success and cancellation**:

```go
// Spy that tracks cancellation
type SpyStore struct {
    response  string
    cancelled bool
    t         *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
    data := make(chan string, 1)

    go func() {
        var result string
        for _, c := range s.response {
            select {
            case <-ctx.Done():
                return
            default:
                result += string(c)
            }
        }
        data <- result
    }()

    select {
    case d := <-data:
        return d, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

**Key testing practices with context**:
- Use `context.WithCancel` to create cancellable contexts in tests
- Use `time.AfterFunc` to trigger cancellation after a delay
- Test both happy path (no cancellation) and sad path (cancellation)
- Verify context is propagated through function calls

### Concurrency Testing Principles

**"Make it work, make it right, make it fast"**:
1. First, make tests pass
2. Then, refactor for good design
3. Finally, optimize for performance (use benchmarks to prove improvement)

**Never optimize prematurely**:
> "Premature optimization is the root of all evil" - Donald Knuth

Only optimize after you have:
- Working code (tests pass)
- Clean code (refactored)
- Evidence it's too slow (benchmarks)

**Race conditions**:
- Occur when output depends on uncontrollable timing/sequence
- Use `go test -race` to detect them
- Fix with channels, mutexes, or better design
- Channels for passing ownership, mutexes for managing state

## Property-Based Testing

Go has built-in support for property-based testing via `testing/quick`.

**What is property-based testing?**

Instead of testing specific examples, you test properties that should always be true for any input.

**Example**: Testing that converting to Roman and back gives original number

```go
import "testing/quick"

func TestPropertiesOfConversion(t *testing.T) {
    assertion := func(arabic uint16) bool {
        if arabic > 3999 {
            return true  // Skip invalid inputs
        }
        roman := ConvertToRoman(arabic)
        fromRoman := ConvertToArabic(roman)
        return fromRoman == arabic
    }

    if err := quick.Check(assertion, nil); err != nil {
        t.Error("failed checks", err)
    }
}
```

**How it works**:
- `quick.Check` runs your function with many random inputs (default 100)
- If function returns `false`, the test fails
- Shows you which input caused the failure

**Configuring iterations**:

```go
if err := quick.Check(assertion, &quick.Config{
    MaxCount: 1000,
}); err != nil {
    t.Error("failed checks", err)
}
```

**Benefits**:
- Forces you to think deeply about domain rules
- Finds edge cases you didn't think of
- Tests many scenarios with less code
- Built into standard library

**When to use**:
- You can express domain rules as properties
- Want to test against many inputs
- Complement to example-based tests

**Property examples**:
- Reversing a string twice gives original: `reverse(reverse(s)) == s`
- Adding then subtracting gives original: `(n + x) - x == n`
- Sorting is idempotent: `sort(sort(list)) == sort(list)`

## Testing File Operations

**Use `testing/fstest` for testing file system interactions**:

```go
import "testing/fstest"

func TestNewBlogPosts(t *testing.T) {
    fs := fstest.MapFS{
        "hello world.md":  {Data: []byte("hi")},
        "hello-world2.md": {Data: []byte("hola")},
    }

    posts := blogposts.NewPostsFromFS(fs)

    if len(posts) != len(fs) {
        t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
    }
}
```

**Benefits of fstest.MapFS**:
- Simpler than maintaining test files on disk
- Faster execution
- Easy to test edge cases (permission errors, missing files)
- No cleanup needed

**Design for testability**:

Accept `fs.FS` interface instead of file paths:

```go
// Good - testable
func NewPostsFromFS(filesystem fs.FS) []Post {
    // ...
}

// Bad - hard to test
func NewPostsFromFS(folderPath string) []Post {
    // ...
}
```

This lets users:
- Use real file system in production
- Use test file system in tests
- Use embedded file system (embed.FS)
- Use zip files (zip.Reader)

**Test at package boundary**:

Use `package mypackage_test` to test as a consumer would:

```go
package blogposts_test  // Not package blogposts

import (
    "testing"
    "github.com/you/blogposts"
)
```

This enforces testing through public API only.

## Advanced Testing Strategies

### Acceptance Tests

**What are acceptance tests?**

Black-box tests that exercise your system as a user would, without access to internals.

**Benefits**:
- When they pass, you know the entire system works
- More accurate and faster than manual testing
- Act as verified documentation of system behavior
- No mocking - tests use real implementations

**Drawbacks vs unit tests**:
- Expensive to write
- Take longer to run
- Harder to debug when they fail
- Don't give feedback on internal code quality
- Not all scenarios are practical to test

**When to use**:
- Testing end-to-end critical paths
- Verifying system integration
- Confidence before shipping
- Should be minority of tests (see Test Pyramid)

**Example**: Testing graceful shutdown of HTTP server

```go
func TestGracefulShutdown(t *testing.T) {
    // Build and run actual program
    cleanup, sendInterrupt, err := LaunchTestProgram("8080")
    if err != nil {
        t.Fatal(err)
    }
    defer cleanup()

    // Start HTTP request
    responseChan := make(chan *http.Response)
    go func() {
        resp, _ := http.Get("http://localhost:8080/slow")
        responseChan <- resp
    }()

    // Send SIGTERM before response completes
    time.Sleep(5 * time.Millisecond)
    sendInterrupt()

    // Verify we still get response despite shutdown
    select {
    case resp := <-responseChan:
        if resp.StatusCode != http.StatusOK {
            t.Errorf("expected OK, got %d", resp.StatusCode)
        }
    case <-time.After(2 * time.Second):
        t.Error("did not get response - server did not shutdown gracefully")
    }
}
```

### Test Pyramid

Balance your testing strategy:

```
        /\
       /  \  Acceptance Tests (few, slow, high confidence)
      /____\
     /      \
    / Unit   \ Integration Tests (some, medium speed)
   /  Tests   \
  /____________\ Unit Tests (many, fast, specific)
```

- **Many unit tests**: Fast feedback, test small units
- **Some integration tests**: Test components working together
- **Few acceptance tests**: End-to-end critical paths only

**Aim for**:
- Fast feedback loops (< 10 seconds for unit tests)
- Fast lead time (< 10 minutes from commit to deploy)
- High confidence without slow test suites

### Fakes vs Mocks/Stubs/Spies

**The problem with mocks and stubs**:

Mocks and stubs encode assumptions about dependencies that may not be validated:

```go
// In test - you assume API returns this format
stubAPI := &StubAPI{
    response: `{"name": "Bob", "age": 30}`,
}
```

If the real API changes format, your tests pass but production fails!

**Test doubles comparison**:

1. **Stub**: Returns canned data
   ```go
   type StubStore struct {
       data string
   }
   func (s *StubStore) Get() string { return s.data }
   ```

2. **Spy**: Records how it was called
   ```go
   type SpyStore struct {
       calls []string
   }
   func (s *SpyStore) Get(key string) string {
       s.calls = append(s.calls, key)
       return ""
   }
   ```

3. **Mock**: Panics if called incorrectly
   ```go
   mockStore.Expect("key1").Return("value1")
   // Panics if called with anything other than "key1"
   ```

4. **Fake**: Real implementation optimized for testing
   ```go
   type FakeStore struct {
       data map[string]string
   }
   func (f *FakeStore) Get(key string) string {
       return f.data[key]
   }
   func (f *FakeStore) Set(key, value string) {
       f.data[key] = value
   }
   ```

**Why fakes are better**:

1. **Stateful and realistic**:
   - Fakes behave like real dependencies
   - Can test complex scenarios with multiple interactions
   - Tests read more naturally

2. **Reusable across tests**:
   - One fake serves many tests
   - Mocks/stubs often created per-test

3. **Better for local development**:
   - Run app locally without real database/API
   - Faster and more reliable than real dependencies

4. **Validated behavior**:
   - Fakes can be tested with contract tests
   - Ensures fake behaves like real implementation

### Contract Tests

**Problem**: How do you know your fake actually behaves like the real thing?

**Solution**: Contract tests verify both implementations satisfy the same interface.

```go
// Contract test that runs against BOTH real and fake implementations
func TestStorageContract(t *testing.T) {
    tests := []struct {
        name    string
        storage Storage
    }{
        {"Real PostgreSQL", NewPostgresStorage(connString)},
        {"Fake in-memory", NewFakeStorage()},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test behavior both must satisfy
            tt.storage.Set("key", "value")
            got := tt.storage.Get("key")

            if got != "value" {
                t.Errorf("got %q want %q", got, "value")
            }

            // Both must handle missing keys the same way
            got = tt.storage.Get("nonexistent")
            if got != "" {
                t.Errorf("expected empty string for missing key")
            }
        })
    }
}
```

**Benefits of contract tests**:
- Fake and real implementation stay in sync
- If real API changes, contract test fails, you update fake
- Gives confidence fake is accurate
- Documents expected behavior of the interface

**Testing strategy with fakes**:
1. Write contract tests for the interface
2. Implement fake that passes contract tests
3. Use fake in unit tests (fast, reliable)
4. Use real implementation in acceptance tests (confidence)
5. When real implementation changes, contract test fails, update fake

### Scaling Acceptance Tests

**The problem with traditional acceptance tests**:
- Slow to run
- Brittle and flaky
- Tightly coupled to implementation details
- Expensive to maintain
- Break when you change markup/API structure (even if behavior unchanged)

**Root cause**: Not separating essential complexity (behavior) from accidental complexity (HTTP, databases, UI).

**Solution**: Separate Specifications, DSL, and Drivers

```
Specification (What) → Driver (How) → System
```

**Essential complexity**: Domain rules that exist independent of computers
- "Greet a user by name"
- "Withdraw money from account if sufficient funds"

**Accidental complexity**: Computer-specific stuff
- HTTP endpoints, JSON, databases, markup
- Changes frequently without changing behavior

### Specifications Pattern

Create reusable specifications decoupled from implementation:

```go
// specifications/greet.go
package specifications

type Greeter interface {
    Greet(name string) (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
    got, err := greeter.Greet("Mike")
    assert.NoError(t, err)
    assert.Equal(t, got, "Hello, Mike")
}
```

This specification can test:
- HTTP API (via HTTP driver)
- gRPC service (via gRPC driver)
- Domain code directly (via adapter)
- Web UI (via browser driver)

**All with the same specification!**

### Drivers

Drivers implement the DSL interface for specific systems:

```go
// adapters/httpserver/driver.go
type Driver struct {
    BaseURL string
    Client  *http.Client
}

func (d Driver) Greet(name string) (string, error) {
    res, err := d.Client.Get(d.BaseURL + "/greet?name=" + name)
    if err != nil {
        return "", err
    }
    defer res.Body.Close()
    greeting, err := io.ReadAll(res.Body)
    return string(greeting), err
}
```

### Acceptance Test Structure

```go
func TestGreeterServer(t *testing.T) {
    // 1. Start system (e.g., Docker container)
    StartDockerServer(t, port, dockerFile)

    // 2. Create driver for that system
    driver := httpserver.Driver{
        BaseURL: "http://localhost:8080",
        Client: &http.Client{Timeout: 1 * time.Second},
    }

    // 3. Run specification
    specifications.GreetSpecification(t, driver)
}
```

**Benefits**:
- Specifications only change when behavior changes
- Drivers only change when implementation details change
- Can test locally, staging, and production with same specs
- Fast: tests run in milliseconds even for full system

### Testing with Testcontainers

Use testcontainers to run your actual Docker images in tests:

```go
import "github.com/testcontainers/testcontainers-go"

func StartDockerServer(t testing.TB, port, dockerFilePath string) {
    ctx := context.Background()
    req := testcontainers.ContainerRequest{
        FromDockerfile: testcontainers.FromDockerfile{
            Context:    "../../.",
            Dockerfile: dockerFilePath,
        },
        ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
        WaitingFor:   wait.ForListeningPort(nat.Port(port)),
    }
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    assert.NoError(t, err)
    t.Cleanup(func() {
        assert.NoError(t, container.Terminate(ctx))
    })
}
```

**Benefits**:
- Tests the actual image you'll ship
- Runs locally with no special setup
- Automatic cleanup with `t.Cleanup`
- Fast enough for development

### Adapter Pattern for Domain Testing

Reuse specifications for unit testing domain code:

```go
// specifications/adapters.go
type GreetAdapter func(name string) string

func (g GreetAdapter) Greet(name string) (string, error) {
    return g(name), nil
}

// In unit test
func TestGreet(t *testing.T) {
    specifications.GreetSpecification(
        t,
        specifications.GreetAdapter(domain.Greet),
    )
}
```

Adapter "adapts" your function to match the interface specification needs.

### When to Write Acceptance Tests

**Write acceptance tests for**:
- Critical user journeys
- Features stakeholders care deeply about
- End-to-end integration confidence

**Use unit tests for**:
- Edge cases
- Error handling
- Different input combinations
- Internal behavior details

**Test Pyramid reminder**:
- Many unit tests (fast, specific)
- Some integration tests
- Few acceptance tests (slow, high confidence)

### Separating Test Types

Use `testing.Short()` to skip slow tests:

```go
func TestGreeterServer(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping acceptance test")
    }
    // ... acceptance test
}
```

Run fast tests: `go test -short ./...`
Run all tests: `go test ./...`

### Top-Down Development with Acceptance Tests

**The GOOS approach**:

1. Write acceptance test for new feature (will fail for a while)
2. Use acceptance test as "north star" - all work goes toward making it pass
3. Break down into unit tests and implement step by step
4. Refactor as you go
5. When acceptance test passes, feature is done

**Benefits**:
- Work is always validated against real integration
- Reduces wasted effort on wrong approaches
- Keeps you focused on actual goal
- Catches integration issues early

**Avoid "bottom-up"**:
- Building components without integrating them
- Risk they won't work together or solve the actual problem
- Wastes effort on unvalidated ideas

## Refactoring Principles

**Definition**: Refactoring is improving code structure **without changing behavior**.

**Critical rule**: If refactoring, tests should NOT need to change.

If you're changing tests during "refactoring", you're actually doing design changes, not refactoring.

### What You CAN Do While Refactoring

- Introduce private methods/functions
- Extract variables and constants
- Rename symbols (with IDE support)
- Change internals of public methods
- Move code around

### What You CANNOT Do While Refactoring

- Change behavior
- Change method signatures (without IDE refactoring tool)
- Change what tests assert

### Common Refactoring Patterns

**Extract method** (`command+option+m`):
```go
// Before
func CreateWidget(name string) error {
    url := baseURL + "/widgets"
    payload := []byte(`{"name": "` + name + `"}`)
    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
    // ...
}

// After
func CreateWidget(name string) error {
    url := baseURL + "/widgets"
    req, err := http.NewRequest(http.MethodPost, url, createWidgetPayload(name))
    // ...
}
```

**Extract variable** (`command+option+v`):
```go
// Before - magic number
client := http.Client{Timeout: 1 * time.Second}

// After - named constant
const defaultTimeout = 1 * time.Second
client := http.Client{Timeout: defaultTimeout}
```

**Inline variable** (`command+option+n`):
```go
// Before
url := baseURL + "/user/" + id
res, err := client.Get(url)

// After
res, err := client.Get(baseURL + "/user/" + id)
```

### DRY (Don't Repeat Yourself)

**The real goal**: Capture an **idea** in one place, not just reduce lines of code.

**Good DRY**:
```go
const maxRetries = 3

func callAPI() { /* uses maxRetries */ }
func callDB() { /* uses maxRetries */ }
```

**Bad DRY** (creates coupling where it shouldn't exist):
```go
const timeout = 1 * time.Second

apiClient := http.Client{Timeout: timeout}
dbClient := http.Client{Timeout: timeout}  // Maybe DB needs different timeout?
```

**Warning signs of bad DRY**:
- Long, confusing parameter lists
- Making code more complex to avoid duplication
- Coupling unrelated concepts

### Refactoring Habits

**Run tests constantly**:
- After every small refactoring
- Feedback loop should be < 1 second
- If tests fail, revert immediately

**Use source control**:
- Commit frequently when tests pass
- Don't be afraid to experiment
- Easy to revert if refactoring doesn't work out

**Make public methods scannable**:

Public methods should describe WHAT, not HOW:

```go
// Before - too much detail
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    var req Request
    json.Unmarshal(body, &req)
    // lots more code
}

// After - clear steps
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
    req := s.parseRequest(r)
    result := s.processRequest(req)
    s.sendResponse(w, result)
}
```

**Remove comments with code**:

> "Whenever we feel the need to comment something, we write a method instead." - Martin Fowler

```go
// Before
// Create JSON payload with widget name
payload := []byte(`{"name": "` + name + `"}`)

// After
payload := createWidgetPayload(name)
```

### Refactoring vs Design Changes

**Refactoring**: Automated or very safe changes
- Rename symbols
- Extract/inline methods
- Extract variables
- No test changes needed

**Design changes**: Require TDD process
- Changing method signatures
- Changing behavior
- Restructuring types
- Tests need to change

**For design changes, still use TDD**:
1. Comment out all but one test
2. Drive out the change with TDD
3. Work through remaining tests one by one

### When to Refactor

**Always**: After making tests pass (3rd step of TDD)

**Never**: When you "don't have time"

> "Not having enough time usually is a sign that you need to do some refactoring" - Martin Fowler

**Benefits of frequent refactoring**:
- Code stays easy to understand
- Easier to spot design improvements
- Prevents exponential complexity growth
- Increases team velocity long-term

### Don't Ask Permission

Refactoring is part of your job as a professional developer, not optional work you need approval for.

Bad code slows everyone down. Refactoring speeds everyone up.

## Key Takeaways

### The TDD Cycle (Never Skip!)

1. **🔴 Red** - Write test first, see it fail
2. **🟢 Green** - Write minimal code to pass
3. **🔄 Refactor** - Improve design (tests don't change!)

### Core Principles

- **Write test first** - Always, no exceptions. This is what makes it TDD.
- **See it fail** - Verify the error message is clear before implementing
- **Make it pass** - Minimal code to green, nothing more
- **Refactor** - Improve design with test safety net (if tests change, you're not refactoring!)
- **Small steps** - Keep feedback loops tight (< 10 minutes per cycle)

### Listen to Your Tests

Tests are the first user of your code. If testing is hard, using the code will be hard.

**Warning signs**:
- Too many mocks (>3) → Too many dependencies
- Complex test setup → Complex code structure
- Tests break when refactoring → Testing implementation, not behavior
- Hard to write test → Code is hard to use

**Action**: Simplify dependencies, redesign the API, test at higher abstraction level.

### Testing Best Practices

- **Test behavior, not implementation** - Tests should survive refactoring
- **One assertion per test** - Use subtests for different scenarios
- **Keep it simple** - Simple is not easy, but it's valuable
- **Refactor constantly** - Don't let code quality degrade
- **Use your tools** - Learn your IDE's refactoring shortcuts
- **Don't ask permission to refactor** - It's your professional responsibility

### Go-Specific Testing

- Use `go test -race` to catch concurrency bugs
- Use `go test -cover` to find gaps (but don't chase 100%)
- Use `slices.Equal` (type-safe) over `reflect.DeepEqual`
- Export errors as values: `var ErrNotFound = errors.New("not found")`
- Test through interfaces, inject dependencies
- No mocking frameworks needed - simple structs work great

### Remember

> "TDD gives you the fastest feedback possible on your design" - Dave Farley

> "If testing your code is difficult, then using your code is difficult too"

> "The bigger the interface, the weaker the abstraction" - Go Proverb

**Write tests. See them fail. Make them pass. Refactor. Repeat.**

## Source Material

This skill is based on "Learn Go with Tests" by Chris James, a comprehensive resource for learning Go through TDD. The material emphasizes that TDD is fundamentally about getting fast feedback on design decisions and building software that can evolve safely over time.

**Additional resources**:
- Go testing package: https://pkg.go.dev/testing
- Martin Fowler's Refactoring: https://refactoring.com
- Growing Object-Oriented Software, Guided by Tests (GOOS book)
