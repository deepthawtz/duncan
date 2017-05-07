package deployment

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/betterdoctor/duncan/consul"
	"github.com/spf13/viper"
)

// Deployment represents a Marathon deployment
type Deployment struct {
	ID string `json:"id"`
}

// BeginDeploy checks if Consul ACL allows deployments for
func BeginDeploy(app, env string) error {
	var txn []*consul.TxnItem
	txn = append(txn, &consul.TxnItem{
		KV: &consul.KVPair{
			Key:   fmt.Sprintf("deploys/%s/%s/lock", app, env),
			Value: base64.StdEncoding.EncodeToString([]byte("yo")),
			Verb:  "set",
		},
	})
	client := &http.Client{}
	body, err := json.Marshal(txn)
	if err != nil {
		return err
	}
	url := consul.TxnURL()
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to set deployment lock: %s\n%s\n", resp.Status, string(b))
	}
	return nil
}

// FinishDeploy removes deployment lock after a successful deploy
func FinishDeploy(app, env string) error {
	var txn []*consul.TxnItem
	txn = append(txn, &consul.TxnItem{
		KV: &consul.KVPair{
			Key:  fmt.Sprintf("deploys/%s/%s/lock", app, env),
			Verb: "delete",
		},
	})
	client := &http.Client{}
	body, err := json.Marshal(txn)
	if err != nil {
		return err
	}
	url := consul.TxnURL()
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to remove deployment lock in Consul: %s\n%s\n", resp.Status, string(b))
	}
	return nil
}

// Watch watches a Marathon deployment and handles success or failure
func Watch(id string) error {
	if id == "" {
		return fmt.Errorf("did not get a deployment id from Marathon API")
	}
	fmt.Printf("Waiting for deployment id: %s\n", id)

	url := fmt.Sprintf("%s/service/marathon/v2/deployments", viper.GetString("marathon_host"))

	for {
		time.Sleep(5 * time.Second)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		var d []Deployment
		fmt.Println(string(b))
		if err := json.Unmarshal(b, &d); err != nil {
			return err
		}
		for _, x := range d {
			if x.ID == id {
				continue
			}
		}
		fmt.Println("DONE")
		break
	}

	return nil
}

// UpdateReleaseTags updates the deployment git tags in Consul KV registry
// `tags/{app}/{env}/current` points to the currently deployed tag
// `tags/{app}/{env}/previous` points to the previously deployed tag
//
// This structure allows for rollback if a previous tag exists
//
// Returns previously deployed git tag if one has been deployed
func UpdateReleaseTags(app, env, tag string) (string, error) {
	prev, err := CurrentTag(app, env)
	if err != nil {
		return "", err
	}

	if prev == tag {
		return tag, nil
	}
	m := map[string]string{
		"current":  tag,
		"previous": prev,
	}

	var txn []*consul.TxnItem
	for k, v := range m {
		txn = append(txn, &consul.TxnItem{
			KV: &consul.KVPair{
				Key:   fmt.Sprintf("deploys/%s/%s/%s", app, env, k),
				Value: base64.StdEncoding.EncodeToString([]byte(v)),
				Verb:  "set",
			},
		})
	}
	client := &http.Client{}
	body, err := json.Marshal(txn)
	if err != nil {
		return "", err
	}
	url := consul.TxnURL()
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("failed to update release tags in Consul: %s\n%s", resp.Status, string(b))
	}
	return prev, nil
}

// CurrentTag returns the currently deployed git tag for an app and environment
func CurrentTag(app, env string) (string, error) {
	url := consul.CurrentDeploymentTagURL(app, env)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", fmt.Errorf("could not fetch current release tag: %s", resp.Status)
}

// MarathonGroupID returns a Marathon Group id for an app and env
func MarathonGroupID(app, env string) string {
	return "/" + strings.Join([]string{app, env}, "-")
}

// GithubDiffLink returns a GitHub diff link to view deployment changes
func GithubDiffLink(app, prev, tag string) string {
	if prev == tag || prev == "" {
		return "no changes"
	}
	// TODO: handle if YAML is not filled out correctly
	repo := viper.GetStringMapString("repos")[app]
	return fmt.Sprintf("https://github.com/betterdoctor/%s/compare/%s...%s", repo, prev, tag)
}
