package autoscaling

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

type testPolicies struct {
	QueueLengthScaled []*QueueLengthScaled

	CPUScaled []*CPUScaled
}

type QueueLengthScaled struct {
	Name               string
	AppName            string
	AppType            string
	Environment        string
	RedisURL           string
	Queues             string
	MinInstances       int
	MaxInstances       int
	UpThreshold        int
	DownThreshold      int
	ScaleUpBy          int
	ScaleDownBy        int
	CheckFrequencySecs int
	Enabled            bool
}

type CPUScaled struct {
	Name               string
	AppName            string
	AppType            string
	Environment        string
	MinInstances       int
	MaxInstances       int
	UpThreshold        int
	DownThreshold      int
	ScaleUpBy          int
	ScaleDownBy        int
	CheckFrequencySecs int
	Enabled            bool
}

func TestGetPolicies(t *testing.T) {
	tp := &testPolicies{}
	ts := slytheServer(tp, "policies.json.tmpl", true)
	viper.Set("SLYTHE_HOST", ts.URL)
	_, err := GetPolicies("", "")
	if err != nil {
		t.Errorf("expected no error got: %v", err)
	}

	tp.QueueLengthScaled = append(tp.QueueLengthScaled, &QueueLengthScaled{
		Name:               "FooProductionWorker",
		AppName:            "foo",
		AppType:            "worker",
		Environment:        "production",
		MinInstances:       1,
		MaxInstances:       10,
		UpThreshold:        2000,
		DownThreshold:      1000,
		ScaleUpBy:          3,
		ScaleDownBy:        1,
		CheckFrequencySecs: 20,
		RedisURL:           "redis://yo.dawg:6379/2",
		Queues:             "one,two,three",
		Enabled:            true,
	})
	tp.CPUScaled = append(tp.CPUScaled, &CPUScaled{
		Name:               "FooProductionWeb",
		AppName:            "foo",
		AppType:            "web",
		Environment:        "production",
		MinInstances:       1,
		MaxInstances:       10,
		UpThreshold:        70,
		DownThreshold:      5,
		ScaleUpBy:          3,
		ScaleDownBy:        1,
		CheckFrequencySecs: 20,
		Enabled:            true,
	})
	ts = slytheServer(tp, "policies.json.tmpl", true)
	viper.Set("SLYTHE_HOST", ts.URL)
	p, err := GetPolicies("idontexist", "")
	if err != nil {
		t.Errorf("expected no error got: %v", err)
	}
	if len(p.CPUScaled) != 0 || len(p.QueueLengthScaled) != 0 {
		t.Errorf("expected no policies")
	}
	p, err = GetPolicies("foo", "production")
	if err != nil {
		t.Errorf("expected no error got: %v", err)
	}
	if len(p.QueueLengthScaled) != 1 {
		t.Errorf("expected one policy")
	}
	p, err = GetPolicies("", "production")
	if err != nil {
		t.Errorf("expected no error got: %v", err)
	}
	if len(p.QueueLengthScaled)+len(p.CPUScaled) != 2 {
		t.Errorf("expected two policies")
	}
}

func slytheServer(tp *testPolicies, tmpl string, ok bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		f := filepath.Join("testdata", tmpl)
		b, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		j := new(bytes.Buffer)
		template := template.Must(template.New("task_json").Parse(string(b)))
		err = template.Execute(j, tp)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, j.String())
	}))
}
