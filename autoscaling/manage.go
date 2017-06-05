package autoscaling

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/betterdoctor/slythe/policy"
	"github.com/spf13/viper"
)

// GetPolicies returns all autoscaling policies optionally filtering
// if app or env are not empty string
func GetPolicies(app, env string) (*policy.Policies, error) {
	policies := &policy.Policies{}
	resp, err := http.Get(viper.GetString("SLYTHE_HOST") + "/")
	if err != nil {
		return policies, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(policies); err != nil {
		return policies, err
	}
	if app == "" && env == "" {
		return policies, nil
	}
	if env == "" {
		env = "stage|production"
	}
	fp := &policy.Policies{}
	for _, e := range strings.Split(env, "|") {
		for _, cp := range policies.CPUScaled {
			if app == "" && cp.Environment == e {
				fp.CPUScaled = append(fp.CPUScaled, cp)
			}
			if app != "" && cp.AppName == app && cp.Environment == e {
				fp.CPUScaled = append(fp.CPUScaled, cp)
			}
		}
	}
	for _, e := range strings.Split(env, "|") {
		for _, cp := range policies.QueueLengthScaled {
			if app == "" && cp.Environment == e {
				fp.QueueLengthScaled = append(fp.QueueLengthScaled, cp)
			}
			if app != "" && cp.AppName == app && cp.Environment == e {
				fp.QueueLengthScaled = append(fp.QueueLengthScaled, cp)
			}
		}
	}
	return fp, nil
}

// CreateWorkerPolicy creates an autoscaling worker policy
func CreateWorkerPolicy(wp *policy.Worker) error {
	j, err := json.Marshal(wp)
	if err != nil {
		return err
	}
	url := viper.GetString("SLYTHE_HOST") + "/policies/worker"
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(j)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to POST autoscaling policy: %s\n%s", resp.Status, b)
	}
	return nil
}

// UpdateWorkerPolicy updates an autoscaling worker policy
func UpdateWorkerPolicy(wp *policy.Worker) error {
	j, err := json.Marshal(wp)
	if err != nil {
		return err
	}
	url := viper.GetString("SLYTHE_HOST") + "/policies/worker"
	req, _ := http.NewRequest("PUT", url, strings.NewReader(string(j)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to POST autoscaling policy: %s\n%s", resp.Status, b)
	}
	return nil
}

// CreateCPUPolicy creates an autoscaling worker policy
func CreateCPUPolicy(cp *policy.CPU) error {
	j, err := json.Marshal(cp)
	if err != nil {
		return err
	}
	url := viper.GetString("SLYTHE_HOST") + "/policies/cpu"
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(j)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to POST autoscaling policy: %s\n%s", resp.Status, b)
	}
	return nil
}

// UpdateCPUPolicy creates an autoscaling worker policy
func UpdateCPUPolicy(cp *policy.CPU) error {
	j, err := json.Marshal(cp)
	if err != nil {
		return err
	}
	url := viper.GetString("SLYTHE_HOST") + "/policies/cpu"
	req, _ := http.NewRequest("PUT", url, strings.NewReader(string(j)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to POST autoscaling policy: %s\n%s", resp.Status, b)
	}
	return nil
}
