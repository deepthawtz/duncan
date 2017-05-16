package docker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// TagResponse represents a Quay API tags response
type TagResponse struct {
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
}

// VerifyTagExists checks if a docker tag exists for a given repo
func VerifyTagExists(app, tag string) error {
	url := tagsURL(app, tag)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, strings.NewReader(""))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("quay_token")))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to find tag: %s", resp.Status)
	}

	tr := &TagResponse{}
	if err := json.NewDecoder(resp.Body).Decode(tr); err != nil {
		return err
	}
	if len(tr.Tags) == 0 {
		return fmt.Errorf("failed to find tag")
	}
	return nil
}

func tagsURL(app, tag string) string {
	host := viper.GetString("docker_registry_host")
	if host == "" {
		host = "https://quay.io"
	}
	return fmt.Sprintf("%s/api/v1/repository/betterdoctor/%s/tag/?specificTag=%s", host, app, tag)
}
