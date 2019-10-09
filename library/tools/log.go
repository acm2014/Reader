package tools

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

//var logLevel = zap.NewAtomicLevel()

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02] 15:04:05.000"))
}

func init() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 小写编码器
		EncodeTime:     timeEncoder,                      //时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	logger := zap.New(core, caller, development, zap.AddCallerSkip(1))
	log = logger.Sugar()
}

type Logger struct {
}

var Log Logger

func (*Logger) Info(args ...interface{}) {
	args = format(args)
	log.Info(args...)
}

func (*Logger) Warn(args ...interface{}) {
	args = format(args)
	log.Warn(args...)
}

func (*Logger) Debug(args ...interface{}) {
	args = format(args)
	log.Debug(args...)
}

func (*Logger) Error(args ...interface{}) {
	args = format(args)
	log.Error(args...)
}

func (*Logger) InfoFormat(template string, args ...interface{}) {
	log.Infof(template, args...)
}

func (*Logger) WarnFormat(template string, args ...interface{}) {
	log.Warnf(template, args...)
}

func (*Logger) DebugFormat(template string, args ...interface{}) {
	log.Debugf(template, args...)
}

func (*Logger) ErrorFormat(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

//Fatal方法调用后直接os.Exit(1),慎用
func (*Logger) Fatal(args ...interface{}) {
	args = format(args)
	log.Fatal(args...)
}

//Panic方法调用后直接Panic,慎用
func (*Logger) Panic(args ...interface{}) {
	args = format(args)
	log.Error(args)
}
func format(src []interface{}) (dst []interface{}) {
	for _, s := range src {
		dst = append(dst, s)
		dst = append(dst, " ")
	}
	return
}
