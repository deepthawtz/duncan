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
	cases := []struct {
		ok  bool
		exp bool
	}{
		{ok: true, exp: true},
		{ok: false, exp: false},
	}

	for _, test := range cases {
		tv := &TaskVars{TaskName: "dogfood-production-yo"}
		ms := mesosServer(tv, test.ok)
		ts, err := scheduledTasks(ms.URL, tv.TaskName)
		if err != nil && test.ok {
			t.Errorf("expected nil but got: %s", err)
		}
		if err == nil && !test.ok {
			t.Errorf("expected error but got nil")
		}
		if test.ok {
			l := len(ts.Tasks)
			if l != 2 {
				t.Errorf("expected 2 tasks but got: %d", l)
			}
		}
	}
}

func TestLaunchChronosOneOffCommand(t *testing.T) {
	cases := []struct {
		chronosOK bool
		mesosOK   bool
		exp       bool
	}{
		{chronosOK: true, mesosOK: true, exp: true},
		{chronosOK: false, mesosOK: true, exp: false},
		{chronosOK: true, mesosOK: false, exp: false},
	}

	for _, test := range cases {
		tv := &TaskVars{
			App:      "dogfood",
			Env:      "stage",
			Tag:      "1.2.3",
			Mem:      3,
			Command:  "echo yodawg",
			TaskName: "dogfood-stage-echo-yodawg",
		}
		cs := chronosServer(test.chronosOK)
		ms := mesosServer(tv, test.mesosOK)
		h := func(string, string, int) error {
			return nil
		}
		err := launchChronosOneOffCommand(cs.URL, ms.URL, tv, h)
		if err != nil && test.exp {
			t.Errorf("expected nil but got: %s", err)
		}
		if err == nil && !test.exp {
			t.Errorf("expected error but got nil")
		}
	}
}

func TestValidateSchedule(t *testing.T) {
	cases := []struct {
		schedule string
		ok       bool
	}{
		{"R1//PT30M", true},
		{"derp", false},
		{"yo/yo/yo", false},
		{"RR/yo/yo", false},
		{"R1/2014-10-10T18:32:00Z/PT30M", true},
	}

	for _, test := range cases {
		err := validateSchedule(test.schedule)
		if err != nil && test.ok {
			t.Errorf("expected no error got: %s", err)
		}
		if err == nil && !test.ok {
			t.Errorf("expected error but got nil")
		}
	}
}

func chronosServer(success bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if success {
			w.WriteHeader(http.StatusNoContent)

			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func mesosServer(task *TaskVars, success bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !success {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
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
