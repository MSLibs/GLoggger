package test

import (
	"glogger"
	"testing"
)

func TestGLogger(t *testing.T) {
	glog := glogger.CreateLog()
	glog.Info("logging success...")
	// t.Errorf("logging %s", "success")
}
