FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/canchas-api ./cmd/canchas-api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/canchas-api /app/canchas-api
EXPOSE 8080
ENTRYPOINT ["/app/canchas-api"]
