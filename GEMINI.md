# Sale Description Generator

This project is a gRPC-based API designed to generate high-quality, positive sale descriptions for physical music records (e.g., vinyl LPs). It takes user notes about condition and specific grading levels (media/sleeve) and produces a single-paragraph description based on industry standards like Discogs and Goldmine.

## Tech Stack
- **Language**: Go (Golang)
- **API**: gRPC
- **Documentation**: GEMINI.md, README.md
- **Standards**: Discogs Grading Manual, Goldmine Standard

## Project Structure
- `cmd/`: Application entry points (e.g., `cmd/server/main.go`).
- `internal/`: Private code, including:
    - `internal/server/`: gRPC server implementation.
    - `internal/service/`: Core business logic for description generation.
- `api/proto/`: Protocol Buffer definitions (`api/proto/v1/`).
- `api/gen/`: Generated code from proto files.

## Core Rules
- Output descriptions must be approximately one long paragraph.
- Descriptions should maintain a positive tone while being accurate to the provided grading.
- Follow the industry standards for grading terminology.
