package server

import (
	"context"
	"io"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/lawyzheng/lyhook"
	"github.com/pkg/errors"

	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/conf"
	"github.com/lawyzheng/kragins/pkg/buildinfo"
	"github.com/lawyzheng/kragins/pkg/service"
)

var ProviderSet = wire.NewSet(NewServer)

func getEngine(writer io.Writer) *gin.Engine {
	if buildinfo.DevMode() {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = writer
	return gin.Default()
}

func NewServer(c *conf.Config, services []service.IService) (*Server, error) {
	ginWriter, logger, err := getGinWriterAndLogger(c)
	if err != nil {
		return nil, errors.Wrap(err, "parse server config")
	}

	s := &Server{
		engine:   getEngine(ginWriter),
		logger:   logger,
		listen:   c.GetServer().GetAddr(),
		timeout:  c.GetServer().GetTimeout(),
		services: services,
	}
	return s, nil
}

type Server struct {
	engine   *gin.Engine
	logger   lyhook.Logger
	services []service.IService
	timeout  time.Duration
	listen   string
}

func timeOutMiddleWare(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		c.Request = c.Request.Clone(ctx)

		c.Next()
	}
}

func (s *Server) registerService() {
	for _, svc := range s.services {
		name := svc.ServiceName()
		for r, h := range svc.Registration() {
			p := path.Join(name, r.Path())
			if s.timeout != 0 {
				s.engine.Use(timeOutMiddleWare(s.timeout))
			}
			s.engine.Handle(r.Method(), p, h...)
		}
	}
}

func (s *Server) run() error {
	s.logger.Infof("Listening and Serving HTTP on %s", s.listen)
	return s.engine.Run(s.listen)
}

func (s *Server) Run() error {
	s.registerService()
	return s.run()
}

func (s *Server) HandleError(err error) {
	s.logger.Errorf("Starting server encounter with error: %s", err)
	os.Exit(1)
}
