// Package otel contains configuration and setup for telemetry
package otel

import (
	"context"
	"errors"
	"sync"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	ServiceName string `conf:"default:gottl"`
	Enabled     bool   `conf:"default:true"`
	OTLPAddress string `conf:"default:localhost:4317"`
	Secure      bool   `conf:"default:false"`
}

type OtelService struct {
	logger         zerolog.Logger
	cfg            Config
	shutdownFuncs  []func(context.Context) error
	TraceProvider  *sdktrace.TracerProvider
	MetricProvider *sdkmetric.MeterProvider
	Resource       *sdkresource.Resource
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
		return
	}

	os.Resource = os.initResource(ctx)

	conn, err := grpc.NewClient(os.cfg.OTLPAddress,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		os.logger.Fatal().Msgf("Error establishing grpc connection: %v", err)
	}

	os.TraceProvider = os.initTracerProvider(ctx, conn, os.Resource)
	os.MetricProvider = os.initMeterProvider(ctx, conn, os.Resource)

	os.shutdownFuncs = append(os.shutdownFuncs,
		os.TraceProvider.Shutdown,
		os.MetricProvider.Shutdown,
	)
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

func (os *OtelService) initTracerProvider(ctx context.Context, conn *grpc.ClientConn, res *sdkresource.Resource) *sdktrace.TracerProvider {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		os.logger.Fatal().Msgf("new otlp trace exporter failed: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func (os *OtelService) initMeterProvider(ctx context.Context, conn *grpc.ClientConn, res *sdkresource.Resource) *sdkmetric.MeterProvider {
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		os.logger.Fatal().Msgf("new otlp metric http exporter failed: %v", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)
	return mp
}

func (os *OtelService) initResource(ctx context.Context) *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			ctx,
			// temp, get this from the config
			sdkresource.WithAttributes(semconv.ServiceNameKey.String(os.cfg.ServiceName)),
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
