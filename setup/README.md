# Setup

To install go, follow the instructions on [The Official Install Page](https://golang.org/doc/install).

For MacOs, you can alternately run

```
brew install go
```

# IDE

Go can be written in any text editor without too much hassle, but an IDE will make your life a bit easier if you choose to use one. [VSCode](https://code.visualstudio.com/) with the [ms-vscode.go plugin](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go) is the most popular free solution. [Jetbrain's Goland](https://www.jetbrains.com/go/) is another option a lot of go developers love, but is not free.

# Tooling
The setup in this directory is entirely optional. It is not required to get up and running for workshop 0, but will help you to write better and more idiomatic code.

## Basic commands

### Getting a new tool

To get a new tool in go, you typically just need to run something like 

```
GO111MODULE=on go get -u golang.org/x/tools/gopls@v0.1
```

So what's happening here? GO111MODULE=on is telling go to function in module mode rather than gopath mode, allowing us to use branches other than master. The `go get` command will fetch module requested and install it into your $GOPATH. The `-u` is telling go to check and fetch the latest version of the module and its dependencies within the specified version, even if we've already downloaded this module before. The `golang.org/x/tools/gopls` part is the module we're specifically requesting. And the `@v0.1` is telling `go get` that we specifically want the version 0.1, but not specifying a patch number (go will by default select the highest available).

Something interesting to note about `go get` is that it's behavior changes depending on where you run it from. In order to download a tool, you must run it from outside an existing go project. When you run it from inside a project, it will actually add the module as a dependency to the project. This is incredibly useful, but can be frustrating when you don't realize it's happening.

For more information [click me](https://dev.to/maelvls/why-is-go111module-everywhere-and-everything-about-go-modules-24k)

### Viewing The Environment

You can get a good understanding of how go commands are going to run by checking its environment:

```
go env  # prints the entire environment
go env GOPATH GOOS GOARCH # prints the specific values for the go path, os, and architecture used to compile binaries
```

### Checking Documentation

While many IDEs can handle checking documentation for you, it may be helpful to understand how to do it yourself. You can certainly find most documentation on the internet, but you can also check the documentation on most packages by running:

```
go doc strings  # Show the high level documentation for the package
go doc strings.Replace # Show the documentation for the specific method
go doc -all strings # Show all the documentation for the package
go doc -src strings.Replace # Show the source code
```

### Testing

To test code, you can run:

```
go test . # Test the code in the current directory
go test ./... # Test the code in current and sub-directories
go test -race ./... # Test the code with the race flag, looking for race conditions
go test -bench=. ./... # Run benchmarks
```

### Formatting

Go has specific rules around formatting of code, but also comes with a tool to help you make sure you're following those rules. To format your code, you can run:

```
gofmt -l -w -s . # Formats all *.go files in the current directory and subdirectories
```

You can also run `go fmt` which is effectively a wrapper around `gofmt -l -w`, but it doesn't support the `-s` option, which finds places to simplify your code.

The gofmt tool also offers some really cool tricks for refactoring your code. [You can find out more here](https://blog.golang.org/go-fmt-your-code)

## Static Analysis and Linting

Go comes with some tools to detect areas that could introduce bugs into your code through static analysis and linting. Additionally, there are a number of third party modules that can improve that process further, and allow more customization in rules.

### Go Vet

The go vet tool examines source code and reports areas that may be problematic but do not specifically break a compile. To run it:

```
go vet ./...
```

### Golint

While it doesn't come with go explicitly, there is a supported linter than looks for stylistic mistakes. These differ from the analysis mistakes in that they tend to be more focused on making your code readable rather than directly exposing potential bugs.

To get the linter:

```
(cd /tmp && GO111MODULE=on go get golang.org/x/lint/golint)
```

And then it can be run with:

```
golint ./...
```

### Golangci-lint

The most commonly used 3rd party linter is `golangci-lint`, a fork of the once popular but now deprecated `gometalinter`.

For macOs, installing golangci-lint is as simple as

```
brew install golangci-lint
```

On linux:

```
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.21.0
```

To run:

```
golangci-lint linters   # shows the currently enabled and disabled linters
golangci-lint run ./... # runs the linters over the current directory and subdirectories
golangci-lint run --fix ./...   # attempts to fix the problems the lintrs detect
```

One of the nice things about golangci-lint is that it's highly configurable. The configuration for this project can be found in the `.golangci.yaml` file in the top level directory. This can be especially valuable when choosing to ignore specific directories or apply special rules for a project or organization. The project adds new rules regularly, and in the near future will hopefully allow implementation of custom rules (currently on PR).

[Click here for more information](https://github.com/golangci/golangci-lint)

## Debugging

Most IDEs offer debugging built in, but in some situations that may not be sufficient, and not everyone wants to work within the confines of an IDE. [Delve](https://github.com/go-delve/delve) tends to be the debugger that most people use outside of IDEs. It's pretty simple and strait forward to use.

## Public Release

### Go Get

If you're putting your code on a public repo, your code is automatically available via the `go get` tool. By default, consumers will get your master branch, but you can improve that by adding version tags. Similarly, code that you store in private repos can be retrieved by anyone with access via `go get`, but they'll need to list the repo in the GOPRIVATE field of their `go env` to make sure the dependencies are pulled correctly.

### Goreleaser

Another great option for releasing code is [Goreleaser](https://github.com/goreleaser/goreleaser). Using goreleaser will automatically build binaries for several platforms and allow you to upload those binaries to common distribution systems, like brew.


# Additional resources

Most of the commands that go runs directly can offer additional information about how to ues them with `go help`, for example:

```
go help vet
go help test
go help fmt
```

There's a pretty nice cheat sheet of commands [here](https://github.com/fedir/go-tooling-cheat-sheet/blob/master/go-tooling-cheat-sheet.pdf)