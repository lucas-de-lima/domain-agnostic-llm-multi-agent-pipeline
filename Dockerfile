# Estágio 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copia dependências primeiro (cache layer)
COPY go.mod ./
# COPY go.sum ./ # Descomente se tiver go.sum
RUN go mod download

# Copia o código fonte
COPY . .

# Compila apontando para o novo main
RUN CGO_ENABLED=0 GOOS=linux go build -o distiller cmd/distiller/main.go

# Estágio 2: Runtime
FROM python:3.11-alpine

WORKDIR /app

# Instala dependências de runtime
RUN apk add --no-cache ffmpeg && \
    pip install --no-cache-dir yt-dlp

# Copia o binário
COPY --from=builder /app/distiller .

# Cria pasta de output
RUN mkdir -p output/temp

ENTRYPOINT ["./distiller"]
