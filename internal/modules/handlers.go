package modules

import (
	"goTest/internal/infrastructure/component"
	rHandler "goTest/internal/modules/messanger/controller"
)

type Controllers struct {
	rHandler.Messangerer
}

func NewControllers(services *Services, components *component.Components) *Controllers {
	return &Controllers{
		Messangerer: rHandler.NewMessanger(services.Messangerer, components.Logger),
	}
}
