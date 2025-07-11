package tfpacker

import (
	"fmt"
	"os"
	"strings"

	"github.com/jackchuka/tfpacker/cmd/tfpacker/commands"
	"github.com/jackchuka/tfpacker/internal/config"
	"github.com/jackchuka/tfpacker/internal/logger"
	"github.com/jackchuka/tfpacker/internal/parser"
	"github.com/jackchuka/tfpacker/internal/writer"
	"github.com/spf13/cobra"
)

var (
	configPath      string
	dryRun          bool
	verbose         bool
	excludePatterns string
	outputDir       string
	excludes        []string
	log             *logger.Logger
)

func Eexecute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "tfpacker [directory]",
	Short: "Organize Terraform resources into separate files",
	Long:  `tfpacker is a CLI utility that organizes Terraform resources into separate files based on customizable rules.`,
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if excludePatterns != "" {
			excludes = strings.Split(excludePatterns, ",")
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		return run(dir, configPath, outputDir, dryRun, verbose, excludes)
	},
}

func init() {
	rootCmd.Flags().StringVar(&configPath, "config", "tfpacker.config.yaml", "Path to config file")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable detailed logging")
	rootCmd.Flags().StringVar(&excludePatterns, "exclude", "", "Patterns to exclude from processing (comma-separated)")
	rootCmd.Flags().StringVar(&outputDir, "output", "output", "Directory to write output files to")

	rootCmd.AddCommand(commands.NewVersionCommand())
}

func run(
	dir,
	configPath,
	outputDir string,
	dryRun,
	verbose bool,
	excludes []string,
) error {
	log = logger.New(verbose)

	var cfg *config.Config
	var err error

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Debug("Config file %s not found, using default configuration", configPath)
		cfg = config.DefaultConfig()
	} else {
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		log.Debug("Loaded configuration from %s with %d rules", configPath, len(cfg.Rules))
	}

	p := parser.NewParser(excludes, log)

	log.Debug("Parsing Terraform files in %s", dir)

	blocks, err := p.ParseDirectory(dir)
	if err != nil {
		return fmt.Errorf("failed to parse directory: %w", err)
	}

	log.Debug("Found %d blocks in directory", len(blocks))

	w := writer.NewWriter(cfg, outputDir, dryRun, log)
	fileBlockCount, err := w.WriteBlocks(blocks)
	if err != nil {
		return fmt.Errorf("failed to write blocks: %w", err)
	}

	log.Info("%s", w.GenerateSummary(fileBlockCount))

	return nil
}
