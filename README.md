# Releases Parser

This is a Go application that parses releases from a given base URL and provides a web interface to view and download the releases.

## Getting Started

### Prerequisites

- Go (version 1.20.4)

### Installation

1. Clone the repository.
1. Run `go build -o appname cmd/main.go` to build the application.

## Usage
### As binary
1. Start the application by running `GIN_MODE=release ./appname -p 8080`.
1. To enable local caching, add the `-c` flag: `GIN_MODE=release ./appname -p 8080 -c`
1. Access the web interface at `http://localhost:8080`.

### As Docker container
1. Build the container: `docker build -t appname:test -f build/Dockerfile .`
1. Run: `docker run -p 8080:8080 --rm appname:test -c`
1. To save cache state, you can attach a volume to the container: `docker run -p 8080:8080 -d --rm -v appname_cache:/app/cache appname:test -c`

## Endpoints

- `/hc`: Health check endpoint.
- `/`: Home page to view available releases.

## License

This project is licensed under the [MIT License](LICENSE).
