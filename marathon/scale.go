package marathon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/betterdoctor/duncan/deploy"
)

// Scale increases or decreases number of running instances of
// an application within a Marathon Group
func Scale(app, env string, procs []string) error {
	groups, err := listGroups()
	if err != nil {
		return err
	}

	for _, g := range groups.Groups {
		if g.ID == "/"+strings.Join([]string{app, env}, "-") {
			marathonJSON, err := scaledMarathonJSON(&g, app, env, procs)
			if err != nil {
				return err
			}
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", deploymentURL(), strings.NewReader(marathonJSON))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			d := &DeploymentResponse{}
			if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
				return err
			}
			if err := waitForDeployment(d.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func scaledMarathonJSON(group *Group, app, env string, procs []string) (string, error) {
	tag, err := deploy.CurrentTag(app, env, nil)
	if err != nil {
		return "", err
	}
	for _, a := range group.Apps {
		for _, proc := range procs {
			s := strings.Split(proc, "=")
			proc := s[0]
			count, err := strconv.Atoi(s[1])
			if err != nil {
				return "", err
			}

			if count < 0 {
				return "", fmt.Errorf("cannot scale %s below zero", proc)
			}

			if a.ID == "/"+strings.Join([]string{app, env}, "-")+"/"+proc {
				prev := a.Instances
				fmt.Printf("scaling %s from %d to %d\n", proc, prev, count)
				a.Instances = count

				a.UpdateReleaseTag(tag)
			}
		}
	}

	j, err := json.Marshal(group)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
