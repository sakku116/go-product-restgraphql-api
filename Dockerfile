FROM golang:1.22-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o backend .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/backend .
CMD ["./backend", "--ensure-indexes", "--seed-users"]