package orchestrator

import (
	"fmt"
	"strings"

	"github.com/ai4next/pegasus/internal/bus"
	"github.com/ai4next/pegasus/internal/global"
)

func Aggregate(plan Plan, queue bus.TaskQueue) (string, error) {
	var builder strings.Builder
	builder.WriteString("# ")
	builder.WriteString(global.FirstNonEmptyTrimmed(plan.Goal, plan.ID))
	builder.WriteString("\n\n")
	for _, task := range plan.Tasks {
		busTaskID := TaskBusID(&plan, task.ID)
		if busTaskID == "" {
			continue
		}
		result, ok, err := queue.TaskResult(busTaskID)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", fmt.Errorf("task %s has no result", task.ID)
		}
		title := global.FirstNonEmptyTrimmed(task.Title, task.ID)
		builder.WriteString("## ")
		builder.WriteString(title)
		builder.WriteString("\n\n")
		if strings.TrimSpace(result.Result) != "" {
			builder.WriteString(strings.TrimSpace(result.Result))
		} else {
			builder.WriteString(strings.TrimSpace(result.Summary))
		}
		builder.WriteString("\n\n")
	}
	return strings.TrimSpace(builder.String()), nil
}
