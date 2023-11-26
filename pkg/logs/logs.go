package logs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lawyzheng/lyhook"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/lawyzheng/kragins/pkg/buildinfo"
	"github.com/lawyzheng/kragins/pkg/process"
)

var (
	caller = lyhook.NewDefaultCaller().SetIfCall(func(packageName string) bool {
		return strings.HasPrefix(packageName, "github.com/lawyzheng/kragins")
	})
)

func NewLogHook(output interface{}) *lyhook.LyHook {
	return lyhook.NewLyHook(output, lyhook.PickFormatter(buildinfo.DevMode())).SetCaller(caller)
}

func NewLog(output interface{}) *logrus.Logger {
	return lyhook.NewLoggerWithHook(NewLogHook(output))
}

func NewRotateFileLog(dir string, name string, level logrus.Level) (lyhook.Logger, error) {
	p := filepath.Join(dir, name)
	f, err := lyhook.NewRotateFile(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("create log file[%s] error", p))
	}
	logger := NewLog(f).WithFields(map[string]interface{}{
		"pid":  process.GetPID(),
		"file": name,
	})
	logger.Infof("Set Logger Level to [%s]", level.String())
	return logger, nil
}
