# Job 方法

- `CollectionName() string`
  - 返回：Mongo 集合名，固定为 `job`
- `GenerateApplication() *appv1.Application`
  - 输入：无（使用 job 内部字段与全局 `manifestRepo`）
  - 返回：构造好的 Argo CD Application 对象（未提交）
  - Application 
    - 关闭自动sync
    - 使用plugin的方式，名称是plugin
    - parameters 包含 manifest-id，job-id， env

## 相关链接
- `model/job.go`
