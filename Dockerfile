FROM golang:alpine
ENV GO111MODULE=on
#ENV http_proxy="http://www-proxy.mrec.ar:8080"
#ENV https_proxy="http://www-proxy.mrec.ar:8080"
ENV GO_ENV=produccion

# Instala dependencias de compilación para SQLite
RUN apk add --no-cache gcc musl-dev

# Move to working directory /build
WORKDIR /opt/go

#RUN go mod init api-rbac

# Copy the code into the container
COPY . .
RUN go mod tidy
RUN go build -o main .

# Export necessary port
EXPOSE 8229

# Command to run when starting the container
CMD ["./main"]