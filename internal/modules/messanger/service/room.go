package service

import (
	"context"
	"fmt"
	"goTest/internal/infrastructure/component"
	"goTest/internal/modules/messanger/storage"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	CountLastMessages = 50
	CapOfSocketsRoom  = 0
	CapOfWebSockets   = 0
)

type Messangerer interface {
	WebSocketService(ctx context.Context, conn *websocket.Conn, name string, room_id int)
	CreateRoom(ctx context.Context, name string) error
	GetRoomId(ctx context.Context, name string) (int, error)
}

type Messanger struct {
	storage     storage.Messangerer
	roomStorage RoomMap
	logger      *zap.Logger
}

func NewMessanger(storage storage.Messangerer, components *component.Components) *Messanger {
	return &Messanger{
		storage: storage,
		roomStorage: RoomMap{
			rooms: make(map[*websocket.Conn]int),
		},
		logger: components.Logger,
	}
}

func (m *Messanger) WebSocketService(ctx context.Context, conn *websocket.Conn, name string, room_id int) {
	defer func() {
		m.roomStorage.delete(conn)
	}()

	m.roomStorage.set(conn, room_id)

	m.WritePrevMessages(ctx, conn, room_id)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			m.logger.Error(fmt.Sprintf("error occurred while reading message: %s", err.Error()))
			return
		}

		message := append([]byte(name+": "), msg...)
		m.storage.WriteMessage(ctx, string(message), room_id)

		err = m.WriteSocketsFromRoom(ctx, message, room_id)
		if err != nil {
			return
		}
	}
}

func (m *Messanger) WritePrevMessages(ctx context.Context, conn *websocket.Conn, room_id int) error {
	lastMess, err := m.storage.GetLastMessages(ctx, room_id, CountLastMessages)
	if err != nil {
		m.logger.Error(fmt.Sprintf("error to write in the db : %s", err.Error()))
		return err
	}

	if err := conn.WriteMessage(1, lastMess); err != nil {
		m.logger.Error(fmt.Sprintf("error occurred while write messages: %s", err.Error()))
		return err
	}

	return nil
}

func (m *Messanger) WriteSocketsFromRoom(ctx context.Context, message []byte, room_id int) error {

	sockets := m.roomStorage.findAllConnOfSuchRoom(room_id)
	for _, socket := range sockets {
		err := socket.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error occurred while reading message: %s", err.Error()))
			return err
		}
	}
	return nil
}

func (m *Messanger) CreateRoom(ctx context.Context, name string) error {
	return m.storage.CreateRoom(ctx, name)
}

func (m *Messanger) GetRoomId(ctx context.Context, name string) (int, error) {
	return m.storage.GetRoomId(ctx, name)
}
