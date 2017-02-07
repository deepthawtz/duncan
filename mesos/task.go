package mesos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// TaskFramework is the framework name used by duncan run
	TaskFramework = "chronos"

	// TaskRunning represents status for a running task
	TaskRunning = "TASK_RUNNING"

	// TaskStaging represents status for a staged task
	TaskStaging = "TASK_STAGING"

	// TaskFinished represents status for a successfully completed task
	TaskFinished = "TASK_FINISHED"

	// TaskFailed represents status for a failed task
	TaskFailed = "TASK_FAILED"

	// TaskKilled represents status for a prematurely killed task
	TaskKilled = "TASK_KILLED"
)

// Tasks represents Mesos tasks
// used to deserialize a Mesos tasks API response
type Tasks struct {
	Tasks []*Task `json:"tasks"`
}

// TasksFor returns running task for given task app, env, name
func (t *Tasks) TasksFor(name string) (*Tasks, error) {
	tasks := &Tasks{}
	for _, task := range t.Tasks {
		p := strings.Split(task.ID, ":")
		if len(p) > 2 {
			s := p[len(p)-2]
			if name == s {
				tasks.Tasks = append(tasks.Tasks, task)
			}
		}
	}
	return tasks, nil
}

// Len helps Tasks implement the sort.Interface
func (t *Tasks) Len() int {
	return len(t.Tasks)
}

// Less helps Tasks implement the sort.Interface
func (t *Tasks) Less(i, j int) bool {
	ti := t.Tasks[i]
	tj := t.Tasks[j]
	if len(ti.Statuses) == 0 && len(tj.Statuses) >= 1 {
		return false
	}
	if len(tj.Statuses) <= 0 && len(ti.Statuses) >= 1 {
		return true
	}
	return ti.Statuses[len(ti.Statuses)-1].Timestamp < tj.Statuses[len(tj.Statuses)-1].Timestamp
}

// Swap helps Tasks implement the sort.Interface
func (t *Tasks) Swap(i, j int) {
	t.Tasks[i], t.Tasks[j] = t.Tasks[j], t.Tasks[i]
}

// Task repesents a Mesos task
type Task struct {
	ID          string    `json:"id"`
	FrameworkID string    `json:"framework_id"`
	SlaveID     string    `json:"slave_id"`
	State       string    `json:"state"`
	Statuses    []*Status `json:"statuses"`
}

// SlaveIP returns the IP of slave the task is running on
func (t *Task) SlaveIP() (string, error) {
	var out string
	for _, s := range t.Statuses {
		for _, n := range s.Container.NetworkInfos {
			for _, i := range n.IPAddresses {
				out = i.IP
			}
		}
	}
	if out == "" {
		return "", fmt.Errorf("could not find slave IP for %v", t)
	}
	return out, nil
}

// Duration returns the duration a task took to complete
func (t *Task) Duration() (float64, error) {
	s := t.Statuses
	if len(s) < 2 {
		return 0.0, fmt.Errorf("task did not run properly")
	}
	end := s[len(s)-1]
	start := s[len(s)-2]
	if end.State != TaskFinished && end.State != TaskFailed && end.State != TaskKilled {
		return 0.0, fmt.Errorf("task incomplete")
	}
	if start.State != TaskRunning {
		return 0.0, fmt.Errorf("task incomplete")
	}
	dur := end.Timestamp - start.Timestamp
	return dur, nil
}

// LogDirectory returns the Mesos sandbox directory for a task
func (t *Task) LogDirectory(host string) (string, error) {
	if !strings.Contains(host, "://") {
		host = "http://" + host + ":5051"
	}
	url := fmt.Sprintf("%s/state", host)
	for {
		time.Sleep(200 * time.Millisecond)
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to fetch slave state: %s", resp.Status)
		}

		st := &SlaveTasks{}
		if err := json.NewDecoder(resp.Body).Decode(st); err != nil {
			return "", err
		}

		for _, f := range st.CompletedFrameworks {
			if f.Name == TaskFramework {
				for _, e := range f.CompletedExecutors {
					if e.ID == t.ID {
						return e.Directory, nil
					}
				}
			}
		}
		for _, f := range st.Frameworks {
			if f.Name == TaskFramework {
				for _, e := range f.Executors {
					if e.ID == t.ID {
						return e.Directory, nil
					}
				}
			}
		}
	}
}

// Status repesents Mesos task status
type Status struct {
	State     string           `json:"state"`
	Timestamp float64          `json:"timestamp"`
	Container *ContainerStatus `json:"container_status"`
}

// ContainerStatus repesents Mesos task container status
type ContainerStatus struct {
	NetworkInfos []*NetworkInfo `json:"network_infos"`
}

// NetworkInfo repesents Mesos task container network info
type NetworkInfo struct {
	IPAddresses []*IPAddress `json:"ip_addresses"`
}

// IPAddress repesents a Mesos task slave IP
type IPAddress struct {
	IP string `json:"ip_address"`
}
