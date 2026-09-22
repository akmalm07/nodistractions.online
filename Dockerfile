# Cloud Build uses the repository root as its Docker build context.
FROM golang:1.22-alpine AS build

WORKDIR /src

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /api ./cmd/api

FROM gcr.io/distroless/static-debian12

COPY --from=build /api /api

ENV PORT=8080
ENV ENV_MODE=true
ENV SECRET_MANAGER_SECRET=projects/933874694474/locations/us-east4/secrets/nodistractions-secrets

USER nonroot:nonroot
ENTRYPOINT ["/api"]
