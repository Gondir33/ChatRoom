package modules

import (
	"goTest/internal/infrastructure/component"
	rService "goTest/internal/modules/messanger/service"
	"goTest/internal/storages"
)

type Services struct {
	rService.Messangerer
}

func NewServices(storages *storages.Storages, components *component.Components) *Services {
	return &Services{
		Messangerer: rService.NewMessanger(storages.Messangerer, components),
	}
}
