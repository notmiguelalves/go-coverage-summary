FROM golang:1.22.5-alpine

WORKDIR /usr/src/cov

COPY main.go ./main.go
COPY go.mod ./go.mod

RUN go mod download && go mod verify

COPY . .
RUN go build -o=/usr/local/bin/cov main.go

ENTRYPOINT ["cov"]
