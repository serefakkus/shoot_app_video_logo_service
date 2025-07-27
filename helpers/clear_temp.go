package helpers

import (
	"github.com/serefakkus/shot_app_video_logo_service/configs"
	"os"
)

// ClearTempFiles temp klasöründeki tüm geçici dosyaları siler uygulama başlatıldığında çağrılır
func ClearTempFiles() {
	_ = os.RemoveAll(configs.TempPath)
	_ = os.MkdirAll(configs.TempPath, 0755)
}
