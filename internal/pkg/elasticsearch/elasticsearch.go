// Package elasticsearch provides the elasticsearch client.
package elasticsearch

import (
	"emperror.dev/errors"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

// NewElasticClient creates a new elasticsearch client.
func NewElasticClient(cfg *ElasticOptions) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.URL},
	})
	if err != nil {
		return nil, errors.WrapIf(err, "v8.elasticsearch")
	}

	return es, nil
}
