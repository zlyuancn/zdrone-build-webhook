/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/5/17
   Description :
-------------------------------------------------
*/

package notifer

import (
    "errors"
    "fmt"
    "io/ioutil"

    "github.com/zlyuancn/zdingtalk"

    "github.com/zlyuancn/zdrone-build-webhook/config"
    "github.com/zlyuancn/zdrone-build-webhook/logger"
    "github.com/zlyuancn/zdrone-build-webhook/message"
    "github.com/zlyuancn/zdrone-build-webhook/template"
)

var _ INotifer = (*DingtalkNotifer)(nil)

var dingtalk_start_template string
var dingtalk_end_template string

type DingtalkNotifer struct {
    dt *zdingtalk.DingTalk
}

func NewDingtalkNotifer() INotifer {
    m := &DingtalkNotifer{}

    if dingtalk_start_template == "" {
        if config.Config.DingtalkStartTemplate == "" {
            config.Config.DingtalkStartTemplate = "conf/dingtask_start_template.md"
        }
        if config.Config.DingtalkEndTemplateFile == "" {
            config.Config.DingtalkEndTemplateFile = "conf/dingtask_end_template.md"
        }
        dingtalk_start_template = m.loadTemplate(config.Config.DingtalkStartTemplate)
    }
    if dingtalk_end_template == "" {
        dingtalk_end_template = m.loadTemplate(config.Config.DingtalkEndTemplateFile)
    }

    if config.Config.DingtalkAccessToken == "" {
        logger.Warn("未设置DingtalkAccessToken")
        return m
    }

    m.dt = zdingtalk.NewDingTalk(config.Config.DingtalkAccessToken).SetSecret(config.Config.DingtalkSecret)
    return m
}

func (m *DingtalkNotifer) loadTemplate(file string) string {
    body, err := ioutil.ReadFile(file)
    if err != nil {
        logger.Fatalf("无法加载模板文件: %s: %s", file, err.Error())
    }
    return string(body)
}

func (m *DingtalkNotifer) Name() string {
    return "dingtalk"
}

func (m *DingtalkNotifer) Notify(msg *message.Msg) error {
    if m.dt == nil {
        return errors.New("未创建Dingtalk实例")
    }
    return m.dt.Send(m.makeDingtalkMsg(msg), config.Config.NotifyRetry)
}

func (m *DingtalkNotifer) makeDingtalkMsg(msg *message.Msg) *zdingtalk.Msg {
    title := fmt.Sprintf("[%s] #%d %s", msg.StatusDesc, msg.TaskNum, msg.RepoName)
    text := m.makeContext(msg)
    btns := []zdingtalk.Button{
        {
            Title:     "更改记录",
            ActionURL: msg.CommitUrl,
        },
        {
            Title:     "任务构建信息",
            ActionURL: msg.TaskUrl,
        },
    }
    return zdingtalk.NewCustomCard(title, text, btns...).HorizontalButton()
}

func (m *DingtalkNotifer) makeContext(msg *message.Msg) string {
    if msg.Status == "start" {
        return template.Render(dingtalk_start_template, msg)
    }
    return template.Render(dingtalk_end_template, msg)
}
