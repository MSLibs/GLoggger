package test

import (
	"context"
	"testing"

	"github.com/MSLibs/glogger"
)

var glog = glogger.CreateLog(glogger.GLoggerConfig{})

func TestGLogger(t *testing.T) {
	glog.Info("logging success...")
	// t.Errorf("logging %s", "success")
}

func TestGLoggerFormat(t *testing.T) {
	glog.Infof("logging format...")
}

func TestGLoggerSetContext_V2(t *testing.T) {
	ctx := context.WithValue(context.Background(), "requestId", "jasiudhasuidhuaisdhuaisdhiuasdhui")
	ctx = context.WithValue(ctx, "platformId", "PC")
	ctx = context.WithValue(ctx, "userflag", "185236523365")
	ctx = context.WithValue(ctx, glogger.Method, "mymethod")
	ctx = context.WithValue(ctx, glogger.UserAgent, "chome65")
	ctx = context.WithValue(ctx, glogger.Url, "localhost/demo")
	ctx = context.WithValue(ctx, glogger.SourceIP, "localhost")
	ctx = context.WithValue(ctx, glogger.ServerIP, "192.168.3.67")
	glog.WithInfof(&ctx, "logging context v2 info with context")
	glog.WithWarnf(&ctx, "logging context v2 warn with context")
}

func TestGLoggerSetContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "requestId", "jasiudhasuidhuaisdhuaisdhiuasdhui")
	ctx = context.WithValue(ctx, "platformId", "PC")
	ctx = context.WithValue(ctx, "userflag", "185236523365")
	log := glog.SetContext(&ctx)
	log.Infof("logging context...")
}
