FROM golang:1.21.5 AS build

WORKDIR /app
RUN apt-get update && apt-get install -y git

# Copy go mod and sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

COPY . /app


RUN CGO_ENABLED=0 go build -o myapp



FROM alpine:latest AS runtime

COPY --from=build /app /app

WORKDIR /app

# Set the entrypoint for the container to run the binary
ENTRYPOINT ["/app/myapp"]
CMD ["-b 0.0.0.0"]