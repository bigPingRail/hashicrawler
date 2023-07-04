FROM golang:1.20.4-bullseye as builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY templates ./templates
COPY utils ./utils
RUN go build -ldflags "-linkmode external -extldflags -static" -o hashicrawler cmd/main.go && chmod +x hashicrawler

FROM debian:11-slim
ARG INTERNAL=false
RUN if [ "$INTERNAL" = "true" ]; then \
        echo "Acquire { https::Verify-Peer false }" >> /etc/apt/apt.conf.d/99verify-peer.conf && \
        echo "deb [trusted=yes] https://nexus.ittests.ru/repository/bullseye-proxy/ bullseye main" > /etc/apt/sources.list; \
    fi 
RUN apt update && apt -y install ca-certificates curl && apt clean
RUN mkdir -m 0755 -p /app/cache && chown 1000:1000 /app/cache
USER 1000
WORKDIR /app
COPY --from=builder /src/hashicrawler /app
COPY --from=builder /src/templates /app/templates
HEALTHCHECK --interval=30s --timeout=30s --start-period=60s --retries=3 CMD curl --fail localhost:8080/hc || exit 1
STOPSIGNAL SIGINT 
EXPOSE 8080
ENTRYPOINT ["/app/hashicrawler"]
