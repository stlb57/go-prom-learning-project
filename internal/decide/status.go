package decide

type Status string

const (
	Healthy  Status = "HEALTHY"
	Degraded Status = "DEGRADED"
	Unstable Status = "UNSTABLE"
)

type Decision struct {
	Status Status
	Reason string
}
