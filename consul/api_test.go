package consul

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestEnvMap(t *testing.T) {
	data := map[string][]byte{
		"YODAWG_LEVEL": []byte("9000"),
		"FOO_ENABLED":  []byte("true"),
	}
	var kvs []KVPair
	for k, v := range data {
		key := fmt.Sprintf("env/yo/stage/%s", k)
		value := base64.StdEncoding.EncodeToString(v)
		kvp := &KVPair{
			Key:   key,
			Value: value,
		}
		kvs = append(kvs, *kvp)
	}

	m := envMap(kvs)

	if len(m) != len(data) {
		t.Errorf("expected %v got %v", len(data), len(m))
	}

	for _, kvp := range kvs {
		p := strings.Split(kvp.Key, "/")
		key := p[len(p)-1]
		value, _ := base64.StdEncoding.DecodeString(kvp.Value)
		if m[key] != string(value) {
			t.Errorf("expected %s but got %s", m[key], value)
		}
	}
}
