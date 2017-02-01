package chronos

// Template for Chronos API version >= 3.0.0
var taskTemplate = `
{
  "name": "{{.TaskName}}",
  "description": "one-off duncan run ({{.Command}})",
  "schedule": "R1//PT30M",
  "retries": 0,
  "owner": "ops@betterdoctor.com",
  "container": {
    "type": "DOCKER",
    "network": "HOST",
    "image": "quay.io/betterdoctor/{{.App}}:{{.Tag}}"
  },
  "command": "envconsul -config envconsul-{{.Env}}.hcl {{.Command}}",
  "cpus": "1.0",
  "mem": "1024",
  "fetch": [
    {
      "uri": "https://s3.amazonaws.com/betterdoctor-operations-qhtumyvauxvxorwmeujn/Configs/docker.tar.gz",
      "cache": false,
      "extract": true,
      "executable": false
    }
  ]
}
`
