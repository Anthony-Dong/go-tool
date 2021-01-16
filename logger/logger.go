package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/anthony-dong/go-tool/util"
)

type (
	Level int
)

var (
	levelMap = map[string]Level{
		"fatal": LevelFatal,
		"error": LevelError,
		"warn":  LevelWarning,
		"info":  LevelInfo,
		"debug": LevelDebug,
	}
	LogLevelToString = func() string {
		list, _ := util.GetMapKeysToString(levelMap)
		return util.ToCliMultiDescString(list)
	}
)

const (
	LevelFatal Level = iota + 1
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

type Logger interface {
	SetLevel(level string)
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool
	IsFatalEnabled() bool
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

func NewStdOutLogger(op ...Option) Logger {
	logger := newDefaultLogger()
	for _, elem := range op {
		elem(logger)
	}
	if !strings.HasSuffix(logger.prefix, " ") {
		logger.prefix = logger.prefix + " "
	}
	logger.Logger = log.New(logger.out, logger.prefix, logger.flag)
	return logger
}

const (
	TimeAndFileLoggerFormat int = log.Lshortfile | log.LstdFlags
	TimeLoggerFormat        int = log.LstdFlags
)

type StdOutLogger struct {
	level  Level
	prefix string
	out    io.Writer
	*log.Logger
	caller int
	flag   int
}

type Option func(*StdOutLogger)

var (
	NameOp = func(name string) Option {
		return func(logger *StdOutLogger) {
			logger.prefix = name
		}
	}
	FlagOp = func(flag int) Option {
		return func(logger *StdOutLogger) {
			logger.flag = flag
		}
	}
	CallerOp = func(caller int) Option {
		return func(logger *StdOutLogger) {
			logger.caller = caller
		}
	}
	OutOp = func(out io.Writer) Option {
		return func(logger *StdOutLogger) {
			logger.out = out
		}
	}
	LevelOp = func(level Level) Option {
		return func(logger *StdOutLogger) {
			logger.level = level
		}
	}
)

func newDefaultLogger() *StdOutLogger {
	return &StdOutLogger{
		level:  LevelDebug,
		caller: 3,
		flag:   TimeAndFileLoggerFormat,
		out:    os.Stdout,
		prefix: "[GO-TOOL]",
	}
}
func (s *StdOutLogger) output(level Level, str string) {
	if level > s.level {
		return
	}
	formatStr := ""
	switch level {
	case LevelFatal:
		formatStr = "\033[35m[FATAL]\033[0m " + str
	case LevelError:
		formatStr = "\033[31m[ERROR]\033[0m " + str
	case LevelWarning:
		formatStr = "\033[33m[WARN]\033[0m " + str
	case LevelInfo:
		formatStr = "\033[32m[INFO]\033[0m " + str
	case LevelDebug:
		formatStr = "\033[36m[DEBUG]\033[0m " + str
	}
	_ = s.Output(s.caller, formatStr)
}

func (s *StdOutLogger) IsDebugEnabled() bool {
	return !(LevelDebug > s.level)
}

func (s *StdOutLogger) IsInfoEnabled() bool {
	return !(LevelInfo > s.level)
}

func (s *StdOutLogger) IsWarnEnabled() bool {
	return !(LevelWarning > s.level)
}

func (s *StdOutLogger) IsErrorEnabled() bool {
	return !(LevelError > s.level)
}
func (s *StdOutLogger) IsFatalEnabled() bool {
	return !(LevelFatal > s.level)
}

func (s *StdOutLogger) Debugf(format string, v ...interface{}) {
	s.output(LevelDebug, fmt.Sprintf(format, v...))
}

func (s *StdOutLogger) Infof(format string, v ...interface{}) {
	s.output(LevelInfo, fmt.Sprintf(format, v...))
}
func (s *StdOutLogger) Warnf(format string, v ...interface{}) {
	s.output(LevelWarning, fmt.Sprintf(format, v...))
}
func (s *StdOutLogger) Errorf(format string, v ...interface{}) {
	s.output(LevelError, fmt.Sprintf(format, v...))
}
func (s *StdOutLogger) Fatalf(format string, v ...interface{}) {
	s.output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(-1)
}

func (s *StdOutLogger) SetLevel(level string) {
	s.level = levelMap[level]
}
