package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

type LogEntry struct {
	Timestamp   time.Time   `json:"timestamp"`
	Method      string      `json:"method"`
	URI         string      `json:"uri"`
	Status      int         `json:"status"`
	Latency     string      `json:"latency"`
	RequestBody interface{} `json:"request_body,omitempty"`
	Error       string      `json:"error,omitempty"`
	UserAgent   string      `json:"user_agent,omitempty"`
	RemoteIP    string      `json:"remote_ip,omitempty"`
}

func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			var requestBody interface{}
			if c.Request().Body != nil {
				bodyBytes, _ := io.ReadAll(c.Request().Body)
				c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				if len(bodyBytes) > 0 {
					json.Unmarshal(bodyBytes, &requestBody)
				}
			}

			err := next(c)

			latency := time.Since(start)

			logEntry := LogEntry{
				Timestamp:   start,
				Method:      c.Request().Method,
				URI:         c.Request().RequestURI,
				Status:      c.Response().Status,
				Latency:     latency.String(),
				RequestBody: requestBody,
				UserAgent:   c.Request().UserAgent(),
				RemoteIP:    c.RealIP(),
			}

			if err != nil {
				logEntry.Error = err.Error()
			}

			if logBytes, marshalErr := json.Marshal(logEntry); marshalErr == nil {
				log.Printf("API_LOG: %s", string(logBytes))
			}

			return err
		}
	}
}
