package metrics

import "github.com/prometheus/client_golang/prometheus"

func RegisterOperatorMetrics() (*Metrics, error) {
	operatorErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "integreatly_tutorial_webapp_operator_reconcile_errors_total",
		Help: "Number of errors that occurred while reconciling the integreatly tutorial webapp deployment",
	})

	err := prometheus.Register(operatorErrors)
	if err != nil {
		return nil, err
	}

	return &Metrics{operatorErrors: operatorErrors}, nil
}
