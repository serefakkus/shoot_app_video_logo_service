package configs

import "os"

const CircularLogoPath = "./assets/logo_circular.png"
const CircularBackgroundPath = "./assets/circular_bg.png"
const (
	TempPath   = "./temp/"
	OutputPath = "./output/"
)

var WebhookURL string

func InitConfigs() {
	// logo dosyasının varlığını kontrol et
	if _, err := os.Stat(CircularLogoPath); os.IsNotExist(err) {
		panic("logo_circular.png file does not exist")
	}

	if _, err := os.Stat(CircularBackgroundPath); os.IsNotExist(err) {
		panic("circular_bg.png file does not exist")
	}

	//webhook_url ortam değişkenini oku
	WebhookURL = os.Getenv("WEBHOOK_URL")
	if WebhookURL == "" {
		panic("WEBHOOK_URL is empty")
	}
}
