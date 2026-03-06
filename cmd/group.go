package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/boscod/responsewatch-cli/internal/api"
	"github.com/boscod/responsewatch-cli/internal/models"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	groupName  string
	groupPhone string
	groupPICs  []string
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"grp"},
	Short:   "Vendor group management",
	Long:    `Create, view, update, and manage vendor groups.`,
}

// groupListCmd represents the group list command
var groupListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all vendor groups",
	Long:    `List all your vendor groups.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var groups []models.VendorGroup
		if err := client.Get("/api/vendor-groups", &groups, true); err != nil {
			return fmt.Errorf("failed to list groups: %w", err)
		}

		if len(groups) == 0 {
			formatter.PrintInfo("No vendor groups found")
			return nil
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(groups)
		}

		// Print table
		headers := []string{"ID", "NAME", "PHONE", "PICS", "CREATED"}
		rows := [][]string{}

		for _, g := range groups {
			pics := strings.Join(g.PICs, ", ")
			if len(pics) > 30 {
				pics = pics[:27] + "..."
			}

			rows = append(rows, []string{
				fmt.Sprintf("%d", g.ID),
				g.Name,
				g.Phone,
				pics,
				g.CreatedAt[:10],
			})
		}

		color.Cyan("\nVendor Groups\n")
		formatter.PrintTable(headers, rows)
		fmt.Println()

		return nil
	},
}

// groupGetCmd represents the group get command
var groupGetCmd = &cobra.Command{
	Use:     "get [ID]",
	Aliases: []string{"show", "view"},
	Short:   "Get group details",
	Long:    `Get detailed information about a specific vendor group.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var group models.VendorGroup
		if err := client.Get("/api/vendor-groups/"+id, &group, true); err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(group)
		}

		// Print details
		color.Cyan("\nGroup Details\n")
		color.Cyan("=============\n")

		data := map[string]string{
			"ID":      fmt.Sprintf("%d", group.ID),
			"Name":    group.Name,
			"Phone":   group.Phone,
			"Created": group.CreatedAt,
		}

		formatter.PrintKV(data)

		fmt.Println("\nPICs:")
		for i, pic := range group.PICs {
			fmt.Printf("  %d. %s\n", i+1, pic)
		}
		fmt.Println()

		return nil
	},
}

// groupCreateCmd represents the group create command
var groupCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new vendor group",
	Long:    `Create a new vendor group with name, phone, and PICs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)

		// Interactive mode if flags not provided
		if groupName == "" {
			fmt.Print("Group Name (required): ")
			name, _ := reader.ReadString('\n')
			groupName = strings.TrimSpace(name)
		}

		if groupName == "" {
			return fmt.Errorf("group name is required")
		}

		if groupPhone == "" {
			fmt.Print("Phone (optional): ")
			phone, _ := reader.ReadString('\n')
			groupPhone = strings.TrimSpace(phone)
		}

		if len(groupPICs) == 0 {
			fmt.Println("Enter PIC names (empty line to finish):")
			for {
				fmt.Print("  - ")
				pic, _ := reader.ReadString('\n')
				pic = strings.TrimSpace(pic)
				if pic == "" {
					break
				}
				groupPICs = append(groupPICs, pic)
			}
		}

		req := models.CreateVendorGroupRequest{
			Name:  groupName,
			Phone: groupPhone,
			PICs:  groupPICs,
		}

		var created models.VendorGroup
		if err := client.Post("/api/vendor-groups", req, &created, true); err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Group created: %s", created.Name))
		return nil
	},
}

// groupUpdateCmd represents the group update command
var groupUpdateCmd = &cobra.Command{
	Use:     "update [ID]",
	Aliases: []string{"edit"},
	Short:   "Update a vendor group",
	Long:    `Update an existing vendor group.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Get current group
		var current models.VendorGroup
		if err := client.Get("/api/vendor-groups/"+id, &current, true); err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		reader := bufio.NewReader(os.Stdin)
		req := models.UpdateVendorGroupRequest{}

		// Name
		if cmd.Flags().Changed("name") {
			req.Name = groupName
		} else {
			fmt.Printf("Name [%s]: ", current.Name)
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			if name != "" {
				req.Name = name
			}
		}

		// Phone
		if cmd.Flags().Changed("phone") {
			req.Phone = groupPhone
		} else {
			fmt.Printf("Phone [%s]: ", current.Phone)
			phone, _ := reader.ReadString('\n')
			phone = strings.TrimSpace(phone)
			if phone != "" {
				req.Phone = phone
			}
		}

		// PICs
		if len(groupPICs) > 0 {
			req.PICs = groupPICs
		} else {
			fmt.Printf("\nCurrent PICs: %s\n", strings.Join(current.PICs, ", "))
			fmt.Println("Enter new PIC names (empty line to keep current, 'clear' to remove all):")
			var newPICs []string
			for {
				fmt.Print("  - ")
				pic, _ := reader.ReadString('\n')
				pic = strings.TrimSpace(pic)
				if pic == "" {
					break
				}
				if pic == "clear" {
					newPICs = []string{}
					break
				}
				newPICs = append(newPICs, pic)
			}
			if len(newPICs) > 0 {
				req.PICs = newPICs
			}
		}

		var updated models.VendorGroup
		if err := client.Put("/api/vendor-groups/"+id, req, &updated, true); err != nil {
			return fmt.Errorf("failed to update group: %w", err)
		}

		formatter.PrintSuccess("Group updated successfully")
		return nil
	},
}

// groupDeleteCmd represents the group delete command
var groupDeleteCmd = &cobra.Command{
	Use:     "delete [ID]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a vendor group",
	Long:    `Delete a vendor group by its ID.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		// Validate ID is numeric
		if _, err := strconv.Atoi(id); err != nil {
			return fmt.Errorf("invalid group ID: %s", id)
		}

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Confirm deletion
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to delete group %s? [y/N]: ", id)
		confirm, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			formatter.PrintInfo("Deletion cancelled")
			return nil
		}

		if err := client.Delete("/api/vendor-groups/"+id, nil, true); err != nil {
			return fmt.Errorf("failed to delete group: %w", err)
		}

		formatter.PrintSuccess("Group deleted successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)

	groupCmd.AddCommand(groupListCmd)
	groupCmd.AddCommand(groupGetCmd)
	groupCmd.AddCommand(groupCreateCmd)
	groupCmd.AddCommand(groupUpdateCmd)
	groupCmd.AddCommand(groupDeleteCmd)

	// Flags for create
	groupCreateCmd.Flags().StringVar(&groupName, "name", "", "Group name")
	groupCreateCmd.Flags().StringVar(&groupPhone, "phone", "", "Phone number")
	groupCreateCmd.Flags().StringSliceVar(&groupPICs, "pics", []string{}, "PIC names (comma-separated)")

	// Flags for update
	groupUpdateCmd.Flags().StringVar(&groupName, "name", "", "Group name")
	groupUpdateCmd.Flags().StringVar(&groupPhone, "phone", "", "Phone number")
	groupUpdateCmd.Flags().StringSliceVar(&groupPICs, "pics", []string{}, "PIC names (comma-separated)")
}
