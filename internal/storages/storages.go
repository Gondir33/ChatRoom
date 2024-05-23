package storages

import (
	rStorage "goTest/internal/modules/messanger/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storages struct {
	rStorage.Messangerer
}

func NewStorages(pool *pgxpool.Pool) *Storages {
	return &Storages{
		Messangerer: rStorage.NewMessanger(pool),
	}
}
