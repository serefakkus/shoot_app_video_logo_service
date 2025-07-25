package handlers

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"shoot_app_video_logo_service/configs"
)

func processVideoAndNotify(videoID, originalFile, outputFile string) {
	log.Printf("[%s] Starting video processing...", videoID)

	// Opaklığı %50 olarak ayarlanmış logo için FFmpeg komutu
	opacity := "0.5" // Bu değeri dinamik olarak da ayarlayabilirsiniz.

	// her 5 saniyede bir logoyu farklı bir konumda göstermek için FFmpeg filter'ı
	filter := fmt.Sprintf("[1:v]colorchannelmixer=aa=%s[logo_transparent];[0:v][logo_transparent]overlay=x='if(or(between(mod(t,20),0,5),between(mod(t,20),15,20)),10,W-w-10)':y='if(or(between(mod(t,20),0,5),between(mod(t,20),10,15)),10,H-h-10)'", opacity)

	cmd := exec.Command("ffmpeg",
		"-i", originalFile,
		"-i", configs.LogoPath, // Logo dosyasının yolu
		"-filter_complex", filter,
		"-codec:a", "copy",
		outputFile,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("[%s] FFmpeg error: %s\n%s", videoID, err.Error(), stderr.String())
		// Dosyaları temizle
		err := os.Remove(originalFile)
		if err != nil {
			log.Printf("[%s] Failed to remove original file: %s", videoID, err.Error())
		}
		return
	}

	log.Printf("[%s] Video processing finished. Output: %s", videoID, outputFile)

	// Webhook gönder
	err = sendWebhook(videoID)
	if err != nil {
		log.Printf("[%s] Failed to send webhook: %s", videoID, err.Error())
	} else {
		log.Printf("[%s] Webhook sent successfully.", videoID)
	}

	// Geçici dosyaları temizle
	err = os.Remove(originalFile)
	if err != nil {
		return
	}
}
