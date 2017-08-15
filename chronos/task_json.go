package chronos

// Template for Chronos API version >= 3.0.0
var taskTemplate = `
{
  "name": "{{.TaskName}}",
  "description": "one-off duncan run ({{.Command}})",
  "schedule": "{{.Schedule}}",
  "scheduleTimeZone": "PST",
  "retries": 0,
  "container": {
    "type": "DOCKER",
    "network": "HOST",
    "image": "{{.DockerRepoPrefix}}/{{.App}}:{{.Tag}}"
  },
  "command": "envconsul -config envconsul-{{.Env}}.hcl {{.Command}}",
  "cpus": "{{.CPU}}",
  "mem": "{{.Mem}}",
  "fetch": [
    {
      "uri": "{{.DockerConfURL}}",
      "cache": false,
      "extract": true,
      "executable": false
    }
  ]
}
`
