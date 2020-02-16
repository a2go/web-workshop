# 2. Shutdown

- Catch OS signals
- Gracefully shutdown HTTP server
- Add sleep to demonstrate a long-running request


## Links:

- [The complete guide to Go net/http timeouts](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
- [So you want to expose Go on the Internet](https://blog.cloudflare.com/exposing-go-on-the-internet/)


## File Changes:

```
Modified main.go
```
---
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
	t.Run("Signal Waiter graceful shutdown", func(t *testing.T) {
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

		if finished {
			t.Error("Waiter Exit before signal sent")
		}

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
