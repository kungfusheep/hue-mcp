package main

import (
	"context"
	"crypto/tls"
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
		server.WithResourceCapabilities(true, false),
	)

	// Register tools
	registerLightTools(srv, hueClient)
	registerGroupTools(srv, hueClient)
	registerSceneTools(srv, hueClient)
	registerEffectTools(srv, hueClient)
	registerSystemTools(srv, hueClient)
	registerRoomTools(srv, hueClient)
	registerSensorTools(srv, hueClient)
	registerEntertainmentTools(srv, hueClient)
	registerBatchTools(srv, hueClient)
	registerEventTools(srv, hueClient)
	registerCRUDTools(srv, hueClient)

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
	// Get supported effects dynamically
	ctx := context.Background()
	supportedEffects, err := client.GetAllSupportedEffects(ctx)
	if err != nil {
		log.Printf("Warning: Could not get supported effects, using defaults: %v", err)
		supportedEffects = effects.GetAllEffects()
	}

	// Set effect on light
	lightEffectTool := mcp.NewTool("light_effect",
		mcp.WithDescription("Set a dynamic effect on a light"),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("The ID of the light")),
		mcp.WithString("effect", mcp.Required(), 
			mcp.Description("Effect to apply"),
			mcp.Enum(supportedEffects...),
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
			mcp.Enum(supportedEffects...),
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

// registerRoomTools adds room and zone control tools
func registerRoomTools(srv *server.MCPServer, client *hue.Client) {
	// List rooms
	listRoomsTool := mcp.NewTool("list_rooms",
		mcp.WithDescription("List all rooms with their lights"),
	)
	srv.AddTool(listRoomsTool, mcpserver.HandleListRooms(client))

	// List zones
	listZonesTool := mcp.NewTool("list_zones",
		mcp.WithDescription("List all zones"),
	)
	srv.AddTool(listZonesTool, mcpserver.HandleListZones(client))

	// List devices
	listDevicesTool := mcp.NewTool("list_devices",
		mcp.WithDescription("List all devices with their details"),
	)
	srv.AddTool(listDevicesTool, mcpserver.HandleListDevices(client))

	// Get device details
	getDeviceTool := mcp.NewTool("get_device",
		mcp.WithDescription("Get detailed information about a device"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("The ID of the device")),
	)
	srv.AddTool(getDeviceTool, mcpserver.HandleGetDevice(client))
}

// registerSensorTools adds sensor reading tools
func registerSensorTools(srv *server.MCPServer, client *hue.Client) {
	// Motion sensors
	listMotionTool := mcp.NewTool("list_motion_sensors",
		mcp.WithDescription("List all motion sensors and their states"),
	)
	srv.AddTool(listMotionTool, mcpserver.HandleListMotionSensors(client))

	// Temperature sensors
	listTempTool := mcp.NewTool("list_temperature_sensors",
		mcp.WithDescription("List all temperature sensors and their readings"),
	)
	srv.AddTool(listTempTool, mcpserver.HandleListTemperatureSensors(client))

	// Light level sensors
	listLightLevelTool := mcp.NewTool("list_light_level_sensors",
		mcp.WithDescription("List all light level sensors and their readings"),
	)
	srv.AddTool(listLightLevelTool, mcpserver.HandleListLightLevelSensors(client))

	// Buttons
	listButtonsTool := mcp.NewTool("list_buttons",
		mcp.WithDescription("List all buttons (dimmer switches) and their last events"),
	)
	srv.AddTool(listButtonsTool, mcpserver.HandleListButtons(client))
}

// registerEntertainmentTools adds entertainment configuration tools
func registerEntertainmentTools(srv *server.MCPServer, client *hue.Client) {
	// List entertainment configurations
	listEntTool := mcp.NewTool("list_entertainment",
		mcp.WithDescription("List all entertainment configurations"),
	)
	srv.AddTool(listEntTool, mcpserver.HandleListEntertainment(client))

	// Start entertainment
	startEntTool := mcp.NewTool("start_entertainment",
		mcp.WithDescription("Start entertainment mode for a configuration"),
		mcp.WithString("config_id", mcp.Required(), mcp.Description("The ID of the entertainment configuration")),
	)
	srv.AddTool(startEntTool, mcpserver.HandleStartEntertainment(client))

	// Stop entertainment
	stopEntTool := mcp.NewTool("stop_entertainment",
		mcp.WithDescription("Stop entertainment mode for a configuration"),
		mcp.WithString("config_id", mcp.Required(), mcp.Description("The ID of the entertainment configuration")),
	)
	srv.AddTool(stopEntTool, mcpserver.HandleStopEntertainment(client))
}

// registerBatchTools adds batch request capability for efficiency
func registerBatchTools(srv *server.MCPServer, client *hue.Client) {
	// Batch commands
	batchTool := mcp.NewTool("batch_commands",
		mcp.WithDescription("Execute multiple commands in a single request for efficiency. Commands format: JSON array of {action, target_id, value, duration}"),
		mcp.WithString("commands", mcp.Required(), mcp.Description("JSON array of commands: [{\"action\":\"light_on\",\"target_id\":\"light_id\"}, {\"action\":\"light_brightness\",\"target_id\":\"light_id\",\"value\":\"75\"}]")),
		mcp.WithNumber("delay_ms", mcp.Description("Delay between commands in milliseconds (default: 100)")),
	)
	srv.AddTool(batchTool, mcpserver.HandleBatchCommands(client))
}

// registerEventTools adds event streaming tools
func registerEventTools(srv *server.MCPServer, client *hue.Client) {
	// Initialize event manager
	mcpserver.InitEventManager(client)
	
	// Start event stream
	startEventTool := mcp.NewTool("start_event_stream",
		mcp.WithDescription("Start real-time event streaming from Hue bridge"),
		mcp.WithString("filter", mcp.Description("Comma-separated event types to filter (e.g., 'light,motion,button')")),
	)
	srv.AddTool(startEventTool, mcpserver.HandleStartEventStream(client))
	
	// Stop event stream
	stopEventTool := mcp.NewTool("stop_event_stream",
		mcp.WithDescription("Stop the event stream"),
	)
	srv.AddTool(stopEventTool, mcpserver.HandleStopEventStream(client))
	
	// Get recent events
	recentEventsTool := mcp.NewTool("get_recent_events",
		mcp.WithDescription("Get recent events from the stream"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of events to return (default: 50)")),
		mcp.WithString("type", mcp.Description("Filter by event type (e.g., 'light', 'motion', 'button')")),
	)
	srv.AddTool(recentEventsTool, mcpserver.HandleGetRecentEvents(client))
	
	// Get stream status
	streamStatusTool := mcp.NewTool("get_event_stream_status",
		mcp.WithDescription("Get the current status of the event stream"),
	)
	srv.AddTool(streamStatusTool, mcpserver.HandleGetEventStreamStatus(client))
}

// registerCRUDTools adds create, update, delete tools
func registerCRUDTools(srv *server.MCPServer, client *hue.Client) {
	// Scene CRUD
	createSceneFromStateTool := mcp.NewTool("create_scene_from_state",
		mcp.WithDescription("Create a new scene capturing current light states"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name for the scene")),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group/room ID to capture")),
	)
	srv.AddTool(createSceneFromStateTool, mcpserver.HandleCreateSceneFromState(client))
	
	updateSceneTool := mcp.NewTool("update_scene",
		mcp.WithDescription("Update a scene's metadata"),
		mcp.WithString("scene_id", mcp.Required(), mcp.Description("Scene ID to update")),
		mcp.WithString("name", mcp.Description("New name for the scene")),
		mcp.WithNumber("speed", mcp.Description("Transition speed (0.0-1.0)")),
	)
	srv.AddTool(updateSceneTool, mcpserver.HandleUpdateScene(client))
	
	deleteSceneTool := mcp.NewTool("delete_scene",
		mcp.WithDescription("Delete a scene"),
		mcp.WithString("scene_id", mcp.Required(), mcp.Description("Scene ID to delete")),
	)
	srv.AddTool(deleteSceneTool, mcpserver.HandleDeleteScene(client))
	
	// Group management
	addLightToGroupTool := mcp.NewTool("add_light_to_group",
		mcp.WithDescription("Add a light to a group/room"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("Light ID to add")),
	)
	srv.AddTool(addLightToGroupTool, mcpserver.HandleAddLightToGroup(client))
	
	removeLightFromGroupTool := mcp.NewTool("remove_light_from_group",
		mcp.WithDescription("Remove a light from a group/room"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("Group ID")),
		mcp.WithString("light_id", mcp.Required(), mcp.Description("Light ID to remove")),
	)
	srv.AddTool(removeLightFromGroupTool, mcpserver.HandleRemoveLightFromGroup(client))
	
	// Zone CRUD
	createZoneTool := mcp.NewTool("create_zone",
		mcp.WithDescription("Create a new zone with specified lights"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name for the zone")),
		mcp.WithString("light_ids", mcp.Required(), mcp.Description("Comma-separated light IDs")),
	)
	srv.AddTool(createZoneTool, mcpserver.HandleCreateZone(client))
	
	updateZoneTool := mcp.NewTool("update_zone",
		mcp.WithDescription("Update a zone"),
		mcp.WithString("zone_id", mcp.Required(), mcp.Description("Zone ID to update")),
		mcp.WithString("name", mcp.Description("New name for the zone")),
		mcp.WithString("light_ids", mcp.Description("Comma-separated light IDs to set")),
	)
	srv.AddTool(updateZoneTool, mcpserver.HandleUpdateZone(client))
	
	deleteZoneTool := mcp.NewTool("delete_zone",
		mcp.WithDescription("Delete a zone"),
		mcp.WithString("zone_id", mcp.Required(), mcp.Description("Zone ID to delete")),
	)
	srv.AddTool(deleteZoneTool, mcpserver.HandleDeleteZone(client))
	
	// Room update
	updateRoomTool := mcp.NewTool("update_room",
		mcp.WithDescription("Update a room's name"),
		mcp.WithString("room_id", mcp.Required(), mcp.Description("Room ID to update")),
		mcp.WithString("name", mcp.Required(), mcp.Description("New name for the room")),
	)
	srv.AddTool(updateRoomTool, mcpserver.HandleUpdateRoom(client))
}