package session_test

import (
	"path/filepath"
	"testing"

	"github.com/ai4next/pegasus/internal/config"
	"github.com/ai4next/pegasus/internal/global"
	pegasussession "github.com/ai4next/pegasus/internal/session"
	adksession "google.golang.org/adk/session"
)

func TestExtractFileReferences(t *testing.T) {
	got := pegasussession.ExtractFileReferences(`read @main.go and @"docs/design doc.md", ignore https://example.com/@x and duplicate @main.go.`)
	want := []string{"main.go", "docs/design doc.md"}
	if len(got) != len(want) {
		t.Fatalf("refs = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("refs = %#v, want %#v", got, want)
		}
	}
}

func TestExtractSessionReferences(t *testing.T) {
	got := pegasussession.ExtractSessionReferences(`continue from [session:past role:user] historical cache decision`)
	if len(got) != 1 || got[0].SessionID != "past" || got[0].Role != pegasussession.MessageUser || got[0].Preview != "historical cache decision" {
		t.Fatalf("refs = %#v", got)
	}
}

func TestRecordPromptReferences(t *testing.T) {
	dir := t.TempDir()
	global.SetConfig(&config.Config{Workspace: dir})
	t.Cleanup(func() { global.SetConfig(nil) })
	svc, err := pegasussession.NewService()
	if err != nil {
		t.Fatal(err)
	}
	extended := svc.(*pegasussession.Service)
	if _, err := svc.Create(t.Context(), &adksession.CreateRequest{AppName: "app", UserID: "user", SessionID: "1"}); err != nil {
		t.Fatal(err)
	}
	mainPath := filepath.Join(dir, "main.go")
	designPath := filepath.Join(dir, "docs/design doc.md")

	counts := pegasussession.RecordPromptReferences(extended, "app", "user", "1", dir, `review @main.go and @"docs/design doc.md" from [session:past role:user] cache decision`)
	if counts.Files != 2 || counts.Sessions != 1 {
		t.Fatalf("counts = %#v", counts)
	}
	files, err := extended.SessionFiles("app", "user", "1")
	if err != nil {
		t.Fatal(err)
	}
	byPath := make(map[string]pegasussession.SessionFile)
	for _, file := range files {
		byPath[file.Path] = file
	}
	if byPath[mainPath].ReadCount != 1 || byPath[designPath].ReadCount != 1 {
		t.Fatalf("files = %#v", files)
	}
	refs, err := extended.SessionReferences("app", "user", "1")
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) != 1 || refs[0].SessionID != "past" {
		t.Fatalf("refs = %#v", refs)
	}
}
