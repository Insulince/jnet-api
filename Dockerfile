FROM golang:alpine as build
MAINTAINER justinreusnow@gmail.com
WORKDIR /go/src/github.com/Insulince/jnet-api
RUN apk add --update git
COPY . .
RUN go get ./cmd/api
RUN go build -o ./jnet-api ./cmd/api;

FROM alpine:latest as deploy
WORKDIR /root
COPY --from=build /go/src/github.com/Insulince/jnet-api/jnet-api ./jnet-api
EXPOSE 8080
CMD ["./jnet-api"]
