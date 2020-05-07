# HTTP Handlers Revisited

This book already has a chapter on [testing a HTTP handler](http-server.md) but this will feature a broader discussion on designing them, so they are simple to test.

Testing HTTP handlers seems to be a recurring question in the Go community, and I think it points to a wider problem of people misunderstanding how to design them.

So often people's difficulties with testing stems from the design of their code rather than the actual writing of tests. As I stress so often in this book:

> If your tests are causing you pain, listen to that signal and think about the design of your code.

## An example

[Santosh Kumar tweeted me](https://twitter.com/sntshk/status/1255559003339284481)

> How do I test a http handler which has mongodb dependency?

Here is the code

```go
func Registration(w http.ResponseWriter, r *http.Request) {
	var res model.ResponseResult
	var user model.User

	w.Header().Set("Content-Type", "application/json")

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	defer r.Body.Close()

	// check if there is proper json body or error
	if err := jsonDecoder.Decode(&user); err != nil {
		res.Error = err.Error()
		// return 400 status codes
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Connect to mongodb
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	// Check if username already exists in users datastore, if so, 400
	// else insert user right away
	collection := client.Database("test").Collection("users")
	filter := bson.D{{"username", user.Username}}
	var foundUser model.User
	err = collection.FindOne(context.TODO(), filter).Decode(&foundUser)
	if foundUser.Username == user.Username {
		res.Error = UserExists
		// return 400 status codes
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		res.Error = err.Error()
		// return 400 status codes
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	user.Password = string(pass)

	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		res.Error = err.Error()
		// return 400 status codes
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// return 200
	w.WriteHeader(http.StatusOK)
	res.Result = fmt.Sprintf("%s: %s", UserCreated, insertResult.InsertedID)
	json.NewEncoder(w).Encode(res)
	return
}
```

Let's just list all the things this one function has to do:

1. Do some HTTP stuff, send headers, status codes, etc.
2. Decode the request's body into a `User`
3. Connect to a database (and all the details around that)
4. Query the database and applying some business logic depending on the result
5. Generate a password
6. Insert a record

This is too much.

## What is a HTTP Handler and what should it do ?

Forgetting specific Go details for a moment, no matter what language I've worked in what has always served me well is the separation of concerns and the single responsibility principle.

This can be quite tricky to apply depending on the problem you're solving. What exactly _is_ a responsibility?

The lines can blur depending on how abstractly you're thinking and sometimes your first guess might not be right. That's what refactoring is for!

Thankfully with HTTP handlers I feel like I have a pretty good idea what they should do, no matter what project I've worked on:

1. Accept a Request, parse and validate it.
2. Call some `ServiceThing` to do `ImportantBusinessLogic` with the stuff I got from step 1.
3. Send an appropriate `HTTP` response depending on what `ServiceThing` returns.

I'm not saying every HTTP handler _ever_ should have roughly this shape, but 99 times out of 100 that seems to be the case for me.

When you force this separation of concerns testing these handlers becomes a breeze but more importantly testing `ImportantBusinessLogic` no longer has to concern itself with `HTTP`, you can just test the important stuff!

By separating these concerns it also means you can use `ImportantBusinessLogic` in other contexts without having to modify it, neat.

## Go's Handlers

[`http.HandlerFunc`](https://golang.org/pkg/net/http/#HandlerFunc)

> The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers.

`type HandlerFunc func(ResponseWriter, *Request)`

Reader, take a breath and look at the code above. What do you notice?

**It is a function that takes some arguments**

There's no framework magic, no annotations, no magic beans, nothing.

It's just a function, and we know how to test functions. It fits in nicely with the commentary above.

- It takes a `Request` which is just a bundle of data for us to inspect, parse and validate.
- It has a `ResponseWriter` which lets us write `HTTP` stuff.

### Super basic example test

```go
func Teapot(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusTeapot)
}

func TestTeapotHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	Teapot(res, req)

	if res.Code != http.StatusTeapot {
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusTeapot)
	}
}
```

To test our function, we _call_ it.

`http.ResponseWriter` is an interface which means we can pass in a `httptest.ResponseRecorder` which our function will use to write `HTTP` stuff. The recorder will record (or _spy_ on) what was sent, and then we can make some assertions in our test.

## Wrapping up

Testing Go's HTTP handlers is not challenging, but designing good software can be!

People make the mistake of thinking HTTP handlers are special and throw out good software engineering practices when writing them which then makes them challenging to test.

Reiterating again; **Go's http handlers are just functions**. If you write them like you would other functions, with clear responsibilities, and a good separation of concerns you will have no trouble testing them, and your codebase will be healthier for it.

### Notes: Things to talk about

- Complaints about examples being too simplistic; well design your code so it is simplistic
- They're just functions! Show them being tested in a simple scenario
- Sep of concerns, iterate with some kind of DB
- What http handlers _should_ be responsible for
