package logger

import (
	"testing"
)

func TestNewStdOutLogger(t *testing.T) {
	logger := NewStdOutLogger(LevelOp(LevelFatal))
	Infof("Infof")
	Debugf("Debugf")
	Errorf("Errorf")
	Warnf("Warnf")
	t.Log(IsDebugEnabled())
	t.Log(IsInfoEnabled())
	//logger.Fatalf("Fatalf")
}
