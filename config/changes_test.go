package config

import (
	"fmt"
	"os/user"
	"testing"
)

func TestChanges(t *testing.T) {
	u, _ := user.Current()
	updated := map[string][]string{
		"FOO": []string{"yabba", "dooo"},
	}
	added := map[string][]string{
		"BAR": []string{"new-thing"},
	}
	deleted := map[string][]string{
		"BAZ": []string{},
	}
	cases := []struct {
		cmd     string
		changes map[string][]string
		output  string
	}{
		{cmd: "env", changes: map[string][]string{}, output: ""},
		{cmd: "secrets", changes: map[string][]string{}, output: ""},
		{cmd: "env", changes: updated, output: fmt.Sprintf("env updated by %s:\n`FOO` updated from `yabba` => `dooo`\n", u.Username)},
		{cmd: "secrets", changes: updated, output: fmt.Sprintf("secrets updated by %s:\n`FOO` updated\n", u.Username)},
		{cmd: "env", changes: added, output: fmt.Sprintf("env updated by %s:\n`BAR` set to `new-thing`\n", u.Username)},
		{cmd: "secrets", changes: added, output: fmt.Sprintf("secrets updated by %s:\n`BAR` added\n", u.Username)},
		{cmd: "env", changes: deleted, output: fmt.Sprintf("env updated by %s:\n`BAZ` deleted\n", u.Username)},
		{cmd: "secrets", changes: deleted, output: fmt.Sprintf("secrets updated by %s:\n`BAZ` deleted\n", u.Username)},
	}

	for _, test := range cases {
		msg := Changes(test.cmd, test.changes)
		if msg != test.output {
			t.Errorf("expected %s but got %s", test.output, msg)
		}
	}
}
