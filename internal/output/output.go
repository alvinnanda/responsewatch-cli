package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Formatter handles output formatting
type Formatter struct {
	Format string // table | json
	Color  bool
}

// NewFormatter creates a new formatter
func NewFormatter(format string, useColor bool) *Formatter {
	return &Formatter{
		Format: format,
		Color:  useColor,
	}
}

// PrintJSON prints data as JSON
func (f *Formatter) PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintSuccess prints a success message
func (f *Formatter) PrintSuccess(message string) {
	if f.Color {
		color.Green("✓ %s", message)
	} else {
		fmt.Printf("✓ %s\n", message)
	}
}

// PrintError prints an error message
func (f *Formatter) PrintError(message string) {
	if f.Color {
		color.Red("✗ %s", message)
	} else {
		fmt.Printf("✗ %s\n", message)
	}
}

// PrintWarning prints a warning message
func (f *Formatter) PrintWarning(message string) {
	if f.Color {
		color.Yellow("⚠ %s", message)
	} else {
		fmt.Printf("⚠ %s\n", message)
	}
}

// PrintInfo prints an info message
func (f *Formatter) PrintInfo(message string) {
	if f.Color {
		color.Cyan("ℹ %s", message)
	} else {
		fmt.Printf("ℹ %s\n", message)
	}
}

// Table creates a new table writer
func (f *Formatter) Table() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	
	if f.Color {
		table.SetHeaderColor(
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		)
	}
	
	return table
}

// PrintTable prints data as a table
func (f *Formatter) PrintTable(headers []string, rows [][]string) {
	if f.Format == "json" {
		// Convert to JSON format
		data := make([]map[string]string, len(rows))
		for i, row := range rows {
			item := make(map[string]string)
			for j, header := range headers {
				if j < len(row) {
					item[strings.ToLower(strings.ReplaceAll(header, " ", "_"))] = row[j]
				}
			}
			data[i] = item
		}
		f.PrintJSON(data)
		return
	}

	table := f.Table()
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	
	table.AppendBulk(rows)
	table.Render()
}

// PrintKV prints key-value pairs
func (f *Formatter) PrintKV(data map[string]string) {
	if f.Format == "json" {
		f.PrintJSON(data)
		return
	}

	maxKeyLen := 0
	for k := range data {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	for k, v := range data {
		if f.Color {
			color.Cyan(fmt.Sprintf("%-*s", maxKeyLen, k))
			fmt.Printf(" : %s\n", v)
		} else {
			fmt.Printf("%-*s : %s\n", maxKeyLen, k, v)
		}
	}
}

// StatusColor returns a colored status string
func (f *Formatter) StatusColor(status string) string {
	if !f.Color {
		return status
	}

	switch strings.ToLower(status) {
	case "waiting":
		return color.YellowString(status)
	case "in_progress", "in progress":
		return color.BlueString(status)
	case "done", "completed":
		return color.GreenString(status)
	case "error", "failed":
		return color.RedString(status)
	default:
		return status
	}
}

// DurationString formats duration in human-readable format
func DurationString(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%dm", seconds/60)
	}
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	if minutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// TruncateString truncates a string to max length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
