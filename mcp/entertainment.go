package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/kungfusheep/hue-mcp/hue"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// HandleListEntertainment returns a handler for listing entertainment configurations
func HandleListEntertainment(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configs, err := client.GetEntertainmentConfigurations(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list entertainment configurations: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d entertainment configurations:\n", len(configs)))
		for _, config := range configs {
			result.WriteString(fmt.Sprintf("- %s (ID: %s)\n", config.Metadata.Name, config.ID))
			result.WriteString(fmt.Sprintf("  Type: %s\n", config.ConfigurationType))
			result.WriteString(fmt.Sprintf("  Status: %s\n", config.Status))
			result.WriteString(fmt.Sprintf("  Channels: %d\n", len(config.Channels)))
			result.WriteString(fmt.Sprintf("  Lights: %d\n", len(config.LightServices)))
			
			if config.ActiveStreamer != nil {
				result.WriteString(fmt.Sprintf("  Active Streamer: %s\n", config.ActiveStreamer.RID))
			}
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleStartEntertainment returns a handler for starting entertainment mode
func HandleStartEntertainment(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		configID, ok := args["config_id"].(string)
		if !ok {
			return mcp.NewToolResultError("config_id is required"), nil
		}

		err := client.StartEntertainment(ctx, configID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to start entertainment: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Entertainment mode started for configuration %s", configID)), nil
	}
}

// HandleStopEntertainment returns a handler for stopping entertainment mode
func HandleStopEntertainment(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		configID, ok := args["config_id"].(string)
		if !ok {
			return mcp.NewToolResultError("config_id is required"), nil
		}

		err := client.StopEntertainment(ctx, configID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to stop entertainment: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Entertainment mode stopped for configuration %s", configID)), nil
	}
}