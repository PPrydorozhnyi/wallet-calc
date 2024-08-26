### Step 1: Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Install protobuf compiler
RUN apt-get update && apt-get install -y protobuf-compiler

# Install Go protobuf plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy the Go module files and download dependencies
COPY ./src/go.* ./
RUN go mod download

# Copy the application source code and build the binary
COPY ./src ./

# Build proto
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --proto_path=. \
    --proto_path=$GOPATH/pkg/mod/github.com/googleapis/googleapis@v0.0.0-20240823220356-a67e27687c1b \
    proto/ledger_record.proto \
    proto/wallet.proto

# Build sources
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp

###
## Step 2: Runtime stage
FROM scratch

# Copy only the binary from the build stage to the final image
COPY --from=builder /app/myapp /

EXPOSE 8081

# Set the entry point for the container
ENTRYPOINT ["/myapp"]