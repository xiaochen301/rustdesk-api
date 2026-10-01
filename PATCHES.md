# PATCHES 台账 — rustdesk-api-fork（管理后台重构 · 后端）

## 项目信息

- 项目：rustdesk-api（apimain）管理后台重构 — 后端配套改造
- 开始：2026-10-02
- 分支：`xc/refactor-admin`（从 `xc/auto-sync` 出发，含设备自动入簿）
- 目标：管理后台入口开关 / 用户类型（托管用户）/ 群组多管理员 / 游离设备 / 角色化菜单配套

## 改造清单（已批准的设计定稿）

1. **管理后台入口开关**（用户点名）——`admin.enable`（config.yaml / env `RUSTDESK_API_ADMIN_ENABLE`）
2. **用户类型**：普通 / 托管（无密码不可登录）/ 系统管理员——`User.Type` 字段
3. **群组增强**：群组模式（集中式/平权式）+ 多管理员（多对多关联表）
4. **游离设备**：hbbs peer 对比 API 台账 → 游离列表 + 收编（绑托管用户）
5. **系统管理员登录拒绝**：`/api/login` 拒绝系统管理员类型
6. **route_names 配套**：新菜单结构的路由名下发

## 债务记录

| # | 日期 | 事项 | 类型 | 文件 | 本金/利息 | 状态 |
|---|------|------|------|------|----------|------|
| 1 | 2026-10-02 | 管理后台开关关闭时，webclient（浏览器远控）与分享功能保持开放——设计边界需在文档中明确（管理后台 vs 用户功能） | 说明 | `router/router.go` | 无息：有意设计，待用户确认边界 | 挂账 |
| 2 | 2026-10-02 | 「管理后台关闭」的完整替代方案（反代层关闭）文档待补——香港机实践为 nginx 层控制 | 债 | 文档 | 低息 | 挂账 |
| 3 | 2026-10-02 | 组集合自动维护未做：群组管理员列表 → 集合规则（type=1 个人规则）的自动映射（含组设备集合的生成逻辑）待实现 | 债 | `service/` | 中息：可见性闭环的关键一环，与前端组管理页面同批 | 挂账 |
| 4 | 2026-10-02 | route_names 配套未做：新菜单结构的路由名下发（角色化菜单）待前端菜单定型后实施 | 债 | `model/user.go` | 中息：与前端 IA 重组同批 | 挂账 |

## 已完成（非债务）

- **管理后台入口开关**（`admin.enable`）：
  - `config/config.go`：`Admin.Enable *bool` + `IsEnabled()`（nil=启用）
  - `http/router/admin.go`：关闭时 `/api/admin/*` 全部返回 403 JSON
  - `http/router/router.go`：关闭时 `/_admin/*` 返回 403 文本
  - `conf/config.yaml`：新增 `enable: true` 配置项（含注释）
  - 客户端协议 `/api/*`（login/heartbeat/sysinfo/ab/audit）完全不受影响
  - 实测：关闭态 403 ✓ / 开启态恢复 ✓ / 客户端协议 200 ✓
- **用户类型（托管用户）**：
  - `model.User.Type`：1=普通 2=托管（默认 1）
  - 托管用户不可登录：客户端 `/api/login` + 管理后台 `/api/admin/login` 双入口拒绝（i18n：ManagedUserCannotLogin）
  - `http/request/admin/user.go`：UserForm 支持 type（默认 1）
  - 实测：创建 type=2 用户 ✓ / 双入口拒绝消息正确 ✓
- **群组增强（Mode + 多管理员）**：
  - `model.Group.Mode`：1=集中式 2=平权式（默认 1）；`AdminIds` 非持久化字段
  - `model.GroupAdmin`：群组管理员多对多表（group_id + user_id）
  - `GroupService.AdminIds/SetAdmins`；Create/Update/Detail/List 全链路支持 admin_ids
  - 可见性映射：集中式=给 N 个管理员各一条个人规则（type=1）；平权式=一条组规则（type=2）——复用现有 CollectionRule 机制
- **游离设备**：
  - `PeerService.ListStray`（user_id=0 的已注册设备）+ `AdoptPeer` 收编（仅限托管用户，校验）
  - 路由：GET `/api/admin/peer/stray` + POST `/api/admin/peer/adopt`
  - 自动脱管：设备实际登录时 `UuidBindUserId` 已有覆盖机制（登录优先）
  - i18n：PeerAlreadyBound / AdoptToManagedOnly
- **系统管理员登录开关**（`app.disable-admin-client-login`，默认 false=允许）：
  - 开启后 IsAdmin 用户不可登录客户端（默认保持既有行为，避免迁移期破坏现有 admin 设备）
  - i18n：AdminCannotLoginClient

## 说明

- 本重构为正规改造（非补丁式），台账记录过程中的临时妥协与技术债
- 与前端台账（rustdesk-api-web/PATCHES.md）配对使用
