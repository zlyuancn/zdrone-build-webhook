package plugin

import (
    "context"

    "github.com/drone/drone-go/drone"
    "github.com/drone/drone-go/plugin/webhook"

    "github.com/zlyuancn/zdrone-build-webhook/logger"
    "github.com/zlyuancn/zdrone-build-webhook/message"
    "github.com/zlyuancn/zdrone-build-webhook/notifer"
)

type plugin struct {
}

func New() webhook.Plugin {
    return &plugin{}
}

func (p *plugin) Deliver(ctx context.Context, req *webhook.Request) error {
    if !Check(req) {
        return nil
    }

    msg, err := message.MakeMsg(req)
    if err != nil {
        logger.Error(err)
    }

    notifer.Notify(msg)
    return nil
}

func Check(req *webhook.Request) bool {
    if req.Event != webhook.EventBuild {
        return false
    }

    switch req.Action {
    case webhook.ActionCreated:
        switch req.Build.Status {
        case drone.StatusPending:
            return true
        case drone.StatusRunning:
            return true
        }
    case webhook.ActionUpdated:
        switch req.Build.Status {
        case drone.StatusPassing:
            return true
        case drone.StatusFailing:
            return true
        case drone.StatusKilled:
            return true
        case drone.StatusError:
            return true
        }
    }
    return false
}
