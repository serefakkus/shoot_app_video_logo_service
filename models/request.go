package models

import (
	"errors"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

type Request struct {
	VideoID string `json:"video_id"`
}

func (r *Request) GetFromJson(c *gin.Context) error {
	err := c.ShouldBindJSON(r)
	if err != nil {
		return errors.New("invalid request data")
	}
	return r.validate()
}

func (r *Request) validate() error {
	if r.VideoID == "" {
		return errors.New("invalid video ID")
	}

	sanitizedFilename := filepath.Base(r.VideoID)
	// Sanitize the video ID to prevent directory traversal attacks
	if sanitizedFilename == "." || sanitizedFilename == ".." {
		return errors.New("invalid video ID")
	}

	return nil
}
