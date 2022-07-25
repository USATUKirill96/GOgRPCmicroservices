FROM golang:1.18
WORKDIR /src/gridgo/
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app locations/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /src/gridgo/app ./
COPY --from=0 /src/gridgo/.env ./
CMD ["./app"]