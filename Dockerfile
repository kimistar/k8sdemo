FROM golang:latest as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM debian:latest
WORKDIR /app
COPY --from=builder /build/app .
ENTRYPOINT ["/app/app"]
