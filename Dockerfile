## Golang
FROM golang:latest AS build-env

RUN go get -u github.com/gorilla/mux

WORKDIR  /go/src/github.com/tusupov/gousers/
ADD  . .

RUN go test -bench=. -v ./...
RUN rm -rf *_test.go

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gousers .
CMD ["./gousers"]

## Main
FROM alpine:latest

WORKDIR /root/
COPY --from=build-env /go/src/github.com/tusupov/gousers/gousers .

CMD ["./gousers"]
