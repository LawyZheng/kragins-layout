package data

import (
	"github.com/lawyzheng/lyhook"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/conf"
	"github.com/lawyzheng/kragins/pkg/logs"
)

type config struct {
	dbDialect gorm.Dialector
	logger    lyhook.Logger
}

func parseConfig(c *conf.Config) (config, error) {
	var dialect gorm.Dialector
	// Or you can user your own database dsn
	dialect = postgres.Open(c.GetData().GetDataBase().GetDSN())

	logger, err := logs.NewRotateFileLog(c.GetLog().GetDir(), "db.log", c.GetLog().GetLevel())
	if err != nil {
		return config{}, err
	}

	return config{
		dbDialect: dialect,
		logger:    logger,
	}, nil
}
