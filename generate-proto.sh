#!/bin/bash

# Generate Protocol Buffer files for Customer Service
# This script generates Go code from .proto files

echo "🔄 Generating Protocol Buffer files for Customer Service..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc could not be found. Please install Protocol Buffers compiler."
    echo "   Installation instructions: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "❌ protoc-gen-go could not be found. Installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "❌ protoc-gen-go-grpc could not be found. Installing..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Set the root directory
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROTO_DIR="$ROOT_DIR/proto"

echo "📂 Root directory: $ROOT_DIR"
echo "📂 Proto directory: $PROTO_DIR"

# Generate customer service protobuf files
echo "🔄 Generating customer.proto..."

protoc \
    --go_out=$ROOT_DIR \
    --go_opt=paths=source_relative \
    --go-grpc_out=$ROOT_DIR \
    --go-grpc_opt=paths=source_relative \
    --proto_path=$PROTO_DIR \
    customer/customer.proto

if [ $? -eq 0 ]; then
    echo "✅ customer.proto generated successfully"
else
    echo "❌ Failed to generate customer.proto"
    exit 1
fi

# List generated files
echo ""
echo "📋 Generated files:"
find $ROOT_DIR/proto -name "*.pb.go" -type f | while read file; do
    echo "   ✓ $(basename $file)"
done

echo ""
echo "🎉 Protocol Buffer generation completed successfully!"
echo ""
echo "💡 Next steps:"
echo "   1. Run 'go mod tidy' to update dependencies"
echo "   2. Build the service with 'go build ./cmd'"
echo "   3. Run tests with 'go test ./...'"
