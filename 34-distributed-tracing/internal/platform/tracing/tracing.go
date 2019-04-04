package tracing

import (
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"
)

type Exporter interface {
	trace.Exporter
	Close() error
}

func NewExporter(name, url, local string) (Exporter, error) {

	localEndpoint, err := openzipkin.NewEndpoint(name, local)
	if err != nil {
		return nil, err
	}
	reporter := zipkinHTTP.NewReporter(url)

	exp := zipkin.NewExporter(reporter, localEndpoint)

	type closer interface {
		Close() error
	}

	ret := struct {
		*zipkin.Exporter
		closer
	}{
		exp,
		reporter,
	}

	return ret, nil
}
