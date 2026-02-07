# Job 字段

## 核心字段
- `application_id` 关联 Application
- `application_name` 应用名
- `project_name` 项目/命名空间
- `manifest_id` 关联 Manifest
- `manifest_name` manifest 名称
- `type` 任务类型：`Install` / `Upgrade` / `Rollback`
- `env` 环境
- `status` 状态：`Pending` / `Running` / `Succeeded` / `Failed` / `RollingBack` / `RolledBack` / `Syncing` / `SyncFailed`

## 相关链接
- `docs/03-DOMAIN-MODEL.md`
