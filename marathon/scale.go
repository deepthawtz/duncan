package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/betterdoctor/duncan/deploy"
)

type scaleEvent map[string]map[string]int

// Scale increases or decreases number of running instances of
// an application within a Marathon Group
func Scale(app, env string, procs []string) (scaleEvent, error) {
	groups, err := listGroups()
	if err != nil {
		return nil, err
	}

	var (
		scaled scaleEvent
		mj     []byte
	)
	for _, g := range groups.Groups {
		if g.ID == "/"+strings.Join([]string{app, env}, "-") {
			tag, err := deploy.CurrentTag(app, env, nil)
			if err != nil {
				return nil, err
			}
			scaled, mj, err = scaledMarathonJSON(&g, app, env, tag, procs)
			if err != nil {
				return nil, err
			}
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", deploymentURL(), bytes.NewReader(mj))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			d := &DeploymentResponse{}
			if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
				return nil, err
			}
			if err := waitForDeployment(d.ID); err != nil {
				return nil, err
			}
		}
	}
	return scaled, nil
}

func scaledMarathonJSON(group *Group, app, env, tag string, procs []string) (scaleEvent, []byte, error) {
	var scaled = make(scaleEvent)
	for _, a := range group.Apps {
		for _, proc := range procs {
			s := strings.Split(proc, "=")
			proc := s[0]
			count, err := strconv.Atoi(s[1])
			if err != nil {
				return nil, []byte(""), err
			}

			if count < 0 {
				return nil, []byte(""), fmt.Errorf("cannot scale %s below zero", proc)
			}

			if a.ID == "/"+strings.Join([]string{app, env}, "-")+"/"+proc {
				prev := a.Instances
				fmt.Printf("scaling %s from %d to %d\n", proc, prev, count)
				a.Instances = count

				scaled[strings.Split(a.ID, "/")[2]] = map[string]int{
					"prev": prev,
					"curr": count,
				}

				a.UpdateReleaseTag(tag)
			}
		}
	}

	j, err := json.Marshal(group)
	return scaled, j, err
}
