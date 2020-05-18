/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/5/17
   Description :
-------------------------------------------------
*/

package notifer

import (
    "encoding/json"
    "strings"

    "github.com/zlyuancn/zsignal"

    "github.com/zlyuancn/zdrone-build-webhook/config"
    "github.com/zlyuancn/zdrone-build-webhook/logger"
    "github.com/zlyuancn/zdrone-build-webhook/message"
)

var notifers []INotifer
var notify_task_queue chan *message.Msg

type INotifer interface {
    Name() string
    Notify(msg *message.Msg) error
}

func Init() {
    notify_task_queue = make(chan *message.Msg, 100)
    zsignal.RegisterOnShutdown(func() {
        close(notify_task_queue)
    })

    go func() {
        for msg := range notify_task_queue {
            notify(msg)
        }
    }()

    if config.Config.Notifer == "" {
        return
    }

    for _, notifer := range strings.Split(config.Config.Notifer, ",") {
        switch notifer {
        case "dingtalk":
            notifers = append(notifers, NewDingtalkNotifer())
        default:
            logger.Warnf("未定义的通告对象: %s", notifer)
            continue
        }
        logger.Debugf("添加<%s>通告者", notifer)
    }
}

func notify(msg *message.Msg) {
    if config.Config.Debug {
        body, _ := json.MarshalIndent(msg, "", "    ")
        logger.Debug(string(body))
    }

    if config.Config.OffCreateNotify && msg.Status == "start" {
        logger.Debugf("跳过<%s>动作的公告", msg.Status)
        return
    }

    for _, notifer := range notifers {
        if err := notifer.Notify(msg); err != nil {
            logger.Warnf("通告<%s>失败: %s", notifer.Name(), err.Error())
        }
    }
}

func Notify(msg *message.Msg) {
    notify_task_queue <- msg
}
