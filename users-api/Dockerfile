# --- build stage ---
FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/users-api ./cmd/users-api

# --- runtime stage ---
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/users-api /app/users-api
EXPOSE 8080
ENTRYPOINT ["/app/users-api"]
