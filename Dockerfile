# Fase de construcción
FROM golang:1.23-alpine AS builder

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
    -trimpath \
    -ldflags="-s -w -extldflags '-static'" \
    -buildvcs=false \
    -o /app/main .

# Fase de ejecución minimalista
FROM alpine:3.19

# Instalar runtime dependencies para SQLite y otras utilidades
RUN apk add --no-cache \   
    ca-certificates \
    curl \
    bash

WORKDIR /app

# Copiar el binario y los archivos RSA necesarios
COPY --from=builder /app/main .
COPY --from=builder /app/config /app/config
COPY --from=builder /app/private.rsa .
COPY --from=builder /app/public.rsa.pub .
COPY --from=builder /app/start.sh .

# Configuración de seguridad
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app && \
    chmod +x /app/start.sh

USER appuser

# Puerto de la aplicación (Railway usará la variable PORT)
EXPOSE 8229

# Comando de ejecución con mejor manejo de errores
CMD ["./start.sh"]





