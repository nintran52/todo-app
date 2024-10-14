# Step 1: Use the official Golang image as the builder
FROM golang:1.23-alpine AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy the rest of the application code
COPY . .

# Step 5: Build the Go binary
RUN go build -o main .

# Step 6: Create a lightweight final image
FROM alpine:latest

# Step 7: Copy the binary from the builder
COPY --from=builder /app/main /main

# Step 8: Expose port 8080
EXPOSE 8080

# Step 9: Run the binary
ENTRYPOINT ["/main"]
