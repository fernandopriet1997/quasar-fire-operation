# Build stage
FROM golang:1.16-alpine as build

WORKDIR ../Go/apps
COPY . .

ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -v -o go-app

# Run stage
FROM alpine:3.11
COPY --from=build ../Go/apps/ app/
CMD ["./app/go-app"]