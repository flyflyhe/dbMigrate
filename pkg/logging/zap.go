package logging

import (
	"github.com/dbMigrate/v2/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

var zapLog *zap.Logger

func initZap(logConfig config.LogConfig) {
	var coreArr []zapcore.Core

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format(time.RFC3339) + "]")
	}

	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()            //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = customTimeEncoder                 //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //按级别显示不同颜色
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder       //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		if logConfig.Debug {
			return lev < zap.ErrorLevel && lev >= zap.DebugLevel
		} else {
			return lev < zap.ErrorLevel && lev > zap.DebugLevel
		}
	})

	log.Println("输出日志目录-----------------------------------------------------")
	log.Println(config.GetApp().Log.Info.Filename)
	//info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.Info.Filename, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    2,                       //文件大小限制,单位MB
		MaxBackups: 100,                     //最大保留日志文件数量
		MaxAge:     30,                      //日志文件保留天数
		Compress:   false,                   //是否压缩处理
	})
	//zapcore.AddSync(os.Stdout)
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	//error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.Error.Filename, //日志文件存放目录
		MaxSize:    1,                        //文件大小限制,单位MB
		MaxBackups: 5,                        //最大保留日志文件数量
		MaxAge:     30,                       //日志文件保留天数
		Compress:   false,                    //是否压缩处理
	})
	//zapcore.AddSync(os.Stdout)
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	options := []zap.Option{zap.AddCaller()}
	if logConfig.Debug {
		options = append(options, zap.AddStacktrace(lowPriority))
	}
	zapLog = zap.New(zapcore.NewTee(coreArr...), options...) //zap.AddCaller()为显示文件名和行号
}
