FROM golang:1.18
WORKDIR /src/gridgo/
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app users/cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate tools/migrate/migrate.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=0 /src/gridgo/app ./
COPY --from=0 /src/gridgo/migrate ./
COPY --from=0 /src/gridgo/.env ./

RUN mkdir -p users/migrations
ADD users/migrations/* users/migrations/

CMD ./migrate up ; ./app