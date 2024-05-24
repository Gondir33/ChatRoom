package storage

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Messangerer interface {
	WriteMessage(ctx context.Context, message string, room int) error
	GetLastMessages(ctx context.Context, room_id int, limit int) ([]byte, error)
	CreateRoom(ctx context.Context, name string) error
	GetRoomId(ctx context.Context, name string) (int, error)
}

type Messanger struct {
	pool    *pgxpool.Pool
	rclient *redis.Client
}

func NewMessanger(pool *pgxpool.Pool, rclient *redis.Client) *Messanger {
	return &Messanger{
		pool:    pool,
		rclient: rclient,
	}
}

func (m *Messanger) WriteMessage(ctx context.Context, message string, room int) error {
	sql := `
		INSERT INTO message (mess, room_id) 
			VALUES ($1,$2)
	`

	_, err := m.pool.Exec(ctx, sql, message, room)
	if err != nil {
		return err
	}

	return nil
}

func (m *Messanger) GetLastMessages(ctx context.Context, room_id int, limit int) ([]byte, error) {
	sql := `
	  SELECT (mess)
	  FROM message
	  WHERE room_id = $1
	  ORDER BY created_at DESC
	  LIMIT $2
	`

	rows, err := m.pool.Query(ctx, sql, room_id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []byte
	var msg string
	for rows.Next() {
		err = rows.Scan(&msg)
		if err != nil {
			return nil, err
		}
		messages = append(messages, []byte(msg+"\n")...)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
func (m *Messanger) CreateRoom(ctx context.Context, name string) error {
	sql := `
	  INSERT INTO room (name)
		VALUES ($1)
	`

	_, err := m.pool.Exec(ctx, sql, name)
	if err != nil {
		return err
	}

	return nil
}

func (m *Messanger) GetRoomId(ctx context.Context, name string) (int, error) {
	var id int

	sql := `
		SELECT id 
		FROM room 
		WHERE name = $1
	`

	err := m.pool.QueryRow(ctx, sql, name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
