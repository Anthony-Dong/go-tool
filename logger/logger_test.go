package logger

import (
	"testing"
)

func TestNewStdOutLogger(t *testing.T) {
	logger := NewStdOutLogger(LevelOp(LevelDebug))
	logger.Infof("Infof")
	logger.Debugf("Debugf")
	logger.Errorf("Errorf")
	logger.Warnf("Warnf")
	logger.Fatalf("Fatalf")
}
