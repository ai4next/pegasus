[English](./README.md) | [简体中文](./README-zh.md)
# Pegasus

![Logo](assets/banner.png)

Pegasus is an AI-native CAD development agent specialized for CAD. It understands design intent, reads and writes project files, runs scripts, connects to CAD/geometry tools, and coordinates geometry, automation, and manufacturing-review experts to produce reproducible CAD and CAD-software work.

## 💡 Design Philosophy

- **CAD-native over generic chat.** Pegasus organizes work around units, constraints, feature history, geometry kernels, file formats, manufacturing constraints, and verifiable artifacts.

- **Reproducible over one-off edits.** Design changes should become scripts, parameters, tests, and source files; generated STEP/STL/DXF/renders should be reproducible from source.

- **Specialists over one giant model.** Pegasus seeds `geometry-kernel`, `cad-automation`, and `manufacturing-review` experts, then can evolve expert prompts and durable memory after completed work.

- **Boundaries over vibes.** Pegasus memory, expert memory, expert souls, evolver memory, session, and audit logs are separate so learning does not pollute projects or specialist responsibilities.

- **Verification over visual hallucination.** CAD results should be checked with geometry assertions, metadata, dimensions/tolerances, generated files, and rendered previews whenever possible.

---

## 🚀 Quick Start

Install the latest release on Linux or macOS:

```bash
curl -fsSL https://raw.githubusercontent.com/ai4next/pegasus/main/install.sh | sh
```

Windows PowerShell:

```powershell
iwr https://raw.githubusercontent.com/ai4next/pegasus/main/install.ps1 -useb | iex
```

```bash
# Create and edit config
pg init

# Set your API key
export OPENAI_API_KEY=sk-...

# Start the terminal UI
pg

# Or run a single prompt
pg run "What's in this directory?"
```

Install a specific release or use a user-writable directory:

```bash
VERSION=v0.0.1 INSTALL_DIR="$HOME/.local/bin" sh -c "$(curl -fsSL https://raw.githubusercontent.com/ai4next/pegasus/main/install.sh)"
```

## ✨ Features

- **Multi-model support** — Gemini (Vertex AI), OpenAI, DeepSeek, Claude, Ollama, and any OpenAI-compatible API
- **Built-in tools** — OS-aware command execution, file read/write/patch, user interaction, memory search, expert delegation, suitable for driving FreeCAD/OpenCascade/Blender automation scripts
- **MCP server integration** — plug in any MCP-compatible tool server via config (stdin/stdout transport)
- **Instant-messaging integration** — run a long-lived server that connects Pegasus to Telegram, Feishu/Lark, WeCom, Weixin, QQ, DingTalk, Slack, Discord, LINE, and Weibo
- **Persistent session** — SQLite-backed session/message store with compact `U/A/T/O` evolution logs, automatic compaction, file revision tracking, and session export/import
- **Runtime audit** — Events (tool calls, text delta, errors, evolutions) streamed to a queryable JSONL audit log
- **In-process task queue** — expert and orchestration tasks use a Go channel queue inside each Pegasus process, so multiple local Pegasus instances do not contend for a shared queue database
- **Flat-file memory** — global facts (L1) and SOP files (L2) stored directly in the workspace
- **Plan-Execute agent loop** — every agent is assembled as `planner -> loop(executor -> replanner)`, so requests are planned, executed step by step, and replanned until completion or the iteration limit
- **CAD expert delegation** — seeded `geometry-kernel`, `cad-automation`, and `manufacturing-review` experts with isolated memory and persistent sessions
- **Layered self-evolution** — agent evolver improves Pegasus/experts from completed session; meta evolver improves only the evolution process from evolver session
- **Plugin system** — unified run/model/tool logging and session reaper
- **Terminal UI** — dark theme, Emacs-style keybindings, sidebar, and dialog system
- **Hook system** — 11 lifecycle event hooks (before/after run, tool, model, etc.) with external script execution via JSON stdin/stdout protocol
- **Skill system** — filesystem-based skills auto-loaded via ADK skilltoolset, compatible with SKILL.md-based workflows, supports multiple skill paths

## ⌨️ Commands

| Command | Description |
|---------|-------------|
| `pg` | Start interactive terminal chat |
| `pg run "prompt"` | Run a single prompt, print response |
| `pg run -f prompt.txt` | Run a prompt from a file |
| `pg run -p "hello"` | Run with `--prompt` flag |
| `pg reflect` | Start autonomous idle-watch + scheduler mode |
| `pg im serve` | Run the instant-messaging integration server |
| `pg init` | Create `config.yaml` from the embedded example template |
| `pg configure` | Show or initialize configuration |
| `pg toolsets` | List configured ADK Skill and MCP toolsets |
| `pg session list` | List persistent session |
| `pg session show <id>` | Show session messages |
| `pg session last` | Show the most recently updated session |
| `pg session search <query>` | Search persisted session messages |
| `pg session files <id>` | Show session working files |
| `pg session history <id>` | Show session file revision history |
| `pg session diff <id> <path>` | Show file revision diff |
| `pg session revert <id> <path>` | Revert a file to its previous revision |
| `pg session export <id>` | Export session (markdown/json/jsonl) |
| `pg session import <path>` | Import a session export |
| `pg session compact <id>` | Compact older session context into a summary |
| `pg session delete <id>` | Delete a persistent session |
| `pg session rename <id> <title>` | Rename a session |
| `pg session queue <id>` | Inspect queued prompts for a session |
| `pg session storage` | Inspect persistent session storage stats |
| `pg session storage gc` | Remove orphaned file revision snapshots |
| `pg runtime events` | List runtime audit events |
| `pg runtime summary` | Summarize runtime audit events |

## ⚙️ Configuration

See `config.example.yaml` for all options. Key settings:

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

# Skill system — multiple paths supported
skills:
  enabled: true
  paths:
    - ${HOME}/.pg/skills
    - ./skills

# MCP server integration
mcp:
  servers:
    - name: my-server
      enabled: true
      command: npx
      args: [-y, "@modelcontextprotocol/server-filesystem", /tmp]
      tools: []                 # empty = all tools; specify names to filter

# Session management
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

`model.headers` is optional and is forwarded with every model request, which is useful for custom OpenAI-compatible gateways.

Environment variables override config: `PEGASUS_MODEL_PROVIDER=openai`, `PEGASUS_MODEL_API_KEY=sk-...`, etc.

`bus.queue` is intentionally in-process. It is used for local async delegate/orchestration work inside the current Pegasus process and is not persisted or shared across simultaneously running Pegasus processes. `bus.audit_log` is the durable JSONL event mirror.

### Instant Messaging

Enable one or more `im.platforms` entries, set the required platform credentials, then run:

```bash
pg im serve --config config.yaml
```

`im serve` is a long-lived server process. Run it under your preferred process manager, or in a shell background job:

```bash
nohup pg im serve --config config.yaml > ~/.pg/runtime/im.log 2>&1 &
```

Example:

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

See `config.example.yaml` for Telegram, Feishu/Lark, WeCom, Weixin, QQ, QQ official bot, DingTalk, Slack, Discord, LINE, and Weibo examples.

For Weixin personal accounts, use QR setup first:

```bash
pg im weixin setup
```

The command prints a QR code in the terminal, waits for phone confirmation, prints the token and account values, then exits. Add the printed values to your `weixin` entry in `im.platforms` before running `pg im serve`.

## 🛠️ Tools

| Tool | Description |
|------|-------------|
| `exec` | Execute a shell command using bash, sh, or PowerShell based on the current OS |
| `read` | Read file lines |
| `write` | Write files |
| `patch` | Replace one exact text match in a file |
| `ask` | Interrupt to ask the user a question |
| `delegate` | Delegate a task to an expert; use `mode=sync` for an immediate result or `mode=async` to enqueue work |

`delegate` is loaded dynamically when at least one expert is available. Experts are stored under `state/{expert_name}`; the directory name is the expert name and `soul.md` is the expert's system prompt.

## 🧠 Agent Runtime

Pegasus builds the main agent and every expert with the same Plan-Execute structure:

```text
{name}                         # sequential root
├── {name}_planner              # produces the initial plan
└── {name}_plan_execute_loop     # bounded loop
    ├── {name}_executor         # executes the first unfinished plan step
    └── {name}_replanner        # evaluates progress, updates the plan, or exits the loop
```

The planner stores the current plan in session state. The executor receives that plan plus normal runtime context and tools, executes only the first unfinished step, and stores the step result. The replanner reads the current plan and latest executor result, then either emits an adjusted plan or calls `exit_loop` when the task is complete. Text events keep the ADK author and event id so callers can display planner/replanner progress while collecting the final result from the last executor event.

## 🔌 Hooks & Skills

### Hooks

Place executable scripts in `hooks/<event>/` directories. They receive JSON context via stdin and return JSON via stdout.

```
hooks/
├── before_run/          # Before agent run
├── after_run/           # After agent run
├── before_tool/         # Before tool execution
├── after_tool/          # After tool execution
├── before_model/        # Before LLM call
├── after_model/         # After LLM call
├── before_agent/        # Before agent execution
├── after_agent/         # After agent execution
├── on_user_message/     # On user message
├── on_model_error/      # On model error
└── on_tool_error/       # On tool error
```

Example script (`hooks/before_tool/audit.sh`):

```bash
#!/bin/sh
# stdin: {"event":"before_tool","tool_name":"write","tool_args":{...}}
echo '{"allow": true}'
# Return {"allow": false, "reason": "..."} to block the tool
```

### Skills

Create skill directories under `skills/` or any configured skill path. Each skill is a `SKILL.md` file with YAML frontmatter.

```
skills/
└── code-review/
    ├── SKILL.md           # Required: YAML frontmatter + Markdown instructions
    └── references/        # Optional: reference docs
```

Example (`skills/code-review/SKILL.md`):

```markdown
---
name: code-review
description: Professional code review for PRs and changes
allowed-tools: [read, patch, web_scan]
---

You are a code review expert. Focus on:
1. Security — OWASP Top 10, injection vulnerabilities
2. Correctness — logic errors, edge cases
3. Maintainability — naming, separation of concerns
```

### MCP Servers

Pegasus supports any MCP-compatible server via stdin/stdout transport. Configure servers in `config.yaml`:

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

Use `pg toolsets` to verify configured servers and their available tools.

## 📁 Project Structure

```
pegasus/
├── main.go                          # Entry point
├── internal/
│   ├── agent/
│   │   ├── agent.go                 # Agent factory entrypoint and shared wiring
│   │   ├── orchestrator.go          # Plan-Execute ADK agent tree assembly
│   │   ├── builtin.go               # Executor prompt/context/tool preparation
│   │   ├── context.go               # Context builder for agent runs
│   │   ├── tool.go                  # Dynamic built-in/toolset assembly
│   │   └── evolver.go               # Agent/meta evolution agent factories
│   ├── prompt/                      # Embedded prompt templates
│   │   └── template/                # Markdown prompt templates
│   ├── config/                      # YAML + env config (viper), embedded config.example.yaml
│   ├── cli/                         # Cobra CLI commands (init, run, reflect, im, configure, toolsets, session, runtime)
│   ├── tui/                         # Terminal UI
│   │   ├── tui.go                   # Compatibility wrapper
│   │   ├── app/                     # Model, runtime, session, commands, dialogs, layout
│   │   ├── components/              # Chat, input line, toolbar, sidebar renderers
│   │   └── styles/                  # Dark theme, icons, color themes
│   ├── model/                       # Multi-provider LLM factory
│   ├── memory/                      # Flat-file memory plus search service (L1 facts, L2 SOP files)
│   ├── session/                     # Persistent session manager with compact logs, file tracking, references
│   ├── store/
│   │   ├── db/                      # GORM/SQLite models, DBRegistry, memory index, session/mailbox stores
│   │   └── fs/                      # File-backed stores such as compact session logs
│   ├── runtime/                     # Run streaming, session compaction, loop detection
│   ├── bus/                         # In-process event broker, channel task queue, audit mirror
│   ├── im/                          # Instant-messaging integration
│   ├── plugin/                      # Plugin registry + built-ins
│   ├── hook/                        # Hook manager + script runner
│   ├── reflect/                     # Autonomous idle watcher + scheduler
│   └── expert/                      # Directory-backed expert registry (`soul.md`)
├── hooks/                            # Hook scripts (convention-based, 11 event dirs)
├── skills/                           # Skill definitions (ADK skilltoolset)
├── config.example.yaml              # Symlink to internal/config/config.example.yaml
├── config/tasks/                    # Sample scheduler task definitions
├── data/
│   └── session/                    # Session history
├── go.mod
└── go.sum
```

## 📂 Runtime Directory

All runtime data is stored under `workspace` in `config.yaml`. If omitted, it defaults to `$HOME/.pg`. The directory is created on first run:

```
~/.pg/                                    # workspace (default: $HOME/.pg)
├── config.yaml                           # User configuration (created by `pg init` or `pg configure`)
├── tui.log                               # Terminal UI runtime log (redirected for display safety)
├── memory/                               # Flat-file memory by agent
│   ├── pegasus/
│   │   ├── l1.toml                       # L1 global facts
│   │   └── l2/                           # L2 SOP files (*.md)
│   └── {expert_name}/
│       ├── l1.toml
│       └── l2/
├── session/                              # Compact session logs and snapshots by agent
│   ├── pegasus/
│   │   ├── <id>.log
│   │   └── snapshots/
│   └── {expert_name}/
│       └── <id>.log
├── state/                                # Agent state stores and souls
│   ├── state.db                          # Global DB: cross-owner memory search index and internal mailbox state
│   ├── pegasus/
│   │   └── state.db
│   └── {expert_name}/
│       ├── soul.md                       # Expert system prompt
│       └── state.db
├── bus/
│   └── events.jsonl                      # Unified bus/runtime audit mirror; task queue is in-process
├── hooks/                                # Hook event scripts (11 lifecycle events)
└── skills/                               # Skill definitions (SKILL.md)
```

## 🏗️ Build

Install from GitHub Releases:

```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/ai4next/pegasus/main/install.sh | sh

# Windows PowerShell
iwr https://raw.githubusercontent.com/ai4next/pegasus/main/install.ps1 -useb | iex
```

Build from source:

```bash
go build -o pg .
./pg --help
```

Requires Go 1.26+.

## ⭐ Community & Support

If this project helped you, please consider leaving a **Star!** 🙏

You're also welcome to join the **Pegasus Experience & Exchange Community** for discussion, feedback, and co-building 👏

<div align="center">
  <table>
    <tr>
      <td align="center"><strong>WeChat Group</strong><br/><img src="assets/wechat_group.jpg" alt="WeChat Group QR" width="240"/></td>
    </tr>
  </table>
</div>

## 📄 License

MIT, See [`LICENSE`](LICENSE) for full text.

---

## 📈 Star History

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
