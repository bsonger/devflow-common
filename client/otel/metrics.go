package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func InitMetrics() error {
	// Prometheus Exporter
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}

	// MeterProvider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)

	otel.SetMeterProvider(provider)
	return nil
}
