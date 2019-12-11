package deployment

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/betterdoctor/duncan/consul"
	"github.com/spf13/viper"
)

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
