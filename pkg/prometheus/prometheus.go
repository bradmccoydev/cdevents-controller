package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"
)

func PushGaugeMetric(logger *zap.Logger, metricName string, metricValue float64, contextId string, contextType string, subjectType string) {
	pushGatewayURL := "http://prometheus-pushgateway.observability:9091"

	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricName,
		Help: "CDEvent metric",
	},
		[]string{"context_id", "context_type", "subject_type"},
	)

	promLabels := prometheus.Labels{
		"context_id":   contextId,
		"context_type": contextType,
		"subject_type": subjectType,
	}

	gauge.With(promLabels).Set(metricValue)

	if err := push.New(pushGatewayURL, "cdevents").Collector(gauge).Push(); err != nil {
		logger.Error("Failed to list Result object:", zap.Error(err))
	}

	logger.Info("Metrics pushed successfully")
}
