package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IRouter interface {
	Method() string
	Path() string
}

func Post(path string) IRouter {
	return &router{
		method: http.MethodPost,
		path:   path,
	}
}

func Get(path string) IRouter {
	return &router{
		method: http.MethodGet,
		path:   path,
	}
}

type router struct {
	method string
	path   string
}

func (r router) Method() string {
	return r.method
}

func (r router) Path() string {
	return r.path
}

type Handlers map[IRouter]gin.HandlersChain

type IService interface {
	ServiceName() string
	Registration() Handlers
}
