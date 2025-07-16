package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/kungfusheep/hue-mcp/hue"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// HandleListRooms returns a handler for listing rooms
func HandleListRooms(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		rooms, err := client.GetRooms(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list rooms: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d rooms:\n", len(rooms)))
		for _, room := range rooms {
			result.WriteString(fmt.Sprintf("- %s (ID: %s)\n", room.Metadata.Name, room.ID))
			
			// List lights in the room
			for _, child := range room.Children {
				if child.RType == "light" {
					result.WriteString(fmt.Sprintf("  └─ Light: %s\n", child.RID))
				}
			}
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleListZones returns a handler for listing zones
func HandleListZones(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		zones, err := client.GetZones(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list zones: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d zones:\n", len(zones)))
		for _, zone := range zones {
			result.WriteString(fmt.Sprintf("- %s (ID: %s)\n", zone.Metadata.Name, zone.ID))
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleListDevices returns a handler for listing devices
func HandleListDevices(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		devices, err := client.GetDevices(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list devices: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d devices:\n", len(devices)))
		for _, device := range devices {
			result.WriteString(fmt.Sprintf("- %s (%s): %s\n", 
				device.Metadata.Name, 
				device.ProductData.ProductName,
				device.ID))
			
			if device.PowerState != nil {
				result.WriteString(fmt.Sprintf("  Power: %s", device.PowerState.PowerState))
				if device.PowerState.BatteryLevel > 0 {
					result.WriteString(fmt.Sprintf(", Battery: %.0f%%", device.PowerState.BatteryLevel))
				}
				result.WriteString("\n")
			}
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleGetDevice returns a handler for getting device details
func HandleGetDevice(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		deviceID, ok := args["device_id"].(string)
		if !ok {
			return mcp.NewToolResultError("device_id is required"), nil
		}

		device, err := client.GetDevice(ctx, deviceID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get device: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Device: %s\n", device.Metadata.Name))
		result.WriteString(fmt.Sprintf("Model: %s\n", device.ProductData.ModelID))
		result.WriteString(fmt.Sprintf("Product: %s\n", device.ProductData.ProductName))
		result.WriteString(fmt.Sprintf("Manufacturer: %s\n", device.ProductData.ManufacturerName))
		result.WriteString(fmt.Sprintf("Type: %s\n", device.ProductData.ProductArchetype))
		result.WriteString(fmt.Sprintf("Software Version: %s\n", device.ProductData.SoftwareVersion))
		
		if device.PowerState != nil {
			result.WriteString(fmt.Sprintf("Power State: %s\n", device.PowerState.PowerState))
			if device.PowerState.BatteryLevel > 0 {
				result.WriteString(fmt.Sprintf("Battery Level: %.0f%%\n", device.PowerState.BatteryLevel))
				result.WriteString(fmt.Sprintf("Battery State: %s\n", device.PowerState.BatteryState))
			}
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}