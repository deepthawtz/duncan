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

// AllowedToManage checks if user is able to manage/deploy
// an app/env
func AllowedToManage(app, env string) (bool, error) {
	var txn []*consul.TxnItem
	txn = append(txn, &consul.TxnItem{
		KV: &consul.KVPair{
			Key:   fmt.Sprintf("deploys/%s/%s/acl_check", app, env),
			Value: base64.StdEncoding.EncodeToString([]byte("yes")),
			Verb:  "set",
		},
	})
	txn = append(txn, &consul.TxnItem{
		KV: &consul.KVPair{
			Key:  fmt.Sprintf("deploys/%s/%s/acl_check", app, env),
			Verb: "delete",
		},
	})
	client := &http.Client{}
	body, err := json.Marshal(txn)
	if err != nil {
		return false, err
	}
	url := consul.TxnURL()
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
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
		if err := json.Unmarshal(b, &d); err != nil {
			return err
		}
		if len(d) == 0 {
			break
		}
		for _, x := range d {
			if x.ID == id {
				continue
			} else {
				fmt.Println("DONE")
				return nil
			}
		}
	}
	fmt.Println("DONE")
	return nil
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
	org := viper.GetString("github_org")
	if org == "" {
		return "no github_org set: cannot generate diff link"
	}
	repo := viper.GetStringMapString("repos")[app]
	if repo == "" {
		repo = app
	}
	return fmt.Sprintf("https://github.com/%s/%s/compare/%s...%s", org, repo, prev, tag)
}
