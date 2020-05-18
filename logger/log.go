/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/4/1
   Description :
-------------------------------------------------
*/

package logger

import (
    "github.com/zlyuancn/zlog"

    "github.com/zlyuancn/zdrone-build-webhook/config"
)

var Log = zlog.DefaultLogger

func Init() {
    conf := zlog.DefaultConfig
    conf.CallerSkip = 1
    conf.IsTerminal = false

    if !config.Config.Debug {
        conf.Level = "info"
        conf.ShowInitInfo = false
        conf.ShowFileAndLinenum = false
    } else {
        conf.Level = "debug"
        conf.ShowInitInfo = true
        conf.ShowFileAndLinenum = true
    }

    if config.Config.LogPath != "" {
        conf.Path = config.Config.LogPath
        conf.WriteToFile = true
    }
    Log = zlog.New(conf)
}

func Debug(v ...interface{}) {
    Log.Debug(v...)
}
func Info(v ...interface{}) {
    Log.Info(v...)
}
func Warn(v ...interface{}) {
    Log.Warn(v...)
}
func Error(v ...interface{}) {
    Log.Error(v...)
}
func DPanic(v ...interface{}) {
    Log.DPanic(v...)
}
func Panic(v ...interface{}) {
    Log.Panic(v...)
}
func Fatal(v ...interface{}) {
    Log.Fatal(v...)
}

func Debugf(format string, v ...interface{}) {
    Log.Debugf(format, v...)
}
func Infof(format string, v ...interface{}) {
    Log.Infof(format, v...)
}
func Warnf(format string, v ...interface{}) {
    Log.Warnf(format, v...)
}
func Errorf(format string, v ...interface{}) {
    Log.Errorf(format, v...)
}
func DPanicf(format string, v ...interface{}) {
    Log.DPanicf(format, v...)
}
func Panicf(format string, v ...interface{}) {
    Log.Panicf(format, v...)
}
func Fatalf(format string, v ...interface{}) {
    Log.Fatalf(format, v...)
}
