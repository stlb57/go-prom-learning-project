package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const windowSize = 5

type Snapshot struct {
	RequestRate float64
	ErrorRatio  float64
	P50Latency  float64
	P95Latency  float64
	SLOBurnRate float64
	Goroutines  float64
	MemoryBytes float64
}

var snap_slice []Snapshot

func extractScalar(value model.Value) float64 {
	vector, ok := value.(model.Vector)
	if !ok || len(vector) == 0 {
		return 0
	}
	return float64(vector[0].Value)
}

const URL = "http://localhost:9090"

func observe(v1api v1.API) Snapshot {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _, err := v1api.Query(ctx, "up", time.Now(), v1.WithTimeout(5*time.Second))
	fmt.Printf("Result:\n%v\n", result)

	reqRateVal, _, err := v1api.Query(ctx, "rate(worker_requests_total[1m])", time.Now())
	if err != nil {
		fmt.Println("reqRate error:", err)
	}

	errRatioVal, _, _ := v1api.Query(ctx,
		"rate(worker_request_errors_total[1m]) / rate(worker_requests_total[1m])",
		time.Now(),
	)

	p50Val, _, _ := v1api.Query(ctx,
		"histogram_quantile(0.5, rate(worker_request_latency_seconds_bucket[5m]))",
		time.Now(),
	)

	p95Val, _, _ := v1api.Query(ctx,
		"histogram_quantile(0.95, rate(worker_request_latency_seconds_bucket[5m]))",
		time.Now(),
	)

	gorVal, _, _ := v1api.Query(ctx, "worker_goroutines", time.Now())
	memVal, _, _ := v1api.Query(ctx, "worker_memory_bytes", time.Now())

	sloBurnVal, _, _ := v1api.Query(ctx,
		"(rate(worker_request_errors_total[5m]) / rate(worker_requests_total[5m])) / 0.001",
		time.Now(),
	)

	return Snapshot{
		RequestRate: extractScalar(reqRateVal),
		ErrorRatio:  extractScalar(errRatioVal),
		P50Latency:  extractScalar(p50Val),
		P95Latency:  extractScalar(p95Val),
		SLOBurnRate: extractScalar(sloBurnVal),
		Goroutines:  extractScalar(gorVal),
		MemoryBytes: extractScalar(memVal),
	}

}

func main() {
	client, err := api.NewClient(api.Config{
		Address: URL,
	})
	if err != nil {
		fmt.Println("Error creating client")
		os.Exit(1)
	}
	v1api := v1.NewAPI(client)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		snap := observe(v1api)
		if len(snap_slice) == windowSize {
			snap_slice = snap_slice[1:]
		}
		snap_slice = append(snap_slice, snap)
		fmt.Printf("%+v\n", snap)
	}

}
