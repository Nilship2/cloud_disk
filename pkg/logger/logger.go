// pkg/logger/logger.go
package logger

import (
	"os"

	"cloud-disk/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init(logConfig *config.LogConfig) error {
	var err error

	switch logConfig.Format {
	case "json":
		Logger, err = initJSONLogger(logConfig)
	case "console":
		fallthrough
	default:
		Logger, err = initConsoleLogger(logConfig)
	}

	if err != nil {
		return err
	}

	// 替换全局logger
	zap.ReplaceGlobals(Logger)

	return nil
}

func initJSONLogger(config *config.LogConfig) (*zap.Logger, error) {
	level := getZapLevel(config.Level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	configs := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	if config.Output == "file" && config.FilePath != "" {
		// 文件输出
		writer := getLogWriter(config.FilePath)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			writer,
			level,
		)
		return zap.New(core, configs...), nil
	}

	// 控制台输出
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	return zap.New(core, configs...), nil
}

func initConsoleLogger(config *config.LogConfig) (*zap.Logger, error) {
	level := getZapLevel(config.Level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	configs := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	return zap.New(core, configs...), nil
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func getLogWriter(filePath string) zapcore.WriteSyncer {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(file)
}

// 便捷方法
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

func Infof(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Sugar().Infof(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Sugar().Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Sugar().Fatalf(format, args...)
	}
}
