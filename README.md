# bigbluebutton-cli

A command line interface for the [BigBlueButton API](https://bbbserver.com/bigbluebutton-api-guide). Manage your [BigBlueButton](https://bigbluebutton.org) server meetings, recordings, and more from your terminal.

Provided by [bbbserver](https://bbbserver.com).

## Installation

Download the latest build for your platform:

| Platform | Download |
|---|---|
| Linux (x86_64) | [bigbluebutton-cli-linux-amd64](https://raw.githubusercontent.com/invokablegmbh/bigbluebutton-cli/refs/heads/master/dist/bigbluebutton-cli-linux-amd64) |
| Linux (ARM64) | [bigbluebutton-cli-linux-arm64](https://raw.githubusercontent.com/invokablegmbh/bigbluebutton-cli/refs/heads/master/dist/bigbluebutton-cli-linux-arm64) |
| Windows (x86_64) | [bigbluebutton-cli-windows-amd64.exe](https://raw.githubusercontent.com/invokablegmbh/bigbluebutton-cli/refs/heads/master/dist/bigbluebutton-cli-windows-amd64.exe) |

```bash
# Example: Linux x86_64
chmod +x bigbluebutton-cli-linux-amd64
sudo mv bigbluebutton-cli-linux-amd64 /usr/local/bin/bigbluebutton-cli
```

### Build from source

Requires Go 1.22+.

```bash
git clone https://github.com/invokablegmbh/bigbluebutton-cli.git
cd bigbluebutton-cli
make build        # builds ./bigbluebutton-cli
make build-all    # cross-compiles to dist/
```

## Configuration

`bigbluebutton-cli` resolves the BigBlueButton server URL and shared secret in this order:

1. **Command line flags** `--url` and `--secret`
2. **Environment variables** `BBB_URL` and `BBB_SECRET`
3. **Config file** `/etc/bigbluebutton-cli/config.yaml` or `~/.bigbluebutton-cli/config.yaml`
4. **Autodetection** from a local BigBlueButton installation (`/etc/bigbluebutton/bbb-web.properties`)
5. **Interactive wizard** (prompts you if running in a terminal)

### Config file format

```yaml
url: https://bbb.example.com
secret: your-shared-secret
```

### Environment variables

```bash
export BBB_URL=https://bbb.example.com
export BBB_SECRET=your-shared-secret
```

### Autodetection

If `bigbluebutton-cli` is run on a BigBlueButton server, it automatically reads the URL and secret from `/etc/bigbluebutton/bbb-web.properties`. No configuration needed.

## Usage

Every command works in two modes:

- **Non-interactive**: pass all required parameters as flags. Ideal for scripts and automation.
- **Interactive**: run in a terminal without required flags and `bigbluebutton-cli` will prompt you for them in a wizard-style interface.

### Global flags

```
-u, --url        BigBlueButton server URL
-s, --secret     BigBlueButton shared secret
-v, --verbose    Enable debug output (shows HTTP requests and responses)
-c, --config     Path to config file
```

### Output format

The commands `get-meetings`, `get-meeting-info`, `get-recordings`, and `get-text-tracks` support a `--format` (`-f`) flag to control the output format:

| Format | Description |
|---|---|
| `human` | Formatted table (default) |
| `json` | JSON output |
| `xml` | XML output |

```bash
# Get meetings as JSON
bigbluebutton-cli get-meetings --format json

# Get meeting info as XML
bigbluebutton-cli get-meeting-info --meeting-id standup-001 -f xml

# Get recordings as JSON (useful for piping to jq)
bigbluebutton-cli get-recordings --format json | jq '.[]'
```

### Meetings

```bash
# Create a meeting and join immediately (recommended)
bigbluebutton-cli create-join --name "Team Standup" --meeting-id standup-001 \
  --user-name "Alice" --role MODERATOR

# Create-join with recording and get just the URL
bigbluebutton-cli create-join --name "Webinar" --meeting-id webinar-042 \
  --user-name "Host" --role MODERATOR --record --url-only

# Create a meeting (without joining)
bigbluebutton-cli create --name "Team Standup" --meeting-id standup-001

# Create with options
bigbluebutton-cli create --name "Webinar" --meeting-id webinar-042 \
  --record \
  --mute-on-start \
  --guest-policy ASK_MODERATOR \
  --max-participants 100 \
  --duration 60 \
  --welcome "Welcome to the webinar!" \
  --meeting-layout PRESENTATION_FOCUS

# Generate a join URL (for an existing meeting)
bigbluebutton-cli join --meeting-id standup-001 --name "Alice" --role MODERATOR

# Get just the URL (useful for scripting)
bigbluebutton-cli join --meeting-id standup-001 --name "Bob" --role VIEWER --url-only

# Check if a meeting is running
bigbluebutton-cli is-meeting-running --meeting-id standup-001

# Get detailed meeting info (participants, settings, etc.)
bigbluebutton-cli get-meeting-info --meeting-id standup-001

# List all running meetings
bigbluebutton-cli get-meetings

# End a meeting
bigbluebutton-cli end --meeting-id standup-001

# Insert a presentation into a running meeting
bigbluebutton-cli insert-document --meeting-id standup-001 \
  --document-url https://example.com/slides.pdf \
  --filename "Q4 Review"
```

### Recordings

```bash
# List all recordings
bigbluebutton-cli get-recordings

# Filter by meeting or state
bigbluebutton-cli get-recordings --meeting-id standup-001
bigbluebutton-cli get-recordings --state published --limit 10

# Publish / unpublish
bigbluebutton-cli publish-recordings --record-id rec-abc123 --publish true
bigbluebutton-cli publish-recordings --record-id rec-abc123 --publish false

# Update metadata
bigbluebutton-cli update-recordings --record-id rec-abc123 --meta name="New Title"

# Delete
bigbluebutton-cli delete-recordings --record-id rec-abc123

# List subtitles / captions
bigbluebutton-cli get-text-tracks --record-id rec-abc123

# Upload a subtitle file
bigbluebutton-cli put-text-track --record-id rec-abc123 \
  --kind subtitles --lang en --label "English" \
  --file captions.vtt
```

### Chat

```bash
# Send a system message to a running meeting
bigbluebutton-cli send-chat-message --meeting-id standup-001 \
  --message "Meeting will end in 5 minutes"

# Send with a custom sender name
bigbluebutton-cli send-chat-message --meeting-id standup-001 \
  --message "Welcome!" --user-name "Admin Bot"
```

### Session

```bash
# Get a join URL from a session token
bigbluebutton-cli get-join-url --session-token sess-xyz123
```

### Other

```bash
# Print version
bigbluebutton-cli version

# Get help for any command
bigbluebutton-cli --help
bigbluebutton-cli create --help
```

## Scripting examples

```bash
# Create a meeting and immediately get a moderator link
LINK=$(bigbluebutton-cli create-join --name "Quick Call" --meeting-id quick-001 \
  --user-name "Host" --role MODERATOR --url-only)
echo "Join here: $LINK"

# End all running meetings
bigbluebutton-cli get-meetings 2>/dev/null | grep -oP '^\S+\s+\K\S+' | while read id; do
  bigbluebutton-cli end --meeting-id "$id"
done

# Use with environment variables in CI
BBB_URL=https://bbb.example.com BBB_SECRET=secret bigbluebutton-cli get-meetings
```

## Verbose / debug mode

Add `--verbose` (or `-v`) to any command to see the full HTTP request URLs and API responses:

```bash
bigbluebutton-cli get-meetings --verbose
```

## Shell completion

```bash
# Bash
bigbluebutton-cli completion bash > /etc/bash_completion.d/bigbluebutton-cli

# Zsh
bigbluebutton-cli completion zsh > "${fpath[1]}/_bigbluebutton-cli"

# Fish
bigbluebutton-cli completion fish > ~/.config/fish/completions/bigbluebutton-cli.fish
```

## License

MIT License. See [LICENSE](LICENSE) for details.
