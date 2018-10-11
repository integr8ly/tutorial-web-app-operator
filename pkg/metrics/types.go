package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	operatorErrors prometheus.Counter
}
