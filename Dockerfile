### Step 1: Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy the Go module files and download dependencies
COPY ./src/go.* ./
RUN go mod download

# Copy the application source code and build the binary
COPY ./src ./
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp

###
## Step 2: Runtime stage
FROM scratch

# Copy only the binary from the build stage to the final image
COPY --from=builder /app/myapp /

EXPOSE 8081

# Set the entry point for the container
ENTRYPOINT ["/myapp"]