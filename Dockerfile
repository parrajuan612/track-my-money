# Etapa 1: Compilación
FROM golang:alpine AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache qpdf poppler-utils ca-certificates tzdata

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]