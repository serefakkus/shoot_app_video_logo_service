package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"shoot_app_video_logo_service/configs"
	"shoot_app_video_logo_service/models"
	"sync"
)

var (
	processingVideos = make(map[string]bool)
	videoMutex       = &sync.RWMutex{}
)

func HandleVideoUpload(c *gin.Context) {
	// responsedan gelen video dosyasını al
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video file not provided"})
		return
	}

	request := &models.Request{}
	if err := request.GetFromJson(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Video ID'nin daha önce işlenip işlenmediğini kontrol et
	videoMutex.RLock()
	processingVideos[request.VideoID] = true
	videoMutex.RUnlock()

	originalFilePath := filepath.Join(configs.TempPath, request.VideoID+"_"+file.Filename)
	outputFilePath := filepath.Join(configs.OutputPath, request.VideoID+".mp4")

	// Videoyu geçici olarak kaydet
	if err := c.SaveUploadedFile(file, originalFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video"})
		return
	}

	// İşlemi arka planda başlat
	go processVideoAndNotify(request.VideoID, originalFilePath, outputFilePath)

	// Kullanıcıya hemen yanıt dön
	c.JSON(http.StatusAccepted, gin.H{
		"status":   "processing",
		"video_id": request.VideoID,
	})
}
