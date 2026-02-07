# API 文档

## 目的
- 描述公共库对外暴露的 Go 包与关键接口

## 范围
- `model/*` 与 `client/*` 的主要类型与函数

## 要点
- 以 Go 包为单位的使用入口
- 初始化函数与依赖项
- 稳定字段与枚举

## model 包
核心结构：
- `Application` / `Manifest` / `Job`
- `BaseModel`（通用 Mongo 字段）
- `Service` / `Port` / `ConfigMap` / `EnvVar`
- 配置：`Config` / `MongoConfig` / `OtelConfig` / `Consul` / `Repo`

### model 方法与规范
#### Application
- `CollectionName() string`  
  返回：Mongo 集合名，固定为 `applications`。

#### Manifest
- `CollectionName() string`  
  返回：Mongo 集合名，固定为 `manifests`。
- `GetStep(taskName string) *ManifestStep`  
  输入：`taskName` 任务名。  
  返回：匹配的步骤指针，未找到返回 `nil`。
- `GeneratePipelineRun(pipelineName string, pvc string) *tknv1.PipelineRun`  
  输入：`pipelineName` Pipeline 名称；`pvc` PVC 名称。  
  返回：构造好的 Tekton PipelineRun 对象（未提交）。
- `GeneratePipelineRunParams() []tknv1.Param`  
  输入：无（使用 manifest 内部字段）。  
  返回：PipelineRun 参数列表（包含 `manifest-id` / `git-url` / `git-revision` / `image-registry` / `name` / `image-tag` / `manifest-name`）。
- `GenerateManifestVersion(name string) string`  
  输入：`name` 应用或前缀。  
  返回：`{name}{YYYYMMDDhhmmss}{0-99随机数}` 格式版本名。
- `PatchManifestRequest.IsEmpty() bool`  
  输入：无。  
  返回：`commit_hash` 与 `digest` 均为空时返回 `true`。

#### Job
- `CollectionName() string`  
  返回：Mongo 集合名，固定为 `job`。
- `GenerateApplication() *appv1.Application`  
  输入：无（使用 job 内部字段与全局 `manifestRepo`）。  
  返回：构造好的 Argo CD Application 对象（未提交）。

#### BaseModel
- `GetID() primitive.ObjectID` / `SetID(id primitive.ObjectID)`  
  用途：读写 Mongo `_id`。
- `WithCreateDefault()` / `WithUpdateDefault()`  
  用途：写入 `created_at` / `updated_at` 时间戳。

## client 包
### client/mongo
用途：MongoDB 连接与通用仓储。
入口：
- `InitMongo(ctx, config, logger) (*mongo.Client, error)`
- `Repo`（全局 Repository）
- `Repository` 常用方法：`Create` / `FindByID` / `Update` / `Delete` / `List` / `UpdateOne` / `UpdateMany` / `FindOne` / `Upsert` / `UpdateByID`

方法参数与返回规范：
- `InitMongo(ctx, config, logger)`  
  输入：`config.URI` / `config.DBName`，`logger` 可用即可。  
  返回：Mongo client 或错误；初始化成功会设置全局 `Repo`。
- `Create(ctx, m)`  
  输入：`m` 需实现 `MongoModel`；若 `id` 为空会自动生成。  
  返回：错误（若有）。
- `FindByID(ctx, m, id)` / `FindOne(ctx, m, filter)`  
  输入：目标模型指针 `m`、查询条件。  
  返回：错误；成功会填充 `m`。
- `Update(ctx, m)`  
  输入：带 `id` 的模型；`$set` 整体更新。  
  返回：错误。
- `UpdateOne(ctx, m, filter, update)` / `UpdateMany(ctx, m, filter, update)` / `Upsert(ctx, m, filter, update)` / `UpdateByID(ctx, m, id, update)`  
  输入：Mongo 更新语句 `update`。  
  返回：错误；`UpdateOne/ByID` 在未匹配时会打日志。
- `Delete(ctx, m, id)`  
  输入：`id`；逻辑删除（设置 `deleted` 标志）。  
  返回：错误。
- `List(ctx, m, filter, results)`  
  输入：`results` 必须是 slice 指针。  
  返回：错误；结果按 `created_at` 倒序。

### client/argo
用途：Argo CD Application 创建与更新。  
入口：
- `InitArgoCdClient(config)`  
- `CreateApplication(ctx, app)`  
- `UpdateApplication(ctx, app)`

方法参数与返回规范：
- `InitArgoCdClient(config)`  
  输入：K8s `rest.Config`。  
  返回：错误；成功后设置全局 `ArgoCdClient`。
- `CreateApplication(ctx, app)` / `UpdateApplication(ctx, app)`  
  输入：完整 `appv1.Application`。  
  返回：错误；`Update` 会先 GET 再替换 `Spec` 保持 `ResourceVersion`。

### client/tekton
用途：Tekton Pipeline/PipelineRun 操作与 PVC 维护。  
入口：
- `InitTektonClient(ctx, config, logger)`  
- `GetPipeline(ctx, namespace, name)`  
- `CreatePipelineRun(ctx, namespace, pr)`  
- `CreatePVC(ctx, namespace, pvcName, storageClass, size)`  
- `PatchPVCOwner(ctx, pvc, pr)`

方法参数与返回规范：
- `InitTektonClient(ctx, config, logger)`  
  输入：K8s `rest.Config`、logger。  
  返回：错误；成功后设置全局 `TektonClient`/`KubeClient`。
- `GetPipeline(ctx, namespace, name)`  
  返回：Pipeline 或错误。
- `CreatePipelineRun(ctx, namespace, pr)`  
  输入：已构造 `PipelineRun`。  
  返回：创建后的对象或错误。
- `CreatePVC(ctx, namespace, pvcName, storageClass, size)`  
  输入：`size` 形如 `5Gi`。  
  返回：PVC 或错误。
- `PatchPVCOwner(ctx, pvc, pr)`  
  输入：PVC 与 PipelineRun。  
  返回：错误；为 PVC 追加 OwnerReference。

### client/otel
用途：Tracing/metrics 初始化。  
入口：
- `InitOtel(ctx, config) (shutdown, error)`  
- `InitMetricProvider()`  
- `Start(ctx, tracerName, spanName, opts...)`

方法参数与返回规范：
- `InitOtel(ctx, config)`  
  输入：`config.Endpoint` / `config.ServiceName` 必填。  
  返回：`shutdown(ctx)` 回调与错误。
- `InitMetricProvider()`  
  返回：错误；成功后设置全局 MeterProvider。
- `Start(ctx, tracerName, spanName, opts...)`  
  返回：新的 `context` 与 `Span`。

### client/logging
用途：Zap logger 初始化与上下文注入。  
入口：
- `InitZapLogger(config)`  
- `InjectLogger(ctx, base)`  
- `LoggerFromContext(ctx)`  
- `NewZapAdapter(logger)`

方法参数与返回规范：
- `InitZapLogger(config)`  
  输入：`config.Level`/`config.Format`。  
  返回：无；初始化全局 `Logger`。
- `InjectLogger(ctx, base)`  
  输入：`ctx`、可选 `base`。  
  返回：注入 trace 信息后的新 `ctx`。
- `LoggerFromContext(ctx)`  
  返回：存在则取出，否则返回全局 `Logger`。
- `NewZapAdapter(logger)`  
  返回：`ZapAdapter`（供第三方库使用）。

### client/consul
用途：Consul 配置读取与监听。  
入口：
- `InitConsulClient(consul)`  
- `LoadConsulConfigAndMerge(consul)`  
- `WatchConsul(consul, logger)`

方法参数与返回规范：
- `InitConsulClient(consul)`  
  输入：`consul.Address` 必填。  
  返回：错误；成功后设置全局 `ConsulClient`。
- `LoadConsulConfigAndMerge(consul)`  
  输入：`consul.Key` 作为 KV 路径。  
  返回：错误；成功后合并到全局 `model.C`。
- `WatchConsul(consul, logger)`  
  输入：Consul 配置与 logger。  
  返回：无；内部 goroutine 监听并更新全局配置。

### client/pyroscope
用途：性能剖析上报。  
入口：
- `InitPyroscope(name, address)`

方法参数与返回规范：
- `InitPyroscope(name, address)`  
  输入：应用名与 Pyroscope 地址。  
  返回：无；启动 profiler。

## 版本管理策略
- Go Module 语义化版本（tag）对外发布
- 结构字段新增向后兼容，字段删除需大版本

## 相关链接
- `docs/06-DEV-GUIDE.md`
