package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/MSLibs/glogger"
	"github.com/MSLibs/glogger/utils"
)

func LogRequestHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ri := &glogger.LogPayload{
			Method:    r.Method,
			Url:       r.URL.String(),
			Referer:   r.Header.Get("Referer"),
			UserAgent: r.Header.Get("User-Agent"),
		}
		ri.SourceIP = requestGetRemoteAddress(r)
		// this runs handler h and captures information about
		// HTTP request
		// m := httpsnoop.CaptureMetrics()
		ri.Size = r.ContentLength
		ctx := initLogContext(r, ri)
		// next
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func initLogContext(r *http.Request, info *glogger.LogPayload) context.Context {
	start := time.Now()
	ctx := r.Context()
	ctx = context.WithValue(ctx, glogger.RequestID, info.RequestID)
	ctx = context.WithValue(ctx, glogger.UserFlag, info.UserFlag)
	ctx = context.WithValue(ctx, glogger.PlatformID, info.PlatformID)
	ctx = context.WithValue(ctx, glogger.Referer, info.Referer)
	ctx = context.WithValue(ctx, glogger.UserAgent, info.UserAgent)
	ctx = context.WithValue(ctx, glogger.Size, info.Size)
	ctx = context.WithValue(ctx, glogger.Duration, start)
	ctx = context.WithValue(ctx, glogger.Url, r.URL.String())
	ctx = context.WithValue(ctx, glogger.SourceIP, requestGetRemoteAddress(r))
	if serverip, err := utils.ExternalIP(); err == nil {
		ctx = context.WithValue(ctx, glogger.ServerIP, serverip)
	}
	return ctx
}

func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}
