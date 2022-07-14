package main

import (
	pb "USATUKirill96/gridgo/protobuf"
	"USATUKirill96/gridgo/users/internal"
	"USATUKirill96/gridgo/users/pkg/user"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

func main() {
	// GRPC client setup
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	locations := pb.NewLocationsClient(conn)

	// Database setup
	userRepository := user.NewFakeRepository()

	app := internal.Application{
		UserService: user.Service{Users: &userRepository, Locations: locations},
	}

	// Paths
	r := mux.NewRouter()
	r.HandleFunc("/", app.UpdateLocation)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
	}

	fmt.Println("Server started and running")
	err = srv.ListenAndServe()
	fmt.Println(err)
}
