# Use the official Go image as the base image
FROM golang:1.21-alpine as base

# Set the working directory inside the container
WORKDIR /app
ENV GO_ENV=production

FROM base as build

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application for linux os
# RUN GOARCH=amd64 GOOS=linux go build -o bin/app
RUN go build -o backend

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/backend /app/backend
COPY --from=build /app/.env.prod /app/.env

EXPOSE 5002

# Set the entry point command to run the built binary
CMD ["./backend"]