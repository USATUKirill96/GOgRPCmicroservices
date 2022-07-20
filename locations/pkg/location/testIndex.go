package location

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var ES7TestIndex = "users_test_1.0"

func NewTestIndex() TestIndex {
	// Load .env
	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatal(err)
	}

	// Setup elasticsearch client and repository
	cfg := elasticsearch.Config{Addresses: []string{os.Getenv("ELASTICSEARCH_URL")}}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return TestIndex{es}
}

type TestIndex struct {
	ES *elasticsearch.Client
}

func (ti TestIndex) NewRepository() Repository {
	return *NewRepository(ti.ES, ES7TestIndex)
}

func (ti TestIndex) TearDown() {
	_, err := ti.ES.Indices.Delete([]string{ES7TestIndex})
	if err != nil {
		log.Fatal(err)
	}
}
