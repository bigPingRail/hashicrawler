FROM golang:1.20.4-bullseye as builer
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
RUN go build -ldflags "-linkmode external -extldflags -static" -o hashiparser cmd/parser.go && chmod +x hashiparser

FROM debian:11-slim
RUN apt update && apt -y install ca-certificates && apt clean
USER 1000
WORKDIR /app
COPY --from=builer /src/hashiparser /app
STOPSIGNAL SIGTERM
EXPOSE 8080
ENTRYPOINT ["/app/hashiparser"]
