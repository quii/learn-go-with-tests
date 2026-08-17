# Revisiting time, with testing/synctest

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/synctest)**

In [the chapter on time](time.md) we gave our poker CLI the ability to schedule "blind is now going up" alerts, using `time.AfterFunc`. Testing this was tricky: `time.AfterFunc` runs its callback in its own goroutine after a real duration elapses, and you can't compare functions in Go, so we couldn't easily inspect what had been scheduled.

We solved that with a familiar tool: dependency injection. We defined a `BlindAlerter` interface, and in our tests we swapped the real implementation for a spy that just records what it was asked to schedule, something like this:

```go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type SpyBlindAlerter struct {
	Alerts []struct {
		At     time.Duration
		Amount int
	}
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.Alerts = append(s.Alerts, struct {
		At     time.Duration
		Amount int
	}{at, amount})
}
```

This is a good design, and the tests it enables are fast and reliable. But look closely at what they actually cover. They assert things like "`ScheduleAlertAt` was called with `10 * time.Minute` and `200`". They never let a *real* alerter run. The one piece of code that actually calls `time.AfterFunc` and prints something, the part with the real bug potential, never gets exercised by a test.

That's not an oversight, it's a trade-off. Testing a real `time.AfterFunc`-based alerter properly would mean either waiting real minutes for a test to finish, or shrinking the durations down to milliseconds and hoping the machine running the test isn't too busy to keep up. Neither is appealing. So historically, we just didn't.

As of Go 1.25, there's a third option: [`testing/synctest`](https://pkg.go.dev/testing/synctest). It lets us run real, unmodified code that uses `time.Sleep`, `time.AfterFunc`, and friends, inside a test that has full control over time: no waiting, no flakiness, no shrinking durations and hoping for the best.

In this chapter we'll build that real alerter from scratch, and test it with `synctest`. But let's earn it, by first seeing exactly what it saves us from.

## Just enough information on `testing/synctest`

`testing/synctest` runs a function inside an isolated **bubble**. Within that bubble:

* The `time` package uses a fake clock. It starts at midnight UTC on the 1st of January, 2000, and only moves forward.
* Fake time only advances when every goroutine in the bubble is **durably blocked**: blocked in a way that only another goroutine in the same bubble can unblock it. `time.Sleep`, a blocking receive on a channel created inside the bubble, `sync.Cond.Wait`, and `sync.WaitGroup.Wait` all count. Locking a `sync.Mutex` does not, because mutexes are typically held only briefly.
* `synctest.Test(t, f)` runs `f` in a new bubble and doesn't return until every goroutine spawned inside it has exited. If the bubble ends up durably blocked with no way to make progress, it fails the test as a deadlock rather than hanging forever.
* `synctest.Wait()` blocks the calling goroutine until every *other* goroutine in the bubble is durably blocked, then returns. It's how you let background work settle before making an assertion.

That's enough to get started. We'll pick up a couple of sharper edges as we go.

## Write the test first

Let's start with the obvious thing, before reaching for `synctest` at all: a `BlindAlerter` that, given a duration and an amount, waits for that duration and then writes a message somewhere, tested with real time.

```go
package poker

import (
	"bytes"
	"testing"
	"time"
)

func TestStdOutAlerter(t *testing.T) {
	out := &bytes.Buffer{}
	alerter := StdOutAlerter(out)

	alerter.ScheduleAlertAt(5*time.Second, 100)

	time.Sleep(6 * time.Second)

	want := "Blind is now 100"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}
```

We sleep a little longer than the scheduled alert (6 seconds, not 5) to give the goroutine `time.AfterFunc` spawns a moment to actually run before we check.

## Try to run the test

We haven't written any production code yet, so this won't compile:

```
./blind_alerter_test.go:11:13: undefined: StdOutAlerter
```

## Write enough code to make it pass

```go
package poker

import (
	"fmt"
	"io"
	"time"
)

// BlindAlerter schedules alerts for blind amounts.
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

// BlindAlerterFunc allows you to implement BlindAlerter with a function.
type BlindAlerterFunc func(duration time.Duration, amount int)

// ScheduleAlertAt is BlindAlerterFunc's implementation of BlindAlerter.
func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

// StdOutAlerter returns a BlindAlerterFunc that schedules alerts and prints them to out.
func StdOutAlerter(out io.Writer) BlindAlerterFunc {
	return func(duration time.Duration, amount int) {
		time.AfterFunc(duration, func() {
			fmt.Fprintf(out, "Blind is now %d", amount)
		})
	}
}
```

Run it:

```
=== RUN   TestStdOutAlerter
--- PASS: TestStdOutAlerter (6.00s)
PASS
ok  	github.com/quii/learn-go-with-tests/synctest/v1	6.191s
```

It passes. It also takes six real seconds, for one alert. Our actual game schedules eleven of them, some as far as an hour and forty minutes apart. Nobody is going to wait for that on every test run, so in practice this test would either get skipped, or written with unrealistically tiny durations that only vaguely resemble the real thing. This is the problem `synctest` exists to solve.

## Introducing synctest

Let's wrap the same test in a bubble and see what happens if we're a bit too optimistic about what a "fake clock" does for us:

```go
func TestStdOutAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		out := &bytes.Buffer{}
		alerter := StdOutAlerter(out)

		alerter.ScheduleAlertAt(5*time.Second, 100)

		want := "Blind is now 100"
		if out.String() != want {
			t.Errorf("got %q, want %q", out.String(), want)
		}
	})
}
```

We dropped the sleep entirely: surely the fake clock handles that for us now? It doesn't:

```
=== RUN   TestStdOutAlerter
    blind_alerter_test.go:19: got "", want "Blind is now 100"
--- FAIL: TestStdOutAlerter (0.00s)
FAIL
```

A bubble's fake clock doesn't run itself forward on a timer of its own. It only advances when something in the bubble durably blocks, waiting for it to. We still have to say what we're waiting for; we just don't have to pay for it in real seconds any more. Let's put the sleep back:

```go
func TestStdOutAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		out := &bytes.Buffer{}
		alerter := StdOutAlerter(out)

		alerter.ScheduleAlertAt(5*time.Second, 100)

		time.Sleep(6 * time.Second)

		want := "Blind is now 100"
		if out.String() != want {
			t.Errorf("got %q, want %q", out.String(), want)
		}
	})
}
```

```
=== RUN   TestStdOutAlerter
--- PASS: TestStdOutAlerter (0.00s)
PASS
ok  	github.com/quii/learn-go-with-tests/synctest/v2	0.133s
```

Passes, and instantly: `time.Sleep(6 * time.Second)` inside a bubble costs nothing in real time. Looks done. Let's just double-check it holds up under `-race`, out of habit:

```
go test -race ./...
```

```
==================
WARNING: DATA RACE
Read at 0x00c00010e630 by goroutine 9:
  bytes.(*Buffer).String()
      /usr/local/go/src/bytes/buffer.go:77 +0x174
  github.com/quii/learn-go-with-tests/synctest/v2.TestStdOutAlerter.func1()
      blind_alerter_test.go:20 +0x15c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x164

Previous write at 0x00c00010e630 by goroutine 11:
  bytes.(*Buffer).grow()
      /usr/local/go/src/bytes/buffer.go:143 +0x354
  bytes.(*Buffer).Write()
      /usr/local/go/src/bytes/buffer.go:185 +0xb4
  fmt.Fprintf()
      /usr/local/go/src/fmt/print.go:225 +0x94
  github.com/quii/learn-go-with-tests/synctest/v2.TestStdOutAlerter.func1.BlindAlerterFunc.ScheduleAlertAt.TestStdOutAlerter.func1.StdOutAlerter.1.2()
      blind_alerter.go:26 +0x6c

Goroutine 11 (finished) created at:
  time.goFunc()
      /usr/local/go/src/time/sleep.go:215 +0x40
==================
--- FAIL: TestStdOutAlerter (0.00s)
    testing.go:1617: race detected during execution of test
FAIL
```

Ouch. Notice "Goroutine 11 (finished)": the write genuinely happened before our read, every single time we ran it. Functionally, the test can never actually fail on the assertion. But the race detector isn't asking "did this go wrong this time?"; it's asking "is there anything that *guarantees* it can't?" `time.Sleep` on our goroutine and `time.AfterFunc`'s callback on its own goroutine are just two independently-scheduled goroutines. Nothing about "I slept for a while" promises another goroutine has finished touching shared memory, no matter how generous the sleep. Sleeping longer doesn't fix this, no matter how long; it just makes it fail less often outside of `-race`.

Worth sitting with for a second: this exact race was already lurking in our very first, plain real-time version too, with a real `time.Sleep`. We just never thought to check: who runs `-race`, repeatedly, on a test that already takes six real seconds to run?

## Write the test first

What we need is a real synchronization point: something that doesn't just make the write *probably* happen first, but actually establishes it did. That's exactly what `synctest.Wait()` is for:

```go
func TestStdOutAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		out := &bytes.Buffer{}
		alerter := StdOutAlerter(out)

		alerter.ScheduleAlertAt(5*time.Second, 100)

		time.Sleep(6 * time.Second)
		synctest.Wait()

		want := "Blind is now 100"
		if out.String() != want {
			t.Errorf("got %q, want %q", out.String(), want)
		}
	})
}
```

## Refactor

There's nothing to change in the production code, but let's make sure this one actually holds up, repeatedly, under `-race`:

```
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.153s
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.143s
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.148s
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.147s
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.143s
```

Clean, every time. Very nice.

## Write the test first

There's a second thing worth testing here, and it's the feature that makes `synctest` genuinely better than "just use real time with generous sleeps": proving something *hasn't* happened yet, without guessing at how long "not yet" should mean:

```go
func TestStdOutAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		out := &bytes.Buffer{}
		alerter := StdOutAlerter(out)

		alerter.ScheduleAlertAt(5*time.Second, 100)

		synctest.Wait()
		if out.String() != "" {
			t.Errorf("did not expect anything to be printed yet, got %q", out.String())
		}

		time.Sleep(5 * time.Second)
		synctest.Wait()

		want := "Blind is now 100"
		if out.String() != want {
			t.Errorf("got %q, want %q", out.String(), want)
		}
	})
}
```

Nothing else in the bubble is doing anything at the point we schedule the alert, so `Wait()` should return straight away: nothing to wait *for* yet, so `out` should still be empty. Run it:

```
=== RUN   TestStdOutAlerter
--- PASS: TestStdOutAlerter (0.00s)
PASS
```

Passes. Let's check `-race`, as we now know we should:

```
==================
WARNING: DATA RACE
Write at 0x00c00011c630 by goroutine 11:
  bytes.(*Buffer).grow()
      /usr/local/go/src/bytes/buffer.go:143 +0x354
  bytes.(*Buffer).Write()
      /usr/local/go/src/bytes/buffer.go:185 +0xb4
  fmt.Fprintf()
      /usr/local/go/src/fmt/print.go:225 +0x94
  github.com/quii/learn-go-with-tests/synctest/v4.TestStdOutAlerter.func1.BlindAlerterFunc.ScheduleAlertAt.TestStdOutAlerter.func1.StdOutAlerter.1.2()
      blind_alerter.go:26 +0x6c

Previous read at 0x00c00011c630 by goroutine 9:
  bytes.(*Buffer).String()
      /usr/local/go/src/bytes/buffer.go:77 +0x164
  github.com/quii/learn-go-with-tests/synctest/v4.TestStdOutAlerter.func1()
      blind_alerter_test.go:18 +0x14c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x164

Goroutine 11 (running) created at:
  time.goFunc()
      /usr/local/go/src/time/sleep.go:215 +0x40
==================
--- FAIL: TestStdOutAlerter (0.00s)
    testing.go:1617: race detected during execution of test
FAIL
```

The race is at line 18: our very first check, the one we were confident about, because there was "nothing to wait for yet". Reliably, every run.

The issue is that `Wait()` alone has nothing bounding how far it's willing to let the fake clock run. There's no other goroutine in the bubble at that point, and no upcoming deadline for *us* to wait on. Only the alert's 5-second timer, sitting there in the runtime's timer heap, is left. So the runtime does the only useful thing it can: it fires the timer to make progress, spawning goroutine 11 to run our callback, *while* `Wait()` is still deciding whether to return. Goroutine 11 is genuinely running concurrently with our check, not before it. Sometimes that write finishes first; sometimes it doesn't. The `bytes.Buffer` we're both touching has no lock, so there's nothing making that safe either way.

We could fix this the way we'd fix any data race: wrap `out` in a `sync.Mutex`, or use `sync/atomic`'s `atomic.Pointer[T]` for a lighter touch. But look at what `StdOutAlerter` is actually being asked to do: decide the message *and* perform the side effect of printing it. Nothing about scheduling a poker blind alert requires knowing about `io.Writer`; that's the caller's business. That's the shape of tension this book keeps coming back to: [if your tests are causing you pain, listen to that signal and think about the design of your code](http-handlers-revisited.md). Let's have the alerter just produce the message when it's due, and let the caller decide what to do with it.

## Refactor

Crossing a goroutine boundary is exactly what channels are for. As the Go proverb goes, [don't communicate by sharing memory; share memory by communicating](https://go.dev/blog/codelab-share). Instead of writing into a shared `out`, our alerter can send the finished message down a channel:

```go
func NewAlerter() (BlindAlerterFunc, <-chan string) {
	alerts := make(chan string)

	scheduleAlertAt := func(duration time.Duration, amount int) {
		time.AfterFunc(duration, func() {
			alerts <- fmt.Sprintf("Blind is now %d", amount)
		})
	}

	return scheduleAlertAt, alerts
}
```

`BlindAlerter` and `BlindAlerterFunc` stay exactly as they were. Only `StdOutAlerter` (which, tellingly, no longer had anything to do with stdout) is gone, replaced by `NewAlerter`, handing back both the alerter and a channel to receive from. Whatever wants these alerts printed (`main`, say) can range over that channel and print them; that's no longer this package's problem.

The test gets simpler too:

```go
func TestNewAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		alerter, alerts := NewAlerter()

		alerter.ScheduleAlertAt(5*time.Second, 100)

		select {
		case got := <-alerts:
			t.Fatalf("did not expect an alert yet, got %q", got)
		default:
		}

		got := <-alerts
		want := "Blind is now 100"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
```

Notice what's gone: no `synctest.Wait()` anywhere, no custom type to guard a shared value, no `time.Sleep` at all. We don't need any of them any more.

The `select` with a `default` case, straight out of [the chapter on select](select.md), never blocks: it either takes a ready case or falls through to `default` immediately. At this point virtual time is still sitting at zero, five whole (fake) seconds shy of the alert, so there's nothing to synchronize: of course nothing's arrived yet. And `got := <-alerts` doesn't need a nudge either: it's a plain blocking receive, so the bubble does exactly what it's designed to do: since the only thing anyone in the bubble is waiting on is that timer, fake time jumps straight to the moment it fires, wakes the `time.AfterFunc` goroutine, and unblocks our receive.

Run this with `-race`, repeatedly. It stays green: there's no shared memory left to race over.

Compare this to the `SpyBlindAlerter` approach from the earlier chapter. That's still a perfectly good tool for a different job: it checks *what* gets scheduled (for a given player count, are the right amounts scheduled at the right offsets?) without caring about real timing at all (useful when you're testing arithmetic, not timing). `synctest` didn't replace that need; it fills the coverage gap that a pure "what was I asked to do" spy always leaves behind: does the scheduling mechanism itself actually work?

## Wrapping up

### What we've covered

* `synctest.Test` and the idea of a bubble with an isolated fake clock, one that only advances when something durably blocks, not on its own.
* "Durably blocked": the condition that lets fake time advance, and why `sync.Mutex` deliberately doesn't count (but a channel receive does).
* Sleeping "long enough" is not the same as synchronizing: even a generous, always-in-practice-correct sleep is still a real data race if nothing enforces the ordering.
* `synctest.Wait()` fixes that for a single check, but calling it with nothing else in the bubble to bound it can let the fake clock run further than you intended, which is exactly what happened testing the "nothing has happened yet" case.
* A test surfacing a design problem, not just a bug, and fixing the design instead of reaching for a lock.

### Gotchas to watch for

* If a goroutine in the bubble is still durably blocked when the bubble's root function returns, `synctest.Test` fails the test as a deadlock rather than hanging; make sure background goroutines actually finish.
* Channels, timers, and tickers are tied to the bubble they were created in; using one from outside its bubble panics.
* Real network and file I/O are not durably blocking, so you can't drive them through `synctest`'s fake clock directly; reach for something like `net.Pipe` if you need an in-memory stand-in.
* `time.AfterFunc`'s callback runs in a goroutine of its own, with none of the synchronization guarantees a channel gives you. If it has to touch shared state directly, that state needs its own lock, same as any other concurrent code.

### Additional material

* [Go blog: Testing concurrent code with testing/synctest](https://go.dev/blog/synctest)
* [`testing/synctest` package documentation](https://pkg.go.dev/testing/synctest)
* [Go 1.25 release notes](https://go.dev/doc/go1.25)

### A note on how this chapter was written

This is the first chapter in the book written with AI assistance (Claude). I want to be upfront about that, and about what "assistance" actually meant here, because it wasn't "describe a chapter, get a chapter".

The process looked a lot like the TDD loop this book has been teaching the whole way through: research the real `testing/synctest` documentation and source rather than guessing, spike small throwaway programs to check claims before writing a single word of prose, and treat every claim from the docs as something to verify with a real `go test -race` run, not trust outright. Several of the gotchas in this chapter, the `Wait()` and race detector interaction chief among them, only exist here because a test genuinely failed in a way that wasn't expected, over several rounds of actually running it, and the reason had to be dug into before deciding what to write.

The design itself changed shape partway through, too. The first draft had `StdOutAlerter` writing straight into an `io.Writer`, which is exactly the kind of thing this book has always pushed back on when a test starts hurting: [if your tests are causing you pain, listen to that signal and think about the design of your code](http-handlers-revisited.md). So we did, and ended up with the channel-based version above, and a shorter, better chapter for it.

Everything here was reviewed and edited, and pushed back on more than once: when the tone didn't sound right, when a section had ballooned past what the actual lesson warranted, when an explanation reached for jargon where a shown, real test failure would do a better job. If something in here still reads oddly, that's on me, not the tool.
