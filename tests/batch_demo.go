package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kungfusheep/hue-mcp/hue"
	mcpserver "github.com/kungfusheep/hue-mcp/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	bridgeIP := os.Getenv("HUE_BRIDGE_IP")
	if bridgeIP == "" {
		bridgeIP = "192.168.87.51"
	}

	username := os.Getenv("HUE_USERNAME")
	if username == "" {
		fmt.Println("Please set HUE_USERNAME environment variable")
		return
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client := hue.NewClient(bridgeIP, username, httpClient)
	ctx := context.Background()

	fmt.Println("🔄 Testing Batch Commands")
	fmt.Println("=========================")

	// Find some lights to work with
	lights, err := client.GetLights(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get lights: %v\n", err)
		return
	}

	if len(lights) < 2 {
		fmt.Println("❌ Need at least 2 lights for batch testing")
		return
	}

	light1 := lights[0]
	light2 := lights[1]

	fmt.Printf("🎯 Using lights: %s and %s\n", light1.Metadata.Name, light2.Metadata.Name)

	// Create batch commands
	batchCommands := []mcpserver.BatchCommand{
		{Action: "light_on", TargetID: light1.ID},
		{Action: "light_on", TargetID: light2.ID},
		{Action: "light_brightness", TargetID: light1.ID, Value: "100"},
		{Action: "light_brightness", TargetID: light2.ID, Value: "50"},
		{Action: "light_color", TargetID: light1.ID, Value: "#FF0000"},
		{Action: "light_color", TargetID: light2.ID, Value: "#00FF00"},
	}

	// Convert to JSON
	commandsJSON, err := json.Marshal(batchCommands)
	if err != nil {
		fmt.Printf("❌ Failed to marshal commands: %v\n", err)
		return
	}

	fmt.Printf("📋 Batch commands JSON:\n%s\n", string(commandsJSON))

	// Create a mock request
	mockRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "batch_commands",
			Arguments: map[string]interface{}{
				"commands":  string(commandsJSON),
				"delay_ms":  200,
			},
		},
	}

	// Execute batch
	fmt.Println("\n🚀 Executing batch commands...")
	handler := mcpserver.HandleBatchCommands(client)
	result, err := handler(ctx, mockRequest)
	
	if err != nil {
		fmt.Printf("❌ Batch execution failed: %v\n", err)
		return
	}

	// Extract result text
	var resultText string
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			resultText = textContent.Text
		}
	}

	fmt.Printf("✅ Batch execution result:\n%s\n", resultText)

	// Clean up - restore neutral state
	fmt.Println("🧹 Cleaning up...")
	cleanupCommands := []mcpserver.BatchCommand{
		{Action: "light_brightness", TargetID: light1.ID, Value: "75"},
		{Action: "light_brightness", TargetID: light2.ID, Value: "75"},
		{Action: "light_color", TargetID: light1.ID, Value: "warm"},
		{Action: "light_color", TargetID: light2.ID, Value: "warm"},
	}

	cleanupJSON, _ := json.Marshal(cleanupCommands)
	cleanupRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "batch_commands",
			Arguments: map[string]interface{}{
				"commands":  string(cleanupJSON),
				"delay_ms":  100,
			},
		},
	}

	cleanupResult, err := handler(ctx, cleanupRequest)
	if err != nil {
		fmt.Printf("❌ Cleanup failed: %v\n", err)
	} else {
		var cleanupText string
		if len(cleanupResult.Content) > 0 {
			if textContent, ok := cleanupResult.Content[0].(mcp.TextContent); ok {
				cleanupText = textContent.Text
			}
		}
		fmt.Printf("✅ Cleanup completed:\n%s\n", cleanupText)
	}

	fmt.Println("🎯 Batch command testing complete!")
	fmt.Println("\n📊 SUMMARY:")
	fmt.Println("• Batch commands allow multiple operations in a single MCP request")
	fmt.Println("• Configurable delay between commands reduces bridge load")
	fmt.Println("• JSON format makes it easy for AI agents to construct complex sequences")
	fmt.Println("• Perfect for synchronized effects, scenes, or multi-light animations")
}