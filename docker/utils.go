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
	Tags []Tag `json:"tags"`
}

// Tag represents a Docker repository tag
type Tag struct {
	Name string `json:"name"`
}

// TagExists checks if a docker tag exists for a given repo
func TagExists(app, tag string) bool {
	url := fmt.Sprintf("https://quay.io/api/v1/repository/betterdoctor/%s/tag/?specificTag=%s", app, tag)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, strings.NewReader(""))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("quay_token")))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	tr := &TagResponse{}
	if err := json.NewDecoder(resp.Body).Decode(tr); err != nil {
		fmt.Println(err)
		return false
	}
	if len(tr.Tags) == 0 || tr.Tags[0].Name != tag {
		return false
	}
	return true
}
