package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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
	excludes []string
	log      *logger.Logger
}

func NewParser(excludes []string, log *logger.Logger) *Parser {
	return &Parser{
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

	file, diags := hclsyntax.ParseConfig(src, path, hcl.InitialPos)
	if diags.HasErrors() {
		return nil, fmt.Errorf("error parsing %s: %s", path, diags.Error())
	}
	tokens, _ := hclsyntax.LexConfig(src, path, hcl.InitialPos)

	body := file.Body.(*hclsyntax.Body)
	p.log.Debug("Found %d blocks in %s", len(body.Blocks), path)

	for _, block := range body.Blocks {
		p.log.Debug("Processing block: %s with %d labels", block.Type, len(block.Labels))
		blockType, subType, name := p.processBlock(block)
		if blockType == "" {
			p.log.Debug("Skipping block with empty type")
			continue
		}

		blocks = append(blocks, BlockInfo{
			Type:       blockType,
			SubType:    subType,
			Name:       name,
			SourceFile: path,
			Content:    extractBlockContent(src, tokens, block.Range()),
		})

		p.log.Debug("Added block: %s %s %s", blockType, subType, name)
	}

	return blocks, nil
}

// extractBlockContent returns the source of the block in rng, extended to
// include the comments directly above it and a comment trailing its closing brace.
func extractBlockContent(src []byte, tokens hclsyntax.Tokens, rng hcl.Range) []byte {
	first := sort.Search(len(tokens), func(i int) bool { return tokens[i].Range.Start.Byte >= rng.Start.Byte })
	next := sort.Search(len(tokens), func(i int) bool { return tokens[i].Range.Start.Byte >= rng.End.Byte })

	start := rng.Start.Byte
	for i := first - 1; i >= 0; i-- {
		if tokens[i].Type == hclsyntax.TokenNewline && i > 0 && isBlockComment(tokens[i-1]) {
			i--
		}
		tok := tokens[i]
		if tok.Type != hclsyntax.TokenComment || !startsLine(src, tok.Range.Start.Byte) {
			break
		}
		start = tok.Range.Start.Byte
	}

	end := rng.End.Byte
	if next < len(tokens) && tokens[next].Type == hclsyntax.TokenComment && tokens[next].Range.Start.Line == rng.End.Line {
		end = tokens[next].Range.End.Byte
	}

	return bytes.TrimRight(src[start:end], "\r\n")
}

func isBlockComment(tok hclsyntax.Token) bool {
	return tok.Type == hclsyntax.TokenComment && bytes.HasPrefix(tok.Bytes, []byte("/*"))
}

func startsLine(src []byte, offset int) bool {
	for i := offset - 1; i >= 0; i-- {
		switch src[i] {
		case ' ', '\t':
			continue
		case '\n':
			return true
		default:
			return false
		}
	}
	return true
}

func (p *Parser) processBlock(block *hclsyntax.Block) (string, string, string) {
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
