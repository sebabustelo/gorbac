# Fase de construcción
FROM golang:1.21-alpine AS builder

# Configuración del entorno
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOPROXY=https://proxy.golang.org,direct

# Instalar dependencias de compilación (incluyendo las necesarias para SQLite)
RUN apk add --no-cache \
    gcc \
    musl-dev 

WORKDIR /app

# 1. Copiar solo los archivos de dependencias primero para mejor caché
COPY go.mod go.sum ./

# 2. Descargar dependencias (se cachea esta capa)
RUN go mod download

# 3. Copiar el resto del código fuente
COPY . .

# 4. Construir la aplicación con optimizaciones
RUN go build \
    -ldflags="-s -w" \    
    -o main .

# Fase de ejecución minimalista
FROM alpine:3.19

# Instalar runtime dependencies para SQLite
RUN apk add --no-cache \   
    ca-certificates

WORKDIR /app

# Copiar el binario y los archivos RSA necesarios
COPY --from=builder /app/main .
COPY --from=builder /app/config /app/config
COPY --from=builder /app/private.rsa .
COPY --from=builder /app/public.rsa.pub .

# Configuración de seguridad
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser

# Puerto de la aplicación
EXPOSE 8229

# Comando de ejecución
CMD ["./main"]





