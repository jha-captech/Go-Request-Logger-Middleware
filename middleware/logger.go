package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type logDetails struct {
	status   int
	duration time.Duration
	clientIP string
	method   string
	path     string
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func LoggerColorMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			logger.Info(
				logFormatter(logDetails{
					status:   wrapped.statusCode,
					duration: time.Since(start),
					clientIP: r.RemoteAddr,
					method:   r.Method,
					path:     r.URL.Path,
				}),
			)
		})
	}
}

func logFormatter(detail logDetails) string {
	statusColor := GetStatusCodeColor(detail.status)
	methodColor := getMethodColor(detail.method)
	resetColor := getResetColor()

	return fmt.Sprintf("|%s %3d %s| %13v | %15s |%s %-7s %s %#v",
		statusColor, detail.status, resetColor,
		detail.duration,
		detail.clientIP,
		methodColor, detail.method, resetColor,
		detail.path,
	)
}

func GetStatusCodeColor(code int) string {
	switch {
	case code >= http.StatusContinue && code < http.StatusOK:
		return white
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func getMethodColor(method string) string {
	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func getResetColor() string {
	return reset
}
