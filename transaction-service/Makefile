# Variables
PROTOC := protoc
PROTO_DIR := proto
OUTPUT_DIR := pb

# Ensure the output directory exists
.PHONY: all generate clean

all: generate

generate: $(PROTO_DIR)/*.proto
	@echo "Generating Go files from proto files..."
	@if not exist "$(OUTPUT_DIR)" mkdir "$(OUTPUT_DIR)"
	$(PROTOC) -I $(PROTO_DIR) \
		--go_out=$(OUTPUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUTPUT_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUTPUT_DIR) --grpc-gateway_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

clean:
	@echo "Cleaning generated files..."
	@if exist "$(OUTPUT_DIR)" rmdir /s /q "$(OUTPUT_DIR)"
