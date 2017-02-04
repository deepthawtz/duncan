package logs

// Query represents criteria to insert into an Elasticsearch log query
type Query struct {
	AppName string
	Start   int64
	End     int64
}

const queryTemplate = `
{
  "query": {
    "filtered": {
      "query": {
        "query_string": {
          "analyze_wildcard": true,
          "query": "app_name:{{.AppName}}"
        }
      },
      "filter": {
        "bool": {
          "must": [
            {
              "range": {
                "@timestamp": {
                  "gte": {{.Start}},
                  "lte": {{.End}},
                  "format": "epoch_millis"
                }
              }
            }
          ]
        }
      }
    }
  },
  "size": 500,
  "sort": [
    {
      "@timestamp": {
        "order": "asc",
        "unmapped_type": "boolean"
      }
    }
  ],
  "fields": [
    "*",
    "_source"
  ],
  "fielddata_fields": [
    "@timestamp"
  ]
}
`
