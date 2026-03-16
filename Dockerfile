# Stage 1: Builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pipeline ./cmd/pipeline/

# Stage 2: Runtime
FROM python:3.11-alpine

WORKDIR /app

RUN apk add --no-cache ffmpeg && \
    pip install --no-cache-dir yt-dlp

COPY --from=builder /app/pipeline .
COPY --from=builder /app/specializations ./specializations

RUN mkdir -p output/temp

# ENV: LLM_PROVIDER, LLM_ENDPOINT, LLM_MODEL, LLM_TIMEOUT, OPENAI_API_KEY (if openai)
# SUBTITLE_LANGS (comma-separated, e.g. pt,en), OUTPUT_DIR, TEMP_DIR
ENTRYPOINT ["./pipeline"]
