package middleware

import (
	"context"
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"

	"hyperfulcrum/pkg/logger"
	"hyperfulcrum/pkg/utils"
)

type requestContextKey string

const requestIDKey requestContextKey = "requestID"

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *responseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}

	w.status = status
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Logger.Error(
					"Request panic recovered",
					"request_id", GetRequestID(r),
					"method", r.Method,
					"path", r.URL.Path,
					"panic", recovered,
					"stack", string(debug.Stack()),
				)
				utils.WriteJSONErrorResponse(
					w,
					http.StatusInternalServerError,
					"Internal server error",
					nil,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		writer := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(writer, r)

		logger.Logger.Info(
			"Request completed",
			"request_id", GetRequestID(r),
			"method", r.Method,
			"path", r.URL.Path,
			"status", writer.status,
			"duration", time.Since(startedAt),
		)
	})
}

func RequestTimeout(timeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UUIDPathValidator(parameters ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, parameter := range parameters {
				if _, err := uuid.Parse(r.PathValue(parameter)); err != nil {
					utils.WriteJSONErrorResponse(
						w,
						http.StatusBadRequest,
						parameter+" must be a valid UUID",
						nil,
					)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetRequestID(r *http.Request) string {
	requestID, _ := r.Context().Value(requestIDKey).(string)
	return requestID
}

func JSONErrorStatus(err error) int {
	switch {
	case errors.Is(err, utils.ErrUnsupportedMediaType):
		return http.StatusUnsupportedMediaType
	case errors.Is(err, utils.ErrRequestBodyTooLarge):
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusBadRequest
	}
}
