//go:build integration

package location

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"io"
	"log"
	"strings"
	"testing"
	"time"
)

func assertLocations(l1, l2 Location) (string, bool) {
	formatError := func(f1, f2 interface{}) string {
		return fmt.Sprintf("Asserting Locations: expected %v, got %v", f1, f2)
	}
	if l1.Username != l2.Username {
		return formatError(l1.Username, l2.Username), false
	}
	if l1.Longitude != l2.Longitude {
		return formatError(l1.Longitude, l2.Latitude), false
	}
	if l1.Latitude != l2.Latitude {
		return formatError(l1.Latitude, l2.Latitude), false
	}
	if !l1.Updated.Equal(l2.Updated) {
		return formatError(l1.Updated, l2.Updated), false
	}
	return "", true
}

func TestRepository_Insert(t *testing.T) {

	ti := NewTestIndex("./../../../.env")
	r := ti.NewRepository()
	defer ti.TearDown()

	l := Location{
		Username:  "TestInsert",
		Longitude: 12.3456,
		Latitude:  -12.3456,
		Updated:   time.Now(),
	}
	res, err := r.Insert(l)
	if err != nil {
		t.Error(err)
	}
	message, ok := assertLocations(l, res)
	if !ok {
		t.Errorf("Incorrect function return: %v", message)
	}

	// Assert saved data
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"username": l.Username,
			},
		},
	}

	es7res, err := ti.ES.Search(
		r.es.Search.WithIndex(r.index),
		r.es.Search.WithBody(esutil.NewJSONReader(body)),
		r.es.Search.WithPretty(),
		r.es.Search.WithSort("updated"),
	)
	if err != nil {
		log.Fatalf("ERROR getting response: %s", err)
	}
	defer es7res.Body.Close()

	if es7res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(es7res.Body).Decode(&e); err != nil {
			log.Fatalf("ERROR parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				es7res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	buf := new(strings.Builder)
	io.Copy(buf, es7res.Body)

	var result map[string]interface{}

	// How to parse response: https://www.elastic.co/guide/en/elasticsearch/reference/7.17/search-your-data.html
	json.Unmarshal([]byte(fmt.Sprint(buf)), &result)
	sl := FromMap(result["hits"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{}))
	message, ok = assertLocations(l, sl)
	if !ok {
		t.Errorf("Incorrect data saved: %v", message)
	}
}

func TestRepository_Find(t *testing.T) {

	ti := NewTestIndex("./../../../.env")
	r := ti.NewRepository()
	defer ti.TearDown()

	var day1 time.Time = time.Date(2021, 8, 10, 0, 0, 0, 100, time.Local)
	var day2 time.Time = time.Date(2021, 8, 20, 0, 0, 0, 100, time.Local)
	var day3 time.Time = time.Date(2021, 8, 30, 0, 0, 0, 100, time.Local)

	ls := []Location{
		{"1", "User0", 12.54321, -12.54321, day1},
		{"2", "User1", 12.54321, -12.54321, day1},
		{"3", "User1", 15.54321, -15.54321, day2},
		{"4", "User1", 17.54321, -17.54321, day3},
	}

	for _, l := range ls {
		_, err := r.Insert(l)
		if err != nil {
			log.Fatal(err)
		}
	}

	cases := []struct {
		username string
		after    time.Time
		before   time.Time
		expected int
	}{
		{"User0", day1, day2, 1}, // exclude User1
		{"User1", day1, day3, 3}, // include not-equal time
		{"User1", day2, day3, 2}, // exclude day 1
		{"User1", day3, day3, 1}, // only one location
	}

	for _, tc := range cases {
		locations, err := r.Find(tc.username, tc.after, tc.before)
		if err != nil {
			t.Error(err)
		}
		if len(locations) != tc.expected {
			t.Errorf(
				"find Locations: incorrect number of results. Expected %v, got %v", tc.expected, len(locations),
			)
		}
	}
}
