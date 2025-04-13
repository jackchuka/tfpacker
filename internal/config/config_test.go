package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.yaml")

	configContent := `rules:
  - match_type: "resource"
    sub_type: "aws_s3_bucket"
    output_file: "s3.tf"
  - match_type: "data"
    sub_type: "aws_region"
    output_file: "aws.tf"
  - match_type: "variable"
    name_prefix: "db_"
    output_file: "database.tf"
  - name_regex: "^db_.*|.*_db$"
    ignore_type: true
    output_file: "all_database.tf"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if len(cfg.Rules) != 4 {
		t.Errorf("Expected 4 rules, got %d", len(cfg.Rules))
	}

	if cfg.Rules[0].MatchType != "resource" {
		t.Errorf("Expected MatchType 'resource', got '%s'", cfg.Rules[0].MatchType)
	}
	if cfg.Rules[0].SubType != "aws_s3_bucket" {
		t.Errorf("Expected SubType 'aws_s3_bucket', got '%s'", cfg.Rules[0].SubType)
	}
	if cfg.Rules[0].OutputFile != "s3.tf" {
		t.Errorf("Expected OutputFile 's3.tf', got '%s'", cfg.Rules[0].OutputFile)
	}

	if cfg.Rules[3].NameRegex != "^db_.*|.*_db$" {
		t.Errorf("Expected NameRegex '^db_.*|.*_db$', got '%s'", cfg.Rules[3].NameRegex)
	}
	if !cfg.Rules[3].IgnoreType {
		t.Errorf("Expected IgnoreType to be true, got false")
	}
	if cfg.Rules[3].OutputFile != "all_database.tf" {
		t.Errorf("Expected OutputFile 'all_database.tf', got '%s'", cfg.Rules[3].OutputFile)
	}

	_, err = LoadConfig(filepath.Join(tempDir, "nonexistent.yaml"))
	if err == nil {
		t.Error("Expected error when loading non-existent config, got nil")
	}

	invalidConfigPath := filepath.Join(tempDir, "invalid_config.yaml")
	err = os.WriteFile(invalidConfigPath, []byte("invalid: yaml: content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid test config file: %v", err)
	}

	_, err = LoadConfig(invalidConfigPath)
	if err == nil {
		t.Error("Expected error when loading invalid config, got nil")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if len(cfg.Rules) != 0 {
		t.Errorf("Expected 0 rules in default config, got %d", len(cfg.Rules))
	}
}

func TestMatchRule(t *testing.T) {
	cfg := &Config{
		Rules: []Rule{
			{
				MatchType:  "resource",
				SubType:    "aws_s3_bucket",
				OutputFile: "s3.tf",
			},
			{
				MatchType:  "data",
				SubType:    "aws_region",
				OutputFile: "aws.tf",
			},
			{
				MatchType:  "variable",
				NamePrefix: "db_",
				OutputFile: "database.tf",
			},
			{
				MatchType:  "resource",
				SubType:    "aws_dynamodb_table",
				NamePrefix: "prod_",
				OutputFile: "prod_dynamodb.tf",
			},
			{
				NameRegex:  "^db_.*|.*_db$",
				IgnoreType: true,
				OutputFile: "all_database.tf",
			},
			{
				MatchType:  "resource",
				NameRegex:  "^api_.*",
				OutputFile: "api_resources.tf",
			},
			{
				IgnoreType: true,
				NameRegex:  "^prod_.*",
				OutputFile: "all_prod.tf",
			},
		},
	}

	tests := []struct {
		name      string
		matchType string
		subType   string
		blockName string
		want      string
	}{
		{
			name:      "Match resource by type",
			matchType: "resource",
			subType:   "aws_s3_bucket",
			blockName: "my_bucket",
			want:      "s3.tf",
		},
		{
			name:      "Match data by type",
			matchType: "data",
			subType:   "aws_region",
			blockName: "current",
			want:      "aws.tf",
		},
		{
			name:      "Match variable by prefix",
			matchType: "variable",
			subType:   "",
			blockName: "db_password",
			want:      "database.tf",
		},
		{
			name:      "Match resource by type and prefix",
			matchType: "resource",
			subType:   "aws_dynamodb_table",
			blockName: "prod_table",
			want:      "prod_dynamodb.tf",
		},
		{
			name:      "Match resource by name regex with db prefix",
			matchType: "resource",
			subType:   "aws_rds_instance",
			blockName: "db_main",
			want:      "all_database.tf", // Should match the regex rule with ignore_type
		},
		{
			name:      "Match variable by name regex with db suffix",
			matchType: "variable",
			subType:   "",
			blockName: "main_db",
			want:      "all_database.tf", // Should match the regex rule with ignore_type
		},
		{
			name:      "Match resource by name regex with api prefix",
			matchType: "resource",
			subType:   "aws_api_gateway",
			blockName: "api_gateway",
			want:      "api_resources.tf", // Should match the regex rule without ignore_type
		},
		{
			name:      "Match data by name regex with prod prefix",
			matchType: "data",
			subType:   "aws_vpc",
			blockName: "prod_vpc",
			want:      "all_prod.tf", // Should match the ignore_type rule
		},
		{
			name:      "Default resource grouping",
			matchType: "resource",
			subType:   "aws_lambda_function",
			blockName: "my_function",
			want:      "resource_aws_lambda_function.tf",
		},
		{
			name:      "Default module grouping",
			matchType: "module",
			subType:   "",
			blockName: "vpc",
			want:      "modules.tf",
		},
		{
			name:      "Default variable grouping",
			matchType: "variable",
			subType:   "",
			blockName: "region",
			want:      "variables.tf",
		},
		{
			name:      "Default output grouping",
			matchType: "output",
			subType:   "",
			blockName: "vpc_id",
			want:      "outputs.tf",
		},
		{
			name:      "Default locals grouping",
			matchType: "locals",
			subType:   "",
			blockName: "",
			want:      "locals.tf",
		},
		{
			name:      "Default provider grouping",
			matchType: "provider",
			subType:   "",
			blockName: "aws",
			want:      "providers.tf",
		},
		{
			name:      "Other block type",
			matchType: "terraform",
			subType:   "",
			blockName: "",
			want:      "terraform.tf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.MatchRule(tt.matchType, tt.subType, tt.blockName)
			if got != tt.want {
				t.Errorf("MatchRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
