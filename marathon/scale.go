package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Scale increases or decreases number of running instances of
// an application within a Marathon Group
func Scale(group *Group, rules map[string]int) (string, error) {
	mj, err := scaledMarathonJSON(group, rules)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", updateGroupURL(), bytes.NewReader(mj))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to update Marathon group: %s", resp.Status)
	}
	d := &deploymentResponse{}
	if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
		return "", err
	}
	return d.ID, nil
}

func scaledMarathonJSON(group *Group, rules map[string]int) ([]byte, error) {
	var out []byte
	for _, a := range group.Apps {
		for proc, count := range rules {
			if a.InstanceType() == proc {
				a.Instances = count
			}
		}
	}

	out, err := json.Marshal(group)
	return out, err
}
