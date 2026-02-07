# 开发指南

## 目的
- 新人快速上手公共库开发

## 范围
- 依赖、构建、使用方式

## 要点
- Go 版本与依赖管理
- 本地验证与发布流程
- 目录约定

## 环境要求
- Go 1.20+（以 `go.mod` 为准）

## 本地开发
- 拉取依赖：`go mod tidy`
- 运行测试：`go test ./...`
- 格式化：`gofmt -w .`

## 目录约定
- `model/` 领域结构与配置
- `client/` 外部系统客户端封装
- `docs/` 文档

## 使用方式
在下游服务中引入模块并按需初始化：
- 日志：`client/logging`
- 观测：`client/otel` / `client/pyroscope`
- 外部系统：`client/mongo` / `client/argo` / `client/tekton` / `client/consul`
