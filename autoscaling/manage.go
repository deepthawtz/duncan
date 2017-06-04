package autoscaling

import (
	"encoding/json"
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
