package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	httpTestWriter := httptest.NewRecorder()
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No specific path",
			args: args{
				w: httpTestWriter,
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Path: "/",
					},
				},
			},
			want: "You asked to GET /\n",
		},
		{
			name: "Posting to /foo/bar",
			args: args{
				w: httpTestWriter,
				r: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Path: "/foo/bar",
					},
				},
			},
			want: "You asked to POST /foo/bar\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Echo(tt.args.w, tt.args.r)
			actual, _ := ioutil.ReadAll(httpTestWriter.Body)
			assertStatus(t, httpTestWriter.Code, http.StatusOK)
			assertResponseBody(t, tt.want, string(actual))
		})
	}
}

func TestWaiter(t *testing.T) {
	t.Run("Signal runserver graceful shutdown", func(t *testing.T) {
		var finished bool
		// Get the operating system process
		proc, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Fatal(err)
		}
		// Discard noisy logs
		logger := log.New(ioutil.Discard, "", log.LstdFlags)
		go func() {
			runServer(logger)
			finished = true
		}()

		if finished {
			t.Error("runServer Exit before signal sent")
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
			t.Error("runServer Did Not Exit")
		}
	})
}

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Echo)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	want := "text/plain; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != want {
		t.Errorf("handler returned wrong status code: got %v want %v",
			contentType, want)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
