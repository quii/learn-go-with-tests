# :bookmark_tabs: Intro

- 此專案目的為藉由測試驅動開發(TDD, Test-Driven Development )的方式，先寫測試再開始的方式進行
- 此專案使用 go 官方的標準測試套件 ( [testing package](https://pkg.go.dev/testing) )，當然實務上可能友會搭配其他好用的第三方套件 ( [stretchr/testify](https://github.com/stretchr/testify) )
- TDD 流程 :

  - 寫測試程式，並執行測試程式，讓他不通過，檢查錯誤訊息是否有意義 ( Write a failing test )
  - 再回去撰寫程式碼，讓測試通過 ( Make the test pass )
  - 重構程式碼 ( Refactor )

    ![TDD](https://marsner.com/wp-content/uploads/test-driven-development-TDD.png)

# :triangular_ruler: 測試規則

1. 必須是 `\_test.go` 結尾 ( 範例 : hello_world_test.go )
2. 測試函數的名稱必須是 `TestHelloWorld` 開頭 ( 範例 : Test )
3. 測試函數只能傳入 `t *testing.T`

# :bug: 錯誤處理方式

### 1. `t.Fail()`

- 是導致當前測試失敗並立即停止執行的函數，不會返回錯誤或提供有關失敗的任何其他信息

### 2. `t.Errorf("\<format error message\>")`

- 跟 `t.Fail()` 類似，但是允許描述失敗的錯誤訊息
- 格式化訊息方式跟 `fmt.Printf("\<format message\>")` 一樣

# :computer: 在測試程式檔案中，將重複程式碼寫到新的函數

- 可使用 `t.Helper()` 在新的函數中，告訴測試套件此函數是 helper，在測試失敗時，錯誤訊息的行數會落在呼叫此函數的地方，這樣方便除錯
- 範例 :

  ```go
  // 傳入 testing.TB 原因是，此介面可讓我們使用 *testing.T 或 *testing.B
  // *testing.B 是用於測試效能(performance)用的
  func assertCorrectMessage(t testing.TB, got, want string) {
    t.Helper()
    if got != want {
      t.Errorf("got %q want %q", got, want)
    }
  }
  ```

# :test_tube: 批次測試不同測試案例

- 使用 `t.Run("\<測試內容說明\>", func(t *testing.T){...})` 方式，可在測試函數中，一次測試多種案例

  ```go
  testcase := []struct {
      name string
      test func(t *testing.T)
    }{
      {
      name: "..."
      test: func( t *testing.T){...}
      },
      {
        ...
      },
      ...
    }

  for i := range testcase {
    tc := testcase[i]
    t.Run(tc.name, tc.test)
  }
  ```

# 🗺 導覽

### Go fundamentals

1. [Install Go](install-go.md) - Set up environment for productivity.
2. [Hello, world](hello-world.md) - Declaring variables, constants, if/else statements, switch, write your first go program and write your first test. Sub-test syntax and closures.
3. [Integers](integers.md) - Further Explore function declaration syntax and learn new ways to improve the documentation of your code.
4. [Iteration](iteration.md) - Learn about `for` and benchmarking.
5. [Arrays and slices](arrays-and-slices.md) - Learn about arrays, slices, `len`, varargs, `range` and test coverage.
6. [Structs, methods & interfaces](structs-methods-and-interfaces.md) - Learn about `struct`, methods, `interface` and table driven tests.
7. [Pointers & errors](pointers-and-errors.md) - Learn about pointers and errors.
8. [Maps](maps.md) - Learn about storing values in the map data structure.
9. [Dependency Injection](dependency-injection.md) - Learn about dependency injection, how it relates to using interfaces and a primer on io.
10. [Mocking](mocking.md) - Take some existing untested code and use DI with mocking to test it.
11. [Concurrency](concurrency.md) - Learn how to write concurrent code to make your software faster.
12. [Select](select.md) - Learn how to synchronise asynchronous processes elegantly.
13. [Reflection](reflection.md) - Learn about reflection
14. [Sync](sync.md) - Learn some functionality from the sync package including `WaitGroup` and `Mutex`
15. [Context](context.md) - Use the context package to manage and cancel long-running processes
16. [Intro to property based tests](roman-numerals.md) - Practice some TDD with the Roman Numerals kata and get a brief intro to property based tests
17. [Maths](math.md) - Use the `math` package to draw an SVG clock
18. [Reading files](reading-files.md) - Read files and process them
19. [Templating](html-templates.md) - Use Go's html/template package to render html from data, and also learn about approval testing
20. [Generics](generics.md) - Learn how to write functions that take generic arguments and make your own generic data-structure
21. [Revisiting arrays and slices with generics](revisiting-arrays-and-slices-with-generics.md) - Generics are very useful when working with collections. Learn how to write your own `Reduce` function and tidy up some common patterns.

### Build an application

Now that you have hopefully digested the _Go Fundamentals_ section you have a solid grounding of a majority of Go's language features and how to do TDD.

This next section will involve building an application.

Each chapter will iterate on the previous one, expanding the application's functionality as our product owner dictates.

New concepts will be introduced to help facilitate writing great code but most of the new material will be learning what can be accomplished from Go's standard library.

By the end of this, you should have a strong grasp as to how to iteratively write an application in Go, backed by tests.

* [HTTP server](http-server.md) - We will create an application which listens to HTTP requests and responds to them.
* [JSON, routing and embedding](json.md) - We will make our endpoints return JSON and explore how to do routing.
* [IO and sorting](io.md) - We will persist and read our data from disk and we'll cover sorting data.
* [Command line & project structure](command-line.md) - Support multiple applications from one code base and read input from command line.
* [Time](time.md) - using the `time` package to schedule activities.
* [WebSockets](websockets.md) - learn how to write and test a server that uses WebSockets.

### Testing fundamentals

Covering other subjects around testing.

* [Introduction to acceptance tests](intro-to-acceptance-tests.md) - Learn how to write acceptance tests for your code, with a real-world example for gracefully shutting down a HTTP server
* [Scaling acceptance tests](scaling-acceptance-tests.md) - Learn techniques to manage the complexity of writing acceptance tests for non-trivial systems.
* [Working without mocks, stubs and spies](working-without-mocks.md) - Learn about how to use fakes and contracts to create more realistic and maintainable tests.
* [Refactoring Checklist](refactoring-checklist.md) - Some discussion on what refactoring is, and some basic tips on how to do it.

### Questions and answers

I often run in to questions on the internets like

> How do I test my amazing function that does x, y and z

If you have such a question raise it as an issue on github and I'll try and find time to write a short chapter to tackle the issue. I feel like content like this is valuable as it is tackling people's _real_ questions around testing.

* [OS exec](os-exec.md) - An example of how we can reach out to the OS to execute commands to fetch data and keep our business logic testable/
* [Error types](error-types.md) - Example of creating your own error types to improve your tests and make your code easier to work with.
* [Context-aware Reader](context-aware-reader.md) - Learn how to TDD augmenting `io.Reader` with cancellation. Based on [Context-aware io.Reader for Go](https://pace.dev/blog/2020/02/03/context-aware-ioreader-for-golang-by-mat-ryer)
* [Revisiting HTTP Handlers](http-handlers-revisited.md) - Testing HTTP handlers seems to be the bane of many a developer's existence. This chapter explores the issues around designing handlers correctly.


# :link: Reference

- [Learn Go with tests](https://quii.gitbook.io/learn-go-with-tests/)
