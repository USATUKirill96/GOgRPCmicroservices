proto:
	protoc --go_out=. --go_opt=paths=source_relative \
               --go-grpc_out=. --go-grpc_opt=paths=source_relative \
               protobuf/*.proto

run-locations:
	go run locations/cmd/main.go

run-users:
	go run users/cmd/main.go

migrate-users-up:
	go run tools/migrate/migrate.go up

test:
	go test -v ./...

test-integration:
	go test -v ./... --tags=integration
