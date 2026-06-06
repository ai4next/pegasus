package tool

import (
	pegasussession "github.com/ai4next/pegasus/internal/session"
	adktool "google.golang.org/adk/tool"
)

func recordFileRead(tctx adktool.Context, path string) {
	if tctx == nil {
		return
	}
	pegasussession.AddFileRead(tctx.Actions(), path)
}

func recordFileWrite(tctx adktool.Context, path string) {
	if tctx == nil {
		return
	}
	pegasussession.AddFileWrite(tctx.Actions(), path)
}

func recordFileRevision(tctx adktool.Context, path, action, before, after string, beforeMissing bool) {
	if tctx == nil {
		return
	}
	pegasussession.AddFileRevision(tctx.Actions(), pegasussession.FileRevisionNote{
		Path:          path,
		Action:        action,
		Before:        before,
		After:         after,
		BeforeMissing: beforeMissing,
	})
}
