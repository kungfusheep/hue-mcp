package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/kungfusheep/hue-mcp/hue"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// HandleCreateSceneFromState creates a scene from current light states
func HandleCreateSceneFromState(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		name, ok := args["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		
		groupID, ok := args["group_id"].(string)
		if !ok || groupID == "" {
			return mcp.NewToolResultError("group_id is required"), nil
		}
		
		scene, err := client.CreateSceneFromCurrentState(ctx, name, groupID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create scene: %v", err)), nil
		}
		
		return mcp.NewToolResultText(fmt.Sprintf("Scene '%s' created successfully with ID: %s", name, scene.ID)), nil
	}
}

// HandleUpdateScene updates a scene
func HandleUpdateScene(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		sceneID, ok := args["scene_id"].(string)
		if !ok || sceneID == "" {
			return mcp.NewToolResultError("scene_id is required"), nil
		}
		
		update := hue.SceneUpdate{}
		
		if name, ok := args["name"].(string); ok && name != "" {
			update.Metadata = &hue.Metadata{Name: name}
		}
		
		if speed, ok := args["speed"].(float64); ok {
			update.Speed = &speed
		}
		
		err := client.UpdateScene(ctx, sceneID, update)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update scene: %v", err)), nil
		}
		
		return mcp.NewToolResultText("Scene updated successfully"), nil
	}
}

// HandleDeleteScene deletes a scene
func HandleDeleteScene(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		sceneID, ok := args["scene_id"].(string)
		if !ok || sceneID == "" {
			return mcp.NewToolResultError("scene_id is required"), nil
		}
		
		err := client.DeleteScene(ctx, sceneID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete scene: %v", err)), nil
		}
		
		return mcp.NewToolResultText(fmt.Sprintf("Scene %s deleted successfully", sceneID)), nil
	}
}

// HandleAddLightToGroup adds a light to a group
func HandleAddLightToGroup(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		groupID, ok := args["group_id"].(string)
		if !ok || groupID == "" {
			return mcp.NewToolResultError("group_id is required"), nil
		}
		
		lightID, ok := args["light_id"].(string)
		if !ok || lightID == "" {
			return mcp.NewToolResultError("light_id is required"), nil
		}
		
		err := client.AddLightToGroup(ctx, groupID, lightID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to add light to group: %v", err)), nil
		}
		
		return mcp.NewToolResultText(fmt.Sprintf("Light %s added to group %s", lightID, groupID)), nil
	}
}

// HandleRemoveLightFromGroup removes a light from a group
func HandleRemoveLightFromGroup(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		groupID, ok := args["group_id"].(string)
		if !ok || groupID == "" {
			return mcp.NewToolResultError("group_id is required"), nil
		}
		
		lightID, ok := args["light_id"].(string)
		if !ok || lightID == "" {
			return mcp.NewToolResultError("light_id is required"), nil
		}
		
		err := client.RemoveLightFromGroup(ctx, groupID, lightID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to remove light from group: %v", err)), nil
		}
		
		return mcp.NewToolResultText(fmt.Sprintf("Light %s removed from group %s", lightID, groupID)), nil
	}
}

// HandleCreateZone creates a new zone
func HandleCreateZone(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		name, ok := args["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		
		// Get light IDs
		lightIDsStr, ok := args["light_ids"].(string)
		if !ok || lightIDsStr == "" {
			return mcp.NewToolResultError("light_ids is required (comma-separated)"), nil
		}
		
		lightIDs := strings.Split(lightIDsStr, ",")
		var children []hue.ResourceIdentifier
		for _, id := range lightIDs {
			id = strings.TrimSpace(id)
			if id != "" {
				children = append(children, hue.ResourceIdentifier{
					RID:   id,
					RType: "light",
				})
			}
		}
		
		if len(children) == 0 {
			return mcp.NewToolResultError("at least one light ID is required"), nil
		}
		
		zoneCreate := hue.ZoneCreate{
			Type: "zone",
			Metadata: hue.Metadata{
				Name: name,
			},
			Children: children,
		}
		
		zone, err := client.CreateZone(ctx, zoneCreate)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create zone: %v", err)), nil
		}
		
		return mcp.NewToolResultText(fmt.Sprintf("Zone '%s' created with ID: %s", name, zone.ID)), nil
	}
}

// HandleUpdateZone updates a zone
func HandleUpdateZone(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		zoneID, ok := args["zone_id"].(string)
		if !ok || zoneID == "" {
			return mcp.NewToolResultError("zone_id is required"), nil
		}
		
		update := hue.ZoneUpdate{}
		
		if name, ok := args["name"].(string); ok && name != "" {
			update.Metadata = &hue.Metadata{Name: name}
		}
		
		if lightIDsStr, ok := args["light_ids"].(string); ok && lightIDsStr != "" {
			lightIDs := strings.Split(lightIDsStr, ",")
			var children []hue.ResourceIdentifier
			for _, id := range lightIDs {
				id = strings.TrimSpace(id)
				if id != "" {
					children = append(children, hue.ResourceIdentifier{
						RID:   id,
						RType: "light",
					})
				}
			}
			update.Children = children
		}
		
		err := client.UpdateZone(ctx, zoneID, update)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update zone: %v", err)), nil
		}
		
		return mcp.NewToolResultText("Zone updated successfully"), nil
	}
}

// HandleDeleteZone deletes a zone
func HandleDeleteZone(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		zoneID, ok := args["zone_id"].(string)
		if !ok || zoneID == "" {
			return mcp.NewToolResultError("zone_id is required"), nil
		}
		
		err := client.DeleteZone(ctx, zoneID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete zone: %v", err)), nil
		}
		
		return mcp.NewToolResultText(fmt.Sprintf("Zone %s deleted successfully", zoneID)), nil
	}
}

// HandleUpdateRoom updates a room's metadata
func HandleUpdateRoom(client *hue.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		roomID, ok := args["room_id"].(string)
		if !ok || roomID == "" {
			return mcp.NewToolResultError("room_id is required"), nil
		}
		
		name, ok := args["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		
		update := hue.RoomUpdate{
			Metadata: &hue.Metadata{Name: name},
		}
		
		err := client.UpdateRoom(ctx, roomID, update)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update room: %v", err)), nil
		}
		
		return mcp.NewToolResultText(fmt.Sprintf("Room renamed to '%s'", name)), nil
	}
}