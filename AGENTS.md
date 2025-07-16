# AGENTS.md - Development Guidelines for Hue MCP Server

## Build & Test Commands
- **Build**: `go build .` or `go build -o hue-mcp`
- **Test all**: `go test -v ./...`
- **Test single package**: `go test -v ./hue` or `go test -v ./mcp`
- **Run**: `go run .` (requires HUE_BRIDGE_IP and HUE_USERNAME env vars)

## Code Style & Conventions
- **Package structure**: `/hue` (API client), `/mcp` (MCP handlers), `/effects` (constants)
- **Imports**: Standard library first, then third-party, then local packages with blank lines between groups
- **Naming**: Use camelCase for unexported, PascalCase for exported. Prefer descriptive names over abbreviations
- **Error handling**: Always check errors, wrap with context using `fmt.Errorf("description: %w", err)`
- **Types**: Define structs in `types.go`, use pointer receivers for methods that modify state
- **Context**: Always pass `context.Context` as first parameter to functions making HTTP calls

## Testing
- Test files use `_test.go` suffix and live alongside source files
- Use table-driven tests with subtests for multiple scenarios
- Mock HTTP responses for API client tests
- Test both success and error cases

## Dependencies
- Uses Go 1.24.0 with `github.com/mark3labs/mcp-go` for MCP server functionality
- HTTP client with TLS skip verification for Hue bridge self-signed certificates
- No external linting tools configured - follow standard Go conventions