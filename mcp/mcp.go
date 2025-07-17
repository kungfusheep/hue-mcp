package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kungfusheep/hue/effects"
	"github.com/kungfusheep/hue/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Light control handlers

// HandleLightOn returns a handler for turning a light on
func HandleLightOn(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		err := hueClient.TurnOnLight(ctx, lightID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to turn on light: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Light %s turned on", lightID)), nil
	}
}

// HandleLightOff returns a handler for turning a light off
func HandleLightOff(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		err := hueClient.TurnOffLight(ctx, lightID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to turn off light: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Light %s turned off", lightID)), nil
	}
}

// HandleLightBrightness returns a handler for setting light brightness
func HandleLightBrightness(hueClient *client.Client) server.ToolHandlerFunc {
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

		err := hueClient.SetLightBrightness(ctx, lightID, brightness)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set brightness: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Light %s brightness set to %.0f%%", lightID, brightness)), nil
	}
}

// HandleLightColor returns a handler for setting light color
func HandleLightColor(hueClient *client.Client) server.ToolHandlerFunc {
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

		err := hueClient.SetLightColor(ctx, lightID, hexColor)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set color: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Light %s color set to %s", lightID, color)), nil
	}
}

// HandleLightEffect returns a handler for setting light effects
func HandleLightEffect(hueClient *client.Client) server.ToolHandlerFunc {
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

		err := hueClient.SetLightEffect(ctx, lightID, effect, duration)
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
func HandleGroupOn(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		groupID, ok := args["group_id"].(string)
		if !ok {
			return mcp.NewToolResultError("group_id is required"), nil
		}

		err := hueClient.TurnOnGroup(ctx, groupID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to turn on group: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Group %s turned on", groupID)), nil
	}
}

// HandleGroupOff returns a handler for turning a group off
func HandleGroupOff(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		groupID, ok := args["group_id"].(string)
		if !ok {
			return mcp.NewToolResultError("group_id is required"), nil
		}

		err := hueClient.TurnOffGroup(ctx, groupID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to turn off group: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Group %s turned off", groupID)), nil
	}
}

// HandleGroupBrightness returns a handler for setting group brightness
func HandleGroupBrightness(hueClient *client.Client) server.ToolHandlerFunc {
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

		err := hueClient.SetGroupBrightness(ctx, groupID, brightness)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set brightness: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Group %s brightness set to %.0f%%", groupID, brightness)), nil
	}
}

// HandleGroupColor returns a handler for setting group color
func HandleGroupColor(hueClient *client.Client) server.ToolHandlerFunc {
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

		err := hueClient.SetGroupColor(ctx, groupID, hexColor)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set color: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Group %s color set to %s", groupID, color)), nil
	}
}

// HandleGroupEffect returns a handler for setting group effects
func HandleGroupEffect(hueClient *client.Client) server.ToolHandlerFunc {
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

		err := hueClient.SetGroupEffect(ctx, groupID, effect, duration)
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
func HandleListScenes(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scenes, err := hueClient.GetScenes(ctx)
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
func HandleActivateScene(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		sceneID, ok := args["scene_id"].(string)
		if !ok {
			return mcp.NewToolResultError("scene_id is required"), nil
		}

		err := hueClient.ActivateScene(ctx, sceneID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to activate scene: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Scene %s activated", sceneID)), nil
	}
}

// HandleCreateScene returns a handler for creating a scene
func HandleCreateScene(hueClient *client.Client) server.ToolHandlerFunc {
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
		sceneCreate := client.SceneCreate{
			Type: "scene",
			Metadata: client.Metadata{
				Name: name,
			},
			Group: client.ResourceIdentifier{
				RID:   groupID,
				RType: "grouped_light",
			},
			Actions: []client.SceneAction{}, // Would need to capture current states
		}

		scene, err := hueClient.CreateScene(ctx, sceneCreate)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create scene: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Scene '%s' created with ID: %s", name, scene.ID)), nil
	}
}

// System handlers

// HandleListLights returns a handler for listing lights
func HandleListLights(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lights, err := hueClient.GetLights(ctx)
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
func HandleListGroups(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		groups, err := hueClient.GetGroups(ctx)
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
func HandleGetLightState(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		light, err := hueClient.GetLight(ctx, lightID)
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
func HandleBridgeInfo(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bridge, err := hueClient.GetBridge(ctx)
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
func HandleIdentifyLight(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		lightID, ok := args["light_id"].(string)
		if !ok {
			return mcp.NewToolResultError("light_id is required"), nil
		}

		err := hueClient.IdentifyLight(ctx, lightID)
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

// BatchCommand represents a single command in a batch
type BatchCommand struct {
	Action   string  `json:"action"`
	TargetID string  `json:"target_id"`
	Value    string  `json:"value,omitempty"`
	Duration float64 `json:"duration,omitempty"`
}

// HandleBatchCommands executes multiple commands in sequence
func HandleBatchCommands(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		// Get commands JSON string
		commandsJSON, ok := args["commands"].(string)
		if !ok {
			return mcp.NewToolResultError("commands JSON array is required"), nil
		}
		
		// Parse commands
		var commands []map[string]interface{}
		if err := json.Unmarshal([]byte(commandsJSON), &commands); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse commands JSON: %v", err)), nil
		}
		
		// Get delay between commands (default 100ms)
		delayMs := 100
		if d, ok := args["delay_ms"].(float64); ok {
			delayMs = int(d)
		}
		
		// Get async flag (default true for non-blocking)
		async := true
		if a, ok := args["async"].(bool); ok {
			async = a
		}
		
		// Check for cache_name to save this scene
		cacheName, _ := args["cache_name"].(string)
		cacheDescription, _ := args["cache_description"].(string)
		
		// If cache_name provided, save the scene
		if cacheName != "" {
			err := globalSceneCache.SaveScene(cacheName, commands, delayMs, cacheDescription)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to cache scene: %v", err)), nil
			}
			log.Printf("Cached scene '%s' with %d commands", cacheName, len(commands))
		}
		
		// Generate batch ID for tracking
		batchID := fmt.Sprintf("batch_%d_%d", time.Now().Unix(), len(commands))
		
		if async {
			// Execute asynchronously - return immediately
			go ExecuteBatchAsync(ctx, hueClient, commands, delayMs, batchID)
			
			responseMsg := fmt.Sprintf("Batch started asynchronously with ID: %s\nCommands: %d\nDelay between commands: %dms", 
				batchID, len(commands), delayMs)
			
			if cacheName != "" {
				responseMsg = fmt.Sprintf("Creating and caching atmosphere: %s...\n%s", cacheName, responseMsg)
			}
			
			return mcp.NewToolResultText(responseMsg), nil
		} else {
			// Execute synchronously
			log.Printf("Starting synchronous batch %s with %d commands", batchID, len(commands))
			
			results := ExecuteBatch(ctx, hueClient, commands, delayMs)
			
			// Summarize results
			successful := 0
			failed := 0
			for _, result := range results {
				if result.Success {
					successful++
				} else {
					failed++
				}
			}
			
			responseMsg := fmt.Sprintf("Batch completed: %d successful, %d failed\nBatch ID: %s", 
				successful, failed, batchID)
			
			if cacheName != "" {
				responseMsg = fmt.Sprintf("Created and cached atmosphere: %s\n%s", cacheName, responseMsg)
			}
			
			return mcp.NewToolResultText(responseMsg), nil
		}
	}
}

// executeBatchCommand executes a single command within a batch
func executeBatchCommand(ctx context.Context, hueClient *client.Client, action, targetID, value string, duration int) (string, error) {
	switch action {
	case "light_on":
		err := hueClient.TurnOnLight(ctx, targetID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Light %s turned on", targetID), nil

	case "light_off":
		err := hueClient.TurnOffLight(ctx, targetID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Light %s turned off", targetID), nil

	case "light_brightness":
		brightness, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return "", fmt.Errorf("invalid brightness value: %s", value)
		}
		if brightness < 0 || brightness > 100 {
			return "", fmt.Errorf("brightness must be between 0 and 100")
		}
		err = hueClient.SetLightBrightness(ctx, targetID, brightness)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Light %s brightness set to %.0f%%", targetID, brightness), nil

	case "light_color":
		if value == "" {
			return "", fmt.Errorf("color value is required")
		}
		hexColor := namedColorToHex(value)
		if hexColor == "" {
			hexColor = value
		}
		if !isValidHexColor(hexColor) {
			return "", fmt.Errorf("invalid color format: %s", value)
		}
		err := hueClient.SetLightColor(ctx, targetID, hexColor)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Light %s color set to %s", targetID, value), nil

	case "light_effect":
		if value == "" {
			return "", fmt.Errorf("effect value is required")
		}
		err := hueClient.SetLightEffect(ctx, targetID, value, duration)
		if err != nil {
			return "", err
		}
		desc := effects.GetDescription(value)
		result := fmt.Sprintf("Light %s effect set to %s - %s", targetID, value, desc)
		if duration > 0 {
			result += fmt.Sprintf(" (duration: %d seconds)", duration)
		}
		return result, nil

	case "group_on":
		err := hueClient.TurnOnGroup(ctx, targetID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Group %s turned on", targetID), nil

	case "group_off":
		err := hueClient.TurnOffGroup(ctx, targetID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Group %s turned off", targetID), nil

	case "group_brightness":
		brightness, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return "", fmt.Errorf("invalid brightness value: %s", value)
		}
		if brightness < 0 || brightness > 100 {
			return "", fmt.Errorf("brightness must be between 0 and 100")
		}
		err = hueClient.SetGroupBrightness(ctx, targetID, brightness)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Group %s brightness set to %.0f%%", targetID, brightness), nil

	case "group_color":
		if value == "" {
			return "", fmt.Errorf("color value is required")
		}
		hexColor := namedColorToHex(value)
		if hexColor == "" {
			hexColor = value
		}
		if !isValidHexColor(hexColor) {
			return "", fmt.Errorf("invalid color format: %s", value)
		}
		err := hueClient.SetGroupColor(ctx, targetID, hexColor)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Group %s color set to %s", targetID, value), nil

	case "group_effect":
		if value == "" {
			return "", fmt.Errorf("effect value is required")
		}
		err := hueClient.SetGroupEffect(ctx, targetID, value, duration)
		if err != nil {
			return "", err
		}
		desc := effects.GetDescription(value)
		result := fmt.Sprintf("Group %s effect set to %s - %s", targetID, value, desc)
		if duration > 0 {
			result += fmt.Sprintf(" (duration: %d seconds)", duration)
		}
		return result, nil

	case "activate_scene":
		err := hueClient.ActivateScene(ctx, targetID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Scene %s activated", targetID), nil

	case "identify_light":
		err := hueClient.IdentifyLight(ctx, targetID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Light %s is blinking for identification", targetID), nil

	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}
}

// BatchResult represents the result of a batch command
type BatchResult struct {
	Success bool
	Message string
	Error   error
}

// ExecuteBatch executes batch commands synchronously and returns results
func ExecuteBatch(ctx context.Context, client *client.Client, commands []map[string]interface{}, delayMs int) []BatchResult {
	results := make([]BatchResult, 0, len(commands))
	
	for i, cmd := range commands {
		// Extract command parameters
		action, _ := cmd["action"].(string)
		targetID, _ := cmd["target_id"].(string)
		value, _ := cmd["value"].(string)
		duration := 0
		if d, ok := cmd["duration"].(float64); ok {
			duration = int(d)
		}
		
		// Execute the command
		result, err := executeBatchCommand(ctx, client, action, targetID, value, duration)
		if err != nil {
			results = append(results, BatchResult{
				Success: false,
				Message: fmt.Sprintf("Command %d (%s): %v", i, action, err),
				Error:   err,
			})
		} else {
			results = append(results, BatchResult{
				Success: true,
				Message: result,
				Error:   nil,
			})
		}
		
		// Add delay between commands (except for the last one)
		if i < len(commands)-1 && delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}
	
	return results
}

// ExecuteBatchAsync executes batch commands asynchronously (exported for testing)
func ExecuteBatchAsync(ctx context.Context, client *client.Client, commands []map[string]interface{}, delayMs int, batchID string) {
	// Create a new context that won't be cancelled by the parent
	asyncCtx := context.Background()
	
	// Log batch start
	log.Printf("Starting async batch %s with %d commands", batchID, len(commands))
	
	// Process each command
	for i, cmd := range commands {
		// Check if context was cancelled
		select {
		case <-ctx.Done():
			log.Printf("Batch %s cancelled at command %d", batchID, i)
			return
		default:
		}
		
		// Extract command parameters
		action, _ := cmd["action"].(string)
		targetID, _ := cmd["target_id"].(string)
		value, _ := cmd["value"].(string)
		duration := 0
		if d, ok := cmd["duration"].(float64); ok {
			duration = int(d)
		}
		
		// Execute the command
		result, err := executeBatchCommand(asyncCtx, client, action, targetID, value, duration)
		if err != nil {
			log.Printf("Batch %s - Command %d (%s) failed: %v", batchID, i, action, err)
		} else {
			log.Printf("Batch %s - Command %d: %s", batchID, i, result)
		}
		
		// Add delay between commands (except for the last one)
		if i < len(commands)-1 && delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}
	
	log.Printf("Batch %s completed", batchID)
}