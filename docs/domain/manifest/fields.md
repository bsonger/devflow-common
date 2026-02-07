# Manifest 字段

## 核心字段
- `application_id` 关联 Application
- `name` manifest 名称/版本
- `application_name` 应用名
- `branch` Git 分支
- `git_repo` 仓库地址
- `commit_hash` 提交哈希
- `replica` 副本数
- `digest` 镜像 digest（如 `sha256:...`）
- `type` 发布类型：`normal` / `canary` / `blue-green`
- `config_maps` 配置映射（见 ConfigMap）
- `service` 服务端口（见 Service/Port）
- `internet` 访问类型：`internal` / `external`
- `envs` 环境变量：`map[string][]EnvVar`
- `pipeline_id` Tekton PipelineRun ID
- `steps` 步骤状态数组（见 ManifestStep）
- `status` 状态：`Pending` / `Running` / `Succeeded` / `Failed`

## 步骤状态（ManifestStep）
- `task_name` 任务名
- `task_run` Tekton TaskRun 名
- `status`：`Pending` / `Running` / `Succeeded` / `Failed`
- `start_time` / `end_time`
- `message` 失败/提示信息

## 相关链接
- `docs/03-DOMAIN-MODEL.md`
