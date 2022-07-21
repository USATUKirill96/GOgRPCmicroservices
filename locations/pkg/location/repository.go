package location

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"io"
	"log"
	"strings"
	"time"
)

func NewRepository(es *elasticsearch.Client, index string) *Repository { return &Repository{index, es} }

type Repository struct {
	index string
	es    *elasticsearch.Client
}

func (r Repository) Insert(l Location) (Location, error) {
	data, err := json.Marshal(l)
	if err != nil {
		log.Println(err)
		return l, err
	}

	req := esapi.IndexRequest{
		Index: r.index,
		Body:  bytes.NewReader(data),
	}

	res, err := req.Do(context.Background(), r.es)
	if err != nil {
		log.Println(err)
		return l, err
	}
	if res.IsError() {
		err := errors.New(fmt.Sprintf("[%s] ERROR indexing document", res.Status()))
		return l, err
	}
	return l, nil
}

func (r Repository) Find(u string, after, before time.Time) ([]Location, error) {

	// How to build queries https://www.elastic.co/guide/en/elasticsearch/reference/7.17/query-dsl-bool-query.html
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"username": u,
						},
					},
					{
						"range": map[string]interface{}{
							"updated": map[string]interface{}{
								"gte": after,
								"lte": before,
							},
						},
					},
				},
			},
		},
	}

	res, err := r.es.Search(
		r.es.Search.WithIndex(r.index),
		r.es.Search.WithBody(esutil.NewJSONReader(body)),
		r.es.Search.WithPretty(),
		r.es.Search.WithSort("updated"),
	)
	if err != nil {
		log.Fatalf("ERROR getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("ERROR parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	buf := new(strings.Builder)
	io.Copy(buf, res.Body)

	var result map[string]interface{}
	var locations []Location

	// How to parse response: https://www.elastic.co/guide/en/elasticsearch/reference/7.17/search-your-data.html
	json.Unmarshal([]byte(fmt.Sprint(buf)), &result)
	for _, m := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		locations = append(locations, FromMap(m.(map[string]interface{})))
	}
	return locations, nil
}
