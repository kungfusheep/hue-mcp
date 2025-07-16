package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

	fmt.Println("🚀 Advanced Batch Command Demo")
	fmt.Println("==============================")

	// Find office lights
	lights, err := client.GetLights(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get lights: %v\n", err)
		return
	}

	var officeLights []hue.Light
	// Include the Hue Play lights (Playbars) in the office setup
	officeNames := []string{"Office 1", "Office 2", "Office 3", "Office 4", "Petes Office Lamp", "Hue Play 1", "Hue Play 2"}
	
	for _, light := range lights {
		for _, name := range officeNames {
			if light.Metadata.Name == name {
				officeLights = append(officeLights, light)
				break
			}
		}
	}

	if len(officeLights) == 0 {
		fmt.Println("❌ No office lights found")
		return
	}

	fmt.Printf("🎯 Found %d office lights:\n", len(officeLights))
	for _, light := range officeLights {
		fmt.Printf("  • %s (ID: %s)\n", light.Metadata.Name, light.ID)
	}

	// Test 1: Synchronized lighting effect
	fmt.Println("\n🎬 Test 1: Synchronized Candle Effect")
	fmt.Println("=====================================")
	
	var candleCommands []mcpserver.BatchCommand
	for _, light := range officeLights {
		candleCommands = append(candleCommands, 
			mcpserver.BatchCommand{Action: "light_on", TargetID: light.ID},
			mcpserver.BatchCommand{Action: "light_brightness", TargetID: light.ID, Value: "80"},
			mcpserver.BatchCommand{Action: "light_effect", TargetID: light.ID, Value: "candle", Duration: 0},
		)
	}

	candleJSON, _ := json.Marshal(candleCommands)
	fmt.Printf("📋 Batch: %d commands to create synchronized candle effect\n", len(candleCommands))
	
	executeCommands(ctx, client, string(candleJSON), 150)
	
	fmt.Println("   ✅ All office lights should now have candle effect!")
	fmt.Println("   Watch for 8 seconds...")
	time.Sleep(8 * time.Second)

	// Test 2: Color wave effect
	fmt.Println("\n🌈 Test 2: Color Wave Effect")
	fmt.Println("============================")
	
	colors := []string{"#FF0000", "#FF7F00", "#FFFF00", "#00FF00", "#0000FF"}
	var waveCommands []mcpserver.BatchCommand
	
	for i, light := range officeLights {
		if i < len(colors) {
			waveCommands = append(waveCommands, 
				mcpserver.BatchCommand{Action: "light_effect", TargetID: light.ID, Value: "no_effect"},
				mcpserver.BatchCommand{Action: "light_color", TargetID: light.ID, Value: colors[i]},
				mcpserver.BatchCommand{Action: "light_brightness", TargetID: light.ID, Value: "100"},
			)
		}
	}

	waveJSON, _ := json.Marshal(waveCommands)
	fmt.Printf("📋 Batch: %d commands to create color wave\n", len(waveCommands))
	
	executeCommands(ctx, client, string(waveJSON), 300)
	
	fmt.Println("   ✅ Office lights should now show a color wave!")
	fmt.Println("   Watch for 5 seconds...")
	time.Sleep(5 * time.Second)

	// Test 3: Fire effect with staggered timing
	fmt.Println("\n🔥 Test 3: Staggered Fire Effect")
	fmt.Println("================================")
	
	var fireCommands []mcpserver.BatchCommand
	for _, light := range officeLights {
		fireCommands = append(fireCommands, 
			mcpserver.BatchCommand{Action: "light_effect", TargetID: light.ID, Value: "fire", Duration: 0},
			mcpserver.BatchCommand{Action: "light_brightness", TargetID: light.ID, Value: "90"},
		)
	}

	fireJSON, _ := json.Marshal(fireCommands)
	fmt.Printf("📋 Batch: %d commands with 500ms stagger for fire effect\n", len(fireCommands))
	
	executeCommands(ctx, client, string(fireJSON), 500)
	
	fmt.Println("   ✅ Office lights should now have staggered fire effect!")
	fmt.Println("   Watch for 8 seconds...")
	time.Sleep(8 * time.Second)

	// Test 4: Cleanup and restore
	fmt.Println("\n🧹 Test 4: Cleanup and Restore")
	fmt.Println("===============================")
	
	var cleanupCommands []mcpserver.BatchCommand
	for _, light := range officeLights {
		cleanupCommands = append(cleanupCommands, 
			mcpserver.BatchCommand{Action: "light_effect", TargetID: light.ID, Value: "no_effect"},
			mcpserver.BatchCommand{Action: "light_brightness", TargetID: light.ID, Value: "60"},
			mcpserver.BatchCommand{Action: "light_color", TargetID: light.ID, Value: "warm"},
		)
	}

	cleanupJSON, _ := json.Marshal(cleanupCommands)
	fmt.Printf("📋 Batch: %d commands to restore neutral state\n", len(cleanupCommands))
	
	executeCommands(ctx, client, string(cleanupJSON), 100)
	
	fmt.Println("   ✅ Office lights restored to neutral state!")

	// Summary
	fmt.Println("\n🎯 Advanced Batch Demo Complete!")
	fmt.Println("=================================")
	fmt.Println("✅ Synchronized candle effect across all office lights")
	fmt.Println("✅ Color wave effect with individual light colors")
	fmt.Println("✅ Staggered fire effect with custom timing")
	fmt.Println("✅ Efficient cleanup with batch restore")
	fmt.Println("\n💡 Key Benefits:")
	fmt.Println("• Single MCP request handles complex multi-light sequences")
	fmt.Println("• Configurable delays prevent bridge rate limiting")
	fmt.Println("• Perfect for AI agents creating dynamic lighting scenes")
	fmt.Println("• Much more efficient than individual MCP calls")
}

func executeCommands(ctx context.Context, client *hue.Client, commandsJSON string, delayMs int) {
	mockRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "batch_commands",
			Arguments: map[string]interface{}{
				"commands":  commandsJSON,
				"delay_ms":  float64(delayMs),
			},
		},
	}

	handler := mcpserver.HandleBatchCommands(client)
	result, err := handler(ctx, mockRequest)
	
	if err != nil {
		fmt.Printf("❌ Batch execution failed: %v\n", err)
		return
	}

	// Extract and display summary
	var resultText string
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			resultText = textContent.Text
		}
	}

	// Just show the summary line
	lines := strings.Split(resultText, "\n")
	if len(lines) > 0 {
		fmt.Printf("   %s\n", lines[0])
	}
}