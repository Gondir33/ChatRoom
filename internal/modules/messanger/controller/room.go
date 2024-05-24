package controller

import (
	"context"
	"fmt"
	"goTest/internal/modules/messanger/service"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Messangerer interface {
	WebSocketHandler(w http.ResponseWriter, r *http.Request)
}

const (
	rules = `
	For Create Room type like this : "create room <name>"
	For Connection Room type like this : "conn room <name>"	
`
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
)

type Messanger struct {
	service service.Messangerer
	logger  *zap.Logger
}

func NewMessanger(service service.Messangerer, logger *zap.Logger) *Messanger {
	return &Messanger{
		service: service,
		logger:  logger,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  ReadBufferSize,
	WriteBufferSize: WriteBufferSize,
}

func (m *Messanger) WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		m.logger.Debug(fmt.Sprintf("can not to create websocket connection:%s", err.Error()))
		return
	}
	defer conn.Close()

	//print rules
	conn.WriteMessage(websocket.TextMessage, []byte(rules))
	//for futute timeouts if need will use context
	ctx := r.Context()

	name, room_id, err := m.Cmds(ctx, conn)
	if err != nil {
		m.logger.Debug(fmt.Sprintf("Error with Cmds : %s", err.Error()))
		return
	}

	//for garbage collector
	// if Use the r.Context, don't use "go", because r.Context dead, and programm dead btw...
	m.service.WebSocketService(ctx, conn, name, room_id)
}

func (m *Messanger) Cmds(ctx context.Context, conn *websocket.Conn) (string, int, error) {
	var room_id int

	for room_id == 0 {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			m.logger.Debug(fmt.Sprintf("error occurred while reading message : %s", err.Error()))
			return "", 0, err
		}
		cmd := string(msg)
		room_id, err = m.RulesRoom(ctx, cmd)
		if err != nil {
			return "", 0, err
		}
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte("Enter tha name:\n"))
	if err != nil {
		m.logger.Debug(fmt.Sprintf("error while write in message :%s", err.Error()))
		return "", 0, err
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		m.logger.Debug(fmt.Sprintf("error occurred while reading message :%s", err.Error()))
		return "", 0, err
	}

	return string(msg), room_id, err
}

func (m *Messanger) RulesRoom(ctx context.Context, text string) (int, error) {
	parts := strings.Split(text, " ")

	for i := 0; i+2 < len(parts); i += 2 {

		if parts[i] == "create" && parts[i+1] == "room" {

			err := m.service.CreateRoom(ctx, parts[i+2])
			if err != nil {
				return 0, err
			}

		} else if parts[i] == "conn" && parts[i+1] == "room" {

			id, err := m.service.GetRoomId(ctx, parts[i+2])
			if err != nil {
				return 0, err
			}
			return id, nil
		}
	}

	return 0, nil
}
