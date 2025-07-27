package configs

import "os"

const LogoPath = "./logo.png"
const (
	TempPath   = "./temp/"
	OutputPath = "./output/"
)

var WebhookURL string

func InitConfigs() {
	// logo dosyasının varlığını kontrol et
	if _, err := os.Stat(LogoPath); os.IsNotExist(err) {
		panic("logo.png file does not exist")
	}

	//webhook_url ortam değişkenini oku
	WebhookURL = os.Getenv("WEBHOOK_URL")
	if WebhookURL == "" {
		panic("WEBHOOK_URL is empty")
	}
}
