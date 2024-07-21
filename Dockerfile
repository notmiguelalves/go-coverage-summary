FROM golang:1.22.5-alpine

COPY main.go main.go
COPY go.mod go.mod

RUN go mod tidy
RUN go build -o=cov main.go

ENTRYPOINT ["cov"]
