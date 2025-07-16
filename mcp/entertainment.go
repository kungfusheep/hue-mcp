package mcp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

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

// Global entertainment streamer management
var (
	activeStreamers = make(map[string]*hue.EntertainmentStreamer)
	streamersMutex  sync.RWMutex
)

// HandleStartStreaming starts UDP streaming for an entertainment configuration
func HandleStartStreaming(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		configID, ok := args["config_id"].(string)
		if !ok || configID == "" {
			return mcp.NewToolResultError("config_id is required"), nil
		}

		// Check if streamer already exists
		streamersMutex.RLock()
		_, exists := activeStreamers[configID]
		streamersMutex.RUnlock()

		if exists {
			return mcp.NewToolResultText(fmt.Sprintf("Streaming already active for configuration %s", configID)), nil
		}

		// Create new streamer
		streamer, err := hue.NewEntertainmentStreamer(client, configID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create streamer: %v", err)), nil
		}

		// Set update rate if provided
		if rateStr, ok := args["update_rate_ms"].(string); ok {
			if rate, err := strconv.Atoi(rateStr); err == nil && rate > 0 {
				streamer.SetUpdateRate(time.Duration(rate) * time.Millisecond)
			}
		}

		// Start streaming
		err = streamer.Start(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to start streaming: %v", err)), nil
		}

		// Store streamer
		streamersMutex.Lock()
		activeStreamers[configID] = streamer
		streamersMutex.Unlock()

		return mcp.NewToolResultText(fmt.Sprintf("UDP streaming started for configuration %s", configID)), nil
	}
}

// HandleStopStreaming stops UDP streaming for an entertainment configuration
func HandleStopStreaming(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		configID, ok := args["config_id"].(string)
		if !ok || configID == "" {
			return mcp.NewToolResultError("config_id is required"), nil
		}

		streamersMutex.Lock()
		streamer, exists := activeStreamers[configID]
		if exists {
			delete(activeStreamers, configID)
		}
		streamersMutex.Unlock()

		if !exists {
			return mcp.NewToolResultText(fmt.Sprintf("No active streaming for configuration %s", configID)), nil
		}

		err := streamer.Stop(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to stop streaming: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("UDP streaming stopped for configuration %s", configID)), nil
	}
}

// HandleSendColors sends color updates to streaming lights
func HandleSendColors(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		configID, ok := args["config_id"].(string)
		if !ok || configID == "" {
			return mcp.NewToolResultError("config_id is required"), nil
		}

		colorsStr, ok := args["colors"].(string)
		if !ok || colorsStr == "" {
			return mcp.NewToolResultError("colors is required (format: 'lightID1:r,g,b;lightID2:r,g,b')"), nil
		}

		streamersMutex.RLock()
		streamer, exists := activeStreamers[configID]
		streamersMutex.RUnlock()

		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("No active streaming for configuration %s", configID)), nil
		}

		// Parse colors
		updates, err := parseColorUpdates(colorsStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse colors: %v", err)), nil
		}

		// Send colors
		err = streamer.SendColors(updates)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to send colors: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Sent color updates to %d lights", len(updates))), nil
	}
}

// HandleStreamingStatus gets the status of all active streamers
func HandleStreamingStatus(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		streamersMutex.RLock()
		defer streamersMutex.RUnlock()

		if len(activeStreamers) == 0 {
			return mcp.NewToolResultText("No active streaming sessions"), nil
		}

		result := "Active Streaming Sessions:\n"
		for configID, streamer := range activeStreamers {
			result += fmt.Sprintf("- Configuration: %s\n", configID)
			lights := streamer.GetLights()
			if lights != nil {
				result += fmt.Sprintf("  Lights: %d\n", len(lights))
				for _, light := range lights {
					result += fmt.Sprintf("    - %s (%s)\n", light.RID, light.RType)
				}
			}
			result += "\n"
		}

		return mcp.NewToolResultText(result), nil
	}
}

// HandleRainbowEffect creates a rainbow effect on streaming lights
func HandleRainbowEffect(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		configID, ok := args["config_id"].(string)
		if !ok || configID == "" {
			return mcp.NewToolResultError("config_id is required"), nil
		}

		durationStr, ok := args["duration"].(string)
		if !ok || durationStr == "" {
			durationStr = "10" // Default 10 seconds
		}

		duration, err := strconv.Atoi(durationStr)
		if err != nil || duration <= 0 {
			return mcp.NewToolResultError("duration must be a positive integer (seconds)"), nil
		}

		streamersMutex.RLock()
		streamer, exists := activeStreamers[configID]
		streamersMutex.RUnlock()

		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("No active streaming for configuration %s", configID)), nil
		}

		// Get lights
		lights := streamer.GetLights()
		if len(lights) == 0 {
			return mcp.NewToolResultError("No lights found in configuration"), nil
		}

		// Start rainbow effect
		go runRainbowEffect(streamer, lights, time.Duration(duration)*time.Second)

		return mcp.NewToolResultText(fmt.Sprintf("Rainbow effect started for %d seconds", duration)), nil
	}
}

// parseColorUpdates parses color updates from string format
func parseColorUpdates(colorsStr string) ([]hue.EntertainmentUpdate, error) {
	var updates []hue.EntertainmentUpdate
	
	pairs := strings.Split(colorsStr, ";")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid color format: %s", pair)
		}
		
		lightID := strings.TrimSpace(parts[0])
		colorStr := strings.TrimSpace(parts[1])
		
		// Parse RGB values
		rgbParts := strings.Split(colorStr, ",")
		if len(rgbParts) != 3 {
			return nil, fmt.Errorf("invalid RGB format: %s", colorStr)
		}
		
		r, err := strconv.Atoi(strings.TrimSpace(rgbParts[0]))
		if err != nil || r < 0 || r > 255 {
			return nil, fmt.Errorf("invalid red value: %s", rgbParts[0])
		}
		
		g, err := strconv.Atoi(strings.TrimSpace(rgbParts[1]))
		if err != nil || g < 0 || g > 255 {
			return nil, fmt.Errorf("invalid green value: %s", rgbParts[1])
		}
		
		b, err := strconv.Atoi(strings.TrimSpace(rgbParts[2]))
		if err != nil || b < 0 || b > 255 {
			return nil, fmt.Errorf("invalid blue value: %s", rgbParts[2])
		}
		
		// Convert to 16-bit values
		red, green, blue := hue.RGBToUint16(uint8(r), uint8(g), uint8(b))
		
		updates = append(updates, hue.EntertainmentUpdate{
			LightID: lightID,
			Red:     red,
			Green:   green,
			Blue:    blue,
		})
	}
	
	return updates, nil
}

// runRainbowEffect runs a rainbow effect on the given lights
func runRainbowEffect(streamer *hue.EntertainmentStreamer, lights []hue.ResourceIdentifier, duration time.Duration) {
	start := time.Now()
	ticker := time.NewTicker(50 * time.Millisecond) // 20fps
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if time.Since(start) >= duration {
				return
			}
			
			// Calculate rainbow colors
			progress := float64(time.Since(start)) / float64(duration)
			var updates []hue.EntertainmentUpdate
			
			for i, light := range lights {
				// Create rainbow effect with phase offset for each light
				hueValue := (progress + float64(i)*0.1) * 360
				for hueValue >= 360 {
					hueValue -= 360
				}
				
				r, g, b := hsvToRGB(hueValue, 1.0, 1.0)
				red, green, blue := hue.FloatRGBToUint16(r, g, b)
				
				updates = append(updates, hue.EntertainmentUpdate{
					LightID: light.RID,
					Red:     red,
					Green:   green,
					Blue:    blue,
				})
			}
			
			streamer.SendColors(updates)
		}
	}
}

// hsvToRGB converts HSV color to RGB
func hsvToRGB(h, s, v float64) (float64, float64, float64) {
	c := v * s
	x := c * (1 - abs(mod(h/60, 2) - 1))
	m := v - c
	
	var r, g, b float64
	
	if h >= 0 && h < 60 {
		r, g, b = c, x, 0
	} else if h >= 60 && h < 120 {
		r, g, b = x, c, 0
	} else if h >= 120 && h < 180 {
		r, g, b = 0, c, x
	} else if h >= 180 && h < 240 {
		r, g, b = 0, x, c
	} else if h >= 240 && h < 300 {
		r, g, b = x, 0, c
	} else {
		r, g, b = c, 0, x
	}
	
	return r + m, g + m, b + m
}

// Helper functions
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func mod(x, y float64) float64 {
	return x - y*float64(int(x/y))
}