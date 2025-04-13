package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackchuka/tfpacker/internal/config"
	"github.com/jackchuka/tfpacker/internal/logger"
	"github.com/jackchuka/tfpacker/internal/parser"
)

type Writer struct {
	cfg       *config.Config
	dryRun    bool
	log       *logger.Logger
	outputDir string
}

func NewWriter(cfg *config.Config, output string, dryRun bool, log *logger.Logger) *Writer {
	return &Writer{
		cfg:       cfg,
		dryRun:    dryRun,
		log:       log,
		outputDir: output,
	}
}

func (w *Writer) WriteBlocks(blocks []parser.BlockInfo) (map[string]int, error) {
	fileBlocksRaw := make(map[string][][]byte)
	fileBlockCount := make(map[string]int)

	w.log.Debug("WriteBlocks: Processing %d blocks", len(blocks))

	for _, block := range blocks {
		outputFile := w.cfg.MatchRule(block.Type, block.SubType, block.Name)

		fileBlocksRaw[outputFile] = append(fileBlocksRaw[outputFile], block.Content)
		fileBlockCount[outputFile]++
	}

	w.log.Debug("WriteBlocks: Writing to %d files", len(fileBlocksRaw))

	for filename, blocks := range fileBlocksRaw {
		if w.dryRun {
			w.log.Info("[DRY RUN] Would write %s (%d blocks)", filename, len(blocks))
			continue
		}

		if err := w.writeRawFile(filename, blocks); err != nil {
			return nil, fmt.Errorf("failed to write file %s: %w", filename, err)
		}

		w.log.Debug("Wrote %s (%d blocks)", filename, len(blocks))
	}

	return fileBlockCount, nil
}

func (w *Writer) writeRawFile(filename string, blocks [][]byte) error {
	var outputPath string
	if filepath.IsAbs(filename) {
		outputPath = filename
	} else {
		outputPath = filepath.Join(w.outputDir, filename)
	}

	w.log.Debug("writeRawFile: Writing %d blocks to %s", len(blocks), outputPath)

	dir := filepath.Dir(outputPath)
	if dir != "." {
		w.log.Debug("writeRawFile: Creating directory %s", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	var content []byte
	for i, block := range blocks {
		if i > 0 {
			content = append(content, '\n', '\n')
		}
		content = append(content, block...)
	}

	w.log.Debug("writeRawFile: File content length: %d bytes", len(content))

	err := os.WriteFile(outputPath, content, 0644)
	if err != nil {
		w.log.Error("writeRawFile: Error writing file: %v", err)
		return err
	}

	w.log.Debug("writeRawFile: Successfully wrote file %s", outputPath)
	return nil
}

func (w *Writer) GenerateSummary(fileBlockCount map[string]int) string {
	summary := "Summary of changes:\n"

	totalFiles := len(fileBlockCount)
	totalBlocks := 0

	for file, count := range fileBlockCount {
		totalBlocks += count
		summary += fmt.Sprintf("  - %s: %d blocks\n", file, count)
	}

	summary += fmt.Sprintf("\nTotal: %d blocks in %d files\n", totalBlocks, totalFiles)

	if w.dryRun {
		summary = "[DRY RUN] " + summary
	}

	return summary
}
