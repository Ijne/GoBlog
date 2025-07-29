package websockets

import (
	"log"
	"sync"

	"github.com/Ijne/core-api_app/internal/tools"
	"golang.org/x/net/websocket"
)

var WS_server = NewWSServer()

type Client struct {
	Conn *websocket.Conn
}

type WSServer struct {
	clients map[int32]*Client
	mu      sync.Mutex
}

func NewWSServer() *WSServer {
	return &WSServer{
		clients: make(map[int32]*Client),
	}
}

func (s *WSServer) HandleWS(ws *websocket.Conn) {
	defer ws.Close()

	user, err := tools.GetUserClaimsFromCookie(ws.Request())
	if err != nil {
		log.Fatal("Can't get user seesion to make ws connection")
	}

	s.mu.Lock()
	s.clients[user.ID] = &Client{ws}
	s.mu.Unlock()

	log.Printf("Новое подключение: %s %s", ws.RemoteAddr(), ws.Request().UserAgent())

	defer func() {
		s.mu.Lock()
		delete(s.clients, user.ID)
		s.mu.Unlock()
		ws.Close()
	}()

	select {}
}

func (s *WSServer) Broadcast(msg string, to int32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, exists := s.clients[to]; exists {
		if err := websocket.Message.Send(client.Conn, msg); err != nil {
			log.Printf("Ошибка отправки: %v", err)
		}
	}
}
