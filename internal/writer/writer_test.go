package writer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jackchuka/tfpacker/internal/config"
	"github.com/jackchuka/tfpacker/internal/logger"
	"github.com/jackchuka/tfpacker/internal/parser"
)

func manyBlocks() []parser.BlockInfo {
	var blocks []parser.BlockInfo
	for i := range 20 {
		blocks = append(blocks, parser.BlockInfo{
			Type:    "resource",
			SubType: fmt.Sprintf("type_%02d", i),
			Name:    "x",
			Content: []byte(fmt.Sprintf(`resource "type_%02d" "x" {}`, i)),
		})
	}
	return blocks
}

func TestWriteBlocksTrailingNewline(t *testing.T) {
	out := t.TempDir()
	w := NewWriter(config.DefaultConfig(), out, false, logger.New(false))
	blocks := []parser.BlockInfo{
		{Type: "variable", Name: "a", Content: []byte(`variable "a" {}`)},
		{Type: "variable", Name: "b", Content: []byte(`variable "b" {}`)},
	}
	if _, err := w.WriteBlocks(blocks); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(filepath.Join(out, "variables.tf"))
	if err != nil {
		t.Fatal(err)
	}
	want := "variable \"a\" {}\n\nvariable \"b\" {}\n"
	if string(got) != want {
		t.Errorf("variables.tf = %q, want %q", got, want)
	}
}

func TestWriteBlocksDryRunSorted(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(false)
	log.SetOutput(&buf)
	w := NewWriter(config.DefaultConfig(), t.TempDir(), true, log)
	if _, err := w.WriteBlocks(manyBlocks()); err != nil {
		t.Fatal(err)
	}

	var files []string
	for line := range strings.Lines(buf.String()) {
		_, rest, ok := strings.Cut(line, "Would write ")
		if !ok {
			continue
		}
		files = append(files, strings.Fields(rest)[0])
	}
	if len(files) != 20 {
		t.Fatalf("got %d dry-run lines, want 20", len(files))
	}
	if !slices.IsSorted(files) {
		t.Errorf("dry-run lines not sorted: %v", files)
	}
}

func TestGenerateSummarySorted(t *testing.T) {
	w := NewWriter(config.DefaultConfig(), "", false, logger.New(false))
	counts := map[string]int{}
	for _, b := range manyBlocks() {
		counts[w.cfg.MatchRule(b.Type, b.SubType, b.Name)]++
	}

	var files []string
	for line := range strings.Lines(w.GenerateSummary(counts)) {
		if f, ok := strings.CutPrefix(line, "  - "); ok {
			files = append(files, strings.Split(f, ":")[0])
		}
	}
	if len(files) != 20 {
		t.Fatalf("got %d summary lines, want 20", len(files))
	}
	if !slices.IsSorted(files) {
		t.Errorf("summary lines not sorted: %v", files)
	}
}
