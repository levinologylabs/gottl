// Package otel contains configuration and setup for telemetry
package otel

import (
	"context"
	"errors"
	"sync"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

type Config struct {
	ServiceName string
	Enabled     bool `conf:"default:false"`
	Address     string
}

type OtelService struct {
	logger        zerolog.Logger
	cfg           Config
	shutdownFuncs []func(context.Context) error
	tp            *sdktrace.TracerProvider
	mp            *sdkmetric.MeterProvider
	rp            *sdkresource.Resource
}

var (
	resource          *sdkresource.Resource
	initResourcesOnce sync.Once
)

func NewOtelService(ctx context.Context, logger zerolog.Logger, cfg Config) *OtelService {
	os := &OtelService{
		logger: logger,
		cfg:    cfg,
	}
	os.init(ctx)
	return os
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func (os *OtelService) init(ctx context.Context) {
	// If not enabled return a noop provider
	if !os.cfg.Enabled {
		tp := noop.NewTracerProvider()
		otel.SetTracerProvider(tp)
		os.shutdownFuncs = append(os.shutdownFuncs, func(ctx context.Context) error { return nil })
	}

	os.tp = initTracerProvider(ctx, os.logger)
	os.mp = initMeterProvider(ctx, os.logger)
	os.rp = initResource(ctx)

	os.shutdownFuncs = append(os.shutdownFuncs, os.tp.Shutdown, os.mp.Shutdown)
}

// Shutdown calls cleanup functions registered via shutdownFuncs.
// The errors from the calls are joined.
// Each registered cleanup will be invoked once.
func (os *OtelService) Shutdown(ctx context.Context) error {
	var err error
	for _, fn := range os.shutdownFuncs {
		err = errors.Join(err, fn(ctx))
	}
	os.shutdownFuncs = nil
	return err
}

func initTracerProvider(ctx context.Context, logger zerolog.Logger) *sdktrace.TracerProvider {
	exporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())
	if err != nil {
		logger.Fatal().Msgf("new otlp trace exporter failed: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource(ctx)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func initMeterProvider(ctx context.Context, logger zerolog.Logger) *sdkmetric.MeterProvider {
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		logger.Fatal().Msgf("new otlp metric http exporter failed: %v", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(initResource(ctx)),
	)
	otel.SetMeterProvider(mp)
	return mp
}

func initResource(ctx context.Context) *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			ctx,
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
		)
		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return resource
}
