package main

import (
	cdevents-controller "github.com/bradmccoydev/cdevents-controller/cue/cdevents-controller"
)

app: cdevents-controller.#Application & {
	config: {
		meta: {
			name:      "cdevents-controller"
			namespace: "default"
		}
		image: tag: "0.0.1"
		resources: requests: {
			cpu:    "100m"
			memory: "16Mi"
		}
		hpa: {
			enabled:     true
			maxReplicas: 3
		}
		ingress: {
			enabled:   true
			className: "nginx"
			host:      "cdevents-controller.example.com"
			tls:       true
			annotations: "cert-manager.io/cluster-issuer": "letsencrypt"
		}
		serviceMonitor: enabled: true
	}
}

objects: app.objects
