package zapgorm2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

//New logger
func New(l *zap.SugaredLogger, config gormlogger.Config) gormlogger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = gormlogger.Green + "%s\n" + gormlogger.Reset + gormlogger.Green + "[info] " + gormlogger.Reset
		warnStr = gormlogger.BlueBold + "%s\n" + gormlogger.Reset + gormlogger.Magenta + "[warn] " + gormlogger.Reset
		errStr = gormlogger.Magenta + "%s\n" + gormlogger.Reset + gormlogger.Red + "[error] " + gormlogger.Reset
		traceStr = gormlogger.Green + "%s\n" + gormlogger.Reset + gormlogger.Yellow + "[%s] " + "[%.3fms] " + gormlogger.BlueBold + "[rows:%v]" + gormlogger.Reset + " %s"
		traceWarnStr = gormlogger.Green + "%s " + gormlogger.Yellow + "%s\n" + gormlogger.Reset + gormlogger.RedBold + "[%s] " + "[%.3fms] " + gormlogger.Yellow + "[rows:%v]" + gormlogger.Magenta + " %s" + gormlogger.Reset
		traceErrStr = gormlogger.RedBold + "%s " + gormlogger.MagentaBold + "%s\n" + gormlogger.Reset + gormlogger.Yellow + "[%s] " + "[%.3fms] " + gormlogger.BlueBold + "[rows:%v]" + gormlogger.Reset + " %s"
	}

	return &logger{
		SugaredLogger: l,
		Config:        config,
		infoStr:       infoStr,
		warnStr:       warnStr,
		errStr:        errStr,
		traceStr:      traceStr,
		traceWarnStr:  traceWarnStr,
		traceErrStr:   traceErrStr,
	}
}

type logger struct {
	*zap.SugaredLogger
	gormlogger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, begin.Format("2006-01-02 15:04:05"), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, begin.Format("2006-01-02 15:04:05"), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, begin.Format("2006-01-02 15:04:05"), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, begin.Format("2006-01-02 15:04:05"), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Infof(l.traceStr, utils.FileWithLineNum(), begin.Format("2006-01-02 15:04:05"), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Infof(l.traceStr, utils.FileWithLineNum(), begin.Format("2006-01-02 15:04:05"), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
