package cmd

import (
	"fmt"

	"github.com/boscod/responsewatch-cli/internal/api"
	"github.com/boscod/responsewatch-cli/internal/models"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// adminCmd represents the admin command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin commands",
	Long:  `Administrative commands for managing users (Admin role only).`,
}

// adminUsersCmd represents the admin users command
var adminUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users",
	Long:  `List all registered users in the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var users []models.User
		if err := client.Get("/admin/users", &users, true); err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		if len(users) == 0 {
			formatter.PrintInfo("No users found")
			return nil
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(users)
		}

		// Print table
		headers := []string{"ID", "USERNAME", "EMAIL", "FULL NAME", "ROLE", "ACTIVE"}
		rows := [][]string{}

		for _, u := range users {
			active := "Yes"
			if !u.IsActive {
				active = color.RedString("No")
			}

			rows = append(rows, []string{
				fmt.Sprintf("%d", u.ID),
				u.Username,
				u.Email,
				u.FullName,
				u.Role,
				active,
			})
		}

		color.Cyan("\nAll Users\n")
		formatter.PrintTable(headers, rows)
		fmt.Println()

		return nil
	},
}

// adminUpgradeCmd represents the admin upgrade command
var adminUpgradeCmd = &cobra.Command{
	Use:   "upgrade [USER_ID]",
	Short: "Upgrade user membership",
	Long:  `Upgrade a user's membership to premium.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		if err := client.Post("/admin/users/"+userID+"/upgrade", nil, nil, true); err != nil {
			return fmt.Errorf("failed to upgrade user: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("User %s upgraded successfully", userID))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)

	adminCmd.AddCommand(adminUsersCmd)
	adminCmd.AddCommand(adminUpgradeCmd)
}
