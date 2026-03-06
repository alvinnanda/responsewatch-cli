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
	groupName   string
	groupPhone  string
	groupPICs   []string
	groupPICNames []string
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

		// Build query params
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")
		
		query := fmt.Sprintf("?page=%d&limit=%d", page, limit)

		var result models.VendorGroupListResponse
		if err := client.Get("/vendor-groups"+query, &result, true); err != nil {
			return fmt.Errorf("failed to list groups: %w", err)
		}

		if len(result.VendorGroups) == 0 {
			formatter.PrintInfo("No vendor groups found")
			return nil
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(result)
		}

		// Print table
		headers := []string{"ID", "NAME", "PHONE", "PICS", "CREATED"}
		rows := [][]string{}

		for _, g := range result.VendorGroups {
			// Format PICs
			var picNames []string
			for _, pic := range g.PICs {
				picNames = append(picNames, pic.Name)
			}
			pics := strings.Join(picNames, ", ")
			if len(pics) > 30 {
				pics = pics[:27] + "..."
			}
			if pics == "" {
				pics = "-"
			}

			phone := g.VendorPhone
			if phone == "" {
				phone = "-"
			}

			createdAt := "-"
			if g.CreatedAt != "" {
				createdAt = formatTime(g.CreatedAt)
			}

			rows = append(rows, []string{
				fmt.Sprintf("%d", g.ID),
				g.GroupName,
				phone,
				pics,
				createdAt,
			})
		}

		color.Cyan("\nVendor Groups (Page %d/%d, Total: %d)\n", 
			result.Pagination.Page, result.Pagination.TotalPages, result.Pagination.Total)
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
		if err := client.Get("/vendor-groups/"+id, &group, true); err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(group)
		}

		// Print details
		color.Cyan("\nGroup Details\n")
		color.Cyan("=============\n")

		data := map[string]string{
			"ID":        fmt.Sprintf("%d", group.ID),
			"Name":      group.GroupName,
			"Phone":     group.VendorPhone,
			"Created":   group.CreatedAt,
			"Updated":   group.UpdatedAt,
		}
		
		if data["Phone"] == "" {
			data["Phone"] = "-"
		}

		formatter.PrintKV(data)

		fmt.Println("\nPICs:")
		if len(group.PICs) == 0 {
			fmt.Println("  (none)")
		}
		for i, pic := range group.PICs {
			phone := pic.Phone
			if phone == "" {
				phone = "-"
			}
			fmt.Printf("  %d. %s (Phone: %s)\n", i+1, pic.Name, phone)
		}
		
		// Backward compatibility pic_names
		if len(group.PICNames) > 0 && len(group.PICs) == 0 {
			fmt.Println("\nPIC Names (legacy):")
			for i, name := range group.PICNames {
				fmt.Printf("  %d. %s\n", i+1, name)
			}
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

		// Build PICs from legacy names if provided
		var pics []models.PIC
		if len(groupPICs) == 0 && len(groupPICNames) == 0 {
			fmt.Println("Enter PIC names (empty line to finish):")
			for {
				fmt.Print("  Name: ")
				name, _ := reader.ReadString('\n')
				name = strings.TrimSpace(name)
				if name == "" {
					break
				}
				fmt.Print("  Phone (optional): ")
				phone, _ := reader.ReadString('\n')
				phone = strings.TrimSpace(phone)
				
				pics = append(pics, models.PIC{Name: name, Phone: phone})
			}
		} else if len(groupPICs) > 0 {
			// From --pics flag (comma-separated names only)
			for _, name := range groupPICs {
				pics = append(pics, models.PIC{Name: name})
			}
		} else if len(groupPICNames) > 0 {
			// From --pic-names flag (legacy)
			for _, name := range groupPICNames {
				pics = append(pics, models.PIC{Name: name})
			}
		}

		req := models.CreateVendorGroupRequest{
			GroupName:   groupName,
			VendorPhone: groupPhone,
			PICs:        pics,
		}

		var created models.VendorGroup
		if err := client.Post("/vendor-groups", req, &created, true); err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Group created: %s (ID: %d)", created.GroupName, created.ID))
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
		if err := client.Get("/vendor-groups/"+id, &current, true); err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		reader := bufio.NewReader(os.Stdin)
		req := models.UpdateVendorGroupRequest{}

		// Name
		if cmd.Flags().Changed("name") {
			req.GroupName = groupName
		} else {
			fmt.Printf("Name [%s]: ", current.GroupName)
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			if name != "" {
				req.GroupName = name
			}
		}

		// Phone
		if cmd.Flags().Changed("phone") {
			req.VendorPhone = groupPhone
		} else {
			phone := current.VendorPhone
			if phone == "" {
				phone = "(none)"
			}
			fmt.Printf("Phone [%s]: ", phone)
			newPhone, _ := reader.ReadString('\n')
			newPhone = strings.TrimSpace(newPhone)
			if newPhone != "" {
				req.VendorPhone = newPhone
			}
		}

		// PICs - interactive if not provided via flags
		if !cmd.Flags().Changed("pics") && !cmd.Flags().Changed("pic-names") {
			fmt.Printf("\nCurrent PICs:\n")
			for i, pic := range current.PICs {
				fmt.Printf("  %d. %s\n", i+1, pic.Name)
			}
			fmt.Println("\nEnter new PICs (name,phone per line, empty line to keep current, 'clear' to remove all):")
			
			var newPICs []models.PIC
			for {
				fmt.Print("  Name: ")
				name, _ := reader.ReadString('\n')
				name = strings.TrimSpace(name)
				
				if name == "" {
					break
				}
				
				if name == "clear" {
					newPICs = []models.PIC{}
					break
				}
				
				fmt.Print("  Phone: ")
				phone, _ := reader.ReadString('\n')
				phone = strings.TrimSpace(phone)
				
				newPICs = append(newPICs, models.PIC{Name: name, Phone: phone})
			}
			
			if len(newPICs) > 0 || (len(newPICs) == 0 && len(current.PICs) > 0) {
				req.PICs = newPICs
			}
		} else if len(groupPICs) > 0 {
			for _, name := range groupPICs {
				req.PICs = append(req.PICs, models.PIC{Name: name})
			}
		} else if len(groupPICNames) > 0 {
			for _, name := range groupPICNames {
				req.PICs = append(req.PICs, models.PIC{Name: name})
			}
		}

		var updated models.VendorGroup
		if err := client.Put("/vendor-groups/"+id, req, &updated, true); err != nil {
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

		if err := client.Delete("/vendor-groups/"+id, nil, true); err != nil {
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

	// Flags for list
	groupListCmd.Flags().Int("page", 1, "Page number")
	groupListCmd.Flags().Int("limit", 10, "Items per page")

	// Flags for create
	groupCreateCmd.Flags().StringVar(&groupName, "name", "", "Group name")
	groupCreateCmd.Flags().StringVar(&groupPhone, "phone", "", "Phone number")
	groupCreateCmd.Flags().StringSliceVar(&groupPICs, "pics", []string{}, "PIC names (comma-separated)")
	groupCreateCmd.Flags().StringSliceVar(&groupPICNames, "pic-names", []string{}, "PIC names legacy (comma-separated)")

	// Flags for update
	groupUpdateCmd.Flags().StringVar(&groupName, "name", "", "Group name")
	groupUpdateCmd.Flags().StringVar(&groupPhone, "phone", "", "Phone number")
	groupUpdateCmd.Flags().StringSliceVar(&groupPICs, "pics", []string{}, "PIC names (comma-separated)")
	groupUpdateCmd.Flags().StringSliceVar(&groupPICNames, "pic-names", []string{}, "PIC names legacy (comma-separated)")
}
