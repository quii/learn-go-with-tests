# Refactoring step, starting checklist

Refactoring is a skill that, once practised enough, becomes, in most cases, second nature reasonably easy.

The activity often gets conflated with more significant design changes, but they are separate. Delineating between refactoring and other programming activities is helpful because it allows me to work with clarity and discipline.

## Refactoring vs other activities

Refactoring is just improving existing code and <u>not changing behaviour</u>; therefore, tests shouldn't have to change.

This is why it's the 3rd step of the TDD cycle. Once you have added a behaviour and a test to back it up, refactoring should be an activity which requires no change to your test code. **You're doing something else** if you are "refactoring" some code and having to change tests at the same time.

Many very helpful refactorings are simple to learn and easy to do (your IDE almost entirely automates many) but, over time, become hugely impactful to the quality of our system.

### Other activities, such as "big" design

> So I'm not changing the "real" behaviour, but I must change my tests? What is that?

Let's say you're working on a type and want to improve its code's quality. *Refactoring shouldn't require you to change the tests*, so you can't:

- Change behaviour
- Change method signatures

...as your tests are coupled to those two things, but you can:

- Introduce private methods, fields and even new types & interfaces
- Change the internals of public methods

What if you want to change the signature of a method?

```go
func (b BirthdayGreeter) WishHappyBirthday(age int, firstname, lastname string, email Email) {
	// some fascinating emailing code
}
```

You may feel its argument list is too long and want to bring more cohesion and meaning to the code.

```go
func (b BirthdayGreeter) WishHappyBirthday(person Person)
```

Well, you're **designing** now and must ensure you tread carefully. If you don't do this with discipline, you can make a mess of your code, the test behind it, *and* probably the things that depend on it - remember, it's not just your tests using `WishHappyBirthday`. Hopefully, it's used by "real" code too!

**You should still be able to drive this change with a test first**. You can split hairs over whether this is a "behaviour" change, but you want your method to behave differently.

As this is a behaviour change, apply the TDD process here too. One benefit of TDD is that it gives you a simple, safe, repeatable way of driving behaviour change in your system; why abandon it in these situations just because it *feels* different?

In this case, you'll change your existing tests to use the new type. The iterative, small steps you usually do with TDD to reduce risk and bring discipline & clarity will help you in these situations, too.

Chances are you'll have several tests that call `WishHappyBirthday`; in these scenarios, I'd suggest commenting out all but one of the tests, driving out the change, and then working through the rest of the tests as you see fit.

### Big design

Design can require more significant changes and more extensive conversations and usually has a level of subjectivity to it. Changing the design of parts of your system is usually a longer process than refactoring; nonetheless, you should still endeavour to reduce risk by thinking about how to do it in small steps.

### Seeing the wood for the trees

> [If someone can't **see the wood for the trees** in British English, or can't see the forest for the trees in American English, they are very involved in the details of something and so they do not notice what is important about the thing as a whole.](https://www.collinsdictionary.com/dictionary/english/cant-see-the-wood-for-the-trees)

Talking about the "big" design issues is more accessible when the **underlying code is well-factored**. If you and your colleagues have to spend a significant amount of time mentally parsing a mess of code every time they open a file, what chance do you have to think about the design of the code?

This is why **constant refactoring is so significant in the TDD process**. If we fail to address the minor design issues, we'll find it hard to engineer the overall design of our more extensive system.

Sadly, badly-factored code gets exponentially worse as engineers pile on complexity on top of shaky foundations.

## Starting mental-checklist

**Get in the habit of running through a mental checklist every TDD cycle.** The more you force yourself to practice, the easier it gets. **It is a skill that needs practice.** Remember, each of these changes should not require any change in your tests.

I have included shortcuts for Intellij/Goland, which my colleagues and I use. Whenever I coach a new engineer, I encourage them to try and gain the muscle memory and habit of using these tools to refactor quickly and safely.

### Inline variables

If you create a variable, only for it to be passed on to another method/function:

```go
url := baseURL + "/user/" + id
res, err := client.Get(url)
```

Consider inlining it (`command+option+n`) *unless* the variable name adds significant meaning.

```go
res, err := client.Get(baseURL + "/user/" + id)
```

Don't be _too_ clever about inlining; the goal is not to have zero variables and instead have ridiculous one-liners that no one can read. If you can add significant naming to a value, it might be best to leave it be.

### DRY up values with extract variables

"Don't repeat yourself" (DRY). Using the same value multiple times in a function? Consider extracting and capturing a variable in a meaningful variable name (`command+option+v`).

This helps with readability and makes changing the value easier in future, as you won't have to remember to update multiple occurrences of the same value.

### DRY up stuff in general

[DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself) gets a bad rep these days, with some justification. DRY is one of those concepts that is *too* easy to understand at a superficial level and then gets misapplied.

An engineer can easily take DRY too far, creating baffling, entangled abstractions to save some lines of code rather than the *real* idea of DRY, which is capturing an _idea_ in one place. Reducing the number of lines of code is often a side-effect of DRY, **but it is not the actual goal**.

So yes, DRY can be misapplied, but the extreme opposite of refusing to DRY up anything is also evil. Repeated code adds noise and increases maintenance costs. A refusal to gather related concepts or values into one thing due to fear of DRY misuse causes *different* problems.

So rather than being extremist on either side of "must DRY everything" or "DRY is bad", engage your brain and think about the code you see in front of you. What is repeated? Does it need to be? Does the parameter list look sensible if you encapsulate some repeated code into a method? Does it feel self-documenting and encapsulate the "idea" clearly?

Nine times out of 10, you can look at the argument list of a function, and if it looks messy and confusing, then it is likely to be a poor application of DRY.

If making some code DRY feels hard, you're probably making things more complex; consider stopping.

DRY with care, **but practising this frequently will improve your judgement**. I encourage my colleagues to "just try it" and use source control to get back to safety if it is wrong.

<u>**Trying these things will teach you more than discussing it**</u>, and source control coupled with good automated tests gives you the perfect setup to experiment and learn.

### Extract "Magic" values.

> [Unique values with unexplained meaning or multiple occurrences which could (preferably) be replaced with named constants](https://en.wikipedia.org/wiki/Magic_number_(programming))

Use extract variable (command+option+v) or constant (command+option+c) to give meaning to magic values. This can be seen as the inverse of the inlining refactor. I often find myself "toggling" the code with inline and extract to help me judge what I think reads better.

Remember that extracting repeated values also adds a level of _coupling_. Everything that uses that value is now coupled. Consider the following code:

```go
func main() {
	api1Client := http.Client{
		Timeout: 1 * time.Second,
	}
	api2Client := http.Client{
		Timeout: 1 * time.Second,
	}
	api3Client := http.Client{
		Timeout: 1 * time.Second,
	}
	//etc
}
```

We are setting up some HTTP clients for our application. There are some _magic values_ here, and we could DRY up the `Timeout` by extracting a variable and giving it a meaningful name.

![A screenshot of me extracting variable](https://i.imgur.com/4sgUG7L.png)

Now the code looks like this

```go
func main() {
	timeout := 1 * time.Second
	api1Client := http.Client{
		Timeout: timeout,
	}
	api2Client := http.Client{
		Timeout: timeout,
	}
	api3Client := http.Client{
		Timeout: timeout,
	}
	// etc..
}
```

We no longer have a magic value; we have given it a meaningful name, but we have also made it so all three clients **share the same timeout**. That _may_ be what you want; refactors are quite context-specific, but it's something to be wary of.

If you can use your IDE well, you can do the _inline_ refactor to let the clients have separate `Timeout` values again.

### Make public methods/functions easy to scan

Does your code have excessively long public methods or functions?

Encapsulate the steps in private methods/functions with the extract method (`command+option+m`) refactor.

The code below has some boring, distracting ceremony around creating a JSON string and turning it into an `io.Reader` so that we can `POST` it in an HTTP request.

```go
func (ws *WidgetService) CreateWidget(name string) error {
	url := ws.baseURL + "/widgets"
	payload := []byte(`{"name": "` + name + `"}`)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(payload),
	)
	//todo: handle codes, err etc
}
```

First, use the inline variable refactor (command+option+n) to put the `payload` into the buffer creation.

```go
func (ws *WidgetService) CreateWidget(name string) error {
	url := ws.baseURL + "/widgets"
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer([]byte(`{"name": "`+name+`"}`)),
	)
	// etc
}
```

Now, we can extract the creation of the JSON payload into a function using the extract method refactor (`command+option+m`) to remove the noise from the method.

```go
func (ws *WidgetService) CreateWidget(name string) error {
	url := ws.baseURL + "/widgets"
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		createWidgetPayload(name),
	)
	// etc
}
```

Public methods and functions should describe *what* they do rather than *how* they do it.

> **Whenever I have to think to understand what the code is doing, I ask myself if I can refactor the code to make that understanding more immediately apparent**

-- Martin Fowler

This helps you understand the overall design better, and it then allows you to ask questions about responsibilities:

>  Why does this method do X? Shouldn't that live in Y?

> Why does this method do so many tasks? Can we consolidate this elsewhere?

Private functions and methods are great; they let you wrap up irrelevant how's into whats.

#### But now I don't know how it works!

A common objection to this refactoring, favouring smaller functions and methods composed of others, is that it can make understanding how the code works difficult. My blunt reply to this is

> Have you learned how to navigate codebases using your tooling effectively?

Quite deliberately, as the _writer_ of `CreateWidget`, I do not want the creation of a specific string to be an essential character in the narration of the method. It is distracting, irrelevant noise for the reader 99% of the time.

However, if someone _does_ care, you press `command+b`  (or whatever "navigate to symbol" is for you) on `createWidgetPayload` ... and read it. Press `command+left-arrow` to go back again.

### Move value creation to construction time.

Methods often have to create value and use them, like the `url` in our `CreateWidget` method from before.

```go
type WidgetService struct {
	baseURL string
	client  *http.Client
}

func NewWidgetService(baseURL string) *WidgetService {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	return &WidgetService{baseURL: baseURL, client: &client}
}

func (ws *WidgetService) CreateWidget(name string) error {
	url := ws.baseURL + "/widgets"
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		createWidgetPayload(name),
	)
	// etc
}
```

A refactoring technique you could apply here is, if a value is being created **that is not dependant on the arguments to the method**, then you can instead create a _field_ in your type and calculate it in your constructor function.

```go
type WidgetService struct {
	client          *http.Client
	createWidgetURL string
}

func NewWidgetService(baseURL string) *WidgetService {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	return &WidgetService{
		createWidgetURL: baseURL + "/widgets",
		client:          &client,
	}
}

func (ws *WidgetService) CreateWidget(name string) error {
	req, err := http.NewRequest(
		http.MethodPost,
		ws.createWidgetURL,
		createWidgetPayload(name),
	)
	// etc
}
```

By moving them to construction time, you can simplify your methods.

#### Comparing and contrasting `CreateWidget`

Starting with

```go
func (ws *WidgetService) CreateWidget(name string) error {
	url := ws.baseURL + "/widgets"
	payload := []byte(`{"name": "` + name + `"}`)
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(payload),
	)
	// etc
}

```

With a few basic refactors, driven almost entirely using automated tooling, we resulted in

```go
func (ws *WidgetService) CreateWidget(name string) error {
	req, err := http.NewRequest(
		http.MethodPost,
		ws.createWidgetURL,
		createWidgetPayload(name),
	)
	// etc
}
```

This is a small improvement, but it undoubtedly reads better. If you are well-practised, this kind of improvement will barely take you a minute, and so long as you have applied TDD well, you'll have the safety net of tests to ensure you're not breaking anything. These continuous minor improvements are vital to the long-term health of a codebase.

### Try to remove comments.

> A heuristic we follow is that whenever we feel the need to comment something, we write a method instead.

-- Martin Fowler

Again, the extract method refactor can be your friend here.

## Exceptions to the rule

There are improvements you can make to your code that require a change in your tests, which I would still be happy to put into the "refactoring" bucket, even though it breaks the rule.

A simple example would be renaming a public symbol (e.g., a method, type, or function) with `shift+F6`. This will, of course, change the production and test codes.

However, as it is an **automated and safe** change, the risk of going into a spiral of breaking tests and production code that so many fall into with other kinds of *design* changes is minimal.

For that reason, any changes you can safely perform with your IDE/editor, I would still happily call refactoring.

## Use your tools to help you practice refactoring.

- You should run your unit tests every time you do one of these small changes. We invest time in making our code unit-testable, and the feedback loop of a few milliseconds is one of the significant benefits; use it!
- Lean on source control. You shouldn't feel shy about trying out ideas. If you're happy, commit it; if not, revert. This should feel comfortable and easy and not a big deal.
- The better you leverage your unit tests and source control, the easier to *practice* refactoring. Once you master this discipline, **your design skills increase quickly** because you have a reliable and effective feedback loop and safety net.
- Too often in my career, I've heard developers complain about not having time to refactor; unfortunately, it is clear that it takes so much time for them because they don't do it with discipline - and they have not practised it enough.
- Whilst typing is never the bottleneck, you should be able to use whatever editor/IDE you use to make refactoring safely and quickly. For instance, if your tool doesn't let you extract variables at a keystroke, you'll do it less because it's more labour-intensive and risky.

## Don't ask permission to refactor

Refactoring should be a frequent occurrence in your work, something you're doing all the time. It also, shouldn't be a time-sink, especially if it's done little and often.

If you don't refactor, your internal quality will suffer, your team's capacity will drop, and pressure will increase.

Martin Fowler has one more fantastic quote for us.

> Other than when you are very close to a deadline, however, you should not put off refactoring because you haven’t got time. Experience with several projects has shown that a bout of refactoring results in increased productivity. Not having enough time usually is a sign that you need to do some refactoring.

## Wrap up

This is not an extensive list, just a start. Read Martin Fowler's Refactoring book (2nd ed) to become a pro.

Refactoring should be extremely quick and safe when you're well-practised, so there's little excuse not to do it. Too many view refactoring as a decision for others to make rather than a skill to learn to where it's a regular part of your work.

We should always strive to leave code in an *exemplary* state.

Good refactoring leads to code that is easier to understand. An understanding of the code means better designs are easier to spot. It is much harder to find designs in systems with massive functions, needlessly duplicated code, deep nesting, etc. **Frequent, small refactoring is necessary for better design**.

