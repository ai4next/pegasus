[English](./README.md) | [简体中文](./README-zh.md)
# Pegasus

![Logo](assets/banner.png)

Pegasus 是一个 AI native 的 CAD 开发 Agent，专注于 CAD 领域：它能理解设计意图、读写工程文件、运行脚本、接入 CAD/几何工具，并通过几何、自动化、制造评审等专家子 Agent 协同完成可复现的建模和 CAD 软件开发工作。

## 💡 设计哲学

- **CAD 原生胜于通用聊天。** Pegasus 不只是回答问题，而是围绕 CAD 的单位、约束、特征树、几何内核、文件格式、制造约束和可验证产物组织工作。

- **可复现胜于一次性编辑。** 设计修改应尽量沉淀为脚本、参数、测试和源文件；STEP/STL/DXF/渲染图等导出物应能从源文件重新生成。

- **专家协同胜于单体模型。** Pegasus 默认种子专家覆盖 `geometry-kernel`、`cad-automation`、`manufacturing-review`，并可在完成任务后进化专家和长期记忆。

- **边界胜于感觉。** Pegasus 记忆、专家记忆、专家 soul、evolver 记忆、session 和 audit logs 分开存储，保证学习不会污染项目或专家职责。

- **验证胜于视觉幻觉。** 对 CAD 结果优先使用几何检查、元数据、尺寸/公差、导出文件和渲染预览验证，而不是只依赖自然语言描述。

---

## 🚀 快速开始

在 Linux 或 macOS 上安装最新发布版：

```bash
curl -fsSL https://raw.githubusercontent.com/ai4next/pegasus/main/install.sh | sh
```

Windows PowerShell：

```powershell
iwr https://raw.githubusercontent.com/ai4next/pegasus/main/install.ps1 -useb | iex
```

```bash
# 创建并编辑配置
pg init

# 设置 API Key
export OPENAI_API_KEY=sk-...

# 启动终端界面
pg

# 或单次执行
pg run "这个目录里有什么？"
```

也可以指定版本或安装到用户目录：

```bash
VERSION=v0.0.1 INSTALL_DIR="$HOME/.local/bin" sh -c "$(curl -fsSL https://raw.githubusercontent.com/ai4next/pegasus/main/install.sh)"
```

## ✨ 功能特性

- **多模型支持** — Gemini (Vertex AI)、OpenAI、DeepSeek、Claude、Ollama，以及任何兼容 OpenAI 的 API
- **内建工具** — 自动适配操作系统的命令执行、文件读写/补丁、用户交互、记忆搜索、专家委托，适合驱动 FreeCAD/OpenCascade/Blender 等 CAD 自动化脚本
- **MCP Server 集成** — 通过配置接入任意 MCP 兼容工具服务（stdin/stdout transport）
- **即时通信软件接入** — 以常驻 server 方式接入 Telegram、飞书/Lark、企业微信、微信个人号、QQ、钉钉、Slack、Discord、LINE、微博等平台
- **持久会话** — SQLite 支撑的 session/message 存储，配套精简 `U/A/T/O` 进化日志，支持自动压缩、文件 revision tracking、session 导入导出
- **运行时审计** — 工具调用、文本增量、错误、进化等事件流式写入可查询 JSONL audit log
- **进程内任务队列** — 专家与编排任务使用每个 Pegasus 进程内的 Go channel 队列，本机同时启动多个 Pegasus 时不会争抢共享队列数据库
- **扁平文件记忆** — 全局事实（L1）和 SOP 文件（L2）直接存储在 workspace 中
- **Plan-Execute Agent 循环** — 每个 Agent 都按 `planner -> loop(executor -> replanner)` 组装，先规划，再按步骤执行，并在完成或达到迭代上限前持续复盘调整
- **CAD 专家委托** — 默认内置 `geometry-kernel`、`cad-automation`、`manufacturing-review` 三类专家，并支持后续自进化扩展
- **分层自进化** — agent evolver 从完成会话中改进 Pegasus/专家；meta evolver 只从 evolver 会话中改进进化过程本身
- **插件系统** — 统一 run/model/tool 日志与会话回收
- **终端界面** — 暗色主题、Emacs 风格键绑定、侧边栏、Dialog 系统
- **Hook 系统** — 11 种生命周期事件钩子（run/tool/model 等前后），通过 JSON stdin/stdout 协议执行外部脚本
- **Skill 系统** — 基于文件系统的技能自动加载（ADK skilltoolset），兼容基于 `SKILL.md` 的工作流，支持多个 skill path

## ⌨️ 命令

| 命令 | 说明 |
|------|------|
| `pg` | 启动交互式终端聊天 |
| `pg run "提示词"` | 单次执行并打印响应 |
| `pg run -f prompt.txt` | 从文件读取提示词执行 |
| `pg run -p "hello"` | 使用 `--prompt` flag 执行 |
| `pg reflect` | 启动自主空闲监听 + 调度模式 |
| `pg im serve` | 启动即时通信接入 server |
| `pg init` | 从嵌入的示例模板创建 `config.yaml` |
| `pg configure` | 查看或初始化配置 |
| `pg toolsets` | 列出已配置的 ADK Skill 和 MCP toolsets |
| `pg session list` | 列出持久会话 |
| `pg session show <id>` | 查看会话消息 |
| `pg session last` | 查看最近更新的会话 |
| `pg session search <query>` | 搜索持久会话消息 |
| `pg session files <id>` | 查看会话工作文件 |
| `pg session history <id>` | 查看会话文件 revision 历史 |
| `pg session diff <id> <path>` | 查看文件 revision diff |
| `pg session revert <id> <path>` | 将文件回滚到上一 revision |
| `pg session export <id>` | 导出会话（markdown/json/jsonl） |
| `pg session import <path>` | 导入会话导出文件 |
| `pg session compact <id>` | 将较旧会话上下文压缩为 summary |
| `pg session delete <id>` | 删除持久会话 |
| `pg session rename <id> <title>` | 重命名会话 |
| `pg session queue <id>` | 查看会话中排队的 prompts |
| `pg session storage` | 查看持久会话存储统计 |
| `pg session storage gc` | 删除孤立的文件 revision snapshots |
| `pg runtime events` | 列出 runtime audit events |
| `pg runtime summary` | 汇总 runtime audit events |

## ⚙️ 配置

完整配置见 `config.example.yaml`。关键配置如下：

```yaml
model:
  provider: openai          # gemini | openai | deepseek | claude | ollama
  name: gpt-4o
  base_url: https://api.openai.com/v1
  api_key: ${OPENAI_API_KEY}
  headers:
    X-Request-Source: pegasus

tools:
  exec:
    enabled: true
    timeout: 30s

expert:
  max_count: 10

bus:
  audit_log: ${HOME}/.pg/bus/events.jsonl
  queue:
    max_size: 100

# Skill 系统，支持多个路径
skills:
  enabled: true
  paths:
    - ${HOME}/.pg/skills
    - ./skills

# MCP Server 集成
mcp:
  servers:
    - name: my-server
      enabled: true
      command: npx
      args: [-y, "@modelcontextprotocol/server-filesystem", /tmp]
      tools: []                 # 空列表 = 全部工具；也可指定工具名过滤

# 会话管理
session:
  max_turns: 75
  loop_detection:
    enabled: true
    window_size: 10
    max_repeats: 5

plugins:
  - name: logger
    enabled: true
```

`model.headers` 是可选配置，会随每次模型请求一起发送，适合接入需要自定义请求头的 OpenAI-compatible 网关。

环境变量可以覆盖配置：`PEGASUS_MODEL_PROVIDER=openai`、`PEGASUS_MODEL_API_KEY=sk-...` 等。

`bus.queue` 是进程内队列，只用于当前 Pegasus 进程里的本地异步 delegate / orchestrator 工作，不持久化，也不在多个同时运行的 Pegasus 进程之间共享。`bus.audit_log` 是持久化 JSONL 事件审计镜像。

### 即时通信接入

在 `im.platforms` 中启用一个或多个平台，配置对应凭据，然后运行：

```bash
pg im serve --config config.yaml
```

`im serve` 是一个常驻 server 进程。可以交给你习惯的进程管理器托管，也可以用 shell 后台方式运行：

```bash
nohup pg im serve --config config.yaml > ~/.pg/runtime/im.log 2>&1 &
```

示例：

```yaml
im:
  platforms:
    - name: feishu
      enabled: true
      options:
        app_id: ${FEISHU_APP_ID}
        app_secret: ${FEISHU_APP_SECRET}
        domain: feishu
        allow_from: ""

    - name: qq
      enabled: false
      options:
        ws_url: ws://127.0.0.1:3001
        token: ${QQ_ONEBOT_TOKEN}
        allow_from: ""
```

Telegram、飞书/Lark、企业微信、微信个人号、QQ、QQ 官方机器人、钉钉、Slack、Discord、LINE、微博等完整示例见 `config.example.yaml`。

微信个人号通常需要先扫码接入：

```bash
pg im weixin setup
```

该命令会在终端打印二维码，等待手机确认登录，打印 token 和账号信息后退出。把打印出的值填到 `im.platforms` 中的 `weixin` 配置后，再运行 `pg im serve`。

## 🛠️ 工具列表

| 工具 | 说明 |
|------|------|
| `exec` | 根据当前操作系统使用 bash、sh 或 PowerShell 执行 shell 命令 |
| `read` | 读取文件行 |
| `write` | 写入文件 |
| `patch` | 替换文件中的一个精确文本匹配 |
| `ask` | 中断执行并向用户提问 |
| `delegate` | 将任务委托给专家；`mode=sync` 立即返回结果，`mode=async` 将任务入队 |

`delegate` 会在每次调用模型前动态判断，只有启用专家委托且至少存在一个专家时才会加载。专家存储在 `state/{expert_name}` 目录下，目录名就是专家名，`soul.md` 是专家的系统提示词。

## 🧠 Agent Runtime

Pegasus 会用同一套 Plan-Execute 结构构建主 Agent 和每个专家 Agent：

```text
{name}                         # sequential root
├── {name}_planner              # 生成初始计划
└── {name}_plan_execute_loop     # 有上限的循环
    ├── {name}_executor         # 执行当前计划中的第一个未完成步骤
    └── {name}_replanner        # 评估进度、更新计划，或退出循环
```

planner 将当前计划写入 session state。executor 读取计划、普通运行上下文和工具，只执行第一个未完成步骤，并写入步骤结果。replanner 读取当前计划和最新 executor 结果，判断是输出调整后的新计划，还是在任务完成时调用 `exit_loop`。文本事件会保留 ADK author 和 event id，因此调用方可以展示 planner/replanner 进度，同时从最后一次 executor 事件收集最终结果。

## 🔌 Hooks & Skills

### Hooks

将可执行脚本放入 `hooks/<event>/` 目录。脚本通过 stdin 接收 JSON 上下文，并通过 stdout 返回 JSON。

```
hooks/
├── before_run/          # Agent run 前
├── after_run/           # Agent run 后
├── before_tool/         # 工具执行前
├── after_tool/          # 工具执行后
├── before_model/        # LLM 调用前
├── after_model/         # LLM 调用后
├── before_agent/        # Agent 执行前
├── after_agent/         # Agent 执行后
├── on_user_message/     # 用户消息
├── on_model_error/      # 模型错误
└── on_tool_error/       # 工具错误
```

示例脚本（`hooks/before_tool/audit.sh`）：

```bash
#!/bin/sh
# stdin: {"event":"before_tool","tool_name":"write","tool_args":{...}}
echo '{"allow": true}'
# 返回 {"allow": false, "reason": "..."} 可以阻止工具执行
```

### Skills

在 `skills/` 或任意已配置的 skill path 下创建技能目录。每个技能是一个带 YAML frontmatter 的 `SKILL.md` 文件。

```
skills/
└── code-review/
    ├── SKILL.md           # 必需：YAML frontmatter + Markdown 指令
    └── references/        # 可选：参考资料
```

示例（`skills/code-review/SKILL.md`）：

```markdown
---
name: code-review
description: Professional code review for PRs and changes
allowed-tools: [read, patch]
---

你是一个代码审查专家。重点关注：
1. 安全性 —— OWASP Top 10、注入漏洞
2. 正确性 —— 逻辑错误、边界条件
3. 可维护性 —— 命名、职责分离
```

### MCP Servers

Pegasus 支持通过 stdin/stdout transport 接入任意 MCP-compatible server。在 `config.yaml` 中配置：

```yaml
mcp:
  servers:
    - name: filesystem
      command: npx
      args: [-y, "@modelcontextprotocol/server-filesystem", /tmp]
    - name: github
      command: npx
      args: [-y, "@modelcontextprotocol/server-github"]
      tools: [issues, pulls]
```

使用 `pg toolsets` 可以确认当前配置的 servers 及其可用工具。

## 📁 项目结构

```
pegasus/
├── main.go                          # 入口
├── internal/
│   ├── agent/
│   │   ├── agent.go                 # Agent 工厂入口与共享 wiring
│   │   ├── orchestrator.go          # Plan-Execute ADK Agent 树组装
│   │   ├── builtin.go               # Executor prompt/context/tool 准备
│   │   ├── context.go               # Agent run 上下文构建器
│   │   ├── tool.go                  # 动态内建工具/toolset 组装
│   │   └── evolver.go               # Agent/meta 进化 Agent 工厂
│   ├── prompt/                      # 内嵌提示词模板
│   │   └── template/                # Markdown 提示词模板
│   ├── config/                      # YAML + 环境变量配置 (viper)，内嵌 config.example.yaml
│   ├── cli/                         # Cobra CLI 命令 (init, run, reflect, im, configure, toolsets, session, runtime)
│   ├── tui/                         # Terminal UI
│   │   ├── tui.go                   # 兼容 wrapper
│   │   ├── app/                     # Model、runtime、session、commands、dialogs、layout
│   │   ├── components/              # Chat、input、toolbar、sidebar renderers
│   │   └── styles/                  # Dark theme、icons、color themes
│   ├── model/                       # 多 Provider LLM 工厂
│   ├── memory/                      # 扁平文件记忆与搜索服务：事实与 SOP
│   ├── session/                     # 持久 SessionService：compact log、file tracking、references
│   ├── store/
│   │   ├── db/                      # GORM/SQLite 模型、DBRegistry、记忆索引、session/mailbox 存储
│   │   └── fs/                      # 文件型存储，例如精简 session log
│   ├── runtime/                     # run 流式处理、session compact、loop detection
│   ├── bus/                         # 进程内事件 broker、channel 任务队列、审计镜像
│   ├── im/                          # 即时通信接入
│   ├── plugin/                      # 插件注册中心 + 内建插件
│   ├── hook/                        # Hook 管理器 + 脚本执行器
│   ├── reflect/                     # 自主空闲监听 + 调度器
│   └── expert/                      # 基于目录的专家注册中心（`soul.md`）
├── hooks/                            # Hook 脚本目录（约定式，11 个事件子目录）
├── skills/                           # Skill 定义目录（ADK skilltoolset）
├── config.example.yaml              # 指向 internal/config/config.example.yaml 的符号链接
├── config/tasks/                    # 调度任务定义示例
├── data/
│   └── session/                    # Session history
├── go.mod
└── go.sum
```

## 📂 运行时目录

所有运行时数据都存储在 `workspace`（默认为 `~/.pg/`），首次启动时自动创建：

```
~/.pg/                                    # workspace（默认: $HOME/.pg）
├── config.yaml                           # 用户配置（由 `pg init` 或 `pg configure` 创建）
├── tui.log                               # 终端界面 runtime 日志（重定向，避免干扰界面）
├── memory/                               # 按 agent 分区的扁平文件记忆
│   ├── pegasus/
│   │   ├── l1.toml                       # L1 全局事实
│   │   └── l2/                           # L2 SOP 文件（*.md）
│   └── {expert_name}/
│       ├── l1.toml
│       └── l2/
├── session/                              # 按 agent 分区的精简会话日志与 snapshots
│   ├── pegasus/
│   │   ├── <id>.log
│   │   └── snapshots/
│   └── {expert_name}/
│       └── <id>.log
├── state/                                # 按 agent 分区的状态库与 soul
│   ├── state.db                          # 全局 DB：跨 owner 记忆搜索索引与内部 mailbox 状态
│   ├── pegasus/
│   │   └── state.db
│   └── {expert_name}/
│       ├── soul.md                       # 专家系统提示词
│       └── state.db
├── bus/
│   └── events.jsonl                      # 统一 bus/runtime 审计镜像；任务队列在进程内
├── hooks/                                # Hook 事件脚本（11 种生命周期事件）
└── skills/                               # Skill 定义（SKILL.md）
```

## 🏗️ 构建

从 GitHub Releases 安装：

```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/ai4next/pegasus/main/install.sh | sh

# Windows PowerShell
iwr https://raw.githubusercontent.com/ai4next/pegasus/main/install.ps1 -useb | iex
```

从源码构建：

```bash
go build -o pg .
./pg --help
```

需要 Go 1.26+。

## ⭐ 社区与支持

如果这个项目对你有帮助，欢迎点一个 **Star!** 🙏

也欢迎加入 **Pegasus 体验交流社区**，一起交流、反馈、共建 👏

<div align="center">
  <table>
    <tr>
      <td align="center"><strong>微信群</strong><br/><img src="assets/wechat_group.jpg" alt="微信群 二维码" width="240"/></td>
    </tr>
  </table>
</div>

## 📄 许可证

MIT，详见 [`LICENSE`](LICENSE)。

---

## 📈 Star 历史

<div align="center">

<a href="https://star-history.com/#ai4next/pegasus&Date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=ai4next/pegasus&type=Date&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=ai4next/pegasus&type=Date" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=ai4next/pegasus&type=Date" />
  </picture>
</a>

<br/><br/>
</div>
