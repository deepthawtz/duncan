package mesos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"text/template"
)

const mesosTasks = `
{
  "tasks": [
    {
      "id": "ct:1485467406026:0:{{.TaskName}}:",
      "name": "ChronosTask:{{.TaskName}}",
      "framework_id": "8b63f436-482d-44d0-9052-17c23e72ef68-0002",
      "executor_id": "",
      "slave_id": "8b63f436-482d-44d0-9052-17c23e72ef68-S0",
      "state": "TASK_FAILED",
      "resources": {
        "disk": 256,
        "mem": 1024,
        "gpus": 0,
        "cpus": 1
      },
      "statuses": [
        {
          "state": "TASK_RUNNING",
          "timestamp": 1485467407.53938,
          "container_status": {
            "network_infos": [ { "ip_addresses": [ { "ip_address": "172.16.1.145" } ] } ]
          }
        },
        {
          "state": "TASK_FAILED",
          "timestamp": 1485467407.5402,
          "container_status": {
            "network_infos": [ { "ip_addresses": [ { "ip_address": "172.16.1.145" } ] } ]
          }
        }
      ]
    },
    {
      "id": "ct:1485467338370:0:{{.TaskName}}:",
      "name": "ChronosTask:{{.TaskName}}",
      "framework_id": "8b63f436-482d-44d0-9052-17c23e72ef68-0002",
      "executor_id": "",
      "slave_id": "8b63f436-482d-44d0-9052-17c23e72ef68-S0",
      "state": "TASK_FINISHED",
      "resources": {
        "disk": 256,
        "mem": 1024,
        "gpus": 0,
        "cpus": 1
      },
      "statuses": [
        {
          "state": "TASK_RUNNING",
          "timestamp": 1485467346.31432,
          "container_status": {
            "network_infos": [ { "ip_addresses": [ { "ip_address": "172.16.1.145" } ] } ]
          }
        },
        {
          "state": "TASK_FINISHED",
          "timestamp": 1485467346.31504,
          "container_status": {
            "network_infos": [ { "ip_addresses": [ { "ip_address": "172.16.1.145" } ] } ]
          }
        }
      ]
    }
  ]
}
`

type TaskVars struct {
	TaskName string
}

func TestTasksFor(t *testing.T) {
	tv := &TaskVars{TaskName: "dogfood-production-echo-hi"}
	tasks, err := tasks(tv)
	if err != nil {
		t.Error(err)
	}
	if len(tasks.Tasks) != 2 {
		t.Errorf("expected 2 tasks but got %d", len(tasks.Tasks))
	}
	st, err := tasks.TasksFor(tv.TaskName)
	if err != nil {
		t.Errorf("expected nil but got: %s", err)
	}
	if len(st) != 2 {
		t.Errorf("expected 2 tasks but got %d", len(st))
	}
}

func TestSlaveIP(t *testing.T) {
	tv := &TaskVars{TaskName: "dogfood-production-echo-hi"}
	tasks, err := tasks(tv)
	if err != nil {
		t.Error(err)
	}
	for _, task := range tasks.Tasks {
		ip, err := task.SlaveIP()
		if err != nil {
			t.Error(err)
		}
		if ip != "172.16.1.145" {
			t.Errorf("expected slave IP but got %s", ip)
		}
	}
}

func TestDuration(t *testing.T) {
	tv := &TaskVars{TaskName: "dogfood-production-echo-hi"}
	tasks, err := tasks(tv)
	if err != nil {
		t.Error(err)
	}
	for _, task := range tasks.Tasks {
		dur, err := task.Duration()
		if err != nil {
			t.Error(err)
		}
		if dur == 0.0 {
			t.Error("expected to calculate task duration")
		}
	}
}

func tasks(tv *TaskVars) (*Tasks, error) {
	j := new(bytes.Buffer)
	template := template.Must(template.New("task_json").Parse(mesosTasks))
	err := template.Execute(j, tv)
	if err != nil {
		return nil, fmt.Errorf("expected nil but got: %s", err)
	}
	tasks := &Tasks{}
	if err = json.NewDecoder(j).Decode(tasks); err != nil {
		return nil, fmt.Errorf("expected nil but got: %s", err)
	}
	return tasks, nil
}
