package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"shoot_app_video_logo_service/configs"
	"shoot_app_video_logo_service/models"
)

func HandleGetVideo(c *gin.Context) {
	// 1. URL'den dosya adını al
	// Video ID'yi parametrelerden al ve boş olup olmadığını kontrol et

	request := &models.Request{}

	if err := request.GetFromJson(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, exists := processingVideos[request.VideoID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	outputFilePath := filepath.Join(configs.OutputPath, request.VideoID+".mp4")

	// 2. GÜVENLİK: Directory Traversal saldırılarını önlemek için dosya adını temizle.
	// Bu işlem, ".." veya "/" gibi karakterleri kaldırarak sadece dosya adının kendisini bırakır.
	sanitizedFilename := filepath.Base(outputFilePath)
	if sanitizedFilename == "." || sanitizedFilename == ".." {
		c.String(http.StatusBadRequest, "Invalid filename")
		return
	}

	// 3. Dosyanın sunucudaki tam yolunu oluştur
	videoPath := filepath.Join(configs.OutputPath, sanitizedFilename)

	// Bu fonksiyon, dosyanın var olup olmadığını kontrol eder.
	// Eğer dosya yoksa, otomatik olarak 404 Not Found hatası döner.
	// Varsa, Content-Type başlığını ayarlar ve dosyayı stream eder.
	c.File(videoPath)
}
