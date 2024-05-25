package service

import (
	"context"
	"fmt"
	"goTest/internal/infrastructure/component"
	"goTest/internal/modules/messanger/storage"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	CountLastMessages = 50
	CapOfSocketsRoom  = 0
	CapOfWebSockets   = 0
)

type Messangerer interface {
	MessagerService(ctx context.Context, conn *websocket.Conn, name string, room_id int)
	CreateRoom(ctx context.Context, name string) error
	GetRoomId(ctx context.Context, name string) (int, error)
}

type Messanger struct {
	storage storage.Messangerer
	rclient *redis.Client

	logger *zap.Logger
}

func NewMessanger(storage storage.Messangerer, components *component.Components, rclient *redis.Client) *Messanger {
	return &Messanger{
		storage: storage,
		rclient: rclient,
		logger:  components.Logger,
	}
}

func (m *Messanger) WritePrevMessages(ctx context.Context, conn *websocket.Conn, room_id int) error {
	lastMess, err := m.storage.GetLastMessages(ctx, room_id, CountLastMessages)
	if err != nil {
		m.logger.Error(fmt.Sprintf("error to write in the db : %s", err.Error()))
		return err
	}

	lastMess = ReverseStringSlice(lastMess)
	msgs := strings.Join(lastMess, "\n")

	if err := conn.WriteMessage(1, []byte(msgs)); err != nil {
		m.logger.Error(fmt.Sprintf("error occurred while write messages: %s", err.Error()))
		return err
	}

	return nil
}

// how can update:
// make the pipeline to redis
// watch the avito messanger again btw...
// a lot ideas how to solve the problem to production
func (m *Messanger) MessagerService(ctx context.Context, conn *websocket.Conn, name string, room_id int) {

	//write previous message
	m.WritePrevMessages(ctx, conn, room_id)

	go m.ReadMessages(ctx, room_id, conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			m.logger.Error(fmt.Sprintf("error occurred while reading message: %s", err.Error()))
			return
		}

		message := name + " : " + string(msg)

		//write message to db
		m.storage.WriteMessage(ctx, message, room_id)

		m.rclient.Publish(strconv.Itoa(room_id), message)
	}
}

func (m *Messanger) ReadMessages(ctx context.Context, room_id int, conn *websocket.Conn) {
	pubsub := m.rclient.Subscribe(strconv.Itoa(room_id))

	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
	}
}

func (m *Messanger) CreateRoom(ctx context.Context, name string) error {
	return m.storage.CreateRoom(ctx, name)
}

func (m *Messanger) GetRoomId(ctx context.Context, name string) (int, error) {
	return m.storage.GetRoomId(ctx, name)
}
