# Sale Description Generator

The **Sale Description Generator** is a gRPC-powered service designed to create professional, positive, and accurate sale descriptions for physical music records (such as vinyl LPs). Leveraging industry standards, it transforms grading levels and user notes into a polished paragraph suitable for platforms like Discogs or eBay.

## Tech Stack
- **Language**: Go (Golang)
- **API Protocol**: gRPC
- **Serialization**: Protocol Buffers (v3)

## Project Structure
- `cmd/server/`: Main application entry point for the gRPC server.
- `internal/server/`: gRPC server implementation and interface.
- `internal/service/`: Core business logic for description generation.
- `api/proto/v1/`: Protocol Buffer definitions.
- `api/gen/v1/`: Generated gRPC and Protobuf code.

## Getting Started

### Prerequisites
- Go 1.21+
- [buf](https://buf.build/docs/installation) (for protobuf generation, if modifying)

### Running the Server
To start the gRPC server on the default port (`:50051`):
```bash
go run cmd/server/main.go
```

The server includes gRPC reflection, allowing for easy testing with tools like `grpcurl`.

## API Definition
The service contract is defined in [sale_description.proto](api/proto/v1/sale_description.proto).

### Core Method: `GenerateDescription`
Takes a request with:
- `record_title`
- `artist`
- `media_condition` (Grading enum)
- `sleeve_condition` (Grading enum)
- `user_notes`

And returns a single string `description`.

## Grading Standards
The generator is built upon industry-standard grading manuals to ensure accuracy and consistency:
- [Discogs Grading Manual](https://www.discogs.com/selling/resources/how-to-grade-items)
- [Goldmine Standard](https://www.goldminemag.com/collector-resources/record-grading-101/)

The output is carefully crafted to remain positive and professional while faithfully representing the condition provided.

