package handlers

import (
    "github.com/gorilla/websocket"
    "log"
)

type Client struct {
    conn *websocket.Conn
    send chan []byte
}
var newline = []byte{'\n'}

func (c *Client) readPump(hub *Hub) {
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        hub.broadcast <- message
    }
}

func (c *Client) writePump() {
    defer func() {
        c.conn.Close()
    }()
    
    for message := range c.send {
        w, err := c.conn.NextWriter(websocket.TextMessage)
        if err != nil {
            return
        }
        w.Write(message)

        n := len(c.send)
        for i := 0; i < n; i++ {
            w.Write(newline)
            w.Write(<-c.send)
        }

        if err := w.Close(); err != nil {
            return
        }
    }

    c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
