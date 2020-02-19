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
			want: "[{\"name\":\"Comic Books\",\"cost\":50,\"quantity\":42},{\"name\":\"McDonalds Toys\",\"cost\":75,\"quantity\":120}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ListProducts(tt.args.w, tt.args.r)
			actual, _ := ioutil.ReadAll(httpTestWriter.Body)
			assertStatus(t, httpTestWriter.Code, http.StatusOK)
			assertResponseBody(t, string(actual), tt.want)
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
