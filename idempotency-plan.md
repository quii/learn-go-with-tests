# Chapter Plan: Creating Retryable Endpoints (working title)

Other title options considered, any of which would also fit the book's style: "Safe to Call Twice", "Handling Retries Safely", "Making Endpoints Retry-Friendly", "Don't Send It Twice", "When the Client Doesn't Know What Happened". "Creating Retryable Endpoints" is the current pick, deliberately avoiding the term "idempotent" in the title itself even though the concept is named and explained plainly inside the chapter.

## Context for the assistant

This is a planned chapter for "Learn Go with Tests" (quii.gitbook.io), forming part of a new "distributed system patterns" section. This chapter is expected to be written **before** a transactional outbox chapter, which does not exist yet, this chapter comes first, and the outbox pattern is planned as a natural follow-on chapter afterward. The book teaches Go through iterative, TDD-driven chapters: each concept is motivated by a failing test first, the simplest possible change is made to pass it, then refactoring is discussed, often surfacing genuine design trade-offs rather than a single "correct" answer. Code is written in Go. The `testing/synctest` package is used elsewhere in the book to make concurrency and time-based logic deterministic and fast to test, and this chapter is expected to lean on it too.

The chapter's subject: building an endpoint that safely handles being called more than once with the same logical request (the property computer science calls "idempotency", introduced and named plainly in the chapter, just kept out of the title). The worked example throughout is an endpoint that runs some "nebulous business logic" and sends a confirmation email as a side effect.

Core framing to preserve throughout the chapter: **idempotency is fundamentally an information problem, not a concurrency problem.** The client cannot always know whether its own previous request succeeded. Concurrent double-calls are just one visible symptom of that underlying uncertainty; a sequential retry after a lost confirmation is the more common, more realistic trigger, and should be introduced first.

The motivating scenario for the primary test is: **a caller crashed or otherwise failed to durably record that its request succeeded before it could act on that knowledge.** On restart, it retries with the same idempotency key. No HTTP timeout or connection-hijacking tricks are needed to demonstrate this, it's simply two sequential calls with the same key, which keeps the test simple and avoids muddying the chapter with network-transport mechanics. A concurrent-callers scenario is introduced later as a variant of the same underlying problem, not a separate justification for the pattern.

**Relationship to the (not-yet-written) outbox chapter:** worth a brief, accurate mention rather than treating them as sharing storage. The outbox and the idempotency-key store are *not* the same table wearing two hats, they sit on opposite sides of a call: an outbox lives on the *caller's* side (recording "I need to send this, and haven't confirmed it went out"), while the idempotency-key store built in this chapter lives on the *receiver's* side (recording "I've already handled this key"). A single service is often both a sender (with an outbox for what it dispatches) and a receiver (with an idempotency store for what it accepts), but they are two separate tables solving mirrored problems, not one shared store. What they do share is the underlying trick: a unique-constraint-backed dedup table, applied to opposite ends of the same call. That structural symmetry is what makes this chapter and the future outbox chapter a natural pair, and is worth a closing pointer along the lines of "the next chapter looks at the same problem from the caller's side."

Related real-world reference: the IETF Internet-Draft "The Idempotency-Key HTTP Header Field" (`draft-ietf-httpapi-idempotency-key-header`, Standards Track, header name `Idempotency-Key`). As of last check it was still a draft, not a ratified RFC, worth verifying current status before publishing. Useful details from it:
- Recommends UUID or similar random identifier as the key.
- Recommends key expiry/purging.
- Specifies `409 Conflict` (optionally with a `Link` header or RFC 7807 Problem Details body) for a request arriving while another with the same key is still in flight.
- The pattern is implicitly scoped per-operation.
- Originally authored by PayPal engineers; Stripe uses the same header name in practice, making it the de facto standard regardless of RFC status.

---

## Cycle-by-cycle plan

### Cycle 1: the happy path (scaffolding, no failing test yet)

Write the endpoint doing its nebulous business logic plus sending a confirmation email, no idempotency machinery. First test: call it once, assert the email fake recorded one send. Establishes the baseline behaviour to make retry-friendly later.

### Cycle 2: the crash-and-replay test (red)

New test: call the endpoint, then call it again with the same idempotency key, simulating a caller that crashed before it could record the first response. Assert the email fake recorded exactly **one** send.

Fails: two emails sent, nothing reads or checks the key yet.

**Motivates:** the endpoint needs somewhere to record "this key has already been handled."

### Cycle 3: simplest thing that could work (green)

A package-level or struct-held `map[string]bool` of seen keys, check-then-set before doing the work. Passes cycle 2's test with minimal code, in keeping with the book's usual restraint against over-building before a test demands it.

**Refactor-step note to surface in prose (no code change yet):** explicitly call out two known gaps even though the test passes: (1) not safe under concurrency, check and set are separate map operations; (2) doesn't yet replay a *response*, only prevents the duplicate side effect. Naming these now motivates cycles 4 and 5.

### Cycle 4: assert the replayed response, not just "no duplicate" (red)

Tighten cycle 2's test: assert the second call's response body/status matches the first call's exactly, not just that only one email went out.

Fails: the naive map blocks the duplicate side effect but returns nothing meaningful on the second call.

**Motivates:** storing the *result*, not just the fact of having seen the key.

### Cycle 5: extract the interface (refactor-driven design)

Pull the seen-keys-plus-stored-response behaviour behind an interface (e.g. `IdempotencyStore`) before adding more logic inline. Rationale to state explicitly: the map is about to grow real logic (claim semantics, stored responses), and hiding it behind an interface sets up the later Postgres-adapter payoff without a rewrite.

**Interesting refactor options to present, not just pick one:**
- A single two-purpose method (`Claim(key) (claimed bool, existing *Response)`) vs. three smaller methods (`Lookup`, `ReserveInProgress`, `Complete`). The single-method version maps better onto the atomic claim actually needed, avoids reintroducing a check-then-act gap at the interface boundary.
- Whether the store returns the full HTTP response (status, body, allowlisted headers) or a result value the handler re-serializes. Storing the actual response shape is what makes replay *faithful*, ties back to the earlier headers discussion (allowlist `Content-Type`, `Location`, any contractually meaningful custom headers; skip transport-layer headers like `Content-Length`, `Date`).

### Cycle 6: force the concurrent case (red, synctest)

New test: two goroutines call the endpoint with the same key, deliberately synchronized (via synctest or a channel-based rendezvous) so both are *guaranteed* in-flight simultaneously, not just "probably" (a bare sleep-based race is flaky and should be avoided as a teaching example). Assert exactly one email sent, and that the second caller gets a sensible signal rather than silently redoing the work or hanging.

**Motivates:** the naive map (or even the cycle 5 interface, if its implementation still isn't atomic) fails here, revealing the check-then-act race explicitly. Worth showing the failure under `-race` too as a concrete demonstration.

**Deliberate teaching note:** build the naive version with genuinely *no* synchronization first (not even an accidental mutex), so the double-send is guaranteed and honest, before introducing real atomic claim semantics as the fix.

### Cycle 7: real locking semantics in the fake (green)

Fix the in-memory implementation with a mutex around claim, and introduce an explicit "in progress" state (distinct from "seen"/"unseen"), so a concurrent second caller can be told "someone else is already doing this."

**Refactor option, present as genuine design choice:** what does the second, in-flight caller get back?
- Block until the first completes: simplest for the reader, but ties up a goroutine/request.
- Fail fast with a `409`-shaped result: matches the IETF draft's behaviour, teaches the reader the standard's actual guidance, more realistic for a real server not wanting to hold connections open indefinitely.

### Cycle 8: mismatched-payload case (red, stretch/optional)

Test: same key, different request body. Assert the endpoint rejects rather than replaying the first response or processing the new body.

**Motivates:** storing a request hash alongside the key and comparing on lookup. Flag as optional depending on chapter length, important to the pattern but not essential to the core narrative.

### Cycle 9: swap the fake for a Postgres-shaped adapter (refactor, no new test behaviour)

Because cycles 5–7 pushed everything behind the interface, this is presented as "the tests don't change, only the implementation does." Show a Postgres adapter using `INSERT ... ON CONFLICT` (or `SELECT ... FOR UPDATE`) to get the atomicity the in-memory mutex was standing in for. State explicitly that any adapter satisfying the interface's claim semantics works, most ACID databases would do.

**Open decision to make once, consistently, for the whole distributed-systems section (not just this chapter):** whether to actually run this adapter against real Postgres in the book's test suite (e.g. via `testcontainers-go`) or leave it as an illustrative, untested sample. Since this chapter is likely to be written first, whatever's decided here effectively sets the section's infra story for the not-yet-written outbox chapter too, worth deciding deliberately rather than by default.

### Cycle 10 (closing, design discussion, not test-driven): mandatory vs. optional key

No new failing test. Walk through making the `Idempotency-Key` header mandatory for this endpoint (reject with `400`/`422` if absent) vs. optional (proceed without protection if absent).

**Framing:** this is a deliberate per-endpoint policy choice, not a universal rule.
- **Optional** suits operations where duplicates are cheap or self-correcting (e.g. naturally idempotent operations like "set status to X", or low-stakes events). Simpler contract for callers who don't need the guarantee.
- **Mandatory** suits expensive or irreversible side effects (this chapter's email-send example arguably belongs here, even though optional was fine as a gentler starting point) — the server enforces the guarantee rather than trusting callers to opt in. Stripe's mutating endpoints work this way for money-related operations, a good real-world anchor.

State the trade-off explicitly: simpler caller contract vs. the server no longer being able to promise the guarantee unconditionally.

---

## Closing notes: idempotency as an enabling precondition for other patterns

Frame idempotency not as a peer to other resilience patterns but as the **precondition that makes them safe to use at all**.

- **Circuit breakers:** exist to fail fast and let callers retry once a downstream service recovers. That advice is only safe if the retried call can't cause harm on top of whatever partially happened before the breaker tripped. Without idempotency, a circuit breaker just defers the duplicate-send problem rather than solving it. Line to land: *a circuit breaker without an idempotent target is retrying blind.*
- **Retries and backoff generally:** the most direct link. Any retry-with-backoff strategy (an HTTP client wrapper, a queue consumer, or the dispatcher in the upcoming outbox chapter) is fundamentally a retry loop, and every retry loop implicitly assumes calling again is safe. "Just add retries" is bad advice on its own and only becomes good advice once idempotency is established as a property of the target. Worth a forward pointer here: the next chapter builds exactly such a retry loop (the outbox dispatcher) from the caller's side, and it leans on the receiver-side guarantee built in this chapter to be safe.
- **At-least-once messaging / load-balancer failover generally:** the umbrella case. Message queues, failover to a different backend after a failed health check, service-mesh retry policies, all quietly rely on the same guarantee.

**Suggested closing line for the chapter:** at-least-once delivery is the easy half of building reliable distributed systems; idempotent receivers are what make it usable. Most of the interesting resilience patterns (circuit breakers, retries, failover) are really just different strategies for deciding *when* to exercise a guarantee that idempotency is what actually provides. This is also the pitch for why the pattern earns its own chapter in a dedicated section rather than being a footnote elsewhere.