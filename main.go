package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/kungfusheep/hue-mcp/effects"
	"github.com/kungfusheep/hue-mcp/hue"
	mcpserver "github.com/kungfusheep/hue-mcp/mcp"
)

func main() {
	// Get configuration from environment
	bridgeIP := os.Getenv("HUE_BRIDGE_IP")
	if bridgeIP == "" {
		bridgeIP = "192.168.87.51" // Default from handover doc
	}

	username := os.Getenv("HUE_USERNAME")
	if username == "" {
		log.Fatal("HUE_USERNAME environment variable is required")
	}

	// Create HTTP client that skips certificate verification for self-signed certs
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Initialize Hue client
	hueClient := hue.NewClient(bridgeIP, username, httpClient)

	// Test connection
	ctx := context.Background()
	if err := hueClient.TestConnection(ctx); err != nil {
		log.Fatalf("Failed to connect to Hue bridge: %v", err)
	}

	// Create MCP server
	srv := server.NewMCPServer(
		"Philips Hue v2 MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true),
	)

	// Register tools
	registerLightTools(srv, hueClient)
	registerGroupTools(srv, hueClient)
	registerSceneTools(srv, hueClient)
	registerEffectTools(srv, hueClient)
	registerSystemTools(srv, hueClient)

	// Start server in stdio mode for Claude Desktop
	log.Println("Starting Hue MCP server...")
	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// registerLightTools adds individual light control tools
func registerLightTools(srv *server.MCPServer, client *hue.Client) {
	// Light on/off
	lightOnTool := mcp.NewTool("light_on",
		mcp.WithDescription("Turn a light on"),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("The ID of the light")),
	)
	srv.AddTool(lightOnTool, mcpserver.HandleLightOn(client))

	lightOffTool := mcp.NewTool("light_off",
		mcp.WithDescription("Turn a light off"),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("The ID of the light")),
	)
	srv.AddTool(lightOffTool, mcpserver.HandleLightOff(client))

	// Brightness control
	brightnessTool := mcp.NewTool("light_brightness",
		mcp.WithDescription("Set light brightness"),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("The ID of the light")),
		mcp.WithNumber("brightness", mcp.Required(), mcp.Description("Brightness percentage (0-100)")),
	)
	srv.AddTool(brightnessTool, mcpserver.HandleLightBrightness(client))

	// Color control
	colorTool := mcp.NewTool("light_color",
		mcp.WithDescription("Set light color"),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("The ID of the light")),
		mcp.WithString("color", mcp.Required(), mcp.Description("Color as hex code (e.g., #FF0000) or color name")),
	)
	srv.AddTool(colorTool, mcpserver.HandleLightColor(client))
}

// registerGroupTools adds group control tools
func registerGroupTools(srv *server.MCPServer, client *hue.Client) {
	// Group on/off
	groupOnTool := mcp.NewTool("group_on",
		mcp.WithDescription("Turn a group of lights on"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("The ID of the group")),
	)
	srv.AddTool(groupOnTool, mcpserver.HandleGroupOn(client))

	groupOffTool := mcp.NewTool("group_off",
		mcp.WithDescription("Turn a group of lights off"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("The ID of the group")),
	)
	srv.AddTool(groupOffTool, mcpserver.HandleGroupOff(client))

	// Group brightness
	groupBrightnessTool := mcp.NewTool("group_brightness",
		mcp.WithDescription("Set group brightness"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Brightness percentage (0-100)")),
	)
	srv.AddTool(groupBrightnessTool, mcpserver.HandleGroupBrightness(client))

	// Group color
	groupColorTool := mcp.NewTool("group_color",
		mcp.WithDescription("Set group color"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("The ID of the group")),
		mcp.WithString("color", mcp.Required(), mcp.Description("Color as hex code or name")),
	)
	srv.AddTool(groupColorTool, mcpserver.HandleGroupColor(client))
}

// registerSceneTools adds scene management tools
func registerSceneTools(srv *server.MCPServer, client *hue.Client) {
	// List scenes
	listScenesTool := mcp.NewTool("list_scenes",
		mcp.WithDescription("List all available scenes"),
	)
	srv.AddTool(listScenesTool, mcpserver.HandleListScenes(client))

	// Activate scene
	activateSceneTool := mcp.NewTool("activate_scene",
		mcp.WithDescription("Activate a scene"),
		mcp.WithString("scene_id", mcp.Required(), mcp.Description("The ID of the scene")),
	)
	srv.AddTool(activateSceneTool, mcpserver.HandleActivateScene(client))

	// Create scene
	createSceneTool := mcp.NewTool("create_scene",
		mcp.WithDescription("Create a new scene from current light states"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name for the scene")),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group to capture")),
	)
	srv.AddTool(createSceneTool, mcpserver.HandleCreateScene(client))
}

// registerEffectTools adds native effect tools
func registerEffectTools(srv *server.MCPServer, client *hue.Client) {
	// Set effect on light
	lightEffectTool := mcp.NewTool("light_effect",
		mcp.WithDescription("Set a dynamic effect on a light"),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("The ID of the light")),
		mcp.WithString("effect", mcp.Required(), 
			mcp.Description("Effect to apply"),
			mcp.Enum(effects.None, effects.Candle, effects.Fireplace, effects.Colorloop, effects.Sunrise, effects.Sparkle, effects.Glisten, effects.Opal, effects.Prism),
		),
		mcp.WithNumber("duration", mcp.Description("Duration in seconds (0 for infinite)")),
	)
	srv.AddTool(lightEffectTool, mcpserver.HandleLightEffect(client))

	// Set effect on group
	groupEffectTool := mcp.NewTool("group_effect",
		mcp.WithDescription("Set a dynamic effect on a group"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("The ID of the group")),
		mcp.WithString("effect", mcp.Required(),
			mcp.Description("Effect to apply"),
			mcp.Enum(effects.None, effects.Candle, effects.Fireplace, effects.Colorloop, effects.Sunrise, effects.Sparkle, effects.Glisten, effects.Opal, effects.Prism),
		),
		mcp.WithNumber("duration", mcp.Description("Duration in seconds (0 for infinite)")),
	)
	srv.AddTool(groupEffectTool, mcpserver.HandleGroupEffect(client))
}

// registerSystemTools adds system and discovery tools
func registerSystemTools(srv *server.MCPServer, client *hue.Client) {
	// List lights
	listLightsTool := mcp.NewTool("list_lights",
		mcp.WithDescription("List all available lights"),
	)
	srv.AddTool(listLightsTool, mcpserver.HandleListLights(client))

	// List groups
	listGroupsTool := mcp.NewTool("list_groups",
		mcp.WithDescription("List all available groups/rooms"),
	)
	srv.AddTool(listGroupsTool, mcpserver.HandleListGroups(client))

	// Get light state
	getLightStateTool := mcp.NewTool("get_light_state",
		mcp.WithDescription("Get current state of a light"),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("The ID of the light")),
	)
	srv.AddTool(getLightStateTool, mcpserver.HandleGetLightState(client))

	// Bridge info
	bridgeInfoTool := mcp.NewTool("bridge_info",
		mcp.WithDescription("Get bridge information and capabilities"),
	)
	srv.AddTool(bridgeInfoTool, mcpserver.HandleBridgeInfo(client))

	// Identify light
	identifyLightTool := mcp.NewTool("identify_light",
		mcp.WithDescription("Make a light blink to identify it"),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("The ID of the light")),
	)
	srv.AddTool(identifyLightTool, mcpserver.HandleIdentifyLight(client))
}