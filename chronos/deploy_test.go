package chronos

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"text/template"
)

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
	ts, err := scheduledTasks(ms.URL, tv.TaskName)
	if err != nil {
		t.Errorf("expected nil but got: %s", err)
	}
	l := len(ts.Tasks)
	if l != 2 {
		t.Errorf("expected 2 tasks but got: %d", l)
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
		filename := filepath.Join("testdata", "mesos_tasks.json.tmpl")
		as, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		t := template.Must(template.New("task_json").Parse(string(as)))
		err = t.Execute(w, task)
		if err != nil {
			panic(err)
		}
	}))
}
