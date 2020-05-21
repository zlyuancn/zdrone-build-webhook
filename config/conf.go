/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/5/16
   Description :
-------------------------------------------------
*/

package config

import (
    "github.com/zlyuancn/zenvconf"
    "github.com/zlyuancn/zlog"
)

const (
    MyPicRepoUrl    = "https://git.zlyuan.xyz/zlyuancn/zdrone-build-webhook-img.git"
    MyPicRepoBranch = "master"
)

var Config struct {
    Bind   string `env:"DRONE_BIND"`   // bind端口
    Debug  bool   `env:"DRONE_DEBUG"`  // debug模式
    Secret string `env:"DRONE_SECRET"` // webhook秘钥

    DroneServer string `env:"DRONE_SERVER"` // drone服务地址

    LogPath string `env:"LOG_PATH"` // 日志路径

    Notifer         string `env:"NOTIFER"`           // 通告者,多个通告者用半角逗号隔开
    NotifyRetry     int    `env:"NOTIFY_RETRY"`      // 通告失败重试次数
    OffCreateNotify bool   `env:"OFF_CREATE_NOTIFY"` // 关闭创建动作的通告

    DingtalkAccessToken     string `env:"DINGTALK_ACCESSTOKEN"`    // 钉钉access_token
    DingtalkSecret          string `env:"DINGTALK_SECRET"`         // 钉钉secret
    DingtalkStartTemplate   string `env:"DINGTALK_START_TEMPLATE"` // 钉钉消息任务开始模板文件
    DingtalkEndTemplateFile string `env:"DINGTALK_END_TEMPLATE"`   // 钉钉消息任务结束模板文件
}

func Init() {
    Config.NotifyRetry = 2

    err := zenvconf.NewEnvConf().Parse(&Config)
    if err != nil {
        zlog.Fatal("初始化失败", err)
    }

    if Config.Bind == "" {
        Config.Bind = ":80"
    }
    if Config.Secret == "" {
        zlog.Fatal("未设置 Secret")
    }
    if Config.DroneServer == "" {
        zlog.Fatal("未设置 DroneServer")
    }
}
