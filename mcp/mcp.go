package mcp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/kungfusheep/hue-mcp/effects"
	"github.com/kungfusheep/hue-mcp/hue"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Light control handlers

// HandleLightOn returns a handler for turning a light on
func HandleLightOn(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		err := client.TurnOnLight(ctx, lightID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to turn on light: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Light %s turned on", lightID)), nil
	}
}

// HandleLightOff returns a handler for turning a light off
func HandleLightOff(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		err := client.TurnOffLight(ctx, lightID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to turn off light: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Light %s turned off", lightID)), nil
	}
}

// HandleLightBrightness returns a handler for setting light brightness
func HandleLightBrightness(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		brightness, ok := args["brightness"].(float64)
		if !ok {
			return mcp.NewToolResultError("brightness is required"), nil
		}

		if brightness < 0 || brightness > 100 {
			return mcp.NewToolResultError("brightness must be between 0 and 100"), nil
		}

		err := client.SetLightBrightness(ctx, lightID, brightness)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set brightness: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Light %s brightness set to %.0f%%", lightID, brightness)), nil
	}
}

// HandleLightColor returns a handler for setting light color
func HandleLightColor(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		color, ok := args["color"].(string)
		if !ok {
			return mcp.NewToolResultError("color is required"), nil
		}

		// Handle named colors
		hexColor := namedColorToHex(color)
		if hexColor == "" {
			hexColor = color
		}

		// Validate hex color
		if !isValidHexColor(hexColor) {
			return mcp.NewToolResultError("Invalid color format. Use hex code (#RRGGBB) or color name"), nil
		}

		err := client.SetLightColor(ctx, lightID, hexColor)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set color: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Light %s color set to %s", lightID, color)), nil
	}
}

// HandleLightEffect returns a handler for setting light effects
func HandleLightEffect(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		effect, ok := args["effect"].(string)
		if !ok {
			return mcp.NewToolResultError("effect is required"), nil
		}

		// Note: We don't validate effects here anymore since they're dynamically generated
		// The MCP enum validation will handle this at the protocol level

		duration := 0
		if d, ok := args["duration"].(float64); ok {
			duration = int(d)
		}

		err := client.SetLightEffect(ctx, lightID, effect, duration)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set effect: %v", err)), nil
		}

		desc := effects.GetDescription(effect)
		result := fmt.Sprintf("Light %s effect set to %s - %s", lightID, effect, desc)
		if duration > 0 {
			result += fmt.Sprintf(" (duration: %d seconds)", duration)
		}

		return mcp.NewToolResultText(result), nil
	}
}

// Group control handlers

// HandleGroupOn returns a handler for turning a group on
func HandleGroupOn(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		groupID, ok := args["group_id"].(string)
		if !ok {
			return mcp.NewToolResultError("group_id is required"), nil
		}

		err := client.TurnOnGroup(ctx, groupID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to turn on group: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Group %s turned on", groupID)), nil
	}
}

// HandleGroupOff returns a handler for turning a group off
func HandleGroupOff(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		groupID, ok := args["group_id"].(string)
		if !ok {
			return mcp.NewToolResultError("group_id is required"), nil
		}

		err := client.TurnOffGroup(ctx, groupID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to turn off group: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Group %s turned off", groupID)), nil
	}
}

// HandleGroupBrightness returns a handler for setting group brightness
func HandleGroupBrightness(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		groupID, ok := args["group_id"].(string)
		if !ok {
			return mcp.NewToolResultError("group_id is required"), nil
		}

		brightness, ok := args["brightness"].(float64)
		if !ok {
			return mcp.NewToolResultError("brightness is required"), nil
		}

		if brightness < 0 || brightness > 100 {
			return mcp.NewToolResultError("brightness must be between 0 and 100"), nil
		}

		err := client.SetGroupBrightness(ctx, groupID, brightness)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set brightness: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Group %s brightness set to %.0f%%", groupID, brightness)), nil
	}
}

// HandleGroupColor returns a handler for setting group color
func HandleGroupColor(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		groupID, ok := args["group_id"].(string)
		if !ok {
			return mcp.NewToolResultError("group_id is required"), nil
		}

		color, ok := args["color"].(string)
		if !ok {
			return mcp.NewToolResultError("color is required"), nil
		}

		// Handle named colors
		hexColor := namedColorToHex(color)
		if hexColor == "" {
			hexColor = color
		}

		// Validate hex color
		if !isValidHexColor(hexColor) {
			return mcp.NewToolResultError("Invalid color format. Use hex code (#RRGGBB) or color name"), nil
		}

		err := client.SetGroupColor(ctx, groupID, hexColor)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set color: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Group %s color set to %s", groupID, color)), nil
	}
}

// HandleGroupEffect returns a handler for setting group effects
func HandleGroupEffect(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		groupID, ok := args["group_id"].(string)
		if !ok {
			return mcp.NewToolResultError("group_id is required"), nil
		}

		effect, ok := args["effect"].(string)
		if !ok {
			return mcp.NewToolResultError("effect is required"), nil
		}

		// Note: We don't validate effects here anymore since they're dynamically generated
		// The MCP enum validation will handle this at the protocol level

		duration := 0
		if d, ok := args["duration"].(float64); ok {
			duration = int(d)
		}

		err := client.SetGroupEffect(ctx, groupID, effect, duration)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set effect: %v", err)), nil
		}

		desc := effects.GetDescription(effect)
		result := fmt.Sprintf("Group %s effect set to %s - %s", groupID, effect, desc)
		if duration > 0 {
			result += fmt.Sprintf(" (duration: %d seconds)", duration)
		}

		return mcp.NewToolResultText(result), nil
	}
}

// Scene handlers

// HandleListScenes returns a handler for listing scenes
func HandleListScenes(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scenes, err := client.GetScenes(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list scenes: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d scenes:\n", len(scenes)))
		for _, scene := range scenes {
			result.WriteString(fmt.Sprintf("- %s: %s (ID: %s)\n", scene.Metadata.Name, scene.ID, scene.IDV1))
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleActivateScene returns a handler for activating a scene
func HandleActivateScene(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		sceneID, ok := args["scene_id"].(string)
		if !ok {
			return mcp.NewToolResultError("scene_id is required"), nil
		}

		err := client.ActivateScene(ctx, sceneID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to activate scene: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Scene %s activated", sceneID)), nil
	}
}

// HandleCreateScene returns a handler for creating a scene
func HandleCreateScene(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResultError("name is required"), nil
		}

		groupID, ok := args["group_id"].(string)
		if !ok {
			return mcp.NewToolResultError("group_id is required"), nil
		}

		// Create scene
		sceneCreate := hue.SceneCreate{
			Type: "scene",
			Metadata: hue.Metadata{
				Name: name,
			},
			Group: hue.ResourceIdentifier{
				RID:   groupID,
				RType: "grouped_light",
			},
			Actions: []hue.SceneAction{}, // Would need to capture current states
		}

		scene, err := client.CreateScene(ctx, sceneCreate)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create scene: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Scene '%s' created with ID: %s", name, scene.ID)), nil
	}
}

// System handlers

// HandleListLights returns a handler for listing lights
func HandleListLights(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lights, err := client.GetLights(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list lights: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d lights:\n", len(lights)))
		for _, light := range lights {
			status := "off"
			if light.On.On {
				status = fmt.Sprintf("on, brightness: %.0f%%", light.Dimming.Brightness)
			}
			result.WriteString(fmt.Sprintf("- %s (%s): %s (ID: %s, v1: %s)\n", 
				light.Metadata.Name, light.Metadata.Archetype, status, light.ID, light.IDV1))
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleListGroups returns a handler for listing groups
func HandleListGroups(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		groups, err := client.GetGroups(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list groups: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found %d groups:\n", len(groups)))
		for _, group := range groups {
			status := "off"
			if group.On.On {
				status = fmt.Sprintf("on, brightness: %.0f%%", group.Dimming.Brightness)
			}
			result.WriteString(fmt.Sprintf("- %s: %s (ID: %s, v1: %s)\n", 
				group.Metadata.Name, status, group.ID, group.IDV1))
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleGetLightState returns a handler for getting light state
func HandleGetLightState(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		light, err := client.GetLight(ctx, lightID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get light: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Light: %s\n", light.Metadata.Name))
		result.WriteString(fmt.Sprintf("Type: %s\n", light.Metadata.Archetype))
		result.WriteString(fmt.Sprintf("On: %v\n", light.On.On))
		result.WriteString(fmt.Sprintf("Brightness: %.0f%%\n", light.Dimming.Brightness))
		
		if light.Color != nil {
			result.WriteString(fmt.Sprintf("Color XY: (%.3f, %.3f)\n", light.Color.XY.X, light.Color.XY.Y))
		}
		
		if light.ColorTemperature != nil && light.ColorTemperature.MirekValid {
			result.WriteString(fmt.Sprintf("Color Temperature: %d mirek\n", light.ColorTemperature.Mirek))
		}
		
		if light.Effects != nil {
			result.WriteString(fmt.Sprintf("Effect: %s\n", light.Effects.Effect))
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleBridgeInfo returns a handler for getting bridge info
func HandleBridgeInfo(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bridge, err := client.GetBridge(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get bridge info: %v", err)), nil
		}

		var result strings.Builder
		result.WriteString("Hue Bridge Information:\n")
		result.WriteString(fmt.Sprintf("Bridge ID: %s\n", bridge.BridgeID))
		result.WriteString(fmt.Sprintf("Time Zone: %s\n", bridge.TimeZone.TimeZone))
		result.WriteString(fmt.Sprintf("API ID: %s\n", bridge.ID))
		result.WriteString(fmt.Sprintf("V1 ID: %s\n", bridge.IDV1))

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleIdentifyLight returns a handler for identifying a light
func HandleIdentifyLight(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		err := client.IdentifyLight(ctx, lightID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to identify light: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Light %s is blinking for identification", lightID)), nil
	}
}

// Helper functions

func namedColorToHex(color string) string {
	colors := map[string]string{
		"red":     "#FF0000",
		"green":   "#00FF00",
		"blue":    "#0000FF",
		"yellow":  "#FFFF00",
		"cyan":    "#00FFFF",
		"magenta": "#FF00FF",
		"white":   "#FFFFFF",
		"warm":    "#FFA500",
		"cool":    "#ADD8E6",
		"orange":  "#FFA500",
		"purple":  "#800080",
		"pink":    "#FFC0CB",
	}
	
	hex, ok := colors[strings.ToLower(color)]
	if ok {
		return hex
	}
	return ""
}

func isValidHexColor(hex string) bool {
	if !strings.HasPrefix(hex, "#") {
		return false
	}
	
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return false
	}
	
	_, err := strconv.ParseUint(hex, 16, 32)
	return err == nil
}