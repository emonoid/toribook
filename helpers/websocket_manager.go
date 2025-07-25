package helpers

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	clients map[string][]*websocket.Conn
	lock    sync.RWMutex
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients: make(map[string][]*websocket.Conn),
	}
}

func (m *WebSocketManager) AddClient(channel string, conn *websocket.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.clients[channel] = append(m.clients[channel], conn)
}

func (m *WebSocketManager) RemoveClient(channel string, conn *websocket.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	conns := m.clients[channel]
	for i, c := range conns {
		if c == conn {
			m.clients[channel] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}

func (m *WebSocketManager) Broadcast(channel string, data interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	conns := m.clients[channel]
	activeConns := make([]*websocket.Conn, 0, len(conns))

	for _, conn := range conns {
		err := conn.WriteJSON(data)
		if err != nil {
			log.Println("WebSocket write error, removing connection:", err)
			conn.Close()
			continue // skip dead connection
		}
		activeConns = append(activeConns, conn) // keep alive ones
	}

	// Update the list with only active connections
	m.clients[channel] = activeConns
}

