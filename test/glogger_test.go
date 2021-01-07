package test

import (
	"context"
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

func TestGLoggerSetContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "requestId", "jasiudhasuidhuaisdhuaisdhiuasdhui")
	ctx = context.WithValue(ctx, "platformId", "PC")
	ctx = context.WithValue(ctx, "userflag", "185236523365")
	log := glog.SetContext(&ctx)
	log.Infof("logging context...")
}
