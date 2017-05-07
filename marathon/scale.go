package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/betterdoctor/duncan/deployment"
)

// ScaleEvent represents a Marathon Group scale event
//
// One or more containers within a Group may be scaled
// up or down at once
type ScaleEvent map[string]map[string]int

// Scale increases or decreases number of running instances of
// an application within a Marathon Group
func Scale(app, env string, procs []string) (ScaleEvent, error) {
	groups, err := listGroups()
	if err != nil {
		return nil, err
	}

	var (
		scaled ScaleEvent
		mj     []byte
	)
	for _, g := range groups.Groups {
		if g.ID == deployment.MarathonGroupID(app, env) {
			tag, err := deployment.CurrentTag(app, env)
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
			if err := deployment.Watch(d.ID); err != nil {
				return nil, err
			}
		}
	}
	return scaled, nil
}

func scaledMarathonJSON(group *Group, app, env, tag string, procs []string) (ScaleEvent, []byte, error) {
	var scaled = make(ScaleEvent)
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

			if a.ID == fmt.Sprintf("%s/%s", deployment.MarathonGroupID(app, env), proc) {
				prev := a.Instances
				if prev == count {
					return nil, []byte(""), fmt.Errorf("already running %d instances of %s", count, a.ID)
				}
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
