package chronos

import (
	"net/http"
	"net/http/httptest"
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

func TestTaskName(t *testing.T) {
	cases := []struct {
		app, env, in, out string
	}{
		{"dogfood", "production", "rake db:migrate:status", "dogfood-production-rake-db-migrate-status"},
		{"dogfood", "production", "rake -T", "dogfood-production-rake-t"},
		{"dogfood", "production", "rake  foo[123]", "dogfood-production-rake-foo-123"},
		// NOTE: command runner should prevent this from executing but
		//       just testing an extreme
		{"dogfood", "production", "curl -s google.com > yo.txt", "dogfood-production-curl-s-google-com-yo-txt"},
	}

	for _, test := range cases {
		s := taskName(test.app, test.env, test.in)
		if s != test.out {
			t.Errorf("expected '%s' but got '%s'", test.out, s)
		}
	}
}

func TestScheduledTasks(t *testing.T) {
	tv := &TaskVars{TaskName: "dogfood-production-yo"}
	ms := mesosServer(tv)
	tasks, err := scheduledTasks(ms.URL, tv.TaskName)
	if err != nil {
		t.Errorf("expected nil but got: %s", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks but got: %d", len(tasks))
	}
}

// func TestLaunchChronosOneOffCommand(t *testing.T) {
// 	tv := &TaskVars{
// 		App:      "dogfood",
// 		Env:      "stage",
// 		Tag:      "1.2.3",
// 		Command:  "echo yodawg",
// 		TaskName: "dogfood-stage-echo-yodawg",
// 	}
// 	cs := chronosServer(true)
// 	ms := mesosServer(tv)
// 	fmt.Println(ms.URL)
// 	fmt.Println(cs.URL)
// 	err := launchChronosOneOffCommand(cs.URL, ms.URL, tv)
// 	if err != nil {
// 		t.Errorf("expected nil but got: %s", err)
// 	}
// }

func chronosServer(success bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if success {
			w.WriteHeader(http.StatusNoContent)

			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func mesosServer(task *TaskVars) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("task_json").Parse(mesosTasks))
		err := t.Execute(w, task)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
	}))
}
