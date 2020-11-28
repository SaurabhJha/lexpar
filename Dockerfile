FROM golang:1.15.5-alpine3.12

WORKDIR /go/src/lexpar

COPY . .

RUN go build

