package log

import "github.com/anthony-dong/go-tool/logger"

var (
	log = logger.NewStdOutLogger(logger.NameOp("[GO-TOOL]"))
)

var (
	Infof    = log.Infof
	Debugf   = log.Debugf
	Errorf   = log.Errorf
	Warnf    = log.Warnf
	Fatalf   = log.Fatalf
	SetLevel = log.SetLevel
)
