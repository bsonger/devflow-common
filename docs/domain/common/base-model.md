# BaseModel

## 字段
- `id`（Mongo `_id`）
- `created_at` / `updated_at`
- `deleted_at`（软删）

## 方法
- `GetID()` / `SetID(id)` 获取/设置 Mongo ID
- `WithCreateDefault()` / `WithUpdateDefault()` 写入时间戳

## 相关链接
- `model/base.go`
