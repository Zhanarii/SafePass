package handlers

import (
	"log"
	"net/http"
	"time"

	"safepass-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleCameraWebSocket(detector *services.Detector, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cameraID := c.Param("id")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		log.Printf("Camera %s WebSocket connected", cameraID)

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				frame := struct {
					CameraID  string    `json:"camera_id"`
					Timestamp time.Time `json:"timestamp"`
					Frame     string    `json:"frame"`
				}{
					CameraID:  cameraID,
					Timestamp: time.Now(),
					Frame:     "",
				}

				if err := conn.WriteJSON(frame); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
			}
		}
	}
}
