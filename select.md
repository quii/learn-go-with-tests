# Select (WIP)

You have been asked to make a function called `WebsiteRacer` which takes two URLs and "races" them by hitting them with a HTTP GET and returning the URL which returned first. If none of them return within 10 seconds then it should return an `error`

For this we will be using

- `net/http` to make the HTTP calls.
- _Dependency injection_ with _mocking_ to let us control our tests, keep them fast and test edge cases.
- `go` routines.
- `select` to synchronise processes. 