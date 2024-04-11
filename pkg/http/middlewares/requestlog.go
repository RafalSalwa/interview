package middlewares

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
)

func RequestLog(logger *logger.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			le := &logEntry{
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
				le.ServerIP = ipFromHostPort(addr.String())
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

			le.Latency = time.Since(start)
			if rcc.err == nil && rcc.r != nil {
				_, err := io.Copy(io.Discard, rcc)
				if err != nil {
					return
				}
			}
			le.RequestBodySize = rcc.n
			le.Status = w2.code
			if le.Status == 0 {
				le.Status = http.StatusOK
			}
			le.ResponseHeaderSize, le.ResponseBodySize = w2.size()
			if le.RequestURL != "/metrics" {
				logger.Info().
					Time("received_time", le.ReceivedTime).
					Str("method", le.RequestMethod).
					Str("url", le.RequestURL).
					Int64("header_size", le.RequestHeaderSize).
					Int64("body_size", le.RequestBodySize).
					Str("agent", le.UserAgent).
					Str("referer", le.Referer).
					Str("proto", le.Proto).
					Str("remote_ip", le.RemoteIP).
					Str("server_ip", le.ServerIP).
					Int("status", le.Status).
					Int64("resp_header_size", le.ResponseHeaderSize).
					Int64("resp_body_size", le.ResponseBodySize).
					Dur("latency", le.Latency).
					Msg("Request")
			}
			h.ServeHTTP(w2, r2)
		})
	}
}
