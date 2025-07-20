FROM golang:1.23 as builder
RUN rm -rf /app

RUN mkdir /app

WORKDIR /app
COPY ./account-cmd /app/account-cmd
COPY ./shared /app/shared

WORKDIR /app/account-cmd
# Compilar la aplicación
RUN CGO_ENABLED=0 go build -o myapp ./cmd

# Etapa 2: Crear una imagen mínima
FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ=America/Bogota
WORKDIR /app

# Copiar solo los archivos necesarios desde la etapa de construcción
COPY --from=builder /app/account-cmd/myapp /app/myapp

# Comando para ejecutar la aplicación
CMD ["./myapp"]