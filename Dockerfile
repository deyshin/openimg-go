# # Use the official Golang image as the build stage
# FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/wolfi-base as build

# RUN apk update && apk add build-base git openssh go-1.22

# # Set the working directory inside the container
# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the rest of the application code
# COPY . .

# ARG TARGETOS
# ARG TARGETARCH
# RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o openimg-go

# FROM cgr.dev/chainguard/static:latest

# COPY --from=build /app/openimg-go .

# # Command to run the application
# CMD ["/app/openimg-go", "-addr", "0.0.0.0:8080"]
# EXPOSE 8080

# Use Go 1.23 bookworm as base image
FROM golang:1.23-bookworm AS base

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
RUN go build -o go-blog

# Document the port that may need to be published
EXPOSE 8000

# Start the application
CMD ["/build/go-blog"]