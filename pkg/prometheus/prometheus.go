package prometheus

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

func PushGaugeMetric(metricName string, metricValue float64) {
	pushGatewayURL := "http://prometheus-pushgateway.observability:9091/metrics"

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricName,
		Help: "CDEvent metric",
	})

	gauge.Set(metricValue)

	pusher := push.New(pushGatewayURL, "cdevents")

	err := pusher.Collector(gauge)
	if err != nil {
		log.Fatal("Failed to push metric to Pushgateway:", err)
	}

	pushError := pusher.Push()
	if pushError != nil {
		log.Fatal("Failed to push metrics to Pushgateway:", err)
	}

	fmt.Println("Metrics pushed successfully")
}
