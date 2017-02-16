package logs

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var (
	types  []string
	colors = []func(...interface{}) string{
		color.New(color.FgCyan, color.Bold).SprintFunc(),
		color.New(color.FgYellow, color.Bold).SprintFunc(),
		color.New(color.FgMagenta, color.Bold).SprintFunc(),
		color.New(color.FgGreen, color.Bold).SprintFunc(),
		color.New(color.FgRed, color.Bold).SprintFunc(),
		color.New(color.FgBlue, color.Bold).SprintFunc(),
		color.New(color.FgWhite, color.Bold).SprintFunc(),
	}
	appTypes = map[string]string{}
)

// Result represents an Elasticsearch search result
type Result struct {
	Took int `json:"took"`
	Hits struct {
		Total int `json:"total"`
		Hits  []struct {
			Source struct {
				TimeStamp time.Time `json:"@timestamp"`
				AppName   string    `json:"app_name"`
				AppType   string    `json:"app_type"`
				Message   string    `json:"message"`
			} `json:"_source"`
		}
	} `json:"hits"`
}

// Print displays the logs for a log result
func (r *Result) Print(utc bool) {
	for _, h := range r.Hits.Hits {
		s := h.Source
		fmt.Printf("%s [%s] %s\n",
			displayTime(s.TimeStamp, utc),
			colorAppType(s.AppType),
			s.Message,
		)
	}
}

func colorAppType(t string) string {
	if _, ok := appTypes[t]; !ok {
		types = append(types, t)
		for i, x := range types {
			if x == t {
				appTypes[t] = colors[i](t)
			}
		}
	}
	return appTypes[t]
}

func displayTime(t time.Time, utc bool) time.Time {
	var l string
	if utc {
		l = "UTC"
	} else {
		l = "Local"
	}
	loc, err := time.LoadLocation(l)
	if err != nil {
		panic(err)
	}
	return t.In(loc)
}
