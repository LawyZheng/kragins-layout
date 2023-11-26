package server

import (
	"io"
	"path/filepath"

	"github.com/lawyzheng/lyhook"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"

	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/conf"
	"github.com/lawyzheng/kragins/pkg/buildinfo"
	"github.com/lawyzheng/kragins/pkg/logs"
)

func getGinWriterAndLogger(c *conf.Config) (ginWriter io.Writer, logger lyhook.Logger, e error) {
	if buildinfo.DevMode() {
		ginWriter = colorable.NewColorableStdout()
	} else {
		f, err := lyhook.NewRotateFile(filepath.Join(c.GetLog().GetDir(), "gin.log"))
		if err != nil {
			return nil, nil, errors.Wrap(err, "create log file")
		}
		ginWriter = f
		if c.Log.Console {
			ginWriter = io.MultiWriter(ginWriter, colorable.NewColorableStdout())
		}
	}

	logger, err := logs.NewRotateFileLog(c.GetLog().GetDir(), "server.log", c.GetLog().GetLevel())
	if err != nil {
		return nil, nil, err
	}

	return ginWriter, logger, nil
}
