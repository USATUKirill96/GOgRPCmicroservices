package location

import (
	"time"
)

type Location struct {
	ID        string    `json:"-"`
	Username  string    `json:"username"`
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
	Updated   time.Time `json:"updated"`
}

func FromMap(m map[string]interface{}) Location {
	// How to parse response: https://www.elastic.co/guide/en/elasticsearch/reference/7.17/search-your-data.html

	l := Location{}
	l.ID = m["_id"].(string)

	var s map[string]interface{} = m["_source"].(map[string]interface{})
	l.Username = s["username"].(string)
	l.Longitude = s["longitude"].(float64)
	l.Latitude = s["latitude"].(float64)
	u, _ := time.Parse(time.RFC3339, s["updated"].(string))
	l.Updated = u
	return l
}
