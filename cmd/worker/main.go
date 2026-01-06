package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	worker_requests_total = promauto.NewCounter(prometheus.CounterOpts{
		Name: "worker_requests_total",
		Help: "",
	})

	worker_request_errors_total = promauto.NewCounter(prometheus.CounterOpts{
		Name: "worker_request_errors_total",
		Help: "",
	})

	worker_request_latency_seconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "worker_request_latency_seconds",
		Help:    "",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
	})

	worker_goroutines = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "worker_goroutines",
		Help: "",
	})

	worker_memory_bytes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "worker_memory_bytes",
		Help: "",
	})

	worker_uptime_seconds = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "worker_uptime_seconds",
		Help: "",
	})
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	const worker_count = 100

	go func() {
		var mem runtime.MemStats
		for {
			runtime.ReadMemStats(&mem)
			worker_goroutines.Set(float64(runtime.NumGoroutine()))
			worker_memory_bytes.Set(float64(mem.Alloc))
			worker_uptime_seconds.Set(time.Since(start).Seconds())
			time.Sleep(1 * time.Second)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		reqStart := time.Now()
		worker_requests_total.Inc()

		if rand.Float64() < 0.05 {
			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(20 * time.Millisecond)
		}

		fmt.Fprintln(w, "Hello World")
		worker_request_latency_seconds.Observe(time.Since(reqStart).Seconds())
	})

	var wg sync.WaitGroup
	for w := 1; w <= worker_count; w++ {
		wg.Add(1)
		go worker(w, &wg)
	}
	if err := http.ListenAndServe(":2112", nil); err != nil {
		fmt.Printf("server failed: %v\n", err)
		worker_request_errors_total.Inc()
	}
}
