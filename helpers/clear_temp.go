package helpers

import (
	"os"
	"shoot_app_video_logo_service/configs"
)

// ClearTempFiles temp klasöründeki tüm geçici dosyaları siler uygulama başlatıldığında çağrılır
func ClearTempFiles() {
	_ = os.RemoveAll(configs.TempPath)
	_ = os.MkdirAll(configs.TempPath, 0755)
}
