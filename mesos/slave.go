package mesos

// SlaveTasks represents Mesos slave completed tasks
type SlaveTasks struct {
	Frameworks          []*Framework `json:"frameworks"`
	CompletedFrameworks []*Framework `json:"completed_frameworks"`
}

// Framework represents a completed framework on a Mesos slave
type Framework struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Executors          []*Executor `json:"executors,omitempty"`
	CompletedExecutors []*Executor `json:"completed_executors,omitempty"`
}

// Executor represents a completed executor on a Mesos slave
type Executor struct {
	ID        string `json:"id"`
	Directory string `json:"directory"`
}

// Logs represents logs for a Mesos task
type Logs struct {
	Data string `json:"data"`
}
