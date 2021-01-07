package test

import (
	"testing"

	"github.com/MSLibs/glogger"
)

var glog = glogger.CreateLog()

func TestGLogger(t *testing.T) {
	glog.Info("logging success...")
	// t.Errorf("logging %s", "success")
}

func TestGLoggerFormat(t *testing.T) {
	glog.Infof("logging format...")
}
