package main

import (
	pb "USATUKirill96/gridgo/protobuf"
	"USATUKirill96/gridgo/tools/logging"
	"USATUKirill96/gridgo/users/internal"
	"USATUKirill96/gridgo/users/pkg/user"
	"database/sql"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func main() {

	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	cfg := elasticsearch.Config{Addresses: []string{os.Getenv("ELASTICSEARCH_URL")}}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	logger := logging.NewLogger(es)

	// GRPC client setup
	conn, err := grpc.Dial(
		fmt.Sprintf("%v", os.Getenv("LOCATION_SERVICE_GRPC")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	locations := pb.NewLocationsClient(conn)

	// Database setup
	db, err := sql.Open("pgx", os.Getenv("POSTGRES_DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	userRepository := user.NewRepository(db)

	app := internal.Application{
		UserService: user.Service{
			Users:     &userRepository,
			Locations: locations,
		},
		Logger: logger,
	}

	err = app.Serve()
	logger.ERROR(err)
}
