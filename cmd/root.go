package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/boscod/responsewatch-cli/internal/config"
	"github.com/boscod/responsewatch-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	apiURL    string
	outputFmt string
	noColor   bool
	debugMode bool
	version   = "dev" // Set by build flags

	cfg       *config.Config
	formatter *output.Formatter
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "rwcli",
	Short: "ResponseWatch CLI - Manage tickets from your terminal",
	Long: `rwcli is the official command-line interface for ResponseWatch.

It allows you to manage tickets, monitor responses, and handle vendor groups
directly from your terminal.

Get started:
  rwcli login                    # Login to your account
  rwcli request list             # List all your tickets
  rwcli monitor                  # View monitoring dashboard

For more information, visit: https://response-watch.web.app`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// For all commands, load config first
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}

		// Apply command-line overrides
		applyConfigOverrides()

		// Initialize formatter
		formatter = output.NewFormatter(outputFmt, !noColor)

		return nil
	},
}

func applyConfigOverrides() {
	if apiURL != "" {
		cfg.API.BaseURL = apiURL
	}
	if outputFmt != "" {
		cfg.Output.Format = outputFmt
	}
	if noColor {
		cfg.Output.Color = false
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.responsewatch/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API base URL (default: https://response-watch.web.app/api)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "table", "Output format: table|json")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug mode")
}

// GetConfig returns the current configuration
func GetConfig() *config.Config {
	return cfg
}

// GetFormatter returns the current formatter
func GetFormatter() *output.Formatter {
	return formatter
}

// truncateString truncates a string to max length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// formatTime formats a time string (RFC3339) to a shorter format
func formatTime(t string) string {
	if t == "" {
		return "-"
	}
	// If it's RFC3339 format, extract just the date part
	if len(t) >= 10 {
		return t[:10] // Return YYYY-MM-DD
	}
	return t
}

// extractToken parses a token from a URL if present, otherwise returns the input
func extractToken(input string) string {
	if strings.Contains(input, "/t/") {
		parts := strings.Split(input, "/t/")
		if len(parts) > 1 {
			// Extract token and remove any query parameters
			return strings.Split(parts[1], "?")[0]
		}
	}
	// Return as is (could be a raw token, ID, or UUID)
	return input
}

// isNumeric checks if a string consists only of digits
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
