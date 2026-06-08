package agent

import (
	"fmt"
	"log"
	"strings"

	"github.com/ai4next/pegasus/internal/global"
	pegasussession "github.com/ai4next/pegasus/internal/session"
	adkagent "google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	adksession "google.golang.org/adk/session"
	"google.golang.org/genai"
)

const (
	defaultContextMessages        = 50
	contextKeepHeadMessages       = 3
	contextKeepRecentToolResults  = 3
	contextMicroCompactRunes      = 120
	contextToolResultBudgetRunes  = 200000
	contextPersistThresholdRunes  = 30000
	contextAutoCompactLimitRunes  = 50000
	contextAutoSummaryMaxRunes    = 4000
	contextMessagePreviewMaxRunes = 600
)

func instructionProvider(build BuildConfig) func(adkagent.CallbackContext, *model.LLMRequest) (string, error) {
	builder := strings.Builder{}
	return func(ctx adkagent.CallbackContext, req *model.LLMRequest) (string, error) {
		defer builder.Reset()
		builder.WriteString(build.Instruction)
		if build.MemoryService != nil {
			if l0Content := build.MemoryService.GetL0Content(); l0Content != "" {
				builder.WriteString("\n\n")
				builder.WriteString(l0Content)
			}
		}
		return strings.TrimSpace(builder.String()), nil
	}
}

func contentsProvider(build BuildConfig) func(adkagent.CallbackContext, *model.LLMRequest) ([]*genai.Content, error) {
	return func(ctx adkagent.CallbackContext, req *model.LLMRequest) ([]*genai.Content, error) {
		if build.SessionService == nil || build.ContextMessages == 0 {
			return nil, nil
		}
		limit := build.ContextMessages
		if limit < 0 {
			limit = defaultContextMessages
		}
		window, err := loadSessionContext(build.SessionService, ctx.AppName(), ctx.UserID(), ctx.SessionID(), limit)
		if err != nil {
			log.Println("failed to get context window", err)
			return nil, nil
		}
		return sessionContextContents(window), nil
	}
}

func sessionContextContents(window sessionContext) []*genai.Content {
	window = compactSessionContext(window)
	if !hasSessionContext(window) {
		return nil
	}
	builder := strings.Builder{}
	builder.WriteString("## Session Context Usage\n")
	builder.WriteString("- Treat session context as compact hints, not complete history.\n")
	builder.WriteString("- Working files are path/status pointers only; call file tools before relying on file contents.\n")
	builder.WriteString("- Session references are user-selected historical pointers; use their preview as intent, not as full transcript.\n")

	if strings.TrimSpace(window.Summary) != "" {
		builder.WriteString("\n\n## Session Summary\n")
		builder.WriteString(window.Summary)
	}
	if len(window.Messages) > 1 {
		historyMessage := window.Messages[:len(window.Messages)-1]
		builder.WriteString("\n\n## Session Context\n")
		for _, msg := range historyMessage {
			builder.WriteString("- ")
			builder.WriteString(string(msg.Role))
			if msg.ToolName != "" {
				builder.WriteString("/")
				builder.WriteString(msg.ToolName)
			}
			builder.WriteString(": ")
			builder.WriteString(compactContextText(msg.Content, msg.Result, msg.Args))
			builder.WriteByte('\n')
		}
	}
	if len(window.Files) > 0 {
		builder.WriteString("\n\n## Session Working Files\n")
		for _, file := range window.Files {
			builder.WriteString("- ")
			builder.WriteString(file.Path)
			var parts []string
			if file.ReadCount > 0 {
				parts = append(parts, "read")
			}
			if file.WriteCount > 0 {
				parts = append(parts, "modified")
			}
			if len(parts) > 0 {
				builder.WriteString(" (")
				builder.WriteString(strings.Join(parts, ", "))
				builder.WriteString(")")
			}
			builder.WriteByte('\n')
		}
	}
	if len(window.References) > 0 {
		builder.WriteString("\n\n## Session References\n")
		for _, ref := range window.References {
			builder.WriteString("- ")
			builder.WriteString(ref.SessionID)
			if ref.Role != "" {
				builder.WriteString(" [")
				builder.WriteString(string(ref.Role))
				builder.WriteString("]")
			}
			if strings.TrimSpace(ref.Preview) != "" {
				builder.WriteString(": ")
				builder.WriteString(compactContextText(ref.Preview))
			}
			builder.WriteByte('\n')
		}
	}
	return []*genai.Content{genai.NewContentFromText(strings.TrimSpace(builder.String()), genai.RoleUser)}
}

type sessionContext struct {
	Summary     string
	Messages    []pegasussession.Message
	Files       []pegasussession.SessionFile
	References  []pegasussession.SessionReference
	MaxMessages int
}

func loadSessionContext(svc adksession.Service, appName, userID, sessionID string, maxMessages int) (sessionContext, error) {
	messages, err := pegasussession.Messages(svc, appName, userID, sessionID)
	if err != nil {
		return sessionContext{}, err
	}
	window := sessionContext{MaxMessages: maxMessages}
	for _, msg := range messages {
		if msg.Summary {
			window.Summary = msg.Content
			continue
		}
		window.Messages = append(window.Messages, msg)
	}

	files, err := pegasussession.SessionFiles(svc, appName, userID, sessionID)
	if err != nil {
		return sessionContext{}, err
	}
	if maxMessages > 0 && len(files) > maxMessages {
		files = files[:maxMessages]
	}
	window.Files = files

	refs, err := pegasussession.SessionReferences(svc, appName, userID, sessionID)
	if err != nil {
		return sessionContext{}, err
	}
	if maxMessages > 0 && len(refs) > maxMessages {
		refs = refs[:maxMessages]
	}
	window.References = refs
	return window, nil
}

func compactSessionContext(window sessionContext) sessionContext {
	window.Messages = cloneSessionMessages(window.Messages)
	applyToolResultBudget(&window)
	window.Messages = snipContextMessagesPreservingCurrent(window.Messages, effectiveContextMaxMessages(window.MaxMessages))
	microCompactToolResults(window.Messages)
	if estimateSessionContextRunes(window) > contextAutoCompactLimitRunes {
		window = autoCompactSessionContext(window)
	}
	return window
}

func effectiveContextMaxMessages(maxMessages int) int {
	if maxMessages < 0 {
		return defaultContextMessages
	}
	return maxMessages
}

func snipContextMessagesPreservingCurrent(messages []pegasussession.Message, maxMessages int) []pegasussession.Message {
	if maxMessages <= 0 || len(messages) <= maxMessages {
		return messages
	}
	if len(messages) <= 1 {
		return messages
	}
	current := messages[len(messages)-1]
	history := snipContextMessages(messages[:len(messages)-1], maxMessages)
	out := make([]pegasussession.Message, 0, len(history)+1)
	out = append(out, history...)
	out = append(out, current)
	return out
}

func snipContextMessages(messages []pegasussession.Message, maxMessages int) []pegasussession.Message {
	if maxMessages <= 0 || len(messages) <= maxMessages {
		return messages
	}
	if maxMessages < 3 {
		return messages[len(messages)-maxMessages:]
	}
	keepHead := min(contextKeepHeadMessages, maxMessages-2)
	keepTail := maxMessages - keepHead - 1
	snipped := len(messages) - keepHead - keepTail
	marker := pegasussession.Message{
		Role:    pegasussession.MessageUser,
		Content: fmt.Sprintf("[snipped %d older session messages]", snipped),
	}
	out := make([]pegasussession.Message, 0, maxMessages)
	out = append(out, messages[:keepHead]...)
	out = append(out, marker)
	out = append(out, messages[len(messages)-keepTail:]...)
	return out
}

func microCompactToolResults(messages []pegasussession.Message) {
	var toolResults []*pegasussession.Message
	for i := range messages {
		if messages[i].Role == pegasussession.MessageTool && strings.TrimSpace(messages[i].Result) != "" {
			toolResults = append(toolResults, &messages[i])
		}
	}
	if len(toolResults) <= contextKeepRecentToolResults {
		return
	}
	for _, msg := range toolResults[:len(toolResults)-contextKeepRecentToolResults] {
		if strings.Contains(msg.Result, "<persisted-output>") {
			continue
		}
		if len([]rune(msg.Result)) > contextMicroCompactRunes {
			msg.Result = "[Earlier tool result compacted. Re-run the tool if exact output is needed.]"
		}
	}
}

func applyToolResultBudget(window *sessionContext) {
	total := 0
	for _, msg := range window.Messages {
		if msg.Role == pegasussession.MessageTool {
			total += len([]rune(msg.Result))
		}
	}
	if total <= contextToolResultBudgetRunes {
		return
	}
	for i := range window.Messages {
		if total <= contextToolResultBudgetRunes {
			return
		}
		msg := &window.Messages[i]
		if msg.Role != pegasussession.MessageTool || len([]rune(msg.Result)) <= contextPersistThresholdRunes {
			continue
		}
		original := msg.Result
		msg.Result = persistedToolResultPlaceholder(msg, original)
		total = total - len([]rune(original)) + len([]rune(msg.Result))
	}
}

func persistedToolResultPlaceholder(msg *pegasussession.Message, result string) string {
	preview := global.TrimRunes(global.TrimAndCollapse(result), 2000)
	label := global.FirstNonEmpty(msg.ToolID, msg.ToolName, msg.ID, "unknown")
	return fmt.Sprintf("<persisted-output>\nTool result %s exceeded the context budget and was compacted in session context.\nPreview:\n%s\n</persisted-output>", label, preview)
}

func autoCompactSessionContext(window sessionContext) sessionContext {
	summary := strings.Builder{}
	if strings.TrimSpace(window.Summary) != "" {
		summary.WriteString(window.Summary)
		summary.WriteString("\n\n")
	}
	summary.WriteString("Auto-compacted session context from older messages.\n\n")
	for _, msg := range window.Messages {
		line := compactContextText(compactMessageRole(msg), msg.Content, msg.Result, msg.Args)
		if line == "" {
			continue
		}
		next := "- " + line + "\n"
		if len([]rune(summary.String()+next)) > contextAutoSummaryMaxRunes {
			summary.WriteString("- ... session context truncated during auto compaction\n")
			break
		}
		summary.WriteString(next)
	}
	window.Summary = strings.TrimSpace(summary.String())
	if len(window.Messages) > contextKeepRecentToolResults+1 {
		window.Messages = window.Messages[len(window.Messages)-(contextKeepRecentToolResults+1):]
	}
	return window
}

func compactMessageRole(msg pegasussession.Message) string {
	role := string(msg.Role)
	if msg.ToolName != "" {
		role += "/" + msg.ToolName
	}
	return role
}

func estimateSessionContextRunes(window sessionContext) int {
	total := len([]rune(window.Summary))
	for _, msg := range window.Messages {
		total += len([]rune(msg.Content)) + len([]rune(msg.Result)) + len([]rune(msg.Args)) + len([]rune(msg.ToolName)) + 16
	}
	for _, file := range window.Files {
		total += len([]rune(file.Path)) + 16
	}
	for _, ref := range window.References {
		total += len([]rune(ref.SessionID)) + len([]rune(ref.Preview)) + 16
	}
	return total
}

func cloneSessionMessages(messages []pegasussession.Message) []pegasussession.Message {
	if len(messages) == 0 {
		return nil
	}
	out := make([]pegasussession.Message, len(messages))
	copy(out, messages)
	return out
}

func hasSessionContext(window sessionContext) bool {
	return strings.TrimSpace(window.Summary) != "" || len(window.Messages) > 0 || len(window.Files) > 0 || len(window.References) > 0
}

func compactContextText(values ...string) string {
	text := strings.TrimSpace(strings.Join(values, " "))
	text = global.TrimAndCollapse(text)
	return global.TrimRunes(text, contextMessagePreviewMaxRunes)
}
