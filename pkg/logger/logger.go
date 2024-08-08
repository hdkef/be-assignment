package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetReportCaller(true)
}

func LogWarn(ctx context.Context, args ...interface{}) {

	reqID := getRequestID(ctx)
	fileAndLine := getFileAndLine()
	log.WithFields(logrus.Fields{
		"request_id": reqID,
		"source":     fileAndLine,
	}).Warn(args...)
}

func LogError(ctx context.Context, args ...interface{}) {

	reqID := getRequestID(ctx)
	fileAndLine := getFileAndLine()
	log.WithFields(logrus.Fields{
		"request_id": reqID,
		"source":     fileAndLine,
	}).Error(args...)
}

func LogFatal(ctx context.Context, args ...interface{}) {

	reqID := getRequestID(ctx)
	fileAndLine := getFileAndLine()
	log.WithFields(logrus.Fields{
		"request_id": reqID,
		"source":     fileAndLine,
	}).Fatal(args...)
}

func LogInfo(ctx context.Context, args ...interface{}) {

	reqID := getRequestID(ctx)
	fileAndLine := getFileAndLine()
	log.WithFields(logrus.Fields{
		"request_id": reqID,
		"source":     fileAndLine,
	}).Info(args...)
}

const RequestIDKey = "RequestID"

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Read the request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Generate request ID
		requestID := uuid.New().String()

		// Set request ID to context
		c.Set(RequestIDKey, requestID)

		// Create a custom response writer
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log the request
		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency":    latency,
			"client_ip":  c.ClientIP(),
			"request":    string(bodyBytes),
			"response":   blw.body.String(),
		}).Info("HTTP Request")
	}
}

func getRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
