package service

import (
	"github.com/gin-gonic/gin"
	"github.com/lawyzheng/kragins/pkg/service"
)

var (
	_ service.IService = (*HelloService)(nil)
)

type HelloUseCase interface {
	GreetHandler(ctx *gin.Context)
}

func NewHelloService(uc HelloUseCase) *HelloService {
	return &HelloService{
		uc: uc,
	}
}

type HelloService struct {
	uc HelloUseCase
}

// ServiceName is the top level router
func (s *HelloService) ServiceName() string {
	return "kragins"
}

func (s *HelloService) Registration() service.Handlers {
	return service.Handlers{
		// this is the suffix router after service name
		service.Get("/hello"): []gin.HandlerFunc{s.uc.GreetHandler},
	}
}
