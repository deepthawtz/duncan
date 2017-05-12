package mesos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"text/template"
)

type TestTask struct {
	ID       string
	TaskName string
	State    string
}

var (
	allTasks = map[string]*Tasks{}
	tt       = &TestTask{
		TaskName: "dogfood-production-echo-hi",
		ID:       "ct:1485292701088:0:dogfood-production-echo-hi:",
	}
)

func init() {
	for _, x := range []string{"completed", "incomplete"} {
		filename := filepath.Join("testdata", fmt.Sprintf("%s_tasks.json.tmpl", x))
		ct, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		tt.State = x
		t, err := createTasks(tt, string(ct))
		if err != nil {
			panic(err)
		}
		allTasks[x] = t
	}
}

func TestTasksFor(t *testing.T) {
	ct := allTasks["completed"]
	taskCount := len(ct.Tasks)
	if taskCount != 2 {
		t.Errorf("expected 2 tasks but got %d", taskCount)
	}
	st, err := ct.TasksFor(tt.TaskName)
	if err != nil {
		t.Errorf("expected nil but got: %s", err)
	}
	if len(st.Tasks) != 2 {
		t.Errorf("expected 2 tasks but got %d", len(st.Tasks))
	}
}

func TestSlaveIP(t *testing.T) {
	ct := allTasks["completed"]
	for _, task := range ct.Tasks {
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
	ct := allTasks["completed"]
	for _, task := range ct.Tasks {
		dur, err := task.Duration()
		if err != nil {
			t.Error(err)
		}
		if dur == 0.0 {
			t.Error("expected to calculate task duration")
		}
	}
	it := allTasks["incomplete"]
	for _, task := range it.Tasks {
		_, err := task.Duration()
		if err == nil {
			t.Error("expected error")
		}
	}
}

func createTasks(tt *TestTask, templ string) (*Tasks, error) {
	j := new(bytes.Buffer)
	template := template.Must(template.New("task_json").Parse(templ))
	err := template.Execute(j, tt)
	if err != nil {
		return nil, fmt.Errorf("expected nil but got: %s", err)
	}
	tasks := &Tasks{}
	if err = json.NewDecoder(j).Decode(tasks); err != nil {
		return nil, fmt.Errorf("expected nil but got: %s", err)
	}
	return tasks, nil
}

func mesosAgentServer(tt *TestTask) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Join("testdata", fmt.Sprintf("agent_state_%s.json.tmpl", tt.State))
		as, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		t := template.Must(template.New("agent_state").Parse(string(as)))
		err = t.Execute(w, tt)
		if err != nil {
			panic(err)
		}
		return
	}))
}
