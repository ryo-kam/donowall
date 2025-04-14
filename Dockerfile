FROM golang:alpine AS build

WORKDIR /build

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

COPY . .

RUN go build -o /donowall ./src

FROM alpine:latest AS run

WORKDIR /app

# Copy the application executable from the build image
COPY --from=build /donowall ./donowall

EXPOSE 8080
CMD ["./donowall"]