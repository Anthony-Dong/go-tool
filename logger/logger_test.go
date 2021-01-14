package logger

import (
	"testing"
)

func TestNewStdOutLogger(t *testing.T) {
	logger := NewStdOutLogger(LevelOp(LevelFatal))
	logger.Infof("Infof")
	logger.Debugf("Debugf")
	logger.Errorf("Errorf")
	logger.Warnf("Warnf")
	t.Log(logger.IsDebugEnabled())
	t.Log(logger.IsInfoEnabled())
	//logger.Fatalf("Fatalf")
}
