# Go'nun Alpine tabanlı (küçük boyutlu) versiyonunu kullanarak derleme aşamasını başlat
FROM golang:1.24-alpine AS builder

# Uygulama için çalışma dizini oluştur
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Tüm proje dosyalarını kopyala
COPY . .

# Uygulamayı derle. CGO_ENABLED=0 statik bir binary oluşturur.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .


FROM alpine:latest

# Derleme aşamasından sadece çalıştırılabilir dosyayı kopyala
COPY --from=builder /main /main

# Logo dosyasını kopyala.
COPY logo.png /logo.png

# Çıktı ve geçici dosyalar için klasörleri oluştur
RUN mkdir -p /output /temp

# *** ORTAM DEĞİŞKENİNİ (ENVIRONMENT VARIABLE) AYARLA
ENV WEBHOOK_URL="https://example.com/webhook"

EXPOSE 80

CMD ["/main"]