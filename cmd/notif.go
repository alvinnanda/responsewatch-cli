package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/boscod/responsewatch-cli/internal/api"
	"github.com/boscod/responsewatch-cli/internal/models"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// notifCmd represents the notif command
var notifCmd = &cobra.Command{
	Use:     "notif",
	Aliases: []string{"notifications"},
	Short:   "Notification management",
	Long:    `View and manage your notifications.`,
}

// notifListCmd represents the notif list command
var notifListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all notifications",
	Long:    `List all your notifications.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Build query params
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")
		
		query := fmt.Sprintf("?page=%d&limit=%d", page, limit)

		var result models.NotificationListResponse
		if err := client.Get("/notifications"+query, &result, true); err != nil {
			return fmt.Errorf("failed to list notifications: %w", err)
		}

		if len(result.Notifications) == 0 {
			formatter.PrintInfo("No notifications found")
			return nil
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(result)
		}

		// Print table
		headers := []string{"ID", "TYPE", "TITLE", "STATUS", "CREATED"}
		rows := [][]string{}

		for _, n := range result.Notifications {
			status := "Read"
			if !n.IsRead {
				status = color.YellowString("Unread")
			}

			rows = append(rows, []string{
				fmt.Sprintf("%d", n.ID),
				n.Type,
				truncateString(n.Title, 30),
				status,
				formatTime(n.CreatedAt),
			})
		}

		color.Cyan("\nYour Notifications (Page %d/%d, Total: %d)\n",
			result.Pagination.Page, result.Pagination.TotalPages, result.Pagination.Total)
		formatter.PrintTable(headers, rows)
		fmt.Println()

		return nil
	},
}

// notifUnreadCmd represents the notif unread command
var notifUnreadCmd = &cobra.Command{
	Use:   "unread",
	Short: "Count unread notifications",
	Long:  `Display the number of unread notifications.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var result models.UnreadCountResponse
		if err := client.Get("/notifications/unread-count", &result, true); err != nil {
			return fmt.Errorf("failed to get unread count: %w", err)
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(result)
		}

		if result.Count == 0 {
			formatter.PrintInfo("You have no unread notifications")
		} else if result.Count == 1 {
			color.Yellow("📬 You have 1 unread notification")
		} else {
			color.Yellow("📬 You have %d unread notifications", result.Count)
		}
		fmt.Println()

		return nil
	},
}

// notifReadCmd represents the notif read command
var notifReadCmd = &cobra.Command{
	Use:   "read [ID]",
	Short: "Mark notification as read",
	Long:  `Mark a specific notification as read.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var result models.SuccessResponse
		if err := client.Post("/notifications/"+id+"/read", nil, &result, true); err != nil {
			return fmt.Errorf("failed to mark notification as read: %w", err)
		}

		formatter.PrintSuccess("Notification marked as read")
		return nil
	},
}

// notifReadAllCmd represents the notif read-all command
var notifReadAllCmd = &cobra.Command{
	Use:   "read-all",
	Short: "Mark all notifications as read",
	Long:  `Mark all notifications as read.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Confirm
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Mark all notifications as read? [y/N]: ")
		confirm, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			formatter.PrintInfo("Cancelled")
			return nil
		}

		var result models.SuccessResponse
		if err := client.Post("/notifications/read-all", nil, &result, true); err != nil {
			return fmt.Errorf("failed to mark all as read: %w", err)
		}

		formatter.PrintSuccess("All notifications marked as read")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(notifCmd)

	notifCmd.AddCommand(notifListCmd)
	notifCmd.AddCommand(notifUnreadCmd)
	notifCmd.AddCommand(notifReadCmd)
	notifCmd.AddCommand(notifReadAllCmd)

	// Flags for list
	notifListCmd.Flags().Int("page", 1, "Page number")
	notifListCmd.Flags().Int("limit", 20, "Items per page")
}
