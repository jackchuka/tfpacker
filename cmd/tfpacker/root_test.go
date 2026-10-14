package tfpacker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGolden(t *testing.T) {
	dirs, err := filepath.Glob("../../test/*")
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range dirs {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			out := t.TempDir()
			if err := run(dir, filepath.Join(dir, "tfpacker.config.yaml"), out, false, false, nil); err != nil {
				t.Fatal(err)
			}

			golden := filepath.Join(dir, "golden")
			want, err := os.ReadDir(golden)
			if err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadDir(out)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(want) {
				t.Errorf("got %d files, want %d", len(got), len(want))
			}
			for _, e := range want {
				wantContent, err := os.ReadFile(filepath.Join(golden, e.Name()))
				if err != nil {
					t.Fatal(err)
				}
				gotContent, err := os.ReadFile(filepath.Join(out, e.Name()))
				if err != nil {
					t.Errorf("missing %s: %v", e.Name(), err)
					continue
				}
				if string(gotContent) != string(wantContent) {
					t.Errorf("%s:\n%s\nwant:\n%s", e.Name(), gotContent, wantContent)
				}
			}
		})
	}
}
