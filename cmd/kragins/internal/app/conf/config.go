package conf

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/lawyzheng/kragins/pkg/process"
)

type Server struct {
	Addr    string
	Timeout int
}

func (s *Server) GetAddr() string {
	return s.Addr
}

func (s *Server) GetTimeout() time.Duration {
	return time.Second * time.Duration(s.Timeout)
}

type DataBase struct {
	DSN string
}

func (d *DataBase) GetDSN() string {
	return d.DSN
}

type Data struct {
	DataBase DataBase
}

func (d *Data) GetDataBase() *DataBase {
	return &d.DataBase
}

type Log struct {
	Console bool
	Level   string
	Dir     string

	once  sync.Once
	level *logrus.Level
}

func (l *Log) GetConsole() bool {
	return l.Console
}

func (l *Log) GetLevel() logrus.Level {
	return *l.level
}

func (l *Log) GetDir() string {
	if filepath.IsAbs(l.Dir) {
		return l.Dir
	} else {
		return filepath.Join(process.ExecutablePath(), l.Dir)
	}
}

type Config struct {
	Data   Data
	Server Server
	Log    Log
}

func (c *Config) GetData() *Data {
	return &c.Data
}

func (c *Config) GetServer() *Server {
	return &c.Server
}

func (c *Config) GetLog() *Log {
	c.Log.once.Do(func() {
		level := logrus.InfoLevel
		if err := level.UnmarshalText([]byte(c.Log.Level)); err != nil {
			panic(err.Error())
		}
		logrus.SetLevel(level)
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		c.Log.level = &level

		if !c.Log.Console {
			logrus.SetOutput(io.Discard)
		} else {
			logrus.SetOutput(colorable.NewColorableStdout())
		}
	})
	return &c.Log
}

func Read(path string) (c *Config, err error) {
	c = &Config{
		Server: Server{
			Addr: "0.0.0.0:30080",
		},
		Log: Log{
			Console: false,
			Level:   "info",
			Dir:     process.ExecutablePath(),
		},
	}

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return nil, errors.Wrap(err, "read in conf")
	}

	if err := viper.Unmarshal(&c); err != nil {
		return nil, errors.Wrap(err, "unmarshal to conf")
	}

	return c, nil
}
