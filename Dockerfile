# Go'nun resmi imajı Debian temelli
FROM golang:1.24-bookworm AS builder

# Uygulama için bir çalışma dizini oluştur
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

#FFmpeg'i kur
RUN apt-get update && apt-get install -y ffmpeg

COPY . .

# CGO_ENABLED=0 statik bir binary oluşturur.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server .

# Çok daha küçük ve minimal bir imaj. Bu, güvenlik ve boyut için önemlidir.
FROM debian:bookworm-slim

# Güvenlik güncelleştirmelerini yap ve SADECE FFmpeg'i kur.
RUN apt-get update && apt-get install -y ffmpeg && rm -rf /var/lib/apt/lists/*

# Çalışma dizinini ayarla
WORKDIR /app

COPY --from=builder /server /app/server

COPY logo.png /app/logo.png

# *** ORTAM DEĞİŞKENİNİ (ENVIRONMENT VARIABLE) AYARLA
ENV WEBHOOK_URL="localhost/weebhook"

# Uygulamanın video dosyalarını yazacağı klasörler
RUN mkdir -p /app/temp /app/output

EXPOSE 8080

CMD ["/app/server"]
