package data

import (
	"context"
	"time"

	"github.com/google/wire"
	"github.com/lawyzheng/gormlogger"
	"github.com/lawyzheng/lyhook"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/biz"
	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/conf"
)

var ProviderSet = wire.NewSet(NewData, NewRepo)

type Data struct {
	db     *gorm.DB
	logger lyhook.Logger
}

func NewData(c *conf.Config) (*Data, func(), error) {
	config, err := parseConfig(c)
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(config.dbDialect, &gorm.Config{
		Logger: gormlogger.NewLogger(config.logger),
	})
	if err != nil {
		return nil, nil, err
	}
	sql, err := db.DB()
	if err != nil {
		return nil, nil, errors.Wrap(err, "get sql instance error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := sql.PingContext(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "ping database error")
	}

	if err := db.AutoMigrate(&biz.Model{}); err != nil {
		sql.Close()
		return nil, nil, errors.Wrap(err, "migrate database error")
	}

	logger := config.logger
	cleanup := func() {
		sql.Close()
		logger.Info("closing the data resources")
	}

	return &Data{
		db:     db,
		logger: logger,
	}, cleanup, nil
}
