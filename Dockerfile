FROM golang:1.22.5-alpine

COPY main.go /main.go

ENTRYPOINT ["/main.go"]
