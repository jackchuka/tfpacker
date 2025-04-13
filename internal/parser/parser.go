package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	"github.com/jackchuka/tfpacker/internal/logger"
)

type BlockInfo struct {
	Type       string
	SubType    string
	Name       string
	SourceFile string
	// Block content
	Content []byte
}

type Parser struct {
	parser   *hclparse.Parser
	excludes []string
	log      *logger.Logger
}

func NewParser(excludes []string, log *logger.Logger) *Parser {
	return &Parser{
		parser:   hclparse.NewParser(),
		excludes: excludes,
		log:      log,
	}
}

func (p *Parser) ParseDirectory(dir string) ([]BlockInfo, error) {
	var blocks []BlockInfo

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !entry.Type().IsRegular() {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		// Skip non-Terraform files
		if !strings.HasSuffix(path, ".tf") || strings.HasSuffix(path, ".tf.json") {
			continue
		}

		excluded := false
		for _, exclude := range p.excludes {
			if matched, _ := filepath.Match(exclude, entry.Name()); matched {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		fileBlocks, err := p.ParseFile(path)
		if err != nil {
			p.log.Error("Failed to parse file %s: %v", path, err)
			continue
		}

		blocks = append(blocks, fileBlocks...)
	}

	return blocks, nil
}

func (p *Parser) ParseFile(path string) ([]BlockInfo, error) {
	var blocks []BlockInfo

	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	file, diag := p.parser.ParseHCL(src, path)
	if diag.HasErrors() {
		return nil, fmt.Errorf("error parsing %s: %s", path, diag.Error())
	}

	sourceContent := src

	schema, err := tfschema.CoreModuleSchemaForVersion(tfschema.LatestAvailableVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema: %w", err)
	}
	body, diag := file.Body.Content(schema.ToHCLSchema())
	if diag.HasErrors() {
		p.log.Error("Error getting content from %s: %s", path, diag.Error())
	}

	p.log.Debug("Found %d blocks in %s", len(body.Blocks), path)

	for _, block := range body.Blocks {
		p.log.Debug("Processing block: %s with %d labels", block.Type, len(block.Labels))
		blockType, subType, name := p.processBlock(block)
		if blockType == "" {
			p.log.Debug("Skipping block with empty type")
			continue
		}

		// Extract the block content using the extractBlockContent function
		content, err := extractBlockContent(sourceContent, block)
		if err != nil {
			p.log.Error("Failed to extract block content: %v", err)
			continue
		}

		blocks = append(blocks, BlockInfo{
			Type:       blockType,
			SubType:    subType,
			Name:       name,
			SourceFile: path,
			Content:    content,
		})

		p.log.Debug("Added block: %s %s %s", blockType, subType, name)
	}

	return blocks, nil
}

func extractBlockContent(sourceContent []byte, block *hcl.Block) ([]byte, error) {
	startPos := block.DefRange.Start
	startOffset := startPos.Byte

	if startOffset < 0 || startOffset >= len(sourceContent) {
		return nil, fmt.Errorf("invalid start offset: %d", startOffset)
	}

	openBracePos := -1
	for i := startOffset; i < len(sourceContent); i++ {
		if sourceContent[i] == '{' {
			openBracePos = i
			break
		}
	}

	if openBracePos == -1 {
		return nil, fmt.Errorf("opening brace not found")
	}

	braceCount := 1
	endOffset := -1

	for i := openBracePos + 1; i < len(sourceContent); i++ {
		if sourceContent[i] == '{' {
			braceCount++
		} else if sourceContent[i] == '}' {
			braceCount--
			if braceCount == 0 {
				endOffset = i + 1
				break
			}
		}
	}

	if endOffset == -1 {
		return nil, fmt.Errorf("closing brace not found")
	}

	return sourceContent[startOffset:endOffset], nil
}

func (p *Parser) processBlock(block *hcl.Block) (string, string, string) {
	labels := block.Labels
	var subType, name string

	switch block.Type {
	case "resource", "data":
		if len(labels) < 2 {
			return "", "", ""
		}
		subType, name = labels[0], labels[1]
	case "module", "variable", "output", "provider":
		if len(labels) < 1 {
			return "", "", ""
		}
		subType = ""
		name = labels[0]
	case "locals":
		subType = ""
		name = ""
	default:
		if len(labels) > 0 {
			name = labels[0]
		}
	}

	return block.Type, subType, name
}
