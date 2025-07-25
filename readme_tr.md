# Dinamik Filigran Video İşleme API'si Dokümantasyonu

**Sürüm:** 1.0
**Son Güncelleme:** 25 Temmuz 2025

Önemli Not: lütfen webhookURL'yi ortam değişkenine ekle.

## Genel Bakış

Bu servis, HTTP üzerinden yüklenen video dosyalarına, her 5 saniyede bir konum değiştiren dinamik bir logo (filigran) ekler. İşlemler, sunucuyu yormamak adına asenkron olarak (arka planda) yürütülür. Bir video işleme alındığında, işlem tamamlandığında sonucun bildirileceği bir Webhook URL'ine POST isteği gönderilir. Ayrıca işlenmiş videolar sunucudan indirilebilir veya silinebilir.

docker-compose.yml dosyası ile Docker ortamında çalıştırılabilir.

## Temel URL (Base URL)

Tüm API istekleri aşağıdaki temel URL üzerinden yapılmalıdır:

`http://localhost:80` (Yerel geliştirme ortamı için)

## Kimlik Doğrulama

Mevcut sürümde herhangi bir kimlik doğrulama yöntemi (API Key, OAuth vb.) bulunmamaktadır. Endpoint'ler halka açıktır. Güvenlik için API'nin izole bir ağda çalıştırılması veya bir API Gateway arkasına konumlandırılması önerilir.

---

## Endpoint'ler

### 1. Video Yükleme ve İşlemi Başlatma

Bir video yükler ve filigran ekleme işlemini arka planda başlatır. İstek kabul edildiğinde anında yanıt döner.

- **Endpoint:** `POST /add-logo`
- **Açıklama:** Yeni bir video işleme görevi oluşturur.
- **Request Tipi:** `multipart/form-data`

**Request Body:**

| Alan Adı | Tip  | Zorunluluk | Açıklama                                                                |
| :------- | :--- | :--------- | :---------------------------------------------------------------------- |
| `video`  | File | **Gerekli** | Üzerinde işlem yapılacak video dosyası (mp4, mov, avi, vb.).           |
| `video_id` | String | **Gerekli** | İşlem için oluşturulan benzersiz kimlik. |

**Örnek cURL İsteği:**

```bash
curl -X POST \
  http://localhost:80/add-logo \
  -F "video=@/path/to/your/local_video.mp4" \
  -F "video_id=a1b2c3d4-e5f6-7890-1234-567890abcdef"
```

Başarılı Yanıt (202-Accepted):

İsteğin kabul edildiğini ve işlemin arka planda başladığını belirtir.

JSON

{
"status": "processing",
"video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}

Hatalı Yanıtlar:

400 Bad Request (video dosyası eksikse):

JSON

{ "error": "Video file not provided" }
400 Bad Request (video_id eksikse):

JSON

{ "error": "Video ID is required" }
409 Conflict (Aynı video_id ile devam eden bir işlem varsa):

JSON

{ "error": "Video is already being processed" }
2. İşlenmiş Videoyu İndirme
   İşlemi tamamlanmış bir videoyu sunucudan indirir.

Endpoint: POST /get-video

Açıklama: video_id'ye göre işlenmiş videoyu getirir.

Request Tipi: application/json

Request Body:

JSON

{
"video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}

Örnek cURL İsteği:

Bash

curl -X POST \
http://localhost:80/get-video \
-H "Content-Type: application/json" \
-d '{"video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"}' \
--output processed_video.mp4
Başarılı Yanıt (200-OK):

Content-Type: video/mp4

Body: Ham video dosyası verisi.

Hatalı Yanıtlar:

400 Bad Request (video_id eksikse veya geçersizse):

JSON

{ "error": "Video ID is required" }
404 Not Found (Video bulunamazsa veya işlem henüz tamamlanmadıysa).

3. İşlenmiş Videoyu Silme
   İşlemi tamamlanmış bir videoyu sunucudan kalıcı olarak siler.

Endpoint: DELETE /del-video

Açıklama: video_id'ye göre ilgili videoyu ve kaydını siler.

Request Tipi: application/json

Request Body:

JSON

{
"video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
Örnek cURL İsteği:

Bash

curl -X DELETE \
http://localhost:80/del-video \
-H "Content-Type: application/json" \
-d '{"video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"}'
Başarılı Yanıt (200-OK):

JSON

{ "message": "Video deleted successfully" }
Hatalı Yanıtlar:

400 Bad Request (video_id eksikse veya geçersizse).

404 Not Found (Silinecek video bulunamazsa).

500 Internal Server Error (Dosya silinirken bir hata oluşursa).

4. Servis Durum Kontrolü (Ping)
   Servisin ayakta olup olmadığını kontrol etmek için kullanılır.

Endpoint: GET /ping

Örnek cURL İsteği:

Bash

curl http://localhost:80/ping
Başarılı Yanıt (200-OK):

Boş bir body ile 200 durum kodu döner.

Webhook Bildirimi
Video işleme görevi başarıyla tamamlandığında, handlers/video_upload.go içinde tanımlı olan webhookURL adresine aşağıdaki formatta bir HTTP POST isteği gönderir.

Method: POST

Content-Type: application/json

Webhook Body İçeriği:

JSON

{
"video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
Bu bildirim, işlemin tamamlandığını ve videonun artık indirilebilir veya silinebilir durumda olduğunu belirtir.