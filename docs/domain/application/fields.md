# Application 字段

## 核心字段
- `name` 应用名
- `project_name` 项目/命名空间
- `repo_url` 源码仓库地址
- `active_manifest_name` 当前生效 manifest 名称
- `active_manifest_id` 当前生效 manifest ID
- `replica` 默认副本数
- `type` 发布类型：`normal` / `canary` / `blue-green`
- `config_maps` 配置映射（见 ConfigMap）
- `service` 服务端口（见 Service/Port）
- `internet` 访问类型：`internal` / `external`
- `envs` 环境变量：`map[string][]EnvVar`
- `status` 当前状态（来自 Job）：`Running` / `Failed` / `Degraded`

## 相关链接
- `docs/03-DOMAIN-MODEL.md`
