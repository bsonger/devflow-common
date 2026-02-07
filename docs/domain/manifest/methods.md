# Manifest 方法

- `CollectionName() string`
  - 返回：Mongo 集合名，固定为 `manifests`
- `GetStep(taskName string) *ManifestStep`
  - 输入：`taskName` 任务名
  - 返回：匹配步骤指针，未找到返回 `nil`
- `GeneratePipelineRun(pipelineName string, pvc string) *tknv1.PipelineRun`
  - 输入：`pipelineName` Pipeline 名称；`pvc` PVC 名称
  - 返回：构造好的 Tekton PipelineRun 对象（未提交）
- `GeneratePipelineRunParams() []tknv1.Param`
  - 输入：无（使用 manifest 内部字段）
  - 返回：PipelineRun 参数列表
- `GenerateManifestVersion(name string) string`
  - 输入：`name` 应用或前缀
  - 返回：`{name}{YYYYMMDDhhmmss}{0-99随机数}` 版本名
- `PatchManifestRequest.IsEmpty() bool`
  - 输入：无
  - 返回：`commit_hash` 与 `digest` 均为空时为 `true`

## 相关链接
- `model/manifest.go`
