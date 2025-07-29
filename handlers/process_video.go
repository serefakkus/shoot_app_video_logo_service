package handlers

import (
	"bytes"
	"github.com/serefakkus/shot_app_video_logo_service/configs"
	"log"
	"os"
	"os/exec"
)

func processVideoAndNotify(videoID, originalFile, outputFile string) {
	//sol ust kısımda sabit logo kare arkaplan
	//filter := "[1:v]scale=150:-1[logo];color=c=white@0.7:s=170x95[bg];[0:v][bg]overlay=50:50:shortest=1[tmp];[tmp][logo]overlay=60:60"

	// her 5 saniyede bir logoyu farklı bir konumda göstermek için FFmpeg filtresi kare arkaplan. (sol ust, sag alt, sol alt. ***sag üstte tarih kısımı oladuğu için eklenmedi)
	//filter := "[1:v]scale=150:-1[logo];color=c=white@0.7:s=170x95[bg];[bg][logo]overlay=(W-w)/2:(H-h)/2[watermark];[0:v][watermark]overlay=x='if(between(mod(t,15),5,10),W-w-50,50)':y='if(lt(mod(t,15),5),50,H-h-50)':shortest=1"

	//yüksek kaliteli çıktı. Optimize edilmemiş
	//cmd := exec.Command("ffmpeg",
	//	"-i", originalFile,
	//	"-i", configs.LogoPath,
	//	"-filter_complex", filter,
	//	"-codec:a", "copy",
	//	outputFile,
	//)

	// Optimize edilmiş komut.
	//kalite ayari icin -crf = guncelle (daha buyuk sayi daha az kalite , en fazla 51 -en fazla sıkıştırma- ,en az 0 -kayipsiz- olarak ayarlanir. video streaming servisleri icin 28 ideal degerdir.)
	//dairesel arkaplan

	filter := "[1:v]scale=120:-1[logo];" +
		"[2:v]scale=160:-1,colorchannelmixer=aa=1[bg];" +
		"[bg][logo]overlay=(W-w)/2:(H-h)/2[watermark];" +
		"[0:v][watermark]overlay=x='if(between(mod(t,15),5,10),W-w-50,50)':y='if(lt(mod(t,15),5),50,H-h-50)'"

	cmd := exec.Command("ffmpeg",
		"-i", originalFile,
		"-i", configs.CircularLogoPath,
		"-i", configs.CircularBackgroundPath,
		"-filter_complex", filter,
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "28",
		"-c:a", "copy",
		"-shortest",
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

	// Webhook gönder
	err = sendWebhook(videoID)
	if err != nil {
		log.Printf("[%s] Failed to send webhook: %s", videoID, err.Error())
	}

	// Geçici dosyaları temizle
	err = os.Remove(originalFile)
	if err != nil {
		return
	}
}
