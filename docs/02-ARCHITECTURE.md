# 架构说明

## 目的
- 描述公共库的模块边界与依赖关系

## 范围
- Go 包、外部依赖、运行时约束

## 要点
- 模块分层：`model`（结构定义）与 `client`（外部系统适配）
- 公共约束：MongoDB 持久化模型、K8s 生态（ArgoCD/Tekton）
- 初始化顺序：日志/配置 → 客户端 → 领域对象

## 模块划分
### model
职责：定义 application / manifest / job 等核心结构与通用类型。  
依赖：MongoDB bson、K8s 类型（仅在 manifest/job 的方法里引用）。

### client
职责：封装外部系统访问，提供初始化与常用操作。  
子模块：
- `client/mongo`：MongoDB 连接与通用 Repository
- `client/argo`：Argo CD Application 创建/更新
- `client/tekton`：Tekton Pipeline/PipelineRun 访问与 PVC 操作
- `client/otel`：OpenTelemetry tracing/metrics 初始化
- `client/logging`：Zap 日志初始化与上下文注入
- `client/consul`：Consul 配置加载与监听
- `client/pyroscope`：性能剖析上报

## 依赖关系
```
model
  ↘ (被) client/*
client/*
  ↘ 外部系统 SDK（MongoDB / Argo CD / Tekton / OTEL / Consul / Pyroscope）
```

## 运行时约束
- `client/mongo` 需要 MongoDB 可用且 `model.MongoConfig` 正确
- `client/argo` / `client/tekton` 需要 K8s REST config
- `client/otel` 需要 OTEL endpoint 与 service name
- `client/consul` 需要 Consul Address 与 Key

## 相关链接
- `docs/03-DOMAIN-MODEL.md`
- `docs/04-API.md`
