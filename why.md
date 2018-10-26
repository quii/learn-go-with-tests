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

> Any software system used in the real-world must change or become less and less useful in the environment

It feels obvious that a system _has_ to change or it becomes less useful but
how often is this ignored? 

Many teams are incentivised to deliver a project on
a particular date and then moved on to the next project. If the software is
"lucky" there is at least some kind of hand-off to another set of individuals
to maintain it, but didn't write it of course. 

People often concern themselves with trying to pick a framework which will help
them "deliver quickly" but not focusing on the longevity of the system in terms
of how it needs to evolve.

Even if you're an incredible software engineer, you will still fall victim of
not knowing the future needs of your system. As the business changes some of the
brilliant code you wrote is now no longer relevant. **Software must change**

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

We dont want to have to be thinking about lots of things at once because that's when we make mistakes. I've
witnessed so many refactoring endeavours fail because the developers are biting
off more than they can chew. 

When I was doing factorisations in maths classes with pen and paper I would
have to manually check that I hadn't changed the meaning of the expressions in
my head. How do we know we aren't changing behaviour when refactoring when
working with code, especially on a system that is non-trivial?

Those who choose not to write tests will typically be reliaint on manual
testing. For anything other than a small project this will be a tremendous
time-sink and doesn't scale in the long run. 

**In order to safely refactor you need automated tests**, they
provide

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
the software and TDD enforces this approach.

TDD addresses the laws that Lehman talks about and other lessons hard learned
through history by encouraging a methodology of constantly refactoring and delivering
iteratively. 


