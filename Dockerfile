FROM golang:1.21.5 AS build


WORKDIR /app
RUN apt-get update && apt-get install -y git
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
#RUN apk update && apk add git


#RUN go install github.com/cosmtrek/air@latest
RUN go get -d
RUN go build -o myapp