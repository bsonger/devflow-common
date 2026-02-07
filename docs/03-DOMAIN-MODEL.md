# 领域模型

## 目的
- 定义核心概念与关系

## 范围
- 领域对象、状态、约束

## 要点
- 领域对象列表
  - application：应用元信息与部署相关配置
  - manifest：一次发布/构建产物的声明与状态
  - job：发布任务与执行状态
- 状态机或生命周期
- 关键约束与不变量

## 领域对象
- Application：`docs/domain/application/README.md`
- Manifest：`docs/domain/manifest/README.md`
- Job：`docs/domain/job/README.md`

## 公共结构
- `docs/domain/common/README.md`

## 配置结构
- `docs/domain/config/README.md`

## 枚举与常量
- `docs/domain/enums/README.md`

## 相关链接
- `model/application.go`
- `model/manifest.go`
- `model/job.go`
