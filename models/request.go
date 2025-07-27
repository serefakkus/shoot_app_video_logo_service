package models

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path/filepath"
)

type Request struct {
	VideoID string `json:"video_id"`
}

func (r *Request) GetVideoForm(c *gin.Context) (file *multipart.FileHeader, err error) {
	file, err = c.FormFile("video")
	if err != nil {
		err = errors.New("file not provided")
		return
	}

	r.VideoID, _ = c.GetPostForm("video_id")

	return file, r.validate()
}

func (r *Request) GetForm(c *gin.Context) (err error) {
	r.VideoID, _ = c.GetPostForm("video_id")

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
