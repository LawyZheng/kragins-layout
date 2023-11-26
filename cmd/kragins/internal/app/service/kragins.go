package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/biz"
	"github.com/lawyzheng/kragins/pkg/service"
)

var (
	_ service.IService = (*HelloService)(nil)
)

func NewHelloService(uc *biz.HelloUseCase) *HelloService {
	return &HelloService{
		uc: uc,
	}
}

type HelloService struct {
	uc *biz.HelloUseCase
}

// ServiceName is the top level router
func (s *HelloService) ServiceName() string {
	return "kragins"
}

func (s *HelloService) Registration() service.Handlers {
	return service.Handlers{
		// this is the suffix router after service name
		service.Get("/hello"): []gin.HandlerFunc{helloKragins},
	}
}

func helloKragins(ctx *gin.Context) {
	// Define your own gin handler here
	ctx.String(http.StatusOK, "hello, Kragins!")
}
