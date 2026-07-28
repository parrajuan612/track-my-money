# Etapa 1: Compilación del binario de Go (usando la versión más reciente de Alpine Go)
FROM golang:alpine AS builder

# Permitir que Go descargue la versión exacta si go.mod lo requiere
ENV GOTOOLCHAIN=auto

WORKDIR /app

# Copiar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente y compilar
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Etapa 2: Imagen final liviana con las herramientas instaladas
FROM alpine:latest

WORKDIR /app

# Instalar qpdf y poppler-utils (que contiene pdftotext)
RUN apk add --no-cache qpdf poppler-utils ca-certificates

# Copiar el ejecutable compilado desde la etapa anterior
COPY --from=builder /app/main .

# Exponer el puerto del backend
EXPOSE 8080

# Ejecutar el servidor
CMD ["./main"]