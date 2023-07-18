package logs

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger            *zap.SugaredLogger
	lumberjackWrapper *LumberjackWrapper
	logOnce           sync.Once
)

//InitLogger init log
func InitLogger(cfg *LumberjackWrapperConfig) {
	var cores []zapcore.Core
	if !cfg.FileClose {
		createFile(cfg.Path)
		writeSyncer := setLogWriter(cfg)
		encoder := setEncoder()
		fileCore := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
		cores = append(cores, fileCore)
	}
	if !cfg.StdClose {
		consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoderConfig), zapcore.Lock(os.Stdout), zapcore.DebugLevel)
		cores = append(cores, consoleCore)
	}
	core := zapcore.NewTee(cores...)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

//Sync sync log
func Sync() {
	logger.Sync()
	lumberjackWrapper.Sync()
	lumberjackWrapper.Close()
}

//GetLogger ...
func GetLogger() *zap.SugaredLogger {
	return logger
}

func setEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func setLogWriter(cfg *LumberjackWrapperConfig) zapcore.WriteSyncer {
	logOnce.Do(func() {
		lumberjackWrapper = NewLumberjackWrapper(cfg)
	})

	return zapcore.AddSync(lumberjackWrapper)
}

//Error ...
func Error(args ...interface{}) {
	logger.Error(args...)
}

//Errorf ...
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

//Info ...
func Info(args ...interface{}) {
	logger.Info(args...)
}

//Infof ...
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

//Warnf ...
func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

//Fatalf ...
func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

//Infoln ...
func Infoln(args ...interface{}) {
	logger.Info(args...)
}

//WithInfoln ...
func WithInfoln(withFields map[string]string, args interface{}) {
	var fields []interface{}
	for field, value := range withFields {
		fields = append(fields, zap.String(field, value))
	}

	logger.With(fields...).Info(args)
}

func createFile(output string) {
	//创建本地存储目录
	if _, err := os.Stat(output); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(output, os.ModePerm); err != nil {
				return
			}
		}
	}

	if _, err := os.Stat(filepath.Join(output, "info.log")); err != nil {
		if os.IsNotExist(err) {
			os.Create(filepath.Join(output, "info.log"))
		}
	}
}

//RecoveryWithZap ...
func RecoveryWithZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logger.Error("[Recovery from panic]",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
