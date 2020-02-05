package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
