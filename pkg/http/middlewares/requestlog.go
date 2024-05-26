package middlewares

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/RafalSalwa/auth-api/pkg/logger"
)

func RequestLog(logger *logger.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			entry := &logEntry{
				ReceivedTime:      start,
				RequestMethod:     r.Method,
				RequestURL:        r.URL.String(),
				RequestHeaderSize: headerSize(r.Header),
				UserAgent:         r.UserAgent(),
				Referer:           r.Referer(),
				Proto:             r.Proto,
				RemoteIP:          ipFromHostPort(r.RemoteAddr),
			}

			if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
				entry.ServerIP = ipFromHostPort(addr.String())
			}
			body, _ := io.ReadAll(r.Body)
			err := r.Body.Close()
			if err != nil {
				logger.Error().Err(err)
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			r2 := new(http.Request)
			*r2 = *r
			rcc := &readCounterCloser{r: io.NopCloser(bytes.NewBuffer(body))}
			r2.Body = rcc
			w2 := &responseStats{w: w}
			r2.Body = io.NopCloser(bytes.NewBuffer(body))

			entry.Latency = time.Since(start)
			if rcc.err == nil && rcc.r != nil {
				_, err := io.Copy(io.Discard, rcc)
				if err != nil {
					return
				}
			}
			entry.RequestBodySize = rcc.n
			entry.Status = w2.code
			if entry.Status == 0 {
				entry.Status = http.StatusOK
			}
			entry.ResponseHeaderSize, entry.ResponseBodySize = w2.size()
			if entry.RequestURL != "/metrics" {
				logger.Info().
					Time("received_time", entry.ReceivedTime).
					Str("method", entry.RequestMethod).
					Str("url", entry.RequestURL).
					Int64("header_size", entry.RequestHeaderSize).
					Int64("body_size", entry.RequestBodySize).
					Str("agent", entry.UserAgent).
					Str("referer", entry.Referer).
					Str("proto", entry.Proto).
					Str("remote_ip", entry.RemoteIP).
					Str("server_ip", entry.ServerIP).
					Int("status", entry.Status).
					Int64("resp_header_size", entry.ResponseHeaderSize).
					Int64("resp_body_size", entry.ResponseBodySize).
					Dur("latency", entry.Latency).
					Msg("Request")
			}
			h.ServeHTTP(w2, r2)
		})
	}
}
