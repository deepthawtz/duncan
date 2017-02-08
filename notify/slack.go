package notify

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/user"
	"strings"

	"github.com/spf13/viper"
)

// Slack notifies slack with the botname to use and message
func Slack(botname, message string) error {
	url := viper.GetString("slack_webhook_url")
	msg := messageBody(botname, message)
	if url == "" {
		fmt.Println("slack_webhook_url not set, skipping notification...")
		fmt.Println(msg)
		return nil
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(msg))
	if err != nil {
		return fmt.Errorf("slack notification failed: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("slack notification failed: %s", b)
	}

	return nil
}

// ConfigChange notifies team when ENV or secrets are updated
func ConfigChange(typ, app, deployEnv string, updated map[string][]string) {
	if len(updated) == 0 {
		return
	}
	var changes string
	for k, v := range updated {
		if len(v) == 2 {
			if typ == "env" {
				changes += fmt.Sprintf("`%s` updated from `%s` => `%s`\n", k, v[0], v[1])
			} else {
				changes += fmt.Sprintf("`%s` updated\n", k)
			}
		} else if len(v) == 0 {
			changes += fmt.Sprintf("`%s` deleted\n", k)
		} else {
			if typ == "env" {
				changes += fmt.Sprintf("`%s` set to `%s`\n", k, v[0])
			} else {
				changes += fmt.Sprintf("`%s` added\n", k)
			}
		}
	}
	u, _ := user.Current()
	msg := fmt.Sprintf("%s updated by %s:\n%s", typ, u.Username, changes)
	Slack(
		fmt.Sprintf("%s %s", app, deployEnv),
		msg,
	)
}

// Emoji returns an icon to visually distinguish a production
// deployment from a stage deployment
func Emoji(env string) string {
	if env == "production" {
		return ":balloon:"
	}

	return ""
}

func messageBody(botname, message string) string {
	return fmt.Sprintf(`{
  "username": "%s",
  "text": "%s"
}`, botname, message)
}
