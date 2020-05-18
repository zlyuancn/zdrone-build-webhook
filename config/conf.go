/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/5/16
   Description :
-------------------------------------------------
*/

package config

import (
    "github.com/kelseyhightower/envconfig"
    "github.com/zlyuancn/zlog"
)

const (
    MyRepoUrl    = "https://github.com/zlyuancn/zdrone-build-webhook.git"
    MyMainBranch = "master"
)

var Config struct {
    Bind   string `envconfig:"DRONE_BIND"`   // bind端口
    Debug  bool   `envconfig:"DRONE_DEBUG"`  // debug模式
    Secret string `envconfig:"DRONE_SECRET"` // webhook秘钥

    DroneServer string `envconfig:"DRONE_SERVER"` // drone服务地址

    LogPath string `envconfig:"LogPath"` // 日志路径

    Notifer         string `envconfig:"Notify"`          // 通告者,多个通告者用半角逗号隔开
    NotifyRetry     int    `envconfig:"NotifyRetry"`     // 通告失败重试次数
    OffCreateNotify bool   `envconfig:"OffCreateNotify"` // 关闭创建动作的通告

    DingtalkAccessToken     string `envconfig:"DingtalkAccessToken"`     // 钉钉access_token
    DingtalkSecret          string `envconfig:"DingtalkSecret"`          // 钉钉secret
    DingtalkStartTemplate   string `envconfig:"DingtalkStartTemplate"`   // 钉钉消息任务开始模板文件
    DingtalkEndTemplateFile string `envconfig:"DingtalkEndTemplateFile"` // 钉钉消息任务结束模板文件
}

func Init() {
    Config.NotifyRetry = 2

    err := envconfig.Process("", &Config)
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
