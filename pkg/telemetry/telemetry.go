package telemetry

import (
	"net/http"

	"github.com/Shopify/goose/logger"
	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/trace/stdout"
	tracerstdout "go.opentelemetry.io/otel/exporters/trace/stdout"
	"go.opentelemetry.io/otel/plugin/httptrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var log = logger.New("telemetry")

const STDOUT = "STDOUT"

func InitTracer(provider string) func() {
	// Some providers require cleanup
	cleanupFunc := func() {}

	var exporter *tracerstdout.Exporter
	var err error
	switch provider {
	case STDOUT:
		exporter, err = tracerstdout.NewExporter(stdout.Options{PrettyPrint: true})
	default:
		log(nil, nil).WithField("provider", provider).Fatal("Unsuported provider")
	}

	if err != nil {
		log(nil, err).Fatal("failed to initialize exporter")
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log(nil, err)
	}
	global.SetTraceProvider(tp)

	return cleanupFunc
}

func OpenTelemetryMiddleware(next http.Handler) http.Handler {
	tracer := global.Tracer("request")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attrs, entries, spanCtx := httptrace.Extract(r.Context(), r)

		r = r.WithContext(correlation.ContextWithMap(r.Context(), correlation.NewMap(correlation.MapUpdate{
			MultiKV: entries,
		})))
		ctx, span := tracer.Start(
			trace.ContextWithRemoteSpanContext(r.Context(), spanCtx),
			"hello",
			trace.WithAttributes(attrs...),
		)
		defer span.End()
		span.AddEvent(ctx, "handling this...")
		next.ServeHTTP(w, r)
	})
}
