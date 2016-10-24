package notify

import (
	"encoding/json"
	"testing"
)

func TestEmoji(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{in: "production", out: ":balloon:"},
		{in: "stage", out: ""},
	}

	for _, test := range cases {
		e := Emoji(test.in)
		if e != test.out {
			t.Errorf("expected %s but got %s", test.out, e)
		}
	}
}

func TestMessageBody(t *testing.T) {
	type msg struct {
		Username string
		Text     string
	}
	m := &msg{}
	b := messageBody("yo", "yodawg")
	if err := json.Unmarshal([]byte(b), &m); err != nil {
		t.Errorf("expected message body to be valid JSON: %s", err)
	}
	if m.Username != "yo" || m.Text != "yodawg" {
		t.Error("expected JSON to be filled out with provided arguments")
	}
}
