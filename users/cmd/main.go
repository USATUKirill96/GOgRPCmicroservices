package main

import (
	pb "USATUKirill96/gridgo/protobuf"
	"USATUKirill96/gridgo/users/internal"
	"USATUKirill96/gridgo/users/pkg/user"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
)

func main() {

	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// GRPC client setup
	conn, err := grpc.Dial(
		fmt.Sprintf(":%v", os.Getenv("LOCATION_SERVICE_GRPC")),
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
		UserService: user.Service{Users: &userRepository, Locations: locations},
	}

	// Paths
	r := mux.NewRouter()
	r.HandleFunc("/update", app.UpdateLocation)
	r.HandleFunc("/users", app.FindByDistance)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%v", os.Getenv("USER_SERVICE_PORT")),
	}

	fmt.Println("Server started and running")
	err = srv.ListenAndServe()
	fmt.Println(err)
}
