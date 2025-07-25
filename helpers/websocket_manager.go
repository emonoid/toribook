package helpers

import (
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
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, conn := range m.clients[channel] {
		err := conn.WriteJSON(data)
		if err != nil {
			conn.Close()
		}
	}
}
