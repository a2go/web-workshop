package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	logger := log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	err := runServer(logger)
	if err == nil {
		logger.Println("finished clean")
		os.Exit(0)
	} else {
		logger.Printf("Got error: %v", err)
		os.Exit(1)
	}
}

func runServer(logger *log.Logger) error {
	httpServer := NewHTTPServer(logger)
	// make a buffered channel for Signals
	quit := make(chan os.Signal, 1)

	// listen for all interrupt signals, send them to quit channel
	signal.Notify(quit,
		os.Interrupt,    // interrupt = SIGINT = Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)

	// receive signals on quit channel, tell server to shutdown
	go func() {
		//cleanup: on interrupt shutdown webserver
		<-quit
		err := httpServer.Shutdown(context.Background())

		if err != nil {
			logger.Printf("An error occurred on shutdown: %v", err)
		}
	}()

	// listen and serve blocks until error or shutdown is called
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// NewHTTPServer is factory function to initialize a new server
func NewHTTPServer(logger *log.Logger) *http.Server {
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8080"
	}

	s := &ServerHandler{todos: &TodoList{}}
	// pass logger
	s.SetLogger(logger)

	h := &http.Server{
		Addr:         addr,
		Handler:      s,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return h
}

// ServerHandler implements type http.Handler interface, with our logger
type ServerHandler struct {
	logger *log.Logger
	mux    *http.ServeMux
	tmpl   *template.Template
	todos  *TodoList
	once   *sync.Once
}

// SetLogger provides external injection of logger
func (s *ServerHandler) SetLogger(logger *log.Logger) {
	s.logger = logger
}

// ServeHTTP satisfies Handler interface, sets up the Path Routing
func (s *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	once := s.initOnce()
	// on the first request only, lazily initialize
	once.Do(s.RegisterHandlers)
	s.mux.ServeHTTP(w, r)
}

func (s *ServerHandler) initOnce() *sync.Once {
	if s.once == nil {
		s.once = &sync.Once{}
	}
	return s.once
}

func (s *ServerHandler) RegisterHandlers() {
	if s.logger == nil {
		s.logger = log.New(os.Stdout,
			"INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile)
		s.logger.Printf("Default Logger used")
	}
	s.tmpl = template.Must(template.ParseFiles("template.html"))
	s.mux = http.NewServeMux()
	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	s.mux.HandleFunc("/health", HealthCheck)
	s.mux.HandleFunc("/", s.TodoForm)
}

// HealthCheck verifies externally that the program is still responding
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(200)
}

// RedirectToHome Will Log the Request, and respond with a HTTP 303 to redirect to /
func (s *ServerHandler) RedirectToHome(w http.ResponseWriter) {
	w.Header().Add("location", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (s ServerHandler) TodoForm(w http.ResponseWriter, r *http.Request) {
	s.logger.Printf("TodoForm method %v request %v to /", r.Method, r.RequestURI)
	s.logger.Printf("TodoForm Items before %v %v", s.todos.items, s.todos.nextId)
	switch r.Method {
	case http.MethodGet:
		s.logger.Printf("Get Items %v", s.todos.items)
		s.tmpl.Execute(w, s.todos.items)
	case http.MethodPost:
		item := r.FormValue("item")
		s.logger.Printf("Item %v", item)
		id, _ := strconv.ParseInt(item, 10, 64)
		switch r.URL.Path {
		case "/done":
			s.todos.Check(id)
		case "/not-done":
			s.todos.UnCheck(id)
		case "/delete":
			s.todos.Delete(id)
		default:
			s.todos.Add(item)
		}
		s.RedirectToHome(w)
		s.logger.Printf("TodoForm Items after %v %v", s.todos.items, s.todos.nextId)
	}
}

type Todo struct {
	Id        int64
	Title     string
	Completed bool
}

type TodoList struct {
	items  []Todo
	nextId int64
}

func (s *TodoList) Add(name string) {
	s.items = append(s.items, Todo{s.nextId, name, false})
	s.nextId++
}

func (s *TodoList) Check(id int64) {
	for i, item := range s.items {
		if item.Id == id {
			s.items[i].Completed = true
		}
	}
}

func (s *TodoList) UnCheck(id int64) {
	for i, item := range s.items {
		if item.Id == id {
			s.items[i].Completed = false
		}
	}
}

func (s *TodoList) Delete(id int64) {
	var newList []Todo
	for _, item := range s.items {
		if item.Id != id {
			newList = append(newList, item)
		}
	}
	s.items = newList
}
