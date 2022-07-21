package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"net/http"
	"runtime/debug"
	"time"
)

type record struct {
	Level     string `json:"level"` // ERROR/INFO
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"` // GET/POST/UPDATE/DELETE
	Params    string `json:"params"`
	Body      string `json:"body"`
	URL       string `json:"url"`
	Source    string `json:"source"` //ip address
	Message   string `json:"message"`
}

func newRecord(level, message string, r *http.Request) record {
	record := record{
		Level:     level,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	if r != nil {
		record.Method = r.Method
		record.URL = r.URL.Path
		record.Params = r.URL.Query().Encode()
		record.Source = r.RemoteAddr
		if r.Method == "POST" {
			defer r.Body.Close()
			b, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
			}
			record.Body = string(b)
		}
	}
	return record
}

func NewLogger(es *elasticsearch.Client) Logger { return Logger{es} }

type Logger struct {
	es *elasticsearch.Client
}

func (l Logger) insert(r record) {

	currentIndex := fmt.Sprintf("logging_%v", time.Now().Year())

	data, _ := json.Marshal(r)
	req := esapi.IndexRequest{
		Index: currentIndex,
		Body:  bytes.NewReader(data),
	}
	req.Do(context.Background(), l.es)
}

func (l Logger) new(level string, message string, r ...*http.Request) {
	fmt.Println(message)
	var request *http.Request
	if len(r) > 0 {
		request = r[0]
	} else {
		request = nil
	}
	record := newRecord(level, message, request)
	l.insert(record)
}

func (l Logger) ERROR(err error, r ...*http.Request) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	go l.new("ERROR", trace, r...)
}
func (l Logger) INFO(message string, r ...*http.Request) {
	go l.new("INFO", message, r...)
}
