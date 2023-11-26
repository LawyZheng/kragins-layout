package biz

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Model struct {
	User string
}

type Repo interface {
	Add(ctx context.Context, m *Model) (*Model, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, m *Model) (*Model, error)
	Get(ctx context.Context, id int) (*Model, error)
}

func NewHelloUseCase(repo Repo) *HelloUseCase {
	return &HelloUseCase{repo: repo}
}

type HelloUseCase struct {
	repo Repo
}

func (u *HelloUseCase) GreetHandler(ctx *gin.Context) {
	model, _ := u.repo.Get(ctx, 0)
	ctx.String(http.StatusOK, "Hello, %s", model.User)
	ctx.Abort()
	return
}
