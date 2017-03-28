package config

import (
	"fmt"
	"os/user"
)

// Changes provides a human-readable message summarizing what env/secrets
// changes have taken place by which UNIX user
func Changes(typ string, updated map[string][]string) string {
	if len(updated) == 0 {
		return ""
	}
	var changes string
	for k, v := range updated {
		if len(v) == 0 {
			changes += fmt.Sprintf("`%s` deleted\n", k)
		} else if len(v) == 2 {
			if typ == "env" {
				changes += fmt.Sprintf("`%s` updated from `%s` => `%s`\n", k, v[0], v[1])
			} else {
				changes += fmt.Sprintf("`%s` updated\n", k)
			}
		} else {
			if typ == "env" {
				changes += fmt.Sprintf("`%s` set to `%s`\n", k, v[0])
			} else {
				changes += fmt.Sprintf("`%s` added\n", k)
			}
		}
	}
	u, _ := user.Current()
	return fmt.Sprintf("%s updated by %s:\n%s", typ, u.Username, changes)
}
