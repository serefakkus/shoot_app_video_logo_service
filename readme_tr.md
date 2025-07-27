# Dinamik Filigran Video İşleme API'si Dokümantasyonu

**Sürüm:** 1.0
**Son Güncelleme:** 27 Temmuz 2025

## Genel Bakış

Bu servis, HTTP üzerinden yüklenen video dosyalarına her 5 saniyede bir konumunu değiştiren dinamik bir logo (filigran) ekler. Sunucuyu engellememek için işlemler asenkron olarak (arka planda) yürütülür. Bir video işleme tamamlandığında, tamamlandığını bildirmek için önceden tanımlanmış bir Webhook URL'sine bir POST isteği gönderilir. İşlenmiş videolar sunucudan indirilebilir veya silinebilir.

Sağlanan `docker-compose.yml` dosyası kullanılarak bir Docker ortamında çalıştırılabilir.

**Dikkat:** Lütfen `WEBHOOK_URL`'yi bir ortam değişkeni olarak ayarladığınızdan emin olun.

## Temel URL

Tüm API istekleri aşağıdaki temel URL'ye yapılmalıdır:

`http://localhost:8080` (Yerel geliştirme ortamı için)

## Kimlik Doğrulama

Mevcut sürümde herhangi bir kimlik doğrulama yöntemi bulunmamaktadır. Endpoint'ler halka açıktır. Güvenlik için API'nin izole bir ağda çalıştırılması veya bir API Gateway arkasına yerleştirilmesi önerilir.

---

## Endpoint'ler

### 1. Video Yükleme ve İşlemi Başlatma

Bir video yükler ve filigran ekleme işlemini arka planda başlatır. İsteği kabul ettiğinde anında bir yanıt döndürür.

- **Endpoint:** `POST /add-logo`
- **Açıklama:** Yeni bir video işleme görevi oluşturur.
- **Request Tipi:** `multipart/form-data`

**Request Body:**

| Alan Adı | Tip    | Zorunluluk | Açıklama                                                    |
| :------- | :----- | :--------- | :---------------------------------------------------------- |
| `video`  | Dosya  | **Gerekli** | İşlenecek video dosyası (mp4, mov, avi, vb.).               |
| `video_id` | String | **Gerekli** | İşlem için benzersiz bir tanımlayıcı.                        |

**Örnek cURL İsteği:**

```bash
curl -X POST \
  http://localhost:8080/add-logo \
  -F "video=@/path/to/your/local_video.mp4" \
  -F "video_id=a1b2c3d4-e5f6-7890-1234-567890abcdef"
```

**Başarılı Yanıt (202-Accepted):**

İsteğin kabul edildiğini ve işlemin arka planda başladığını belirtir.

```json
{
  "status": "processing",
  "video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
```

**Hata Yanıtları:**

- **400 Bad Request:** Video dosyası veya `video_id` eksikse.
- **500 Internal Server Error:** Videoyu kaydederken bir hata oluşursa.

### 2. İşlenmiş Videoyu İndirme

İşlemi tamamlanmış bir videoyu sunucudan indirir.

- **Endpoint:** `POST /get-video`
- **Açıklama:** `video_id`'sine göre işlenmiş videoyu getirir.
- **Request Tipi:** `multipart/form-data`

**Request Body:**

| Alan Adı | Tip    | Zorunluluk | Açıklama                        |
| :------- | :----- | :--------- | :------------------------------ |
| `video_id` | String | **Gerekli** | Videonun benzersiz tanımlayıcısı. |

**Örnek cURL İsteği:**

```bash
curl -X POST \
  http://localhost:8080/get-video \
  -F "video_id=a1b2c3d4-e5f6-7890-1234-567890abcdef" \
  --output islenmis_video.mp4
```

**Başarılı Yanıt (200-OK):**

- **Content-Type:** `video/mp4`
- **Body:** Ham video dosyası verileri.

**Hata Yanıtları:**

- **400 Bad Request:** `video_id` eksik veya geçersizse.
- **404 Not Found:** Video bulunamazsa veya henüz işlenmemişse.

### 3. İşlenmiş Videoyu Silme

İşlenmiş bir videoyu sunucudan kalıcı olarak siler.

- **Endpoint:** `DELETE /del-video`
- **Açıklama:** Belirtilen videoyu `video_id`'sine göre siler.
- **Request Tipi:** `multipart/form-data`

**Request Body:**

| Alan Adı | Tip    | Zorunluluk | Açıklama                        |
| :------- | :----- | :--------- | :------------------------------ |
| `video_id` | String | **Gerekli** | Videonun benzersiz tanımlayıcısı. |

**Örnek cURL İsteği:**

```bash
curl -X DELETE \
  http://localhost:8080/del-video \
  -F "video_id=a1b2c3d4-e5f6-7890-1234-567890abcdef"
```

**Başarılı Yanıt (200-OK):**

```json
{ "message": "Video deleted successfully" }
```

**Hata Yanıtları:**

- **400 Bad Request:** `video_id` eksik veya geçersizse.
- **404 Not Found:** Silinecek video bulunamazsa.
- **500 Internal Server Error:** Dosya silme sırasında bir hata oluşursa.

### 4. Servis Durum Kontrolü (Ping)

Servisin çalışıp çalışmadığını kontrol etmek için kullanılır.

- **Endpoint:** `GET /ping`

**Örnek cURL İsteği:**

```bash
curl http://localhost:8080/ping
```

**Başarılı Yanıt (200-OK):**

Boş bir gövde ile 200 durum kodu döndürür.

---

## Webhook Bildirimi

Bir video işleme görevi başarıyla tamamlandığında, servis ortam değişkenlerinde tanımlanan `WEBHOOK_URL`'ye bir HTTP POST isteği gönderir.

- **Method:** `POST`
- **Content-Type:** `application/json`

**Webhook Body İçeriği:**

```json
{
  "video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
```

Bu bildirim, işlemin tamamlandığını ve videonun artık indirilebilir veya silinebilir durumda olduğunu belirtir.
