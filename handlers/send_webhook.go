package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/serefakkus/shot_app_video_logo_service/configs"
	"io"
	"net/http"
)

func sendWebhook(videoID string) error {
	payload := map[string]string{
		"video_id": videoID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(configs.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook received non-2xx status code: %d", resp.StatusCode)
	}

	return nil
}
