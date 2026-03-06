# ResponseWatch CLI (rwcli)

Official command-line interface for ResponseWatch - manage tickets from your terminal.

## Features

- 🔐 **Authentication**: Login with email/password, auto token refresh
- 📝 **Request Management**: Create, list, update, delete tickets
- 👥 **Vendor Groups**: Manage vendor groups and PICs
- 📊 **Monitoring**: Static Kanban view of all requests
- 🔔 **Notifications**: View and manage notifications
- 📝 **Notes**: Personal notes with reminders
- 🛡️ **Admin**: User management (admin only)
- 🌐 **Public Actions**: Start/finish requests via public token

## Installation

### macOS / Linux

```bash
curl -sSfL https://raw.githubusercontent.com/alvinnanda/responsewatch-cli/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
iwr -useb https://raw.githubusercontent.com/alvinnanda/responsewatch-cli/main/install.ps1 | iex
```

### Manual

Download the latest binary from [GitHub Releases](https://github.com/alvinnanda/responsewatch-cli/releases) and add to your PATH.

### Build from Source

```bash
git clone https://github.com/alvinnanda/responsewatch-cli.git
cd responsewatch-cli
go build -o rwcli .
```

## Quick Start

```bash
# Login
rwcli login

# List your requests
rwcli request list

# Create a new request
rwcli request create

# View monitoring dashboard
rwcli monitor

# Logout
rwcli logout
```

## Configuration

Configuration is stored in `~/.responsewatch/config.yaml`:

```yaml
api:
  base_url: "https://response-watch.web.app/api"
  timeout: 30
auth:
  token: "..."
  refresh_token: "..."
  expires_at: "..."
output:
  format: "table"  # table | json
  color: true
```

## Commands

### Authentication
- `rwcli login` - Login with email/password
- `rwcli logout` - Logout and clear credentials
- `rwcli me` - View current user profile
- `rwcli profile update` - Update profile
- `rwcli password change` - Change password

### Request Management
- `rwcli request list` - List all requests
- `rwcli request get <ID>` - Get request details
- `rwcli request create` - Create new request (interactive)
- `rwcli request update <ID>` - Update request
- `rwcli request delete <UUID>` - Delete request
- `rwcli request reopen <ID>` - Reopen completed request
- `rwcli request assign <ID>` - Assign vendor/PIC
- `rwcli request stats` - Show statistics
- `rwcli request export` - Export to Excel
- `rwcli request start <TOKEN>` - Start request (public)
- `rwcli request finish <TOKEN>` - Finish request (public)

### Vendor Groups
- `rwcli group list` - List all groups
- `rwcli group get <ID>` - Get group details
- `rwcli group create` - Create new group
- `rwcli group update <ID>` - Update group
- `rwcli group delete <ID>` - Delete group

### Monitoring
- `rwcli monitor` - View Kanban dashboard
- `rwcli monitor public <username>` - View public monitoring

### Notes
- `rwcli note list` - List all notes
- `rwcli note create` - Create new note
- `rwcli note update <ID>` - Update note
- `rwcli note delete <ID>` - Delete note
- `rwcli note reminders` - View upcoming reminders

### Notifications
- `rwcli notif list` - List all notifications
- `rwcli notif unread` - Count unread notifications
- `rwcli notif read <ID>` - Mark notification as read
- `rwcli notif read-all` - Mark all as read

### Admin
- `rwcli admin users` - List all users
- `rwcli admin upgrade <USER_ID>` - Upgrade user membership

## Global Flags

- `--api-url` - Custom API base URL
- `-o, --output` - Output format: `table` or `json`
- `--no-color` - Disable colored output
- `--debug` - Enable debug mode
- `--config` - Custom config file path

## Examples

```bash
# List with filter
rwcli request list --status waiting --limit 10

# Create request with flags
rwcli request create --title "Server Down" --desc "Critical issue" --pin

# Get request as JSON
rwcli request get 123 --output json

# View public monitoring
rwcli monitor public johndoe
```

## Security

- **Token Storage**: Your authentication token is stored locally at `~/.responsewatch/config.yaml` with restricted file permissions (`0600`). Keep this file secure and do not share it.
- **Password Input**: In interactive mode, passwords are never displayed on screen. Avoid using the `--password` flag in scripts as it may be visible in shell history. Use environment variable `RWCLI_PASSWORD` as a safer alternative for automation.
- **HTTPS**: All API communication uses HTTPS by default.

## License

MIT License
