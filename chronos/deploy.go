package chronos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/betterdoctor/duncan/consul"
	"github.com/betterdoctor/duncan/mesos"
	"github.com/spf13/viper"
)

// SlaveTasks represents Mesos slave completed tasks
type SlaveTasks struct {
	Frameworks []*Framework `json:"completed_frameworks"`
}

// Framework represents a completed framework on a Mesos slave
type Framework struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Executors []*Executor `json:"completed_executors"`
}

// Executor represents a completed executor on a Mesos slave
type Executor struct {
	ID        string `json:"id"`
	Directory string `json:"directory"`
}

// TaskVars represents a one-off Chronos task
type TaskVars struct {
	App, Env, Tag, Command, TaskName, DockerRepoPrefix, DockerConfURL string
}

var logsOpened bool

// RunCommand spins up a Chronos task to run the given command and exits
func RunCommand(app, env, cmd string, follow bool) error {
	if !follow {
		logsOpened = true
	}
	tag, err := consul.CurrentTag(app, env)
	if err != nil {
		return err
	}
	prefix := viper.GetString("docker_repo_prefix")
	if prefix == "" {
		return fmt.Errorf("docker_repo_prefix not set in config")
	}
	// TODO: remove when we change how we handle passing docker config to tasks
	confURL := viper.GetString("docker_conf_url")
	if confURL == "" {
		return fmt.Errorf("docker_conf_url not set in config")
	}

	task := &TaskVars{
		App:              app,
		Env:              env,
		Tag:              tag,
		Command:          cmd,
		TaskName:         taskName(app, env, cmd),
		DockerRepoPrefix: prefix,
		DockerConfURL:    confURL,
	}
	chronosURL := fmt.Sprintf("%s/service/chronos/v1/scheduler/iso8601", viper.GetString("chronos_host"))
	mesosURL := fmt.Sprintf("%s/mesos/tasks", viper.GetString("marathon_host"))
	if err := launchChronosOneOffCommand(chronosURL, mesosURL, task); err != nil {
		return err
	}

	return nil
}

// taskName generates a valid Chronos task name based on the app/env/command given
func taskName(app, env, cmd string) string {
	out := []string{app, env}
	re := regexp.MustCompile("[a-zA-Z0-9]*")
	p := strings.Split(cmd, " ")
	for _, c := range p {
		m := re.FindAllString(strings.TrimSpace(c), -1)
		for _, x := range m {
			if x != "" {
				out = append(out, strings.ToLower(x))
			}
		}
	}

	return strings.Join(out, "-")
}

func renderChronosTaskJSON(task *TaskVars) (string, error) {
	t, err := template.New("chronos_task").Parse(taskTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}
	j := new(bytes.Buffer)
	if err := t.Execute(j, task); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}
	return j.String(), nil
}

func launchChronosOneOffCommand(chronosURL, mesosURL string, task *TaskVars) error {
	tasks, err := scheduledTasks(mesosURL, task.TaskName)
	if err != nil {
		return err
	}

	cj, err := renderChronosTaskJSON(task)
	if err != nil {
		return err
	}
	resp, err := http.Post(chronosURL, "application/json", strings.NewReader(cj))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to launch command: %s", body)
	}
	fmt.Printf("executing '%s' in instance of %s:%s (%s)\n", task.Command, task.App, task.Tag, task.Env)
	if err = handleTask(mesosURL, task.TaskName, len(tasks.Tasks)); err != nil {
		return err
	}

	return nil
}

func cleanupTask(name string) error {
	url := fmt.Sprintf("%s/service/chronos/v1/scheduler/job/%s",
		viper.GetString("chronos_host"),
		name,
	)
	req, _ := http.NewRequest("DELETE", url, strings.NewReader(""))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not clean up task: %s", name)
	}
	if resp.StatusCode != http.StatusNoContent {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf(string(b))
	}
	return nil
}

func handleTask(url, name string, count int) error {
	fmt.Println("scheduling task")
	state := map[string]bool{}
	for {
		time.Sleep(1 * time.Second)
		tasks, err := scheduledTasks(url, name)
		if err != nil {
			return err
		}
		if len(tasks.Tasks) == count {
			continue
		}
		if len(tasks.Tasks) == count+1 {
			sort.Sort(tasks)
			task := tasks.Tasks[len(tasks.Tasks)-1]

			if _, ok := state[task.State]; !ok {
				state[task.State] = true
				fmt.Printf("task state: %s\n", task.State)
			}
			switch task.State {
			case mesos.TaskRunning:
				if err := openLogPage(task); err != nil {
					return err
				}
				continue
			case mesos.TaskStaging:
				continue
			case mesos.TaskFinished:
				dur, err := task.Duration()
				if err != nil {
					return err
				}
				if err := openLogPage(task); err != nil {
					return err
				}
				fmt.Printf("\ntask finished: %.02f seconds\n", dur)
				if err := printLogs(task); err != nil {
					return err
				}
				if err := cleanupTask(name); err != nil {
					return err
				}

				return nil
			case mesos.TaskFailed:
				dur, err := task.Duration()
				if err != nil {
					return err
				}

				if err := openLogPage(task); err != nil {
					return err
				}
				if err := printLogs(task); err != nil {
					return err
				}
				if err := cleanupTask(name); err != nil {
					return err
				}
				return fmt.Errorf("\ntask failed: %.02f seconds\n", dur)
			case mesos.TaskKilled:
				dur, err := task.Duration()
				if err != nil {
					return err
				}
				if err := cleanupTask(name); err != nil {
					return err
				}

				return fmt.Errorf("\ntask killed: %.02f seconds\n", dur)
			default:
				return fmt.Errorf("\ntask state unhandled: %s", task.State)
			}
		}
	}
}

func printLogs(t *mesos.Task) error {
	ip, err := t.SlaveIP()
	if err != nil {
		return err
	}
	dir, err := t.LogDirectory(ip)
	if err != nil {
		return fmt.Errorf("cannot fetch logs: %s\n", err)
	}

	var out string
	if t.State == mesos.TaskFinished {
		out = "stdout"
	} else {
		out = "stderr"
	}

	url := fmt.Sprintf("http://%s:5051/files/read?path=%s/%s&offset=0", ip, dir, out)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not fetch logs from: %s\n%s\n", url, resp.Status)
	}
	l := &mesos.Logs{}
	if err := json.NewDecoder(resp.Body).Decode(l); err != nil {
		return err
	}
	fmt.Printf("\n\n")
	lines := strings.Split(l.Data, "\n")
	for _, l := range lines {
		re := regexp.MustCompile("(cpp:|sandbox_directory)")
		if !re.MatchString(l) {
			fmt.Printf("%s\n", l)
		}
	}
	return nil
}

func openLogPage(t *mesos.Task) error {
	if logsOpened {
		return nil
	}
	ip, err := t.SlaveIP()
	if err != nil {
		return err
	}
	dir, err := t.LogDirectory(ip)
	if err != nil {
		return fmt.Errorf("cannot fetch logs: %s\n", err)
	}
	p := strings.Split(dir, "/var/lib/mesos/slave/slaves/")
	slaveID := strings.Split(p[1], "/")[0]
	url := fmt.Sprintf("%s/mesos/#/agents/%s/browse?path=%s", viper.GetString("marathon_host"), slaveID, dir)
	cmd := exec.Command("open", url)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not open link to task logs")
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	logsOpened = true
	return nil
}

// scheduledTasks fetches tasks from Mesos API and returns list of tasks that
// match the given task name
func scheduledTasks(url, name string) (*mesos.Tasks, error) {
	var offset int
	tasks := &mesos.Tasks{}
	for {
		u := fmt.Sprintf("%s?offset=%d", url, offset)
		resp, err := http.Get(u)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		t := &mesos.Tasks{}
		if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
			return nil, err
		}
		for _, tt := range t.Tasks {
			tasks.Tasks = append(tasks.Tasks, tt)
		}
		if len(t.Tasks) >= 100 {
			offset += 100
			continue
		}
		break
	}

	rt, err := tasks.TasksFor(name)
	if err != nil {
		return nil, err
	}

	return rt, nil
}
