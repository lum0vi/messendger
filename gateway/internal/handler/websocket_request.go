package handler

import (
	"context"
	"encoding/json"
	"gateway/internal/models"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (h *Handler) HandleWebsocketRequest(msg *models.MessageDelivery) error {
	socketID, err := h.redisCon.Get(context.Background(), "socket:"+msg.UserID).Result()
	if err != nil {
		logrus.Errorf("Error connecting to redis: %v", err)
		return err
	}
	logrus.Infof("ebsocket Send for SocketID: %v", socketID)
	conn, ok := h.connections[socketID]
	if !ok {
		logrus.Errorf("Error connecting to redis for socket %v", socketID)
		return err
	}
	if conn == nil {
		logrus.Errorf("Error connecting to redis for socket: %v", socketID)
		return err
	}
	out := struct {
		Type string                  `json:"type"`
		Data *models.MessageDelivery `json:"data"`
	}{
		Type: "new_message",
		Data: msg,
	}
	jsMes, err := json.Marshal(out)
	if err != nil {
		logrus.Errorf("Error marshalling json: %v", err)
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, jsMes); err != nil {
		logrus.Errorf("Error sending ping message: %v", err)
		return err
	}

	logrus.Info("Successfully send mess")
	return nil
}
