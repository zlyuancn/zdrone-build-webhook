/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/5/16
   Description :
-------------------------------------------------
*/

package message

import (
    "errors"
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "time"

    "github.com/drone/drone-go/drone"
    "github.com/drone/drone-go/plugin/webhook"

    "github.com/zlyuancn/zdrone-build-webhook/config"
)

const TimeLayout = "2006-01-02 15:04:05"

type Msg struct {
    TaskNum  string `json:"task_num"`  // 任务号
    TaskUrl  string `json:"task_url"`  // 任务跳转url
    RepoName string `json:"repo_name"` // 仓库名
    Branch   string `json:"branch"`    // 分支名
    RepoUrl  string `json:"repo_url"`  // 仓库地址, 转到该分支

    Auther       string `json:"auther"`        // 操作人员
    AutherEmail  string `json:"auther_email"`  // 操作人员邮箱
    AutherAvatar string `json:"auther_avatar"` // 操作人员头像

    Status       string `json:"status"`         // 执行结果
    StatusDesc   string `json:"status_desc"`    // 执行结果描述
    StatusPicUrl string `json:"status_pic_url"` // 执行结果图片url

    StartTime   string `json:"start_time"`   // 开始时间
    EndTime     string `json:"end_time"`     // 结束时间
    ProcessTime string `json:"process_time"` // 处理时间

    CommitMsg string `json:"commit_msg"` // 提交信息
    CommitId  string `json:"commit_id"`  // 提交id
    CommitUrl string `json:"commit_url"` // 提交信息的跳转url

    template_values map[string]string
}

func (m *Msg) Get(key string) string {
    if len(m.template_values) == 0 {
        msg_type := reflect.TypeOf(m).Elem()
        msg_val := reflect.ValueOf(m).Elem()

        field_count := msg_type.NumField()
        m.template_values = make(map[string]string, field_count)
        for i := 0; i < field_count; i++ {
            field := msg_type.Field(i)
            if field.PkgPath != "" {
                continue
            }

            k := field.Tag.Get("template")
            if k == "" {
                k = field.Tag.Get("json")
            }
            if k == "" {
                k = field.Name
            }
            v := msg_val.Field(i).String()
            m.template_values[k] = v
        }
    }

    if v, ok := m.template_values[key]; ok {
        return v
    }
    return fmt.Sprintf("{ %s: undefined }", key)
}

func MakeMsg(req *webhook.Request) (*Msg, error) {
    repo := req.Repo
    build := req.Build
    if repo == nil {
        return nil, errors.New("没有 repo 信息")
    }
    if build == nil {
        return nil, errors.New("没有 build 信息")
    }

    repo_url := repo.HTTPURL
    if build.Source != repo.Branch {
        repo_url = makeBranchUrl(repo.HTTPURL, build.Source)
    }

    start_time := time.Unix(build.Created, 0).Format(TimeLayout)
    end_time := ""
    process_time := ""
    if build.Finished > 0 {
        end_time = time.Unix(build.Finished, 0).Format(TimeLayout)
        process_time = (time.Duration(build.Finished-build.Created) * time.Second).String()
    }

    status := "start"
    status_desc := "开始"
    if req.Action != webhook.ActionCreated {
        switch req.Build.Status {
        case drone.StatusPassing:
            status = "success"
            status_desc = "完成"
        case drone.StatusFailing:
            status = "failure"
            status_desc = "失败"
        case drone.StatusKilled:
            status = "killed"
            status_desc = "删除"
        case drone.StatusError:
            status = "error"
            status_desc = "错误"
        }
    }

    return &Msg{
        TaskNum:      strconv.FormatInt(build.Number, 10),
        RepoName:     repo.Slug,
        Branch:       build.Source,
        RepoUrl:      repo_url,
        Auther:       build.AuthorName,
        AutherEmail:  build.AuthorEmail,
        AutherAvatar: build.AuthorAvatar,
        Status:       status,
        StatusDesc:   status_desc,
        StatusPicUrl: makeStatusPicUrl(status),
        StartTime:    start_time,
        EndTime:      end_time,
        ProcessTime:  process_time,
        CommitMsg:    strings.TrimSpace(build.Message),
        CommitId:     build.After,
        CommitUrl:    makeCommitUrl(repo.HTTPURL, build.After),
        TaskUrl:      makeTaskUrl(repo.Slug, build.Number),
    }, nil
}

// 构建储存库基础url(储存库的url地址)
func makeRepoBaseUrl(repo_url string) string {
    return strings.TrimSuffix(repo_url, ".git")
}

// 构建分支url
func makeBranchUrl(repo_url, branch string) string {
    repo_base_url := makeRepoBaseUrl(repo_url)
    if strings.Contains(repo_base_url, "//github.com/") {
        return fmt.Sprintf("%s/tree/%s", repo_base_url, branch)
    }
    if strings.Contains(repo_base_url, "//gitee.com/") {
        return fmt.Sprintf("%s/tree/%s", repo_base_url, branch)
    }

    return fmt.Sprintf("%s/src/%s", repo_base_url, branch)
}

// 构建资源url
func makeResUrl(repo_url, branch, res string) string {
    repo_base_url := makeRepoBaseUrl(repo_url)
    if strings.Contains(repo_base_url, "//github.com/") {
        raw_base_url := strings.ReplaceAll(repo_base_url, "//github.com/", "//raw.githubusercontent.com/")
        return fmt.Sprintf("%s/%s/%s", raw_base_url, branch, res)
    }
    return fmt.Sprintf("%s/raw/%s/%s", repo_base_url, branch, res)
}

// 构建commit跳转url
func makeCommitUrl(repo_url, commit_id string) string {
    repo_base_url := makeRepoBaseUrl(repo_url)
    return fmt.Sprintf("%s/commit/%s", repo_base_url, commit_id)
}

// 构建任务跳转url
func makeTaskUrl(slug string, task_id int64) string {
    return fmt.Sprintf("%s/%s/%d", config.Config.DroneServer, slug, task_id)
}

// 构建状态图片url
func makeStatusPicUrl(status string) string {
    return makeResUrl(config.MyPicRepoUrl, config.MyPicRepoBranch, fmt.Sprintf("assets/%s.png", status))
}
