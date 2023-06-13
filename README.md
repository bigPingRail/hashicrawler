# Releases Parser

This is a Go application that parses releases from a given base URL and provides a web interface to view and download the releases.

## Getting Started

### Prerequisites

- Go (version 1.20.4)

### Installation

1. Clone the repository.
1. Run `go build -o appname cmd/main.go` to build the application.

## Usage

1. Start the application by running `./appname`.
1. Access the web interface at `http://localhost:8080`.

## Endpoints

- `/hc`: Health check endpoint.
- `/`: Home page to view available release keys.

## License

This project is licensed under the [MIT License](LICENSE).
