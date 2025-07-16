package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/kungfusheep/hue-mcp/hue"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// HandleListMotionSensors returns a handler for listing motion sensors
func HandleListMotionSensors(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sensors, err := client.GetMotionSensors(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list motion sensors: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d motion sensors:\n", len(sensors)))
		for _, sensor := range sensors {
			status := "No motion"
			if sensor.Motion.Motion {
				status = "Motion detected"
			}
			enabled := "enabled"
			if !sensor.Enabled {
				enabled = "disabled"
			}
			result.WriteString(fmt.Sprintf("- %s: %s (%s) (ID: %s)\n", 
				sensor.ID, status, enabled, sensor.IDV1))
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleListTemperatureSensors returns a handler for listing temperature sensors
func HandleListTemperatureSensors(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sensors, err := client.GetTemperatureSensors(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list temperature sensors: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d temperature sensors:\n", len(sensors)))
		for _, sensor := range sensors {
			enabled := "enabled"
			if !sensor.Enabled {
				enabled = "disabled"
			}
			
			tempC := sensor.Temperature.Temperature
			tempF := tempC*9/5 + 32
			
			result.WriteString(fmt.Sprintf("- %s: %.1f°C (%.1f°F) (%s) (ID: %s)\n", 
				sensor.ID, tempC, tempF, enabled, sensor.IDV1))
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleListLightLevelSensors returns a handler for listing light level sensors
func HandleListLightLevelSensors(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sensors, err := client.GetLightLevelSensors(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list light level sensors: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d light level sensors:\n", len(sensors)))
		for _, sensor := range sensors {
			enabled := "enabled"
			if !sensor.Enabled {
				enabled = "disabled"
			}
			
			// Convert light level to lux (approximate)
			lux := float64(sensor.LightLevel.LightLevel)
			
			result.WriteString(fmt.Sprintf("- %s: %.0f lux (%s) (ID: %s)\n", 
				sensor.ID, lux, enabled, sensor.IDV1))
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleListButtons returns a handler for listing buttons (dimmer switches)
func HandleListButtons(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		buttons, err := client.GetButtons(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list buttons: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d buttons:\n", len(buttons)))
		for _, button := range buttons {
			lastEvent := "none"
			if button.Button.ButtonReport != nil {
				lastEvent = button.Button.ButtonReport.Event
			}
			
			result.WriteString(fmt.Sprintf("- %s: Last event: %s (ID: %s)\n", 
				button.Metadata.Name, lastEvent, button.ID))
			
			if len(button.Button.EventValues) > 0 {
				result.WriteString(fmt.Sprintf("  Supported events: %s\n", 
					strings.Join(button.Button.EventValues, ", ")))
			}
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}