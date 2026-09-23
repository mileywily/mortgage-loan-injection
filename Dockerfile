FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/api-bin cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-bin .

EXPOSE 8081
ENV PORT=8081

CMD ["./api-bin"]
