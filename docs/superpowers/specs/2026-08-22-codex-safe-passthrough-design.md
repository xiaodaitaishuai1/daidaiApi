# Codex Safe Passthrough Design

## Goal

让 Codex 频道把下游 `/v1/responses` 请求的原始 body 和非敏感客户端请求头尽可能原样发送到 Codex 上游，同时继续使用频道 OAuth 凭据完成上游认证。

## Architecture

保留现有 Router -> ResponsesHelper -> Codex Adaptor -> shared API request 链路。仅在 Codex responses 模式下改变 body/header 的构造：body 默认读取已缓存的原始请求体；header 从下游复制允许转发的请求头，再由 Codex adaptor 覆盖 `Authorization`、`chatgpt-account-id`、目标媒体类型和流式 Accept。Host、Content-Length 及 hop-by-hop 头不复制。

## Data Flow

1. `ResponsesHelper` 从 `BodyStorage` 读取原始 body，不执行 DTO 重编码、禁用字段删除或参数覆盖。
2. `Codex.Adaptor.SetupRequestHeader` 复制下游非敏感、非 hop-by-hop headers。
3. 频道 OAuth key 解析后覆盖 `Authorization` 与 `chatgpt-account-id`。
4. HTTP client 根据上游 URL 和 body 重新计算 Host、Content-Length/Transfer-Encoding。

## Security Boundaries

- 永不透传下游 `Authorization`、API key headers、Cookie、Host 或 hop-by-hop headers。
- `chatgpt-account-id` 必须来自与 OAuth access token 配对的频道凭据。
- `Content-Type` 固定为 `application/json`，流式请求的 `Accept` 固定为 `text/event-stream`，避免上游严格校验失败。

## Testing

- 验证 Codex responses 默认发送完全相同的 body。
- 验证普通客户端 headers 会透传，敏感/连接类 headers 不会透传。
- 验证频道 OAuth Authorization/account id 覆盖下游值。
- 验证现有 responses/compact 路径和 header override 行为不回归。
