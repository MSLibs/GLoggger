package test

import (
	"testing"

	"github.com/MSLibs/glogger"
)

func TestGLogger(t *testing.T) {
	glog := glogger.CreateLog()
	glog.Info("logging success...")
	// t.Errorf("logging %s", "success")
}
