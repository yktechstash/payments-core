# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/api /app/api
COPY internal/config/prod.yaml /app/internal/config/prod.yaml
ENV CONFIG_PATH=/app/internal/config/prod.yaml
ENV CONFIG_FILE=prod.yaml
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/api"]
