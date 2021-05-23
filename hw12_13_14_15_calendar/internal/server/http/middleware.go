package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"

	log "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)

		m := httpsnoop.CaptureMetrics(next, w, r)
		log.Info(fmt.Sprintf("%s [%s] %s %s %s %d %d %s", getIP(r), t.String(), r.Method, r.URL.String(), r.Proto, 200, m.Duration, r.Header.Get("User-Agent")), nil)
	})
}

func ipFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func getIP(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipFromRemoteAddr(r.RemoteAddr)
	}

	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}

		return parts[0]
	}
	return hdrRealIP
}
