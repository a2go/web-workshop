# Workshop 1
## Getting started by knowing when to quit!
In the course of development, it is important to have a rapid feedback cycle.
When you are building a web application, you want to run it, check it out, and then _**stop**_ it.

You need some way to signal to your program that it should exit. Go by default will kill your program as soon as it receives any quit signals from the operating system.

<details><summary>Example</summary>

Make a file called "main.go" and put this in it:
```go
package main

import (
	"log"
	"os"
	"time"
)

func main() {
	// Set up logging format
	logger := log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Starting up! Press CTRL-C to stop!")

	for {
		time.Sleep(1 * time.Second)
		log.Printf("Tick.")
	}
	logger.Printf("I never get here!")
}
```
When you `go run main.go` in the terminal, if you type "Control-C", your operating system will send an interrupt signal to your program. By default, golang will listen to these signals and stop your program for you.

</details>

This might be fine, but what if our program was in the middle of something important and we only want it to stop when it's safely done?
We need to explicitly listen for that operating signal, and then Go will let our program pick when and how to exit.

<details><summary>Example</summary>

Let's alter our `main.go` file to be this:
```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Starting up! Press CTRL-C to stop!")

	// make a quit channel for operating system signals, buffered to size 1
	quit := make(chan os.Signal, 1)

	// listen for all interrupt signals, send them to quit channel
	signal.Notify(quit,
		os.Interrupt,    // interrupt = SIGINT = Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)
	logger.Printf("Just going to wait here until you press control-C")
	// block, waiting for receive on quit channel
	sig := <-quit
	logger.Printf("Shutting down after receiving %v signal!", sig)
	// 0 means no errors
	os.Exit(0)
}
```
Here we are making a [channel](https://tour.golang.org/concurrency/2), a typed conduit through which you can send and receive values.

This particular channel is typed for Operating System signals, and can hold a buffer of up to 1 quit signal at a time.

The `<-` is the Go channel `receive` operator, and it will block the program there until a signal on the quit channel is sent.

</details>

Once signal.Notify is called, Go's default exit behavior is off, and it will leave it to our program to exit when it can safely do so.

Now how can we test this behavior so we don't accidentally break it when we change other things?
Getting good feedback on changes is most valuable when you’re trying to make a lot of them! Moreover, not only is the system implementation changing, but the very definition of desired behavior often is as well, as the product evolves.

<details><summary>Example</summary>

Let's change our `main.go` again to extract the Waiter function:
```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Starting up! Press CTRL-C to stop!")

	Waiter(logger)
	os.Exit(0)
}

func Waiter(logger *log.Logger) {
	quit := make(chan os.Signal, 1)

	// listen for all interrupt signals, send them to quit channel
	signal.Notify(quit,
		os.Interrupt,    // interrupt = SIGINT = Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)
	logger.Printf("Just going to wait here until you press control-C")
	// block, waiting for receive on quit channel
	sig := <-quit
	logger.Printf("Shutting down after receiving %v signal!", sig)
	// 0 means no errors
}
```
This extracts the waiting functionality into something that can be called from a test, unlike `main()`.

Now add a new file called `main_test.go`:

```go
package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestWaiter(t *testing.T) {
	t.Run("Wait with func", func(t *testing.T) {
		var finished bool
		// Get the operating system process
		proc, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Fatal(err)
		}
		// Discard noisy logs
		logger := log.New(ioutil.Discard, "", log.LstdFlags)
		go func() {
			Waiter(logger)
			finished = true
		}()
		// if we signal too early, Waiter isn't listening yet
		time.Sleep(10 * time.Millisecond)
		//Send the SIGQUIT
		proc.Signal(syscall.SIGQUIT)
		// if we test finished too early, finished may not have been updated yet
		time.Sleep(10 * time.Millisecond)
		//reset signal notification
		signal.Reset()
		if !finished {
			t.Error("Waiter Did Not Exit")
		}
	})
}
```
If we run `go test`, then it will verify that Waiter listens, blocks, receives, and then returns.
</details>

### Making a Todo web application

This branch starts with a simple, basic HTTP Server with graceful shutdown that
provides a Todo List via HTML forms.

Your job is to write tests first and then extend the existing code to satisfy the basic Todo CRUDL API for a new JSON front-end:

| HTTP Method | URI | Content Type |
|---|---|---|
| DELETE | "/api" | "application/json" |
| GET | "/api" | "application/json" |
| OPTIONS | "/api" | "application/json" |
| POST | "/api" | "application/json" |
| DELETE | "/api/:id" | "application/json" |
| GET | "/api/:id" | "application/json" |
| PATCH | "/api/:id" | "application/json" |

You will need to [serialize and deserialize JSON data](https://gobyexample.com/json)
and add some additional routes and handlers.

The application has this basic data model:

```
type Todo struct {
	Id        int64  `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Completed bool   `json:"completed,omitempty"`
}

type TodoList struct {
	items  []Todo
	nextId int64
}
```

Note: Yeah, uuids are way better than integers, but this is simpler for now.

An example Todo List API HTTP Handler could be:
```
// ListTodos is an HTTP Handler for returning a list of Todos.
func ListTodos(w http.ResponseWriter, r *http.Request) {
	list := []Todo{
		{Id: 1, Title: "Bring home milk", false},
		{Id: 2, Title: "Drink more milk", false},
	}

	data, err := json.Marshal(list)
	if err != nil {
		log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing result", err)
	}
}
```

# Workshop 0

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
    http.ListenAndServe(":8080", nil)
}
```

</details>

Run it again and try to connect with a web browser or curl.

```sh
curl http://127.0.0.1:8080/
```

Go is generating a 404 response for us. Use `curl -D -` to see the HTTP headers and notice which headers it created for us.

## Add a response

It is always a good idea to read the documentation for functions that we use. One of Go's strengths is the excellently documented standard library. Read https://golang.org/pkg/net/http/#ListenAndServe

We can use something called the default serve mux. We don't go into what this means in this workshop. We save that for a later workshop. Expand the example on the godocs if you would like to see the next step of our code, or use `http.HandleFunc` and write the code yourself.

## Logging the error

One important part of writting Go is to leave no error unchecked, usually. If a function returns an error, that error should probably be part of control flow and maybe logged.

You may have noticed that `http.ListenAndServe` returns an error. It is probably safe to ignore any error that comes from it, but what if we had fat fingered our listen port as `";8080"` in our code instead of `":8080"`? What do you think would happen?

Lets use the `log` package from the standard library to log this error. `log.Fatal` will end the program when called.

Now with that error we get notified that the listen port is incorrect:

```
2019/12/15 12:54:19 listen tcp: address ;8080: missing port in address
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
    log.Fatal(http.ListenAndServe(":8080", nil))
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
