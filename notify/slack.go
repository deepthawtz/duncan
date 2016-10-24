package notify

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func Slack(botname, message string) error {
	url := viper.GetString("slack_webhook_url")
	msg := messageBody(botname, message)
	if url == "" {
		fmt.Println("slack_webhook_url not set, skipping notification...")
		fmt.Println(msg)
		return nil
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(msg))
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("slack notification failed: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("slack notification failed: %s", b)
	}

	fmt.Println("notified slack")
	return nil
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
