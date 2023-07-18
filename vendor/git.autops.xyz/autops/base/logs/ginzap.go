package logs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// ErrorLogger returns an ErrorLoggerT with parameter gin.ErrorTypeAny
func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

// ErrorLoggerT returns an ErrorLoggerT middleware with the given
// type gin.ErrorType.
func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.Writer.Written() {
			json := c.Errors.ByType(typ).JSON()
			if json != nil {
				c.JSON(-1, json)
			}
		}
	}
}

// Logger prints a logline for each request and measures the time to
// process for a call. It formats the log entries similar to
// http://godoc.org/github.com/gin-gonic/gin#Logger does.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// process request
		c.Next()

		var (
			latency           = time.Since(t)
			clientIP          = c.ClientIP()
			platform, version = c.GetHeader("X-Platform"), c.GetHeader("X-Version")
			method            = c.Request.Method
			statusCode        = c.Writer.Status()
			statusColor       = colorForStatus(statusCode)
			methodColor       = colorForMethod(method)
			path              = c.Request.URL.Path
			raw               = c.Request.URL.RawQuery
			deviceID          = c.GetHeader("X-Request-Id")
			userID            = c.GetString("__mp_user_id__")
		)

		if userID == "" {
			userID = c.GetString("__user_id__")
		}

		if raw != "" {
			path = path + "?" + raw
		}

		logMessage := "[GIN] |" + statusColor + " " + strconv.Itoa(statusCode) + " " + reset + "| " + fmt.Sprintf("%5v", latency) +
			" | i:" + clientIP + " c:" + platform + " v:" + version + " d:" + deviceID + " u:" + userID + "|" + methodColor + " " +
			reset + " " + method + " " + path + " " + c.Errors.String()

		switch {
		case statusCode == 401:
			logger.Warn(logMessage)
		case statusCode >= 400:
			logger.Error(logMessage)
		default:
			logger.Info(logMessage)
		}
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
