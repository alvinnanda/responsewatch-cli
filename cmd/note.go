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

var (
	noteTitle    string
	noteContent  string
	noteColor    string
	noteReminder string
)

// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:     "note",
	Aliases: []string{"notes"},
	Short:   "Note management",
	Long:    `Create, view, update, and manage your personal notes.`,
}

// noteListCmd represents the note list command
var noteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all notes",
	Long:    `List all your personal notes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var notes []models.Note
		if err := client.Get("/notes", &notes, true); err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(notes) == 0 {
			formatter.PrintInfo("No notes found")
			return nil
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(notes)
		}

		// Print table
		headers := []string{"ID", "TITLE", "COLOR", "REMINDER", "CREATED"}
		rows := [][]string{}

		for _, n := range notes {
			reminder := "-"
			if n.Reminder != nil {
				reminder = n.Reminder.Format("2006-01-02")
			}

			rows = append(rows, []string{
				fmt.Sprintf("%d", n.ID),
				truncateString(n.Title, 30),
				n.Color,
				reminder,
				n.CreatedAt.Format("2006-01-02"),
			})
		}

		color.Cyan("\nYour Notes\n")
		formatter.PrintTable(headers, rows)
		fmt.Println()

		return nil
	},
}

// noteCreateCmd represents the note create command
var noteCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new note",
	Long:    `Create a new personal note with optional reminder.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)

		// Interactive mode if flags not provided
		if noteTitle == "" {
			fmt.Print("Title (required): ")
			title, _ := reader.ReadString('\n')
			noteTitle = strings.TrimSpace(title)
		}

		if noteTitle == "" {
			return fmt.Errorf("title is required")
		}

		if noteContent == "" {
			fmt.Println("Content (empty line to finish):")
			var lines []string
			for {
				line, _ := reader.ReadString('\n')
				line = strings.TrimRight(line, "\n")
				if line == "" {
					break
				}
				lines = append(lines, line)
			}
			noteContent = strings.Join(lines, "\n")
		}

		if noteColor == "" {
			fmt.Print("Color (yellow/blue/green/red/purple, default: yellow): ")
			color, _ := reader.ReadString('\n')
			noteColor = strings.TrimSpace(color)
			if noteColor == "" {
				noteColor = "yellow"
			}
		}

		if noteReminder == "" {
			fmt.Print("Reminder (YYYY-MM-DD, optional): ")
			reminder, _ := reader.ReadString('\n')
			noteReminder = strings.TrimSpace(reminder)
		}

		req := models.CreateNoteRequest{
			Title:    noteTitle,
			Content:  noteContent,
			Color:    noteColor,
			Reminder: noteReminder,
		}

		var created models.Note
		if err := client.Post("/notes", req, &created, true); err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Note created: %s", created.Title))
		return nil
	},
}

// noteUpdateCmd represents the note update command
var noteUpdateCmd = &cobra.Command{
	Use:     "update [ID]",
	Aliases: []string{"edit"},
	Short:   "Update a note",
	Long:    `Update an existing note.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Get current note
		var current models.Note
		if err := client.Get("/notes/"+id, &current, true); err != nil {
			return fmt.Errorf("failed to get note: %w", err)
		}

		reader := bufio.NewReader(os.Stdin)
		req := models.UpdateNoteRequest{}

		// Title
		if cmd.Flags().Changed("title") {
			req.Title = noteTitle
		} else {
			fmt.Printf("Title [%s]: ", current.Title)
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			if title != "" {
				req.Title = title
			}
		}

		// Content
		if cmd.Flags().Changed("content") {
			req.Content = noteContent
		} else {
			fmt.Printf("Content [%s]: ", truncateString(current.Content, 30))
			fmt.Println("\n(Enter new content or press Enter to keep current, empty line to finish):")
			var lines []string
			for {
				line, _ := reader.ReadString('\n')
				line = strings.TrimRight(line, "\n")
				if line == "" {
					break
				}
				lines = append(lines, line)
			}
			if len(lines) > 0 {
				req.Content = strings.Join(lines, "\n")
			}
		}

		// Color
		if cmd.Flags().Changed("color") {
			req.Color = noteColor
		}

		// Reminder
		if cmd.Flags().Changed("reminder") {
			req.Reminder = noteReminder
		}

		var updated models.Note
		if err := client.Put("/notes/"+id, req, &updated, true); err != nil {
			return fmt.Errorf("failed to update note: %w", err)
		}

		formatter.PrintSuccess("Note updated successfully")
		return nil
	},
}

// noteDeleteCmd represents the note delete command
var noteDeleteCmd = &cobra.Command{
	Use:     "delete [ID]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a note",
	Long:    `Delete a note by its ID.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Confirm deletion
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to delete note %s? [y/N]: ", id)
		confirm, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			formatter.PrintInfo("Deletion cancelled")
			return nil
		}

		if err := client.Delete("/notes/"+id, nil, true); err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}

		formatter.PrintSuccess("Note deleted successfully")
		return nil
	},
}

// noteRemindersCmd represents the note reminders command
var noteRemindersCmd = &cobra.Command{
	Use:   "reminders",
	Short: "View upcoming reminders",
	Long:  `Display all notes with upcoming reminders.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var notes []models.Note
		if err := client.Get("/notes/reminders", &notes, true); err != nil {
			return fmt.Errorf("failed to get reminders: %w", err)
		}

		if len(notes) == 0 {
			formatter.PrintInfo("No upcoming reminders")
			return nil
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(notes)
		}

		color.Cyan("\nUpcoming Reminders\n")
		color.Cyan("==================\n")

		for _, n := range notes {
			if n.Reminder != nil {
				color.Yellow("📅 %s - %s\n", n.Reminder.Format("2006-01-02"), n.Title)
			}
		}
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)

	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteCreateCmd)
	noteCmd.AddCommand(noteUpdateCmd)
	noteCmd.AddCommand(noteDeleteCmd)
	noteCmd.AddCommand(noteRemindersCmd)

	// Flags for create
	noteCreateCmd.Flags().StringVar(&noteTitle, "title", "", "Note title")
	noteCreateCmd.Flags().StringVar(&noteContent, "content", "", "Note content")
	noteCreateCmd.Flags().StringVar(&noteColor, "color", "", "Note color (yellow/blue/green/red/purple)")
	noteCreateCmd.Flags().StringVar(&noteReminder, "reminder", "", "Reminder date (YYYY-MM-DD)")

	// Flags for update
	noteUpdateCmd.Flags().StringVar(&noteTitle, "title", "", "Note title")
	noteUpdateCmd.Flags().StringVar(&noteContent, "content", "", "Note content")
	noteUpdateCmd.Flags().StringVar(&noteColor, "color", "", "Note color")
	noteUpdateCmd.Flags().StringVar(&noteReminder, "reminder", "", "Reminder date (YYYY-MM-DD)")
}
