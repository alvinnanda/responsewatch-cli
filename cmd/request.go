package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/boscod/responsewatch-cli/internal/api"
	"github.com/boscod/responsewatch-cli/internal/models"
	"github.com/boscod/responsewatch-cli/internal/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	reqTitle       string
	reqDesc        string
	reqGroup       int
	reqPIC         string
	reqRefLink     string
	reqPin         bool
	reqScheduled   bool
	reqScheduledAt string
	reqStatus      string
	reqPage        int
	reqLimit       int
	reqSearch      string
)

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"req"},
	Short:   "Request management",
	Long:    `Create, view, update, and manage your tickets/requests.`,
}

// requestListCmd represents the request list command
var requestListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all requests",
	Long:    `List all your tickets with optional filtering by status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Build query params
		query := fmt.Sprintf("?page=%d&limit=%d", reqPage, reqLimit)
		if reqStatus != "" {
			query += "&status=" + reqStatus
		}
		if reqSearch != "" {
			query += "&search=" + reqSearch
		}

		var result models.ListResponse
		if err := client.Get("/api/requests"+query, &result, true); err != nil {
			return fmt.Errorf("failed to list requests: %w", err)
		}

		requests, ok := result.Data.([]interface{})
		if !ok {
			// Try direct array
			if err := client.Get("/api/requests"+query, &requests, true); err != nil {
				return fmt.Errorf("failed to parse requests: %w", err)
			}
		}

		if len(requests) == 0 {
			formatter.PrintInfo("No requests found")
			return nil
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(result.Data)
		}

		// Print table
		headers := []string{"ID", "TITLE", "STATUS", "GROUP", "PIC", "CREATED"}
		rows := [][]string{}

		for _, r := range requests {
			reqMap, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			id := fmt.Sprintf("%.0f", reqMap["id"])
			title := output.TruncateString(getString(reqMap, "title"), 30)
			status := getString(reqMap, "status")
			group := getString(reqMap, "group_name")
			pic := getString(reqMap, "pic")
			if pic == "" {
				pic = "-"
			}
			created := formatTime(getString(reqMap, "created_at"))

			rows = append(rows, []string{
				id,
				title,
				formatter.StatusColor(status),
				group,
				pic,
				created,
			})
		}

		color.Cyan("\nYour Requests (Total: %d)\n", result.Total)
		formatter.PrintTable(headers, rows)
		fmt.Println()

		return nil
	},
}

// requestGetCmd represents the request get command
var requestGetCmd = &cobra.Command{
	Use:     "get [ID or UUID]",
	Aliases: []string{"show", "view"},
	Short:   "Get request details",
	Long:    `Get detailed information about a specific request.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var req models.Request
		if err := client.Get("/api/requests/"+id, &req, true); err != nil {
			return fmt.Errorf("failed to get request: %w", err)
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(req)
		}

		// Print details
		printRequestDetails(&req)
		return nil
	},
}

// requestCreateCmd represents the request create command
var requestCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new request",
	Long:    `Create a new ticket/request with interactive prompts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)

		// Interactive mode if flags not provided
		if reqTitle == "" {
			fmt.Print("Title (required): ")
			title, _ := reader.ReadString('\n')
			reqTitle = strings.TrimSpace(title)
		}

		if reqTitle == "" {
			return fmt.Errorf("title is required")
		}

		if reqDesc == "" {
			fmt.Print("Description (optional): ")
			desc, _ := reader.ReadString('\n')
			reqDesc = strings.TrimSpace(desc)
		}

		if reqRefLink == "" {
			fmt.Print("Ref Link (optional): ")
			link, _ := reader.ReadString('\n')
			reqRefLink = strings.TrimSpace(link)
		}

		if !cmd.Flags().Changed("pin") {
			fmt.Print("Pin this request? [y/N]: ")
			pin, _ := reader.ReadString('\n')
			reqPin = strings.ToLower(strings.TrimSpace(pin)) == "y"
		}

		// Build request
		req := models.CreateRequestRequest{
			Title:       reqTitle,
			Description: reqDesc,
			RefLink:     reqRefLink,
			IsPinned:    reqPin,
		}

		if reqGroup > 0 {
			req.GroupID = &reqGroup
		}
		if reqPIC != "" {
			req.PIC = reqPIC
		}
		if reqScheduled && reqScheduledAt != "" {
			req.IsScheduled = true
			req.ScheduledDate = reqScheduledAt
		}

		var created models.Request
		if err := client.Post("/api/requests", req, &created, true); err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Request created: %s", created.Title))
		fmt.Printf("\nPublic URL: %s\n", created.PublicURL)
		fmt.Printf("Token: %s\n\n", created.URLToken)

		return nil
	},
}

// requestUpdateCmd represents the request update command
var requestUpdateCmd = &cobra.Command{
	Use:     "update [ID or UUID]",
	Aliases: []string{"edit"},
	Short:   "Update a request",
	Long:    `Update an existing request. Uses interactive prompts for fields not provided as flags.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Get current request
		var current models.Request
		if err := client.Get("/api/requests/"+id, &current, true); err != nil {
			return fmt.Errorf("failed to get request: %w", err)
		}

		reader := bufio.NewReader(os.Stdin)
		req := models.UpdateRequestRequest{}

		// Interactive mode for fields not provided via flags
		if cmd.Flags().Changed("title") {
			req.Title = reqTitle
		} else {
			fmt.Printf("Title [%s]: ", current.Title)
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			if title != "" {
				req.Title = title
			}
		}

		if cmd.Flags().Changed("desc") {
			req.Description = reqDesc
		} else {
			fmt.Printf("Description [%s]: ", current.Description)
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)
			if desc != "" {
				req.Description = desc
			}
		}

		if cmd.Flags().Changed("ref-link") {
			req.RefLink = reqRefLink
		}

		if cmd.Flags().Changed("group") && reqGroup > 0 {
			req.GroupID = &reqGroup
		}

		if cmd.Flags().Changed("pic") {
			req.PIC = reqPIC
		}

		var updated models.Request
		if err := client.Put("/api/requests/"+id, req, &updated, true); err != nil {
			return fmt.Errorf("failed to update request: %w", err)
		}

		formatter.PrintSuccess("Request updated successfully")
		return nil
	},
}

// requestDeleteCmd represents the request delete command
var requestDeleteCmd = &cobra.Command{
	Use:     "delete [UUID]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a request",
	Long:    `Soft delete a request by its UUID.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		uuid := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Confirm deletion
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to delete request %s? [y/N]: ", uuid)
		confirm, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			formatter.PrintInfo("Deletion cancelled")
			return nil
		}

		if err := client.Delete("/api/requests/"+uuid, nil, true); err != nil {
			return fmt.Errorf("failed to delete request: %w", err)
		}

		formatter.PrintSuccess("Request deleted successfully")
		return nil
	},
}

// requestReopenCmd represents the request reopen command
var requestReopenCmd = &cobra.Command{
	Use:   "reopen [ID]",
	Short: "Reopen a completed request",
	Long:  `Reopen a request that has been marked as done.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var req models.Request
		if err := client.Put("/api/requests/"+id+"/reopen", nil, &req, true); err != nil {
			return fmt.Errorf("failed to reopen request: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Request reopened: %s", req.Title))
		return nil
	},
}

// requestAssignCmd represents the request assign command
var requestAssignCmd = &cobra.Command{
	Use:   "assign [ID]",
	Short: "Assign vendor/PIC to a request",
	Long:  `Assign a vendor group or specific PIC to a request.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		req := models.AssignVendorRequest{}
		if reqGroup > 0 {
			req.GroupID = &reqGroup
		}
		if reqPIC != "" {
			req.PIC = reqPIC
		}

		var updated models.Request
		if err := client.Put("/api/requests/"+id+"/assign-vendor", req, &updated, true); err != nil {
			return fmt.Errorf("failed to assign vendor: %w", err)
		}

		formatter.PrintSuccess("Vendor assigned successfully")
		return nil
	},
}

// requestStatsCmd represents the request stats command
var requestStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show request statistics",
	Long:  `Display statistics about your requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		premium, _ := cmd.Flags().GetBool("premium")

		if premium {
			var stats models.RequestStatsPremium
			if err := client.Get("/api/requests/stats/premium", &stats, true); err != nil {
				return fmt.Errorf("failed to get premium stats: %w", err)
			}

			if outputFmt == "json" {
				return formatter.PrintJSON(stats)
			}

			printPremiumStats(&stats)
		} else {
			var stats models.RequestStats
			if err := client.Get("/api/requests/stats", &stats, true); err != nil {
				return fmt.Errorf("failed to get stats: %w", err)
			}

			if outputFmt == "json" {
				return formatter.PrintJSON(stats)
			}

			color.Cyan("\nRequest Statistics\n")
			color.Cyan("==================\n")
			fmt.Printf("Total:      %d\n", stats.Total)
			fmt.Printf("Waiting:    %d\n", stats.Waiting)
			fmt.Printf("In Progress: %d\n", stats.InProgress)
			fmt.Printf("Done:       %d\n", stats.Done)
			fmt.Println()
		}

		return nil
	},
}

// requestExportCmd represents the request export command
var requestExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export requests to Excel",
	Long:  `Download all requests as an Excel file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			outputFile = "requests_export.xlsx"
		}

		// Download file
		resp, err := client.HTTPClient.Get(client.BaseURL + "/api/requests/export")
		if err != nil {
			return fmt.Errorf("failed to export: %w", err)
		}
		defer resp.Body.Close()

		// Save to file
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()

		_, err = file.ReadFrom(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Exported to %s", outputFile))
		return nil
	},
}

// requestStartCmd represents the request start command (public)
var requestStartCmd = &cobra.Command{
	Use:   "start [TOKEN]",
	Short: "Start working on a request (public)",
	Long:  `Mark a request as in-progress using its public token.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]

		client := api.NewClient(cfg)

		reader := bufio.NewReader(os.Stdin)

		// Get PIC name if not provided
		pic, _ := cmd.Flags().GetString("pic")
		if pic == "" {
			fmt.Print("Your Name: ")
			pic, _ = reader.ReadString('\n')
			pic = strings.TrimSpace(pic)
		}

		req := models.PublicRequestAction{
			PIC: pic,
		}

		var result models.Request
		if err := client.Post("/api/public/t/"+token+"/start", req, &result, false); err != nil {
			return fmt.Errorf("failed to start request: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Request started: %s", result.Title))
		return nil
	},
}

// requestFinishCmd represents the request finish command (public)
var requestFinishCmd = &cobra.Command{
	Use:   "finish [TOKEN]",
	Short: "Finish working on a request (public)",
	Long:  `Mark a request as done using its public token.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]

		client := api.NewClient(cfg)

		reader := bufio.NewReader(os.Stdin)

		// Get notes if not provided
		notes, _ := cmd.Flags().GetString("notes")
		if notes == "" {
			fmt.Print("Resolution Notes (optional): ")
			notes, _ = reader.ReadString('\n')
			notes = strings.TrimSpace(notes)
		}

		req := models.PublicRequestAction{
			Notes: notes,
		}

		var result models.Request
		if err := client.Post("/api/public/t/"+token+"/finish", req, &result, false); err != nil {
			return fmt.Errorf("failed to finish request: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Request completed: %s", result.Title))
		return nil
	},
}

// Helper functions
func printRequestDetails(req *models.Request) {
	color.Cyan("\nRequest Details\n")
	color.Cyan("===============\n")

	data := map[string]string{
		"ID":      fmt.Sprintf("%d", req.ID),
		"UUID":    req.UUID,
		"Title":   req.Title,
		"Status":  formatter.StatusColor(req.Status),
		"Group":   req.GroupName,
		"PIC":     req.PIC,
		"Ref Link": req.RefLink,
	}

	if req.StartPIC != "" {
		data["Started By"] = req.StartPIC
	}
	if req.EndPIC != "" {
		data["Finished By"] = req.EndPIC
	}

	formatter.PrintKV(data)

	if req.Description != "" {
		fmt.Printf("\nDescription:\n%s\n", req.Description)
	}

	if req.DurationSeconds != nil && *req.DurationSeconds > 0 {
		fmt.Printf("\nDuration: %s\n", output.DurationString(*req.DurationSeconds))
	}
	if req.ResponseTimeSeconds != nil && *req.ResponseTimeSeconds > 0 {
		fmt.Printf("Response Time: %s\n", output.DurationString(*req.ResponseTimeSeconds))
	}

	fmt.Printf("\nPublic URL: %s\n", req.PublicURL)
	fmt.Println()
}

func printPremiumStats(stats *models.RequestStatsPremium) {
	color.Cyan("\nPremium Statistics\n")
	color.Cyan("==================\n")
	fmt.Printf("Total Requests:      %d\n", stats.Total)
	fmt.Printf("Waiting:             %d\n", stats.Waiting)
	fmt.Printf("In Progress:         %d\n", stats.InProgress)
	fmt.Printf("Done:                %d\n", stats.Done)
	fmt.Printf("Avg Response Time:   %.1f min\n", stats.AvgResponseTimeMinutes)
	fmt.Printf("Avg Duration:        %.1f min\n", stats.AvgDurationMinutes)
	fmt.Println()
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func formatTime(t string) string {
	if len(t) > 10 {
		return t[:10]
	}
	return t
}

func init() {
	rootCmd.AddCommand(requestCmd)

	// Add subcommands
	requestCmd.AddCommand(requestListCmd)
	requestCmd.AddCommand(requestGetCmd)
	requestCmd.AddCommand(requestCreateCmd)
	requestCmd.AddCommand(requestUpdateCmd)
	requestCmd.AddCommand(requestDeleteCmd)
	requestCmd.AddCommand(requestReopenCmd)
	requestCmd.AddCommand(requestAssignCmd)
	requestCmd.AddCommand(requestStatsCmd)
	requestCmd.AddCommand(requestExportCmd)
	requestCmd.AddCommand(requestStartCmd)
	requestCmd.AddCommand(requestFinishCmd)

	// Flags for list
	requestListCmd.Flags().StringVar(&reqStatus, "status", "", "Filter by status: waiting|in_progress|done")
	requestListCmd.Flags().IntVar(&reqPage, "page", 1, "Page number")
	requestListCmd.Flags().IntVar(&reqLimit, "limit", 20, "Items per page")
	requestListCmd.Flags().StringVar(&reqSearch, "search", "", "Search query")

	// Flags for create
	requestCreateCmd.Flags().StringVar(&reqTitle, "title", "", "Request title")
	requestCreateCmd.Flags().StringVar(&reqDesc, "desc", "", "Description")
	requestCreateCmd.Flags().IntVar(&reqGroup, "group", 0, "Vendor group ID")
	requestCreateCmd.Flags().StringVar(&reqPIC, "pic", "", "PIC name")
	requestCreateCmd.Flags().StringVar(&reqRefLink, "ref-link", "", "Reference link")
	requestCreateCmd.Flags().BoolVar(&reqPin, "pin", false, "Pin this request")
	requestCreateCmd.Flags().BoolVar(&reqScheduled, "scheduled", false, "Schedule this request")
	requestCreateCmd.Flags().StringVar(&reqScheduledAt, "scheduled-at", "", "Scheduled date (YYYY-MM-DD)")

	// Flags for update
	requestUpdateCmd.Flags().StringVar(&reqTitle, "title", "", "Request title")
	requestUpdateCmd.Flags().StringVar(&reqDesc, "desc", "", "Description")
	requestUpdateCmd.Flags().StringVar(&reqRefLink, "ref-link", "", "Reference link")
	requestUpdateCmd.Flags().IntVar(&reqGroup, "group", 0, "Vendor group ID")
	requestUpdateCmd.Flags().StringVar(&reqPIC, "pic", "", "PIC name")

	// Flags for assign
	requestAssignCmd.Flags().IntVar(&reqGroup, "group-id", 0, "Vendor group ID")
	requestAssignCmd.Flags().StringVar(&reqPIC, "pic", "", "PIC name")

	// Flags for stats
	requestStatsCmd.Flags().Bool("premium", false, "Show premium statistics")

	// Flags for export
	requestExportCmd.Flags().StringP("output", "o", "requests_export.xlsx", "Output file name")

	// Flags for public actions
	requestStartCmd.Flags().String("pic", "", "PIC name")
	requestFinishCmd.Flags().String("notes", "", "Resolution notes")
}
