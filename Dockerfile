FROM golang:1.22.5-alpine

COPY main.go ./main.go
COPY go.mod ./go.mod

RUN go mod download && go mod verify

COPY . .
RUN go build -o=/sur/local/bin/cov main.go

ENTRYPOINT ["cov"]
