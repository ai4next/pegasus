package app

import (
	"context"
	"strings"
	"time"

	"charm.land/bubbles/v2/textarea"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"

	pegasusagent "github.com/ai4next/pegasus/internal/agent"
	"github.com/ai4next/pegasus/internal/bus"
	"github.com/ai4next/pegasus/internal/config"
	pegasussession "github.com/ai4next/pegasus/internal/session"
	"github.com/ai4next/pegasus/internal/tui/components"
)

type Model struct {
	width          int
	height         int
	agent          agent.Agent
	cfg            *config.Config
	runner         *runner.Runner
	sessionService session.Service
	pluginCfg      runner.PluginConfig
	messages       []components.Message
	input          string
	cursorPos      int
	textarea       textarea.Model
	running        bool
	showWelcome    bool
	sessionID      string
	modelName      string
	cursorRow      int
	cursorCol      int
	scrollOffset   int
	selection      selectionState
	runtimeBroker  bus.Broker
	runtimeCh      <-chan bus.Event
	runtimeCancel  context.CancelFunc
	auditLogger    *bus.AuditLogger
	pulseOn        bool
	currentTool    string
	pendingConfirm *pendingConfirmation
	responseBuffer strings.Builder
	toolStarts     map[string]time.Time
	chatCache      string
	chatLinesCache []string
	chatCacheWidth int
	chatCachePulse bool
	chatCacheDirty bool
	resizeSeq      int
	fileCount      int
	queueCount     int
	sessionTitle   string
	fileChanges    []pegasussession.FileChangeSummary
	sessionDialog  *sessionDialogState
	commandDialog  *commandDialogState
	filePicker     *filePickerState
	searchDialog   *searchDialogState
	toolsets       []pegasusagent.ToolsetDescriptor
	promptHistory  []string
	historyIndex   int
	historyDraft   string
}

type pulseMsg struct{}

type selectionState struct {
	Active bool
	StartX int
	StartY int
	EndX   int
	EndY   int
}

type pendingConfirmation struct {
	ID       string
	ToolName string
	Args     string
}

type sessionDialogState struct {
	Sessions []pegasussession.Metadata
	Selected int
}

type commandDialogState struct {
	Commands []components.CommandDialogItem
	Query    string
	Selected int
}

type filePickerState struct {
	Files    []string
	Query    string
	Selected int
}

type searchDialogState struct {
	Query    string
	Results  []pegasussession.MessageSearchResult
	Selected int
}

type runtimeEventMsg struct {
	Event bus.Event
	OK    bool
}

type resizeRenderMsg struct {
	Seq int
}
