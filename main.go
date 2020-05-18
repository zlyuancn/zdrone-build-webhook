package main

import (
    "net/http"

    "github.com/drone/drone-go/plugin/webhook"
    "github.com/zlyuancn/zsignal"

    _ "github.com/joho/godotenv/autoload"
    "github.com/sirupsen/logrus"

    "github.com/zlyuancn/zdrone-build-webhook/config"
    "github.com/zlyuancn/zdrone-build-webhook/logger"
    "github.com/zlyuancn/zdrone-build-webhook/notifer"
    "github.com/zlyuancn/zdrone-build-webhook/plugin"
)

func main() {
    defer zsignal.Shutdown()

    config.Init()
    logger.Init()
    notifer.Init()

    handler := webhook.Handler(
        plugin.New(),
        config.Config.Secret,
        logrus.StandardLogger(),
    )

    http.Handle("/", handler)
    logger.Infof("服务启动: %s", config.Config.Bind)

    server := &http.Server{Addr: config.Config.Bind}
    zsignal.RegisterOnShutdown(func() {
        _ = server.Close()
    })
    if err := server.ListenAndServe(); err != nil && err == http.ErrServerClosed {
        logger.Fatal(err)
    }
}
