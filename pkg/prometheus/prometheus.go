package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"
)

func PushGaugeMetric(logger *zap.Logger, metricName string, metricValue float64) {
	pushGatewayURL := "http://prometheus-pushgateway.observability:9091/metrics"

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricName,
		Help: "CDEvent metric",
	})

	gauge.Set(metricValue)

	if err := push.New(pushGatewayURL, "cdevents").Collector(gauge).Push(); err != nil {
		logger.Error("Failed to list Result object:", zap.Error(err))
	}

	logger.Info("Metrics pushed successfully")
}
