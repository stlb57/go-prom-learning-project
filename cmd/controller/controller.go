package main

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	pmodel "github.com/prometheus/common/model"

	"worker-metrics/internal/decide"
	"worker-metrics/internal/model"
	"worker-metrics/internal/state"
)

const (
	windowSize = 5
	URL        = "http://localhost:9090"
)

var snaps []model.Snapshot

func extractScalar(v pmodel.Value) float64 {
	vec, ok := v.(pmodel.Vector)
	if !ok || len(vec) == 0 {
		return 0
	}
	return float64(vec[0].Value)
}

func observe(api v1.API) model.Snapshot {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _, _ := api.Query(ctx, "rate(worker_requests_total[1m])", time.Now())
	errs, _, _ := api.Query(ctx,
		"rate(worker_request_errors_total[1m]) / rate(worker_requests_total[1m])",
		time.Now(),
	)
	p95, _, _ := api.Query(ctx,
		"histogram_quantile(0.95, rate(worker_request_latency_seconds_bucket[5m]))",
		time.Now(),
	)

	return model.Snapshot{
		RequestRate: extractScalar(req),
		ErrorRatio:  extractScalar(errs),
		P95Latency:  extractScalar(p95),
	}
}

func main() {
	client, _ := api.NewClient(api.Config{Address: URL})
	api := v1.NewAPI(client)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		snap := observe(api)

		if len(snaps) == windowSize {
			snaps = snaps[1:]
		}
		snaps = append(snaps, snap)

		if d := decide.Evaluate(snaps); d != nil {
			switch d.Status {
			case decide.Healthy:
				state.Set(state.Healthy)
			case decide.Degraded:
				state.Set(state.Degraded)
			case decide.Unstable:
				state.Set(state.Unstable)
			}
		}
	}
}
