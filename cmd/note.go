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
	noteTitle           string
	noteContent         string
	noteBackgroundColor string
	noteTagline         string
	noteRemindAt        string
	noteIsReminder      bool
	noteReminderChannel string
	noteWebhookURL      string
	noteWebhookPayload  string
	noteWhatsAppPhone   string
	noteRequestUUID     string
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

		// Build query params
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")
		search, _ := cmd.Flags().GetString("search")
		
		query := fmt.Sprintf("?page=%d&limit=%d", page, limit)
		if search != "" {
			query += "&search=" + search
		}

		var result models.NoteListResponse
		if err := client.Get("/notes"+query, &result, true); err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(result.Notes) == 0 {
			formatter.PrintInfo("No notes found")
			return nil
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(result)
		}

		// Print table
		headers := []string{"ID", "TITLE", "REMINDER", "LINKED REQUEST", "CREATED"}
		rows := [][]string{}

		for _, n := range result.Notes {
			reminder := "-"
			if n.RemindAt != nil && *n.RemindAt != "" {
				reminder = formatTime(*n.RemindAt)
			}

			linkedReq := "-"
			if n.Request != nil {
				linkedReq = truncateString(n.Request.Title, 20)
			}

			rows = append(rows, []string{
				truncateString(n.ID, 8),
				truncateString(n.Title, 30),
				reminder,
				linkedReq,
				formatTime(n.CreatedAt),
			})
		}

		color.Cyan("\nYour Notes (Page %d, Total: %d)\n", 
			result.Pagination.Page, result.Pagination.Total)
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

		if !cmd.Flags().Changed("reminder") {
			fmt.Print("Set reminder? [y/N]: ")
			ans, _ := reader.ReadString('\n')
			noteIsReminder = strings.ToLower(strings.TrimSpace(ans)) == "y"
		}

		if noteIsReminder && noteRemindAt == "" {
			fmt.Print("Reminder Date/Time (YYYY-MM-DDTHH:MM:SS): ")
			remindAt, _ := reader.ReadString('\n')
			noteRemindAt = strings.TrimSpace(remindAt)
		}

		if noteIsReminder && noteReminderChannel == "" {
			fmt.Print("Reminder Channel (email/webhook/whatsapp) [email]: ")
			ch, _ := reader.ReadString('\n')
			noteReminderChannel = strings.TrimSpace(ch)
			if noteReminderChannel == "" {
				noteReminderChannel = "email"
			}
		}

		if noteBackgroundColor == "" {
			fmt.Print("Background Color (yellow/blue/green/red/purple) [yellow]: ")
			color, _ := reader.ReadString('\n')
			noteBackgroundColor = strings.TrimSpace(color)
			if noteBackgroundColor == "" {
				noteBackgroundColor = "yellow"
			}
		}

		// Build request
		req := models.CreateNoteRequest{
			Title:           noteTitle,
			Content:         noteContent,
			BackgroundColor: noteBackgroundColor,
			Tagline:         noteTagline,
			IsReminder:      noteIsReminder,
		}

		if noteIsReminder && noteRemindAt != "" {
			// Convert to RFC3339 if needed
			req.RemindAt = &noteRemindAt
			req.ReminderChannel = noteReminderChannel
		}

		if noteWebhookURL != "" {
			req.WebhookURL = &noteWebhookURL
		}
		if noteWebhookPayload != "" {
			req.WebhookPayload = &noteWebhookPayload
		}
		if noteWhatsAppPhone != "" {
			req.WhatsAppPhone = &noteWhatsAppPhone
		}
		if noteRequestUUID != "" {
			req.RequestUUID = &noteRequestUUID
		}

		var created models.Note
		if err := client.Post("/notes", req, &created, true); err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Note created: %s (ID: %s)", created.Title, truncateString(created.ID, 8)))
		return nil
	},
}

// noteGetCmd represents the note get command
var noteGetCmd = &cobra.Command{
	Use:     "get [ID]",
	Aliases: []string{"show", "view"},
	Short:   "Get note details",
	Long:    `Get detailed information about a specific note.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		var note models.Note
		if err := client.Get("/notes/"+id, &note, true); err != nil {
			return fmt.Errorf("failed to get note: %w", err)
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(note)
		}

		// Print details
		color.Cyan("\nNote Details\n")
		color.Cyan("============\n")

		data := map[string]string{
			"ID":              truncateString(note.ID, 12),
			"Title":           note.Title,
			"Background Color": note.BackgroundColor,
			"Tagline":         note.Tagline,
			"Is Reminder":     fmt.Sprintf("%t", note.IsReminder),
			"Created":         formatTime(note.CreatedAt),
			"Updated":         formatTime(note.UpdatedAt),
		}
		
		if data["Background Color"] == "" {
			data["Background Color"] = "-"
		}
		if data["Tagline"] == "" {
			data["Tagline"] = "-"
		}
		
		if note.IsReminder {
			data["Reminder Channel"] = note.ReminderChannel
			if note.RemindAt != nil {
				data["Remind At"] = *note.RemindAt
			}
		}

		formatter.PrintKV(data)

		fmt.Printf("\nContent:\n%s\n", note.Content)

		if note.Request != nil {
			fmt.Printf("\nLinked Request:\n")
			fmt.Printf("  Title: %s\n", note.Request.Title)
			fmt.Printf("  UUID: %s\n", note.Request.UUID)
		}
		
		fmt.Println()

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

		// Background Color
		if cmd.Flags().Changed("background-color") {
			req.BackgroundColor = noteBackgroundColor
		}

		// Tagline
		if cmd.Flags().Changed("tagline") {
			req.Tagline = noteTagline
		}

		// Reminder
		if cmd.Flags().Changed("reminder") {
			req.IsReminder = noteIsReminder
			if noteRemindAt != "" {
				req.RemindAt = &noteRemindAt
			}
			if noteReminderChannel != "" {
				req.ReminderChannel = noteReminderChannel
			}
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
		color.Cyan("===================\n")

		for _, n := range notes {
			if n.RemindAt != nil && *n.RemindAt != "" {
				color.Yellow("📅 %s - %s\n", formatTime(*n.RemindAt), n.Title)
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
	noteCmd.AddCommand(noteGetCmd)
	noteCmd.AddCommand(noteUpdateCmd)
	noteCmd.AddCommand(noteDeleteCmd)
	noteCmd.AddCommand(noteRemindersCmd)

	// Flags for list
	noteListCmd.Flags().Int("page", 1, "Page number")
	noteListCmd.Flags().Int("limit", 10, "Items per page")
	noteListCmd.Flags().String("search", "", "Search query")

	// Flags for create
	noteCreateCmd.Flags().StringVar(&noteTitle, "title", "", "Note title")
	noteCreateCmd.Flags().StringVar(&noteContent, "content", "", "Note content")
	noteCreateCmd.Flags().StringVar(&noteBackgroundColor, "background-color", "", "Note background color")
	noteCreateCmd.Flags().StringVar(&noteTagline, "tagline", "", "Note tagline")
	noteCreateCmd.Flags().BoolVar(&noteIsReminder, "reminder", false, "Set as reminder")
	noteCreateCmd.Flags().StringVar(&noteRemindAt, "remind-at", "", "Reminder datetime (RFC3339)")
	noteCreateCmd.Flags().StringVar(&noteReminderChannel, "channel", "", "Reminder channel")
	noteCreateCmd.Flags().StringVar(&noteWebhookURL, "webhook-url", "", "Webhook URL")
	noteCreateCmd.Flags().StringVar(&noteWhatsAppPhone, "whatsapp", "", "WhatsApp phone number")
	noteCreateCmd.Flags().StringVar(&noteRequestUUID, "request-uuid", "", "Link to request UUID")

	// Flags for update
	noteUpdateCmd.Flags().StringVar(&noteTitle, "title", "", "Note title")
	noteUpdateCmd.Flags().StringVar(&noteContent, "content", "", "Note content")
	noteUpdateCmd.Flags().StringVar(&noteBackgroundColor, "background-color", "", "Note background color")
	noteUpdateCmd.Flags().StringVar(&noteTagline, "tagline", "", "Note tagline")
	noteUpdateCmd.Flags().BoolVar(&noteIsReminder, "reminder", false, "Set as reminder")
	noteUpdateCmd.Flags().StringVar(&noteRemindAt, "remind-at", "", "Reminder datetime")
}
