# 限流熔断系统（flow-control）

一个基于纯 Go 标准库（`net/http`，零第三方依赖）的限流与熔断后端，提供限流规则、令牌桶、熔断器、熔断事件、配额、流量统计与告警规则。内置三种限流算法（固定窗口 / 滑动窗口 / 令牌桶）、熔断器状态机（closed → open → half_open → closed）与请求放行决策，并提供原生前端控制台。

## 运行

```bash
cd origin
go run ./cmd/server
# 默认监听 :8080，可用 PORT / ADDR 覆盖
# 可选环境变量：AUTH_TOKEN（启用 Bearer 鉴权）、RATE_LIMIT（每 IP 每分钟限流，默认 200）
```

访问前端控制台：`http://localhost:8080/`（需在 `origin/` 目录下启动）。

## 分层结构

```
origin/
├── cmd/server/main.go        # 入口 + 前端挂载 + 优雅关闭
├── frontend/index.html       # 原生前端控制台（零构建）
├── internal/
│   ├── app/ config/ model/ store/ service/ handler/
└── pkg/ httpx/ idgen/ logger/
```

## 核心概念

- **RateLimitRule**：限流规则，算法 `fixed_window` / `sliding_window` / `token_bucket`，按资源限流。
- **TokenBucket**：令牌桶配置（容量、补充速率）。
- **CircuitBreaker**：熔断器，状态机 `closed → open → half_open → closed`。
- **BreakerEvent**：熔断事件（opened / closed / half_open / rejected）。
- **Quota**：按维度（ip / user / api）的配额。
- **TrafficStats**：资源流量统计（总量 / 放行 / 拒绝）。
- **AlertRule**：告警规则（拒绝率 / 延迟阈值）。

## 限流算法

- **固定窗口**：按 `window_sec` 划分窗口，窗口内计数超过 `limit` 拒绝。
- **滑动窗口**：基于时间戳日志，动态滑动窗口内计数。
- **令牌桶**：容量为 `limit`，按 `limit / window_sec` 速率补充令牌，令牌耗尽拒绝。

## API 一览（核心）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET | /api/rules | 创建 / 列表（resource/algorithm/status/keyword） |
| GET/PUT/DELETE | /api/rules/{id} | 规则详情 / 更新 / 删除 |
| PATCH | /api/rules/{id}/toggle | 启用 / 停用 |
| POST/GET | /api/breakers | 创建 / 列表 |
| GET/PUT/DELETE | /api/breakers/{id} | 熔断器详情 / 更新 / 删除 |
| GET | /api/breaker-events | 熔断事件（breaker_id 筛选） |
| POST/GET | /api/quotas | 创建 / 列表（rule_id 筛选） |
| POST | /api/quotas/consume | 消耗配额 |
| DELETE | /api/quotas/{id} | 删除配额 |
| GET | /api/traffic-stats | 流量统计 |
| GET | /api/traffic-stats/{resource}/summary | 资源流量汇总 |
| POST/GET | /api/alert-rules | 创建 / 列表 |
| GET | /api/alert-rules/evaluate | 评估告警 |
| POST | /api/decide | 放行决策（resource + key） |
| POST | /api/records/success | 记录成功 |
| POST | /api/records/failure | 记录失败 |
| GET | /api/stats/flow | 整体统计 |
| GET | /api/stats/resource?resource= | 单资源统计 |
| GET | /api/export · POST /api/import | 快照导出 / 导入 |
| GET | /healthz · /readyz | 存活 / 就绪探针 |

## 放行决策逻辑

1. 先查熔断器：`open` 且未过 `open_duration_sec` 直接拒绝；超时自动转 `half_open`。
2. 再按资源匹配限流规则，按算法判定是否放行。
3. 同步更新流量统计。

## 统一响应

```json
{ "code": 0, "message": "ok", "data": { } }
```
