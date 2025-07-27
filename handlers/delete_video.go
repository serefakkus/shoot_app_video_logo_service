package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/serefakkus/shot_app_video_logo_service/configs"
	"github.com/serefakkus/shot_app_video_logo_service/models"
	"os"
	"path/filepath"
)

func HandleDeleteVideo(c *gin.Context) {
	// Extract video ID from the request
	request := models.Request{}
	if err := request.GetForm(c); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	videoMutex.RLock()
	_, exists := processingVideos[request.VideoID]
	videoMutex.RUnlock()

	if !exists {
		c.JSON(404, gin.H{"error": "Video not found"})
		return
	}

	// Check if the video file exists
	outputFilePath := filepath.Join(configs.OutputPath, request.VideoID+".mp4")
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "Video not found"})
		return
	}

	err := os.Remove(configs.OutputPath + request.VideoID + ".mp4")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete video"})
		return
	}

	c.JSON(200, gin.H{"message": "Video deleted successfully"})

}
