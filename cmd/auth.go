package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/boscod/responsewatch-cli/internal/api"
	"github.com/boscod/responsewatch-cli/internal/models"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	loginEmail    string
	loginPassword string
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Manage authentication: login, logout, view profile, etc.`,
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to ResponseWatch",
	Long: `Authenticate with your email and password to get an access token.

You can login interactively (recommended) or provide credentials via flags:
  rwcli login
  rwcli login --email user@example.com --password secret

⚠️  Security Note: Using --password flag will expose your password in shell
history. Interactive mode (without flags) is recommended for better security.
You can also use the RWCLI_PASSWORD environment variable as an alternative.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var email, password string

		// Priority: flags > env vars > interactive
		if loginEmail != "" && loginPassword != "" {
			email = loginEmail
			password = loginPassword
		} else {
			// Try environment variables
			if envEmail := os.Getenv("RWCLI_EMAIL"); envEmail != "" && email == "" {
				email = envEmail
			}
			if envPassword := os.Getenv("RWCLI_PASSWORD"); envPassword != "" && password == "" {
				password = envPassword
			}

			// Fall back to interactive mode for missing values
			if email == "" || password == "" {
				reader := bufio.NewReader(os.Stdin)

				if email == "" {
					fmt.Print("Email: ")
					emailInput, err := reader.ReadString('\n')
					if err != nil {
						return err
					}
					email = strings.TrimSpace(emailInput)
				}

				if password == "" {
					fmt.Print("Password: ")
					passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
					if err != nil {
						return err
					}
					fmt.Println()
					password = string(passwordBytes)
				}
			}
		}

		if email == "" || password == "" {
			return fmt.Errorf("email and password are required")
		}

		// Create API client
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		// Call login API
		resp, err := authAPI.Login(email, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		// Save credentials
		if err := authAPI.SaveLogin(resp); err != nil {
			return fmt.Errorf("failed to save credentials: %w", err)
		}

		formatter.PrintSuccess(fmt.Sprintf("Logged in as %s", resp.User.Email))
		return nil
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from ResponseWatch",
	Long:  `Invalidate the current session and remove local credentials.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		// Try to invalidate token on server (ignore errors)
		_ = authAPI.Logout()

		// Clear local credentials
		if err := authAPI.ClearAuth(); err != nil {
			return fmt.Errorf("failed to clear credentials: %w", err)
		}

		formatter.PrintSuccess("Logged out successfully")
		return nil
	},
}

// meCmd represents the me command
var meCmd = &cobra.Command{
	Use:   "me",
	Short: "View current user profile",
	Long:  `Display information about the currently authenticated user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		// Check auth
		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Get user profile
		user, err := authAPI.Me()
		if err != nil {
			return fmt.Errorf("failed to get profile: %w", err)
		}

		if outputFmt == "json" {
			return formatter.PrintJSON(user)
		}

		// Print profile in table format
		fullName := "-"
		if user.FullName != nil {
			fullName = *user.FullName
		}
		org := "-"
		if user.Organization != nil {
			org = *user.Organization
		}
		
		data := map[string]string{
			"ID":           fmt.Sprintf("%d", user.ID),
			"Username":     user.Username,
			"Email":        user.Email,
			"Full Name":    fullName,
			"Organization": org,
			"Role":         user.Role,
			"Active":       fmt.Sprintf("%t", user.IsActive),
			"Created At":   formatTime(user.CreatedAt),
		}

		color.Cyan("\nProfile Information\n")
		color.Cyan("====================\n")
		formatter.PrintKV(data)
		fmt.Println()

		return nil
	},
}

// profileCmd represents the profile command group
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Profile management",
	Long:  `Manage your profile information.`,
}

// profileUpdateCmd represents the profile update command
var profileUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update profile information",
	Long:  `Update your profile information (full name, organization).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		// Check auth
		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)

		// Get current profile for defaults
		currentUser, err := authAPI.Me()
		if err != nil {
			return err
		}

		// Get full name
		currentFullName := ""
		if currentUser.FullName != nil {
			currentFullName = *currentUser.FullName
		}
		fmt.Printf("Full Name [%s]: ", currentFullName)
		fullName, _ := reader.ReadString('\n')
		fullName = strings.TrimSpace(fullName)
		if fullName == "" {
			fullName = currentFullName
		}

		// Get organization
		currentOrg := ""
		if currentUser.Organization != nil {
			currentOrg = *currentUser.Organization
		}
		fmt.Printf("Organization [%s]: ", currentOrg)
		org, _ := reader.ReadString('\n')
		org = strings.TrimSpace(org)
		if org == "" {
			org = currentOrg
		}

		req := models.UpdateProfileRequest{}
		if fullName != "" {
			req.FullName = &fullName
		}
		if org != "" {
			req.Organization = &org
		}

		user, err := authAPI.UpdateProfile(req)
		if err != nil {
			return fmt.Errorf("failed to update profile: %w", err)
		}

		// Update local config
		if user.FullName != nil {
			cfg.User.Name = *user.FullName
			_ = cfg.Save()
		}

		formatter.PrintSuccess("Profile updated successfully")
		return nil
	},
}

// passwordCmd represents the password command
var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Password management",
	Long:  `Change your password.`,
}

// passwordChangeCmd represents the password change command
var passwordChangeCmd = &cobra.Command{
	Use:   "change",
	Short: "Change password",
	Long:  `Change your account password.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)
		authAPI := api.NewAuthAPI(client)

		// Check auth
		if err := authAPI.CheckAuth(); err != nil {
			return err
		}

		// Get current password
		fmt.Print("Current Password: ")
		currentPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()

		// Get new password
		fmt.Print("New Password: ")
		newPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()

		// Confirm new password
		fmt.Print("Confirm New Password: ")
		confirmPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()

		currentPassword := string(currentPasswordBytes)
		newPassword := string(newPasswordBytes)
		confirmPassword := string(confirmPasswordBytes)

		if newPassword != confirmPassword {
			return fmt.Errorf("new passwords do not match")
		}

		if len(newPassword) < 6 {
			return fmt.Errorf("new password must be at least 6 characters")
		}

		if err := authAPI.ChangePassword(currentPassword, newPassword); err != nil {
			return fmt.Errorf("failed to change password: %w", err)
		}

		formatter.PrintSuccess("Password changed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(meCmd)
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileUpdateCmd)
	rootCmd.AddCommand(passwordCmd)
	passwordCmd.AddCommand(passwordChangeCmd)

	// Login flags
	loginCmd.Flags().StringVar(&loginEmail, "email", "", "Email address")
	loginCmd.Flags().StringVar(&loginPassword, "password", "", "Password")
}
