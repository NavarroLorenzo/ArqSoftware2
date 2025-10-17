FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/reservas-api ./cmd/reservas-api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/reservas-api /app/reservas-api
EXPOSE 8080
ENTRYPOINT ["/app/reservas-api"]
