package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/viper"
)

// Stream prints out the logs for the given app and env
func Stream(app, env string, utc bool) {
	t := time.Now().UTC()
	index := fmt.Sprintf("%s-%s-docker-%s", app, env, t.Format("2006.01.02"))
	url := fmt.Sprintf("%s/%s/_search", viper.GetString("elasticsearch_host"), index)
	for {
		q, err := buildQuery(app, env)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, strings.NewReader(q))
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			r := &Result{}
			if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			r.Print(utc)
			time.Sleep(1000 * time.Millisecond)
			continue
		} else {
			fmt.Println(resp.Status)
			fmt.Println(url)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

// buildQuery generates an Elasticsearch query document for an app/env
// and provides a sliding time range with 1 minute lookback
func buildQuery(app, env string) (string, error) {
	end := time.Now().Unix() * 1000
	query := &Query{
		AppName: fmt.Sprintf("%s-%s", app, env),
		Start:   end - 60*1000,
		End:     end,
	}
	t, err := template.New("elasticsearch_query").Parse(queryTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}
	j := new(bytes.Buffer)
	if err := t.Execute(j, query); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}
	return j.String(), nil
}
