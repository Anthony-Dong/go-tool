package log

import (
	logger2 "github.com/anthony-dong/go-tool/commons/logger"
)

var (
	log = logger2.NewStdOutLogger(logger2.NameOp("[GO-TOOL]"))
)

var (
	IsDebugEnabled = log.IsDebugEnabled
	IsInfoEnabled  = log.IsInfoEnabled
	IsWarnEnabled  = log.IsWarnEnabled
	IsErrorEnabled = log.IsErrorEnabled
	IsFatalEnabled = log.IsFatalEnabled
	Infof          = log.Infof
	Debugf         = log.Debugf
	Errorf         = log.Errorf
	Warnf          = log.Warnf
	Fatalf         = log.Fatalf
	SetLevel       = log.SetLevel
)
