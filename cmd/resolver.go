package cmd

import (
	"context"
	"fmt"
	"strings"
	
	"github.com/kungfusheep/hue/client"
)

// resolveLightID takes a name or ID and returns the actual light ID
func resolveLightID(ctx context.Context, nameOrID string) (string, error) {
	// If it looks like a UUID, return it as-is
	if strings.Contains(nameOrID, "-") && len(nameOrID) > 30 {
		return nameOrID, nil
	}
	
	// Otherwise, search for the light by name
	lights, err := hueClient.GetLights(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get lights: %w", err)
	}
	
	// Try exact match first (case-insensitive)
	for _, light := range lights {
		if strings.EqualFold(light.Metadata.Name, nameOrID) {
			return light.ID, nil
		}
	}
	
	// Try partial match
	var matches []struct {
		ID   string
		Name string
	}
	
	searchLower := strings.ToLower(nameOrID)
	for _, light := range lights {
		if strings.Contains(strings.ToLower(light.Metadata.Name), searchLower) {
			matches = append(matches, struct {
				ID   string
				Name string
			}{
				ID:   light.ID,
				Name: light.Metadata.Name,
			})
		}
	}
	
	if len(matches) == 0 {
		return "", fmt.Errorf("no light found matching '%s'", nameOrID)
	}
	
	if len(matches) == 1 {
		return matches[0].ID, nil
	}
	
	// Multiple matches - show them to the user
	return "", fmt.Errorf("multiple lights match '%s':\n%s\nPlease be more specific", 
		nameOrID, formatMatches(matches))
}

// resolveGroupID takes a name or ID and returns the actual group ID
func resolveGroupID(ctx context.Context, nameOrID string) (string, error) {
	// If it looks like a UUID, return it as-is
	if strings.Contains(nameOrID, "-") && len(nameOrID) > 30 {
		return nameOrID, nil
	}
	
	// Search in rooms first (they have names)
	rooms, err := hueClient.GetRooms(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get rooms: %w", err)
	}
	
	// Try exact match first
	for _, room := range rooms {
		if strings.EqualFold(room.Metadata.Name, nameOrID) {
			// Find the grouped_light for this room
			for _, service := range room.Services {
				if service.RType == "grouped_light" {
					return service.RID, nil
				}
			}
		}
	}
	
	// Try partial match
	var matches []struct {
		ID       string
		Name     string
		GroupID  string
	}
	
	searchLower := strings.ToLower(nameOrID)
	for _, room := range rooms {
		if strings.Contains(strings.ToLower(room.Metadata.Name), searchLower) {
			// Find the grouped_light for this room
			groupID := ""
			for _, service := range room.Services {
				if service.RType == "grouped_light" {
					groupID = service.RID
					break
				}
			}
			if groupID != "" {
				matches = append(matches, struct {
					ID      string
					Name    string
					GroupID string
				}{
					ID:      room.ID,
					Name:    room.Metadata.Name,
					GroupID: groupID,
				})
			}
		}
	}
	
	if len(matches) == 0 {
		return "", fmt.Errorf("no room/group found matching '%s'", nameOrID)
	}
	
	if len(matches) == 1 {
		return matches[0].GroupID, nil
	}
	
	// Multiple matches
	var matchInfo []struct {
		ID   string
		Name string
	}
	for _, m := range matches {
		matchInfo = append(matchInfo, struct {
			ID   string
			Name string
		}{
			ID:   m.GroupID,
			Name: m.Name,
		})
	}
	
	return "", fmt.Errorf("multiple rooms match '%s':\n%s\nPlease be more specific", 
		nameOrID, formatMatches(matchInfo))
}

// resolveSceneID takes a name or ID and returns the actual scene ID
func resolveSceneID(ctx context.Context, nameOrID string) (string, error) {
	// If it looks like a UUID, return it as-is
	if strings.Contains(nameOrID, "-") && len(nameOrID) > 30 {
		return nameOrID, nil
	}
	
	// Check if input contains room specifier like "Nightlight:Master Bedroom"
	parts := strings.Split(nameOrID, ":")
	sceneName := strings.TrimSpace(parts[0])
	roomFilter := ""
	if len(parts) == 2 {
		roomFilter = strings.TrimSpace(parts[1])
	}
	
	// Get scenes
	scenes, err := hueClient.GetScenes(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get scenes: %w", err)
	}
	
	// Get rooms and zones for room name lookup
	rooms, err := hueClient.GetRooms(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get rooms: %w", err)
	}
	
	roomIDToName := make(map[string]string)
	for _, room := range rooms {
		roomIDToName[room.ID] = room.Metadata.Name
	}
	
	zones, err := hueClient.GetZones(ctx)
	if err == nil {
		for _, zone := range zones {
			roomIDToName[zone.ID] = zone.Metadata.Name
		}
	}
	
	// Helper to get room name for a scene
	getRoomName := func(scene client.Scene) string {
		if scene.Group.RType == "room" || scene.Group.RType == "zone" {
			return roomIDToName[scene.Group.RID]
		}
		return ""
	}
	
	// If room filter specified, try to find matching scene
	if roomFilter != "" {
		var roomFilterMatches []struct {
			ID       string
			Name     string
			RoomName string
		}
		
		roomFilterLower := strings.ToLower(roomFilter)
		for _, scene := range scenes {
			roomName := getRoomName(scene)
			if strings.EqualFold(scene.Metadata.Name, sceneName) && 
			   strings.Contains(strings.ToLower(roomName), roomFilterLower) {
				roomFilterMatches = append(roomFilterMatches, struct {
					ID       string
					Name     string
					RoomName string
				}{
					ID:       scene.ID,
					Name:     scene.Metadata.Name,
					RoomName: roomName,
				})
			}
		}
		
		if len(roomFilterMatches) == 1 {
			return roomFilterMatches[0].ID, nil
		}
		
		if len(roomFilterMatches) > 1 {
			return "", fmt.Errorf("multiple scenes match '%s' in rooms containing '%s':\n%s\nPlease be more specific", 
				sceneName, roomFilter, formatSceneMatches(roomFilterMatches))
		}
		// If no matches with room filter, continue to show all matches
	}
	
	// Try exact match first (no room filter)
	var exactMatches []struct {
		ID       string
		Name     string
		RoomName string
	}
	
	for _, scene := range scenes {
		if strings.EqualFold(scene.Metadata.Name, sceneName) {
			exactMatches = append(exactMatches, struct {
				ID       string
				Name     string
				RoomName string
			}{
				ID:       scene.ID,
				Name:     scene.Metadata.Name,
				RoomName: getRoomName(scene),
			})
		}
	}
	
	if len(exactMatches) == 1 {
		return exactMatches[0].ID, nil
	}
	
	if len(exactMatches) > 1 {
		// Multiple exact matches - show with room names
		return "", fmt.Errorf("multiple scenes named '%s':\n%s\nSpecify the room like: '%s:Room Name'", 
			sceneName, formatSceneMatches(exactMatches), sceneName)
	}
	
	// Try partial match
	var partialMatches []struct {
		ID       string
		Name     string
		RoomName string
	}
	
	searchLower := strings.ToLower(sceneName)
	for _, scene := range scenes {
		if strings.Contains(strings.ToLower(scene.Metadata.Name), searchLower) {
			partialMatches = append(partialMatches, struct {
				ID       string
				Name     string
				RoomName string
			}{
				ID:       scene.ID,
				Name:     scene.Metadata.Name,
				RoomName: getRoomName(scene),
			})
		}
	}
	
	if len(partialMatches) == 0 {
		return "", fmt.Errorf("no scene found matching '%s'", nameOrID)
	}
	
	if len(partialMatches) == 1 {
		return partialMatches[0].ID, nil
	}
	
	// Multiple matches
	return "", fmt.Errorf("multiple scenes match '%s':\n%s\nPlease be more specific", 
		nameOrID, formatSceneMatches(partialMatches))
}

// formatMatches formats multiple matches for display
func formatMatches(matches []struct {
	ID   string
	Name string
}) string {
	var lines []string
	for _, match := range matches {
		lines = append(lines, fmt.Sprintf("  - %s (ID: %s)", match.Name, match.ID))
	}
	return strings.Join(lines, "\n")
}

// formatSceneMatches formats multiple scene matches with room info
func formatSceneMatches(matches []struct {
	ID       string
	Name     string
	RoomName string
}) string {
	var lines []string
	for _, match := range matches {
		if match.RoomName != "" {
			lines = append(lines, fmt.Sprintf("  - %s (%s) [ID: %s]", match.Name, match.RoomName, match.ID))
		} else {
			lines = append(lines, fmt.Sprintf("  - %s [ID: %s]", match.Name, match.ID))
		}
	}
	return strings.Join(lines, "\n")
}