Demonstrate Go Concurrency Patterns.
<br>
## [Simple Pipeline](https://github.com/AhmadMuzakkir/Go-Concurrency-Patterns/blob/master/simple_pipeline.go)

Demonstrate a simple pipeline construction with explicit cancellation.
<br><br>
## [Channel to Function](https://github.com/AhmadMuzakkir/Go-Concurrency-Patterns/blob/master/functioner.go)

Demonstrate how to convert a channel into a blocking function.

Let's say you have a stream of data in a channel but you don't want to expose the channel across package.
Instead, you want it to be a blocking function. e.g. data, err := get()

The function get() should block until there's a new data or an error.

Error is considered terminal. If there's already an error, get() will immediately return the error.
