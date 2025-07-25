package main

import (
	"github.com/gin-gonic/gin"
	"shoot_app_video_logo_service/configs"
	"shoot_app_video_logo_service/handlers"
	"shoot_app_video_logo_service/helpers"
)

func main() {
	helpers.ClearTempFiles()

	//konfigrasyonu ayarla
	configs.InitConfigs()

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	err := r.SetTrustedProxies(nil)
	if err != nil {
		panic(err.Error())
	}

	r.GET("/ping", handlers.Ping)

	r.POST("/add-logo", handlers.HandleVideoUpload)

	r.POST("/get-video", handlers.HandleGetVideo)

	r.DELETE("/del-video", handlers.HandleDeleteVideo)

	err = r.Run("0.0.0.0:80")
	if err != nil {
		panic(err.Error())
	}
}
