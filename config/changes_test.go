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
		typ     string
		changes map[string][]string
		output  string
	}{
		{typ: "env", changes: map[string][]string{}, output: ""},
		{typ: "secrets", changes: map[string][]string{}, output: ""},
		{typ: "env", changes: updated, output: fmt.Sprintf("env updated by %s:\n`FOO` updated from `yabba` => `dooo`\n", u.Username)},
		{typ: "secrets", changes: updated, output: fmt.Sprintf("secrets updated by %s:\n`FOO` updated\n", u.Username)},
		{typ: "env", changes: added, output: fmt.Sprintf("env updated by %s:\n`BAR` set to `new-thing`\n", u.Username)},
		{typ: "secrets", changes: added, output: fmt.Sprintf("secrets updated by %s:\n`BAR` added\n", u.Username)},
		{typ: "env", changes: deleted, output: fmt.Sprintf("env updated by %s:\n`BAZ` deleted\n", u.Username)},
		{typ: "secrets", changes: deleted, output: fmt.Sprintf("secrets updated by %s:\n`BAZ` deleted\n", u.Username)},
	}

	for _, test := range cases {
		msg := Changes(test.typ, test.changes)
		if msg != test.output {
			t.Errorf("expected %s but got %s", test.output, msg)
		}
	}
}
