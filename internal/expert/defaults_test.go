package expert

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSeedDefaultCADExpertsCreatesSouls(t *testing.T) {
	dir := t.TempDir()

	if err := SeedDefaultCADExperts(dir); err != nil {
		t.Fatal(err)
	}

	for name := range DefaultCADExperts {
		data, err := os.ReadFile(filepath.Join(dir, name, "soul.md"))
		if err != nil {
			t.Fatalf("read seeded soul %s: %v", name, err)
		}
		if len(data) == 0 {
			t.Fatalf("seeded soul %s is empty", name)
		}
	}
}

func TestSeedDefaultCADExpertsDoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	soulPath := filepath.Join(dir, "geometry-kernel", "soul.md")
	if err := os.MkdirAll(filepath.Dir(soulPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(soulPath, []byte("custom soul\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := SeedDefaultCADExperts(dir); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(soulPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "custom soul\n" {
		t.Fatalf("existing soul was overwritten: %q", string(data))
	}
}
