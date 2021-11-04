package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"tiktok-uploader/configs"
)

var ZapLogger *zap.Logger

func Init(fileName string) *zap.SugaredLogger {
	// 通过配置lumberjack.Logger，来设置日志文件的切割
	hook := lumberjack.Logger{
		Filename:   configs.AppConfig.LogPath + fileName, // 日志文件路径
		MaxSize:    64,                                   // 每个日志文件保存的大小 单位:M
		MaxAge:     7,                                    // 文件最多保存多少天
		MaxBackups: 30,                                   // 日志文件最多保存多少个备份
		Compress:   false,                                // 是否压缩
	}
	// 配置zapcore.EncoderConfig，来设置日志的格式
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "log",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
		//ConsoleSeparator: " ",
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}

	// 通过zapcore.NewCore创建一个core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		//zapcore.NewJSONEncoder(encoderConfig),
		// 设置多个输出
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	// 设置初始化字段
	//field := zap.Fields(zap.String("appName", name))
	// 构造日志
	//ZapLogger = zap.New(core, caller, development, field)

	// 通过zap.New创建日志实例
	ZapLogger = zap.New(core, caller, development)
	return ZapLogger.Sugar()

	//ZapLogger.Info("init log module")

	//ZapLogger.Info("无法获取网址",
	//	zap.String("url", "http://www.baidu.com"),
	//	zap.Int("attempt", 3),
	//	zap.Duration("backoff", time.Second))
}