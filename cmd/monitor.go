package cmd

import (
	"fmt"

	"github.com/boscod/responsewatch-cli/internal/api"
	"github.com/boscod/responsewatch-cli/internal/models"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "View monitoring dashboard",
	Long: `Display a static Kanban-style view of all requests.
Shows requests grouped by status: Waiting, In Progress, and Done.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Get monitoring data
		var result models.ListResponse
		if err := client.Get("/api/requests?limit=100", &result, true); err != nil {
			return fmt.Errorf("failed to get requests: %w", err)
		}

		requests, ok := result.Data.([]interface{})
		if !ok {
			// Try direct array
			if err := client.Get("/api/requests?limit=100", &requests, true); err != nil {
				return fmt.Errorf("failed to parse requests: %w", err)
			}
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(result.Data)
		}

		// Categorize requests
		var waiting, inProgress, done []map[string]interface{}

		for _, r := range requests {
			reqMap, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			status := getString(reqMap, "status")
			switch status {
			case "waiting":
				waiting = append(waiting, reqMap)
			case "in_progress":
				inProgress = append(inProgress, reqMap)
			case "done":
				done = append(done, reqMap)
			}
		}

		// Print Kanban view
		printKanbanBoard(waiting, inProgress, done)

		return nil
	},
}

// monitorPublicCmd represents the public monitor command
var monitorPublicCmd = &cobra.Command{
	Use:   "public [username]",
	Short: "View public monitoring for a user",
	Long:  `Display public monitoring data for a specific user (no login required).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		client := api.NewClient(cfg)

		var result models.PublicMonitoringResponse
		if err := client.Get("/api/public/monitoring/"+username, &result, false); err != nil {
			return fmt.Errorf("failed to get public monitoring: %w", err)
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(result)
		}

		color.Cyan("\nPublic Monitoring: %s\n", result.Username)
		color.Cyan("====================\n")

		if len(result.Requests) == 0 {
			formatter.PrintInfo("No public requests found")
			return nil
		}

		// Categorize requests
		var waiting, inProgress, done []models.Request
		for _, r := range result.Requests {
			switch r.Status {
			case "waiting":
				waiting = append(waiting, r)
			case "in_progress":
				inProgress = append(inProgress, r)
			case "done":
				done = append(done, r)
			}
		}

		// Convert to interface for printKanbanBoard
		var w, ip, d []map[string]interface{}
		for _, r := range waiting {
			w = append(w, requestToMap(r))
		}
		for _, r := range inProgress {
			ip = append(ip, requestToMap(r))
		}
		for _, r := range done {
			d = append(d, requestToMap(r))
		}

		printKanbanBoard(w, ip, d)

		return nil
	},
}

func printKanbanBoard(waiting, inProgress, done []map[string]interface{}) {
	// Print header
	fmt.Println()
	color.Yellow("═══════════════════════════════════════════════════════════════════════════")
	color.Yellow("                        📊  MONITORING DASHBOARD")
	color.Yellow("═══════════════════════════════════════════════════════════════════════════")
	fmt.Println()

	// Column width
	colWidth := 35

	// Print headers
	printColumnHeader("⏳ WAITING", len(waiting), color.YellowString, colWidth)
	printColumnHeader("🔵 IN PROGRESS", len(inProgress), color.BlueString, colWidth)
	printColumnHeader("✅ DONE", len(done), color.GreenString, colWidth)
	fmt.Println()

	// Print separator
	sep := ""
	for i := 0; i < 3; i++ {
		sep += repeatString("─", colWidth) + "│"
	}
	fmt.Println(sep)

	// Find max rows
	maxRows := len(waiting)
	if len(inProgress) > maxRows {
		maxRows = len(inProgress)
	}
	if len(done) > maxRows {
		maxRows = len(done)
	}

	// Print rows
	for i := 0; i < maxRows; i++ {
		if i < len(waiting) {
			printCard(waiting[i], colWidth)
		} else {
			printEmptyCard(colWidth)
		}
		fmt.Print("│")

		if i < len(inProgress) {
			printCard(inProgress[i], colWidth)
		} else {
			printEmptyCard(colWidth)
		}
		fmt.Print("│")

		if i < len(done) {
			printCard(done[i], colWidth)
		} else {
			printEmptyCard(colWidth)
		}
		fmt.Println()
	}

	fmt.Println()
}

func printColumnHeader(title string, count int, colorFunc func(string, ...interface{}) string, width int) {
	text := fmt.Sprintf(" %s (%d) ", title, count)
	fmt.Print(padString(text, width))
}

func printCard(req map[string]interface{}, width int) {
	id := fmt.Sprintf("%.0f", req["id"])
	title := getString(req, "title")
	if len(title) > width-4 {
		title = title[:width-7] + "..."
	}

	card := fmt.Sprintf(" #%s %s ", id, title)
	fmt.Print(padString(card, width))
}

func printEmptyCard(width int) {
	fmt.Print(repeatString(" ", width))
}

func padString(s string, width int) string {
	if len(s) > width {
		return s[:width]
	}
	return s + repeatString(" ", width-len(s))
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func requestToMap(r models.Request) map[string]interface{} {
	return map[string]interface{}{
		"id":    float64(r.ID),
		"title": r.Title,
		"status": r.Status,
	}
}

func init() {
	rootCmd.AddCommand(monitorCmd)
	monitorCmd.AddCommand(monitorPublicCmd)
}
