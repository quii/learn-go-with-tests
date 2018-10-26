# Why TDD?

It's difficult to write about TDD without rehashing what others have said but
it helps _me_ to organise my thoughts around the matter. So in a way this is
a selfish endeavour but I do hope this will at least get readers thinking about
TDD and the important role it has in software development. 

The promise of software is that it can change. This is why it is called _soft_
ware, it is malleable compared to hardware. A great engineering team should be
an amazing asset to a company, writing systems that can evolve with a business
to keep delivering value. 

So why are we so bad at it? How many projects do you hear about that outright
fail? Or become "legacy" and have to be entirely re-written (and the re-writes
often fail too!) How does a software system "fail" anyway? Cant it just be
changed until it's correct? That's what we're promised! 

In 1974, a long time before I was born a clever software engineer called Manny
Lehman wrote something called

## The Law of Continuous Change

> Any software system used in the real-world must change or become less and
> less useful in the environment

It feels obvious that a system _has_ to change or it becomes less useful but
how often is this ignored? 

Many teams are incentivised to deliver a project on a particular date and then
moved on to the next project. If the software is "lucky" there is at least some
kind of hand-off to another set of individuals to maintain it, but didn't write
it of course. 

People often concern themselves with trying to pick a framework which will help
them "deliver quickly" but not focusing on the longevity of the system in terms
of how it needs to evolve.

Even if you're an incredible software engineer, you will still fall victim of
not knowing the future needs of your system. As the business changes some of
the brilliant code you wrote is now no longer relevant. **Software must
change**

Lehman was on a roll in the 70s because he gave us another law to chew on.

## The law of increasing complexity

> As a system evolves, its complexity increases _unless work is done to reduce
> it_

(emphasis mine)

What he's saying here is we cant have software teams as blind feature
factories, piling more and more features on to software in the hope it will
survive in the long run. 

We **have** to keep managing the complexity of the system as the knowledge of
our domain changes. 

## Refactoring

There are _many_ facets of software engineering that keeps software malleable,
such as:

- Developer empowerment
- Communication skills
- Architecture
- Observability
- Deployability
- Feedback loops

I am going to focus on refactoring. It's a phrase that get's thrown around
a lot "we need to refactor this" said to a developer on their first day of
programming without a second thought. Where does the phrase come from? How is
refactoring just different from writing code?

### Factorisation

When learning maths at school you probably learned about factorisation. Here's
a very simple example

Calculate `1/2 + 1/4`

To do this you _factorise_ the denominators, turning the expression into 

`2/4 + 1/4` which you can then turn into `3/4`. 

We can take some important lessons from this. When we _factorise the
expression_ we have **not changed the meaning of the expression**. Both of them
equal `3/4` but we have made it easier for us to work with; by changing `1/2`
to `2/4` it fits into our "domain" easier. 

When you refactor your code, you are trying to find ways of making your code
easier to understand and "fit" into your current understanding of what the
system needs to do. Crucially **you should not be changing behaviour**. 

### When refactoring code you must not be changing behaviour

This is very important. If you are changing behaviour at the same time you are
doing _two_ things at once. As software engineers we learn to break systems up
into different files/packages/functions/etc because we know trying to
understand a big blob of stuff is hard. 

We dont want to have to be thinking about lots of things at once because that's
when we make mistakes. I've witnessed so many refactoring endeavours fail
because the developers are biting off more than they can chew. 

When I was doing factorisations in maths classes with pen and paper I would
have to manually check that I hadn't changed the meaning of the expressions in
my head. How do we know we aren't changing behaviour when refactoring when
working with code, especially on a system that is non-trivial?

Those who choose not to write tests will typically be reliaint on manual
testing. For anything other than a small project this will be a tremendous
time-sink and doesn't scale in the long run. 

**In order to safely refactor you need automated tests**, they provide

- Confidence you can reshape code without worrying about changing behaviour
- Documentation for humans as to how the system should behave
- Much faster and more reliable feedback than manual testing
- In order for code to be testable, it _generally_ has to follow best practices
  of single responsibilities, explicit dependencies (i.e no global variables);
properties that also aid in refactoring.

## Why TDD

Some people might take Lehman's quotes about how software has to change and
overthink elaborate designs, wasting lots of time upfront trying to create the
"perfect" extensible system and end up getting it wrong and going nowhere. 

This is the bad old days of software where an analyst team would spend 6 months
writing a requirements document and an architect team would spend another
6 months coming up with a design and a few years later the whole project fails.
I say bad old days but this still happpens! 

Agile teaches us that we need to work iteratively, starting small and evolving
the software so that we get fast feedback on the design of our software and how
it works with real users;  TDD enforces this approach.

TDD addresses the laws that Lehman talks about and other lessons hard learned
through history by encouraging a methodology of constantly refactoring and
delivering iteratively. 

## Common objections with pithy responses

> **Tests _dont_ help me refactor. Every time i refactor loads of tests stop
> passing/compiling**

Remember what refactoring is _supposed_ to be? Just changing the way your
program is expressed, not changing behaviour. Now ask yourself why your tests
are failing. It will be because your tests are **too coupled to implementation
details**. 

You're _probably_ mocking too much and testing irrelevant detail. Remember
a unit test is _not_ only on functions/classes/whatever. A unit of behaviour
can be tested and it may have a number of internal collaborators to make that
behaviour work; just dont test them!

Listen to your tests and act on what they're telling you.

> **I dont like writing tests as I want to explore the design first, then I write
> my tests afterward.** 

It is hard/time-consuming to write your first test; if your first test is "make
a website to rival twitter". 

Irrespective of whether you practice TDD or not it is an important skill as
a software developer to be able to break problems down into small pieces. This
lets us work in a smaller problem space and deliver small pieces of value
quickly, letting us validate our design assumptions as we work. This is all
about learning from the mistakes of the past with too much work on upfront
design.

The beauty of TDD is it forces us to start small (unless you enjoy spending
loads of time writing a big test without the endorphin rush of seeing a test
pass). By starting small it will challenge your assumptions because you'll get
feedback quicker. 

Writing tests after the fact is usually harder and more error prone.
You are more likely to write code that isn't easy to test because your code has
been driven by assumptions in your head rather than tests demanding a specific
behaviour. 

In addition an important step in TDD is the first one; see how your test fails
and see if the error makes sense. This forces you to write ergonomic tests when
they fail so when developers see a failing test they have an easier time
understanding the problem

Too much of my career has been wasted debugging tests that fail with `false was
not true` 

> **It takes too long**

You should read [GeePaw's TDD & The Lump of Coding
Fallacy](http://geepawhill.org/tdd-and-the-lump-of-coding-fallacy/) as it
explains brilliantly why this line of thinking is wrong (at least once you
become proficient with TDD).

If you're too lazy my TL;DR version is

- You dont actually arrive at your desk at 9:30 and constantly write code
   until 5:30
-  What you do is a mixture of. 1) Yes, writing code. 2) Thinking about code,
   studying existing code. 3) Make a change to the code and run it to see what
happens (e.g spin up the server and see what happens, debugging, etc)
- The premise is the tests you write basically are a part of 2 and 3, but make
  it structured and quicker.

People who enjoy TDD will rarely touch a debugger and dont need to spend too
much time spinning up the software and testing it because they have confidence
it works how it should. 

The "studying" part becomes easier because as GeePaw says

> it’s almost like the test code forms a kind of Cliff’s Notes for the shipping
> code. A scaffolding that makes it easier for us to study, and this makes it
> far easier to tell what’s going on. This will cut our code study time in
> about half.

> **All the examples are unrealistic compared to "real" software**

This comes back to being able to break problems down. As you gain
practice with TDD and software development you'll learn how to break down
problems so that they look like the simple examples you learned with.

If a section of your code is too hard to test; it's not "realistic" - it's
poorly written 
