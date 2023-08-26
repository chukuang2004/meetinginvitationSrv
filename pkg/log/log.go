package log

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Level string
}

type Logger struct {
	*logrus.Logger
}

var log *Logger = nil

func GetInstance() *Logger {
	if log == nil {
		log = &Logger{}
	}

	return log
}

func (l *Logger) Init(conf *Config) {

	log := logrus.New()
	log.Out = os.Stdout
	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		level = logrus.InfoLevel
	}

	log.SetLevel(level)

	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		FullTimestamp:   true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := path.Base(frame.File)
			fileLoc := fmt.Sprintf("%s:%d", fileName, frame.Line)

			fun := path.Base(frame.Function)
			return fun, fileLoc
		},
	})

	l.Logger = log
}
