package service

import (
	"github.com/google/wire"

	"github.com/lawyzheng/kragins/pkg/service"
)

var ProviderSet = wire.NewSet(NewServices, NewHelloService)

func NewServices(hs *HelloService) []service.IService {
	svc := make([]service.IService, 0)
	svc = append(svc, hs)
	return svc
}
