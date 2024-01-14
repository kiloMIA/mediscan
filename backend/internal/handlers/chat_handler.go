package handlers

import (
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/sirupsen/logrus"
)

type ChatHandler struct {
    hub *Hub
    lg  *logrus.Logger
}

func NewChatHandler(hub *Hub, lg *logrus.Logger) *ChatHandler {
    return &ChatHandler{hub: hub, lg: lg}
}

func (ch *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        ch.lg.Error("Error during connection upgrade:", err)
        return
    }
    client := &Client{conn: conn, send: make(chan []byte, 256)}
    ch.hub.register <- client

    go client.writePump()
    go client.readPump(ch.hub)
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}