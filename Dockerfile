# Start from the latest golang base image
FROM golang:1.21 AS builder

# Add Maintainer Info
LABEL maintainer="Tomeu Uris tomeu.uris.dev@gmail.com"

RUN apt install ca-certificates gcc

# # Install curl and unzip
# RUN apt install curl unzip git

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY pkg pkg
COPY cmd cmd
COPY api api

# Generate Swagger documentation
RUN swag init --parseDependency --parseInternal -g ./cmd/main.go -o ./docs

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Generate database migrations
RUN go run cmd/migrate/migrate.go



# Start a new stage for development
FROM debian:latest AS development

RUN apt update && apt install -y ca-certificates

# Create a new user and switch to that user
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

# Change to non-root privilege
USER appuser

WORKDIR /home/appuser/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder --chown=appuser:appgroup --chmod=555 /app/main .

# Copy the Swagger documentation from the previous stage
COPY --from=builder --chown=appuser:appgroup --chmod=444 /app/docs ./docs

# Copy database file
COPY --from=builder --chown=appuser:appgroup --chmod=744 /app/database.sqlite ./database.sqlite

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
# CMD ["ls", "-l"]
CMD ["./main"] 


# Start a new stage for production
FROM debian:latest AS production

ENV ENV=prod

RUN apt update && apt install -y ca-certificates

# Create a new user and switch to that user
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

# Change to non-root privilege
USER appuser

WORKDIR /home/appuser/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder --chown=appuser:appgroup --chmod=555 /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 