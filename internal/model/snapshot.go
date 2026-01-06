package model

type Snapshot struct {
	RequestRate float64
	ErrorRatio  float64
	P50Latency  float64
	P95Latency  float64
	SLOBurnRate float64
	Goroutines  float64
	MemoryBytes float64
}
