# syntax=docker/dockerfile:1

# ---------- Etapa 1: build ----------
# Imagen con todo el toolchain de Go (compilador, módulos). Solo se usa para compilar.
FROM golang:1.24-alpine AS build

WORKDIR /src

# 1) Copiar primero solo go.mod y go.sum y bajar dependencias.
#    Así Docker cachea esta capa y no vuelve a descargar todo
#    cada vez que cambia una línea de código.
COPY go.mod go.sum ./
RUN go mod download

# 2) Copiar el resto del código y compilar.
#    CGO_ENABLED=0 -> binario estático (no depende de librerías de C del sistema).
#    -trimpath y -ldflags="-s -w" -> binario más chico y sin rutas de la máquina de build.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/catalog .

# ---------- Etapa 2: runtime ----------
# Imagen mínima: sin shell, sin gestor de paquetes, usuario no-root (UID 65532).
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Del build solo nos llevamos el binario.
COPY --from=build /out/catalog /app/catalog

# No correr nunca como root.
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/catalog"]