# 不再维护的项目, 请转到 https://github.com/zlyuancn/drone-build-notify

# 用于接收drone构建信息的webhook, 并将构建信息推送到钉钉

# 更新你的 Drone 服务配置

```text
DRONE_WEBHOOK_ENDPOINT=your_drone_url
DRONE_WEBHOOK_SECRET=your_webhook_secret
```

# 秘钥创建方法

```console
$ openssl rand -hex 16
```

# docker-compose
```yaml
version: '3'
services:
  zdrone-build-webhook:
    image: zlyuan/zdrone-build-webhook
    restart: unless-stopped
    ports:
      - 80:80
    volumes:
      - /usr/share/zoneinfo/Asia/Shanghai:/etc/localtime
    environment:
      - DRONE_SECRET={{your_webhook_secret}}
      - DRONE_SERVER={{your_drone_url}}
      - NOTIFER=dingtalk
      - DINGTALK_ACCESSTOKEN={{your_dingtalk_access_token}}
      - DINGTALK_SECRET={{your_dingtalk_secret}}
    container_name: zdrone-build-webhook
```

# 环境变量说明

变量名|类型|默认值|说明|示例
--|:-:|:-:|:-:|:-:
DRONE_BIND|string|:80|服务监听地址|:80
DRONE_DEBUG|bool|false|调试模式,会输出额外信息|false
*DRONE_SECRET|string||webhook秘钥|
*DRONE_SERVER|string||drone服务地址|
LOG_PATH|string||日志文件输出目录,不需要预先创建|/var/log/zdrone-build-webhook
NOTIFER|string||通告者,多个通告者用半角逗号隔开|dingtalk
NOTIFY_RETRY|int|2|通告失败重试次数|2
OFF_CREATE_NOTIFY|bool|false|关闭创建动作的通告|false
DINGTALK_ACCESSTOKEN|string||dingtalk通告者的access_token|
DINGTALK_SECRET|string||dingtalk通告者的secret|
DINGTALK_START_TEMPLATE|string|conf/dingtask_start_template.md|钉钉消息任务开始模板文件|
DINGTALK_END_TEMPLATE|string|conf/dingtask_end_template.md|钉钉消息任务结束模板文件|

# 模板语法

```
{{变量名}}
{{ 变量名 }}
```

变量名 | 描述
---|:-:
task_num      |  任务号
task_url      |  任务跳转url
repo_name     |  仓库名
branch        |  分支名
repo_url      |  仓库地址, 转到该分支
auther        |  操作人员
auther_email  |  操作人员邮箱
auther_avatar |  操作人员头像
status        |  执行结果
status_desc   |  执行结果描述
status_pic_url|  执行结果图片url
start_time    |  开始时间
end_time      |  结束时间
process_time  |  处理时间
commit_msg    |  提交信息
commit_id     |  提交id
commit_url    |  提交信息的跳转url
