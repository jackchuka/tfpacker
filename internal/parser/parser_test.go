package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jackchuka/tfpacker/internal/logger"
)

func parseSource(t *testing.T, src string) []BlockInfo {
	t.Helper()
	path := filepath.Join(t.TempDir(), "main.tf")
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	blocks, err := NewParser(nil, logger.New(false)).ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	return blocks
}

func TestParseFileContent(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []string
	}{
		{
			name: "braces in strings",
			src: `resource "null_resource" "a" {
  x = "a}b"
  y = "{"
  z = "${var.v}}"
  after = true
}

variable "next" {}
`,
			want: []string{
				`resource "null_resource" "a" {
  x = "a}b"
  y = "{"
  z = "${var.v}}"
  after = true
}`,
				`variable "next" {}`,
			},
		},
		{
			name: "braces in heredoc",
			src: `locals {
  policy = <<EOT
}
{{
EOT
  after = true
}
`,
			want: []string{
				`locals {
  policy = <<EOT
}
{{
EOT
  after = true
}`,
			},
		},
		{
			name: "braces in comments",
			src: `resource "null_resource" "a" {
  # closing } here
  // another }
  /* and { here */
  after = true
}
`,
			want: []string{
				`resource "null_resource" "a" {
  # closing } here
  // another }
  /* and { here */
  after = true
}`,
			},
		},
		{
			name: "leading comments move with their block",
			src: `# Bucket for logs
// owned by platform
resource "aws_s3_bucket" "logs" {
  bucket = "logs"
}

/*
  Region doc
*/
variable "region" {}
`,
			want: []string{
				`# Bucket for logs
// owned by platform
resource "aws_s3_bucket" "logs" {
  bucket = "logs"
}`,
				`/*
  Region doc
*/
variable "region" {}`,
			},
		},
		{
			name: "detached and trailing comments",
			src: `# file header

variable "a" {} # about a
variable "b" {}
`,
			want: []string{
				`variable "a" {} # about a`,
				`variable "b" {}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := parseSource(t, tt.src)
			if len(blocks) != len(tt.want) {
				t.Fatalf("got %d blocks, want %d", len(blocks), len(tt.want))
			}
			for i, want := range tt.want {
				if got := string(blocks[i].Content); got != want {
					t.Errorf("block %d content:\n%s\nwant:\n%s", i, got, want)
				}
			}
		})
	}
}

func TestParseFileLabels(t *testing.T) {
	blocks := parseSource(t, `resource "aws_s3_bucket" "logs" {}
data "aws_region" "current" {}
module "vpc" {}
locals {}
terraform {}
`)
	want := []struct{ typ, sub, name string }{
		{"resource", "aws_s3_bucket", "logs"},
		{"data", "aws_region", "current"},
		{"module", "", "vpc"},
		{"locals", "", ""},
		{"terraform", "", ""},
	}
	if len(blocks) != len(want) {
		t.Fatalf("got %d blocks, want %d", len(blocks), len(want))
	}
	for i, w := range want {
		b := blocks[i]
		if b.Type != w.typ || b.SubType != w.sub || b.Name != w.name {
			t.Errorf("block %d = (%q, %q, %q), want (%q, %q, %q)", i, b.Type, b.SubType, b.Name, w.typ, w.sub, w.name)
		}
	}
}
