package cmd

import (
	"context"
	"fmt"
	"strings"
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
	
	// Otherwise, search for the scene by name
	scenes, err := hueClient.GetScenes(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get scenes: %w", err)
	}
	
	// Try exact match first
	for _, scene := range scenes {
		if strings.EqualFold(scene.Metadata.Name, nameOrID) {
			return scene.ID, nil
		}
	}
	
	// Try partial match
	var matches []struct {
		ID   string
		Name string
	}
	
	searchLower := strings.ToLower(nameOrID)
	for _, scene := range scenes {
		if strings.Contains(strings.ToLower(scene.Metadata.Name), searchLower) {
			matches = append(matches, struct {
				ID   string
				Name string
			}{
				ID:   scene.ID,
				Name: scene.Metadata.Name,
			})
		}
	}
	
	if len(matches) == 0 {
		return "", fmt.Errorf("no scene found matching '%s'", nameOrID)
	}
	
	if len(matches) == 1 {
		return matches[0].ID, nil
	}
	
	// Multiple matches
	return "", fmt.Errorf("multiple scenes match '%s':\n%s\nPlease be more specific", 
		nameOrID, formatMatches(matches))
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