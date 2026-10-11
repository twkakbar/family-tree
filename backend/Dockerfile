# ============================================================
# Backend – multi-stage build
#   dev:     hot reload via air (compose dev target)
#   build:   compile static binaries
#   runtime: minimal non-root production image
# ============================================================

# --- Development (hot reload) --------------------------------
FROM golang:1.25-alpine AS dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# --- Build ----------------------------------------------------
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=docker
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
      -ldflags "-s -w -X main.version=${VERSION}" \
      -o /out/api ./cmd/api \
    && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" \
      -o /out/migrate ./cmd/migrate

# --- Runtime ---------------------------------------------------
FROM alpine:3.20 AS runtime
RUN apk add --no-cache ca-certificates tzdata su-exec \
    && addgroup -S app && adduser -S -G app app \
    && mkdir -p /app/uploads && chown -R app:app /app
WORKDIR /app
COPY --from=build /out/api /out/migrate ./
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/app/api"]
