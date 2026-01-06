package state

import "sync/atomic"

type Status int32

const (
	Healthy Status = iota
	Degraded
	Unstable
)

var current atomic.Int32

func Set(s Status) {
	current.Store(int32(s))
}

func Get() Status {
	return Status(current.Load())
}
