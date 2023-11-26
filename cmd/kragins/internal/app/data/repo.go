package data

import (
	"context"

	"gorm.io/gorm"

	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/biz"
)

func NewRepo(data *Data) biz.Repo {
	return &repo{
		db: data.db,
	}
}

type repo struct {
	db *gorm.DB
}

func (r *repo) Add(ctx context.Context, m *biz.Model) (*biz.Model, error) {
	return m, nil
}

func (r *repo) Delete(ctx context.Context, id int) error { return nil }

func (r *repo) Update(ctx context.Context, m *biz.Model) (*biz.Model, error) {
	return m, nil
}

func (r *repo) Get(ctx context.Context, id int) (*biz.Model, error) {
	return &biz.Model{
		User: "kragins",
	}, nil
}
