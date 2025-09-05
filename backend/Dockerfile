# Étape de build
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Copie et téléchargement des modules
COPY go.mod go.sum ./
RUN go mod download

# Copie du code source
COPY . .

# Build avec CGO activé (par défaut sur Debian)
RUN go build -o server .

# Image finale minimaliste basée sur Debian
FROM debian:bullseye-slim

WORKDIR /app

# Installer libsqlite3 nécessaire à go-sqlite3
RUN apt-get update && apt-get install -y libsqlite3-0 && rm -rf /var/lib/apt/lists/*

# Copier le binaire depuis l'étape builder
COPY --from=builder /app/server .

# Copier la DB et les migrations dans l'image
COPY db ./db

EXPOSE 8080

CMD ["./server"]
