package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Logger *zap.Logger

func InitLogger(serverName string, stdout bool) *zap.Logger {
	infoHook := lumberjack.Logger{
		Filename:   "./logs/info.log", // 日志文件路径
		MaxSize:    32,                // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                // 日志文件最多保存多少个备份
		MaxAge:     30,                // 文件最多保存多少天
		Compress:   true,              // 是否压缩
		LocalTime:  true,              // 本地时间
	}
	warnHook := lumberjack.Logger{
		Filename:   "./logs/warn.log", // 日志文件路径
		MaxSize:    32,                // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                // 日志文件最多保存多少个备份
		MaxAge:     365,               // 文件最多保存多少天
		Compress:   true,              // 是否压缩
		LocalTime:  true,              // 本地时间

	}
	errorHook := lumberjack.Logger{
		Filename:   "./logs/error.log", // 日志文件路径
		MaxSize:    32,                 // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                 // 日志文件最多保存多少个备份
		MaxAge:     365,                // 文件最多保存多少天
		Compress:   true,               // 是否压缩
		LocalTime:  true,               // 本地时间

	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "line",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(time.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.ErrorLevel
	})

	var (
		infoMultiWriter  zapcore.WriteSyncer
		warnMultiWriter  zapcore.WriteSyncer
		errorMultiWriter zapcore.WriteSyncer
	)

	if stdout {
		infoMultiWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&infoHook))
		warnMultiWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&warnHook))
		errorMultiWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&errorHook))
	} else {
		infoMultiWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&infoHook))
		warnMultiWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&warnHook))
		errorMultiWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&errorHook))
	}
	info := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		infoMultiWriter,                       // 打印到控制台和文件
		infoLevel,                             // 日志级别
	)
	warn := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		warnMultiWriter,                       // 打印到控制台和文件
		warnLevel,                             // 日志级别
	)
	_error := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		errorMultiWriter,                      // 打印到控制台和文件
		errorLevel,                            // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()

	// 当前打印日志的位置意义不大，跳过当前位置，线上上层调用
	//callSkip := zap.AddCallerSkip(1)

	// 调用堆栈trace
	trace := zap.AddStacktrace(zapcore.ErrorLevel)

	// 开启文件及行号
	development := zap.Development()

	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", serverName))

	// 构造日志
	core := zapcore.NewTee(
		info,
		warn,
		_error,
	)
	return zap.New(core, caller, trace, development, filed)
}

// 不区分不同级别日志文件
func FileLogger() *zap.Logger {
    // 动态调整日志级别
    allLevel := zap.NewAtomicLevel()

    hook := lumberjack.Logger{
        Filename:   "./logs/logs.log",
        MaxSize:    1024, // megabytes
        MaxBackups: 3,
        MaxAge:     7,    //days
        Compress:   true, // disabled by default
    }
    w := zapcore.AddSync(&hook)

    allLevel.SetLevel(zap.InfoLevel)
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    core := zapcore.NewCore(
        zapcore.NewConsoleEncoder(encoderConfig),
        w,
		allLevel,
    )

    return zap.New(core)
}
