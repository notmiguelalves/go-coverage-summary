FROM golang:1.22.5-alpine

COPY main.go /main.go

RUN go build -o=/cov /main.go

ENTRYPOINT ["/cov"]
