# 1. Startup

- Startup an HTTP server.
- Respond to all requests with "Hello world"
- Use Postman or similar to make a request.
---
# Workshop 1

In this workshop we will build a simple, basic HTTP server. This assumes you have no experience with Go's net/http library.

## Getting Started

Create a file to hold your main function. It doesn't matter where, the root of this repository is fine. Write a main function and try running it with `go run`. Try building it with `go build`. You may wish to output some text to the console so that you know for sure that things are working.

<details><summary>Example</summary>

```go
package main

import "fmt"

func main() {
       fmt.Println("hello world")
}

```

</details>

## Make an HTTP listener

To make an HTTP server, the `net/http` package needs to be imported. The `ListenAndServe` function in that package starts a TCP listener that responds to HTTP requests.

If you are using Visual Studio Code, you can type `http.ListenAndServe` and press cmd-. to automatically add `net/http` to the imports.

Pick a port on which to listen and pass nil as the handler for now.

<details><summary>Example</summary>

```go
package main

import (
    "net/http"
)

func main() {
    http.ListenAndServe(":8000", nil)
}
```

</details>

Run it again and try to connect with a web browser or curl.

```sh
curl http://127.0.0.1:8000/
```

Go is generating a 404 response for us. Use `curl -D -` to see the HTTP headers and notice which headers it created for us.

## Add a response

It is always a good idea to read the documentation for functions that we use. One of Go's strengths is the excellently documented standard library. Read https://golang.org/pkg/net/http/#ListenAndServe

We can use something called the default serve mux. We don't go into what this means in this workshop. We save that for a later workshop. Expand the example on the godocs if you would like to see the next step of our code, or use `http.HandleFunc` and write the code yourself.

## Logging the error

One important part of writting Go is to leave no error unchecked, usually. If a function returns an error, that error should probably be part of control flow and maybe logged.

You may have noticed that `http.ListenAndServe` returns an error. It is probably safe to ignore any error that comes from it, but what if we had fat fingered our listen port as `";8000"` in our code instead of `":8000"`? What do you think would happen?

Lets use the `log` package from the standard library to log this error. `log.Fatal` will end the program when called.

Now with that error we get notified that the listen port is incorrect:

```
2019/12/15 12:54:19 listen tcp: address ;8000: missing port in address
exit status 1
```

The `log` package is nice enough to timestamp each message.

## Complete

The end program should look something like this:

<details><summary>Expand</summary>

```
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello World")
    }))
    log.Fatal(http.ListenAndServe(":8000", nil))
}
</details>
```

</details>

## Bonus Topics

### What about that weird nil argument to ListenAndServe?

Most experience Go programs that I (Jay) know disagree with that once sentence in the net/http documentation which says `The handler is typically nil, in which case the DefaultServeMux is used.`

Without explaining muxers, it is important to know this about the DefaultServeMux, there can be some unexpected magic associated with that that some gophers feel isn't very Go-like. It may be best to look at an example. Add this to the imports section of our webserver `_ "net/http/pprof"`. Now rebuild, restart and visit the `/debug/pprof/` route.

Some feel that only importing a package shouldn't have side effects and should certainly not be adding a route.

Your homework: read the `net/http` docs to:

- create a mux
- add your handler to that mux
- and use that mux with ListenAndServe
