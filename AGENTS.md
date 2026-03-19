# AGENTS.md

Guide for AI agents working on the bigbluebutton-cli codebase.

## Project overview

A Go CLI tool wrapping the BigBlueButton API. Binary name: `bigbluebutton-cli`. Module path: `github.com/invokablegmbh/bigbluebutton-cli`.

## Build and test commands

```bash
make build          # Build ./bigbluebutton-cli binary
make test           # Run all tests (go test ./... -v -race -count=1)
make build-all      # Cross-compile linux/amd64, linux/arm64, windows/amd64 to dist/
make lint           # Run golangci-lint
make clean          # Remove binaries
go test ./... -count=1   # Quick test run (race detector requires CGO_ENABLED=1)
go build ./...      # Verify compilation
```

Note: The race detector (`-race`) requires `CGO_ENABLED=1`. In environments without cgo, omit the `-race` flag.

## Architecture

### Package layout

```
cmd/           CLI layer (Cobra commands). Each file = one command.
pkg/api/       API client. Pure HTTP + XML/JSON. Zero CLI dependencies.
pkg/config/    Configuration cascade (flags > env > file > autodetect > wizard).
pkg/interactive/  Terminal prompt wrappers (charmbracelet/huh).
internal/testutil/ Test helpers (fake HTTP server, fixtures).
```

### Key design patterns

**Adding a new command** requires touching exactly these places:

1. Create `cmd/<command_name>.go` with a `newXxxCmd() *cobra.Command` function
2. Register it in `cmd/root.go` → `registerCommands()`
3. If the API endpoint is new, add request struct in `pkg/api/requests.go`, response struct in `pkg/api/responses.go`, and the client method in the appropriate `pkg/api/*.go` file (admin/monitoring/recording/engagement)
4. Add a fixture response in `internal/testutil/fixtures.go`
5. Add tests in `cmd/commands_test.go` (command integration) and `pkg/api/*_test.go` (API unit)

**Command structure**: All commands are flat (no subcommand groups). Each command's `RunE` follows this flow:
1. `loadConfig()` — resolves URL + secret via the cascade
2. `requireString()` / `optionalString()` / `optionalBool()` / `optionalInt()` — resolve flags, prompting interactively if needed
3. Build a request struct, call the API client
4. Print output to `os.Stdout`

**Request structs** use pointer types (`*string`, `*bool`, `*int`) for optional fields so zero-value vs unset is distinguishable. Use `api.StringPtr()`, `api.BoolPtr()`, `api.IntPtr()` helpers.

**Config cascade** in `pkg/config/config.go`: the `Load()` function tries each source in order and fills in missing values. The `allowInteractive` param controls whether the wizard fallback is enabled (determined by terminal detection).

**API client** (`pkg/api/client.go`): `doGet` and `doPost` handle checksum calculation, HTTP execution, and verbose logging. Each endpoint method unmarshals the response and checks for API errors. The client accepts an injectable `*http.Client` for testing.

**Checksum**: SHA-256 of `apiCallName + queryString + sharedSecret`. The checksum parameter itself is NOT included in the hash input.

### Testing approach

- **API tests** (`pkg/api/*_test.go`): Use `httptest.NewServer` with canned responses. Test both success and error paths. The helper `newTestServer(t, apiCall, response)` creates a server that verifies the correct API path is called.
- **Config tests** (`pkg/config/*_test.go`): Test each cascade level in isolation using temp files and env vars.
- **Command integration tests** (`cmd/commands_test.go`): Use `executeCommand(t, args...)` which captures `os.Stdout` via pipe (commands write directly to `os.Stdout`, not `cmd.OutOrStdout()`). A fake BBB server provides canned responses.
- **Response parsing tests** (`pkg/api/responses_test.go`): Verify XML/JSON unmarshaling with known payloads.

Important: Commands write to `os.Stdout` directly, not to `cmd.OutOrStdout()`. The test helper `executeCommand` accounts for this by redirecting `os.Stdout` to a pipe.

### Dependencies

| Package | Role |
|---|---|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/spf13/viper` | Configuration (env, file, flags) |
| `github.com/charmbracelet/huh` | Interactive terminal prompts |
| `github.com/stretchr/testify` | Test assertions (used in some tests) |
| `golang.org/x/term` | Terminal detection (`IsTerminal`) |

### File naming conventions

- Command files: `cmd/<verb>_<noun>.go` (e.g., `get_meetings.go`, `publish_recordings.go`)
- API files grouped by domain: `admin.go`, `monitoring.go`, `recording.go`, `engagement.go`
- Tests are co-located: `foo.go` → `foo_test.go`

### Common gotchas

- `NewRootCmd()` creates a fresh command tree each call — required for test isolation since Cobra commands carry state.
- The `cmd` package has an `init()` in both `root.go` (log setup) and `helpers.go` (config prompt registration). Both must execute.
- `pkg/config` avoids importing `pkg/interactive` to prevent cycles. Instead, it uses a `SetPromptFunc` callback that `cmd/helpers.go` wires up.
- Text track endpoints (`getRecordingTextTracks`, `putRecordingTextTrack`) return **JSON**, not XML. All other endpoints return XML.
- `putRecordingTextTrack` uses multipart form POST. `insertDocument` uses XML body POST. All other endpoints use GET.

### BigBlueButton API reference

The full API spec is at https://docs.bigbluebutton.org/development/api/. The 16 endpoints implemented:

| Category | Endpoints |
|---|---|
| Admin | `create`, `create-join`, `join`, `end`, `insertDocument` |
| Monitoring | `isMeetingRunning`, `getMeetingInfo`, `getMeetings` |
| Recording | `getRecordings`, `publishRecordings`, `deleteRecordings`, `updateRecordings`, `getRecordingTextTracks`, `putRecordingTextTrack` |
| Engagement | `sendChatMessage`, `getJoinUrl` |
