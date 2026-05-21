# API Key 对外创建与删除对接方案

本文档说明外部系统如何通过 new-api 的管理接口创建、使用和删除 API Key。

## 1. 对接目标

外部系统通过用户级管理凭证调用 new-api：

1. 获取或配置用户的 `access_token`。
2. 调用创建接口生成 API Key。
3. 使用生成的 API Key 调用 AI 中继接口。
4. 在用户解绑、过期、风控或主动撤销时删除 API Key。

## 2. 基础信息

### 2.1 接口地址

当前示例使用的服务地址：

```text
http://192.168.30.91:8081
```

### 2.2 认证方式

API Key 管理接口使用 `UserAuth`，外部系统调用时推荐使用：

```http
Authorization: Bearer {access_token}
New-Api-User: {user_id}
Content-Type: application/json
```

说明：

- `access_token` 是用户级管理凭证，不是 AI 中继使用的 `sk-xxx`。
- `New-Api-User` 必须等于该 `access_token` 所属用户 ID。
- `Authorization` 也兼容直接传 `{access_token}`，但建议统一使用 `Bearer {access_token}`。
- 创建出来的 API Key 才用于调用 AI 中继接口，格式为 `Bearer sk-{key}`。

## 3. 获取 Access Token

### 3.1 登录获取用户 ID 和 Session

如果外部系统没有可用的 `access_token`，可先通过用户名密码登录获取 Session。

```bash
curl -X POST "http://192.168.30.91:8081/api/user/login" \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "username": "user@example.com",
    "password": "password"
  }'
```

响应示例：

```json
{
  "success": true,
  "message": "",
  "data": {
    "id": 1,
    "username": "user@example.com",
    "display_name": "User",
    "role": 1,
    "status": 1,
    "group": "default"
  }
}
```

外部系统需要保存响应里的 `data.id`，后续作为 `New-Api-User`。

### 3.2 生成 Access Token

```bash
curl -X GET "http://192.168.30.91:8081/api/user/token" \
  -b cookies.txt \
  -H "New-Api-User: 1"
```

响应示例：

```json
{
  "success": true,
  "message": "",
  "data": "access-token-value"
}
```

注意：

- `GET /api/user/token` 会为当前用户生成新的 `access_token`。
- 新生成后，旧的 `access_token` 会被覆盖。
- 建议只在初始化绑定或重置凭证时调用，不要每次创建 API Key 都调用。

## 4. 创建 API Key

### 4.1 接口

```http
POST /api/token/
```

### 4.2 Header

```http
Authorization: Bearer {access_token}
New-Api-User: {user_id}
Content-Type: application/json
```

### 4.3 请求体

```json
{
  "name": "external-app-key",
  "remain_quota": 0,
  "expired_time": -1,
  "unlimited_quota": true,
  "model_limits_enabled": false,
  "model_limits": "",
  "allow_ips": "",
  "group": "",
  "cross_group_retry": false
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `name` | string | 是 | API Key 名称，最长 50 字符。 |
| `remain_quota` | number | 是 | 剩余额度，`unlimited_quota=false` 时生效。 |
| `expired_time` | number | 是 | Unix 秒级时间戳，`-1` 表示永不过期。 |
| `unlimited_quota` | boolean | 是 | 是否无限额度。 |
| `model_limits_enabled` | boolean | 是 | 是否启用模型限制。 |
| `model_limits` | string | 否 | 模型限制配置，空字符串表示不限制。 |
| `allow_ips` | string | 否 | IP 白名单，空字符串表示不限制。 |
| `group` | string | 否 | 使用的分组，空字符串表示默认分组。 |
| `cross_group_retry` | boolean | 否 | auto 分组下是否允许跨组重试。 |

### 4.4 成功响应

```json
{
  "success": true,
  "message": "",
  "data": {
    "id": 123,
    "name": "external-app-key",
    "key": "raw-token-key"
  }
}
```

说明：

- `data.id` 是 API Key 的内部 ID，删除或管理时需要保存。
- `data.key` 是原始 key，不包含 `sk-` 前缀。
- 调用 AI 中继接口时需要使用 `sk-{key}`。

### 4.5 curl 示例

```bash
curl -X POST "http://192.168.30.91:8081/api/token/" \
  -H "Authorization: Bearer access-token-value" \
  -H "New-Api-User: 1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "external-app-key",
    "remain_quota": 0,
    "expired_time": -1,
    "unlimited_quota": true,
    "model_limits_enabled": false,
    "model_limits": "",
    "allow_ips": "",
    "group": "",
    "cross_group_retry": false
  }'
```

## 5. 使用 API Key 调用 AI 中继接口

创建接口返回：

```json
{
  "key": "raw-token-key"
}
```

调用中继接口时需要拼接为：

```http
Authorization: Bearer sk-raw-token-key
```

示例：

```bash
curl -X POST "http://192.168.30.91:8081/v1/chat/completions" \
  -H "Authorization: Bearer sk-raw-token-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [
      {
        "role": "user",
        "content": "Hello"
      }
    ]
  }'
```

## 6. 查询 API Key 列表

删除 API Key 前，外部系统通常应保存创建接口返回的 `id`。如果没有保存，可通过列表接口查询。

### 6.1 接口

```http
GET /api/token/?p=1&size=10
```

### 6.2 Header

```http
Authorization: Bearer {access_token}
New-Api-User: {user_id}
```

### 6.3 curl 示例

```bash
curl -X GET "http://192.168.30.91:8081/api/token/?p=1&size=10" \
  -H "Authorization: Bearer access-token-value" \
  -H "New-Api-User: 1"
```

注意：

- 列表接口返回的 `key` 是脱敏后的值。
- 完整 key 只应在创建时保存；如需再次获取，可使用现有 `POST /api/token/{id}/key` 接口。

## 7. 删除 API Key

### 7.1 接口

```http
DELETE /api/token/{id}
```

其中 `{id}` 是创建 API Key 时返回的 `data.id`。

### 7.2 Header

```http
Authorization: Bearer {access_token}
New-Api-User: {user_id}
```

### 7.3 请求体

无需请求体。

### 7.4 成功响应

```json
{
  "success": true,
  "message": ""
}
```

### 7.5 curl 示例

```bash
curl -X DELETE "http://192.168.30.91:8081/api/token/123" \
  -H "Authorization: Bearer access-token-value" \
  -H "New-Api-User: 1"
```

### 7.6 删除策略建议

建议在以下场景删除 API Key：

- 用户主动解绑第三方系统。
- 用户套餐到期且不再允许调用。
- 外部系统检测到 key 泄露风险。
- 用户账号被禁用或风控命中。
- 外部系统重新发放新 API Key 后，废弃旧 API Key。

如果只是临时停用，也可以使用更新接口将 token 状态置为禁用；删除适用于明确不再使用的 API Key。

## 8. 错误响应与处理

管理接口通常使用以下业务响应结构：

```json
{
  "success": false,
  "message": "error message"
}
```

常见错误：

| 场景 | 可能原因 | 处理建议 |
| --- | --- | --- |
| `401 Unauthorized` | 未登录、缺少 `Authorization`、缺少 `New-Api-User`。 | 检查 Header 是否完整。 |
| `success=false` 且提示 Access Token 无效 | `access_token` 不存在或已被重新生成覆盖。 | 重新生成并更新外部系统保存的凭证。 |
| `success=false` 且提示用户 ID 不匹配 | `New-Api-User` 与 `access_token` 所属用户不一致。 | 使用正确的用户 ID。 |
| `success=false` 且提示令牌数量达到上限 | 用户 API Key 数量超过系统配置。 | 删除无用 key 或调整系统配置。 |
| `success=false` 且提示额度非法 | `remain_quota` 小于 0 或超过上限。 | 修正请求参数。 |
| `403 Forbidden` | 用户被禁用或权限不足。 | 检查用户状态和角色。 |

## 9. 外部系统数据存储建议

外部系统建议保存以下字段：

| 字段 | 说明 |
| --- | --- |
| `user_id` | new-api 用户 ID，用于 `New-Api-User`。 |
| `access_token` | 管理凭证，用于创建和删除 API Key。 |
| `api_key_id` | 创建接口返回的 `data.id`，用于删除。 |
| `api_key` | 建议保存为 `sk-{key}`，用于调用 AI 中继接口。 |
| `api_key_name` | 创建接口返回的 `data.name`。 |
| `created_at` | 外部系统记录的创建时间。 |
| `expired_time` | 与创建请求保持一致。 |

安全要求：

- `access_token` 和 `api_key` 必须加密存储。
- 不要把 `access_token` 下发到浏览器、移动端或第三方客户端。
- 创建接口返回的完整 key 建议只展示一次。
- 日志中不要打印完整 `access_token` 或完整 API Key。

## 10. 推荐完整流程

### 10.1 初始化绑定

1. 用户在 new-api 登录。
2. 外部系统获取用户 ID。
3. 调用 `GET /api/user/token` 生成 `access_token`。
4. 外部系统保存 `user_id` 和 `access_token`。

### 10.2 创建 API Key

1. 外部系统调用 `POST /api/token/`。
2. 保存响应里的 `data.id`。
3. 保存 `sk-{data.key}` 用于后续 AI 调用。
4. 将 key 与外部系统用户、业务空间或应用绑定。

### 10.3 调用 AI 接口

1. 从外部系统读取加密存储的 `sk-{key}`。
2. 调用 `/v1/chat/completions` 等中继接口。
3. 根据业务需要记录用量和错误信息。

### 10.4 删除 API Key

1. 外部系统读取保存的 `api_key_id`。
2. 调用 `DELETE /api/token/{id}`。
3. 删除或标记外部系统保存的 key。
4. 后续请求不再使用该 key。

## 11. Node.js 示例

```js
const baseURL = 'http://192.168.30.91:8081'

async function createApiKey({ accessToken, userId, name }) {
  const response = await fetch(`${baseURL}/api/token/`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${accessToken}`,
      'New-Api-User': String(userId),
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      name,
      remain_quota: 0,
      expired_time: -1,
      unlimited_quota: true,
      model_limits_enabled: false,
      model_limits: '',
      allow_ips: '',
      group: '',
      cross_group_retry: false,
    }),
  })

  const result = await response.json()
  if (!result.success) {
    throw new Error(result.message || 'Create API key failed')
  }

  return {
    id: result.data.id,
    name: result.data.name,
    apiKey: `sk-${result.data.key}`,
  }
}

async function deleteApiKey({ accessToken, userId, apiKeyId }) {
  const response = await fetch(`${baseURL}/api/token/${apiKeyId}`, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${accessToken}`,
      'New-Api-User': String(userId),
    },
  })

  const result = await response.json()
  if (!result.success) {
    throw new Error(result.message || 'Delete API key failed')
  }

  return true
}
```

## 12. 上线检查清单

- 已确认外部系统保存了正确的 `user_id`。
- 已确认 `access_token` 不会暴露到客户端。
- 已确认创建接口响应中的 `key` 会加 `sk-` 后再用于中继请求。
- 已确认外部系统保存了 `api_key_id`，可用于删除。
- 已确认解绑、风控、过期等场景会触发删除或禁用。
- 已确认日志不会输出完整凭证。
- 已确认外部系统具备重试和错误告警机制。
