# Start from the official Go image
FROM golang:1.19-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./services/transaction/cmd/cron/main.go

# Start a new stage from scratch
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Create a shell script to run the main executable in a loop
RUN echo '#!/bin/sh' > run.sh && \
    echo 'while true; do' >> run.sh && \
    echo '    ./main' >> run.sh && \
    echo '    sleep 1' >> run.sh && \
    echo 'done' >> run.sh && \
    chmod +x run.sh

EXPOSE 80

# Command to run the shell script
CMD ["./run.sh"]