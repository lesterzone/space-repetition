# memory
Will review cards based on sample rules of increasing time when answer is
correct and decreasing interval when answer is wrong.

When answer is correct, let's add 50% of time since last execution
When answer is incorrect, let's add 25% of time since last execution.

# Install leaf to watch file changes and run test
- go install github.com/vrongmeal/leaf/cmd/leaf@latest

# watch and run tests
```
leaf
```

# build / install

This is a stateless / single executable and forever runner CLI app.

```shell
go build // or just `go install`
go build -ldflags "-s -w"
```

# TODO

Makefile