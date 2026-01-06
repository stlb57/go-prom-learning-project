package decide

import "worker-metrics/internal/model"

func detectLatencyDegradation(snaps []model.Snapshot) *Decision {
	const (
		thresholdSeconds = 0.300 // 300ms
		requiredHits     = 3
	)

	if len(snaps) < requiredHits {
		return nil
	}
	count := 0
	for i := len(snaps) - 1; i >= 0; i-- {
		if snaps[i].P95Latency > thresholdSeconds {
			count++
			if count == requiredHits {
				return &Decision{
					Status: Degraded,
					Reason: "p95 latency > 300ms for 3 consecutive snapshots",
				}
			}
		} else {
			break
		}
	}

	return nil
}

func detectLatencyRecovery(snaps []model.Snapshot) *Decision {
	const (
		healthyThresholdSeconds = 0.200 // 200ms
		requiredHits            = 4
	)

	if len(snaps) < requiredHits {
		return nil
	}

	count := 0
	for i := len(snaps) - 1; i >= 0; i-- {
		if snaps[i].P95Latency < healthyThresholdSeconds {
			count++
			if count == requiredHits {
				return &Decision{
					Status: Healthy,
					Reason: "p95 latency < 200ms for sustained period",
				}
			}
		} else {
			break
		}
	}

	return nil
}

func detectErrorSpike(snaps []model.Snapshot) *Decision {
	const (
		errorThreshold = 0.02 // 2%
		requiredHits   = 2
	)

	if len(snaps) < requiredHits {
		return nil
	}

	count := 0
	for i := len(snaps) - 1; i >= 0; i-- {
		if snaps[i].ErrorRatio > errorThreshold {
			count++
			if count == requiredHits {
				return &Decision{
					Status: Unstable,
					Reason: "error ratio > 2% for consecutive snapshots",
				}
			}
		} else {
			break
		}
	}

	return nil
}
