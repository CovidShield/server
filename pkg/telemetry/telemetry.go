package telemetry

import (
	"net/http"

	"github.com/Shopify/goose/logger"
	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	metricstdout "go.opentelemetry.io/otel/exporters/metric/stdout"
	"go.opentelemetry.io/otel/exporters/trace/stdout"
	tracerstdout "go.opentelemetry.io/otel/exporters/trace/stdout"
	"go.opentelemetry.io/otel/plugin/httptrace"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var log = logger.New("telemetry")

// Providers

const STDOUT = "STDOUT"
const PROMETHEUS = "PROMETHEUS"

// InitTracer initializes the global trace provider
func InitTracer(provider string) func() {
	// Some providers require cleanup
	cleanupFunc := func() {}

	var exporter *tracerstdout.Exporter
	var err error
	switch provider {
	case STDOUT:
		exporter, err = tracerstdout.NewExporter(stdout.Options{PrettyPrint: true})
	default:
		log(nil, nil).WithField("provider", provider).Fatal("Unsuported trace provider")
	}

	if err != nil {
		log(nil, err).WithField("provider", provider).Fatal("failed to initialize exporter")
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

// InitMeter initializes the global metric progider
func InitMeter(provider string) func() {
	cleanupFunc := func() {}

	var err error
	switch provider {
	case STDOUT:
		var pusher *push.Controller
		pusher, err = metricstdout.InstallNewPipeline(metricstdout.Config{
			Quantiles:   []float64{},
			PrettyPrint: true,
		}, push.WithStateful(false))
		if err != nil {
			break
		}
		cleanupFunc = pusher.Stop
	case PROMETHEUS:
		var exporter *prometheus.Exporter
		exporter, err = prometheus.InstallNewPipeline(prometheus.Config{})
		if err != nil {
			break
		}
		http.HandleFunc("/", exporter.ServeHTTP)
		go func() {
			_ = http.ListenAndServe(":2222", nil)
		}()
	default:
		log(nil, nil).WithField("provider", provider).Fatal("Unsuported metric provider")
	}

	if err != nil {
		log(nil, err).WithField("provider", provider).Fatal("failed to initialize metric stdout exporter")
	}

	initSystemStatsObserver()

	return cleanupFunc
}

// OpenTelemetryMiddleware adds monitoring around HTTP requests
func OpenTelemetryMiddleware(next http.Handler) http.Handler {
	tracer := global.Tracer("covidshield/request")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attrs, entries, spanCtx := httptrace.Extract(r.Context(), r)

		r = r.WithContext(correlation.ContextWithMap(r.Context(), correlation.NewMap(correlation.MapUpdate{
			MultiKV: entries,
		})))
		_, span := tracer.Start(
			trace.ContextWithRemoteSpanContext(r.Context(), spanCtx),
			"HTTP Request",
			trace.WithAttributes(attrs...),
		)
		defer span.End()
		next.ServeHTTP(w, r)
	})
}
