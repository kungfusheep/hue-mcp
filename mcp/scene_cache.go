package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kungfusheep/hue/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// CachedScene represents a stored lighting scene
type CachedScene struct {
	Name        string                   `json:"name"`
	Commands    []map[string]interface{} `json:"commands"`
	DelayMs     int                      `json:"delay_ms"`
	Description string                   `json:"description"`
	CreatedAt   time.Time                `json:"created_at"`
	UsageCount  int                      `json:"usage_count"`
}

// SceneCache manages cached lighting scenes
type SceneCache struct {
	scenes map[string]*CachedScene
	mu     sync.RWMutex
}

// Global scene cache instance
var globalSceneCache = &SceneCache{
	scenes: make(map[string]*CachedScene),
}

// GetSceneCache returns the global scene cache instance
func GetSceneCache() *SceneCache {
	return globalSceneCache
}

// SaveScene stores a scene in the cache
func (sc *SceneCache) SaveScene(name string, commands []map[string]interface{}, delayMs int, description string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Validate scene name
	if name == "" {
		return fmt.Errorf("scene name cannot be empty")
	}
	if len(commands) == 0 {
		return fmt.Errorf("scene must have at least one command")
	}

	sc.scenes[name] = &CachedScene{
		Name:        name,
		Commands:    commands,
		DelayMs:     delayMs,
		Description: description,
		CreatedAt:   time.Now(),
		UsageCount:  0,
	}

	return nil
}

// GetScene retrieves a scene from the cache
func (sc *SceneCache) GetScene(name string) (*CachedScene, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	scene, exists := sc.scenes[name]
	if !exists {
		return nil, fmt.Errorf("scene '%s' not found", name)
	}

	// Increment usage count
	sc.mu.RUnlock()
	sc.mu.Lock()
	scene.UsageCount++
	sc.mu.Unlock()
	sc.mu.RLock()

	return scene, nil
}

// ListScenes returns all cached scenes
func (sc *SceneCache) ListScenes() []*CachedScene {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	scenes := make([]*CachedScene, 0, len(sc.scenes))
	for _, scene := range sc.scenes {
		scenes = append(scenes, scene)
	}

	return scenes
}

// DeleteScene removes a scene from the cache
func (sc *SceneCache) DeleteScene(name string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if _, exists := sc.scenes[name]; !exists {
		return fmt.Errorf("scene '%s' not found", name)
	}

	delete(sc.scenes, name)
	return nil
}

// HandleRecallScene executes a cached scene
func HandleRecallScene(client *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		sceneName, ok := args["scene_name"].(string)
		if !ok {
			return mcp.NewToolResultError("scene_name is required"), nil
		}

		// Get the cached scene
		scene, err := globalSceneCache.GetScene(sceneName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to recall scene: %v", err)), nil
		}

		// Generate batch ID for tracking
		batchID := fmt.Sprintf("recalled_%s_%d", scene.Name, time.Now().Unix())

		// Execute the scene asynchronously
		go ExecuteBatchAsync(ctx, client, scene.Commands, scene.DelayMs, batchID)

		// Format response
		var description string
		if scene.Description != "" {
			description = fmt.Sprintf("\nDescription: %s", scene.Description)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Recalling atmosphere: %s...%s\nCommands: %d\nDelay: %dms\nBatch ID: %s\nUsage count: %d",
			scene.Name, description, len(scene.Commands), scene.DelayMs, batchID, scene.UsageCount)), nil
	}
}

// HandleListCachedScenes lists all cached scenes
func HandleListCachedScenes(client *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scenes := globalSceneCache.ListScenes()

		if len(scenes) == 0 {
			return mcp.NewToolResultText("No cached scenes available"), nil
		}

		// Sort by usage count (most used first)
		// Simple bubble sort for small lists
		for i := 0; i < len(scenes); i++ {
			for j := i + 1; j < len(scenes); j++ {
				if scenes[j].UsageCount > scenes[i].UsageCount {
					scenes[i], scenes[j] = scenes[j], scenes[i]
				}
			}
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Cached scenes (%d):\n\n", len(scenes)))

		for _, scene := range scenes {
			result.WriteString(fmt.Sprintf("📦 %s\n", scene.Name))
			if scene.Description != "" {
				result.WriteString(fmt.Sprintf("   Description: %s\n", scene.Description))
			}
			result.WriteString(fmt.Sprintf("   Commands: %d | Delay: %dms | Used: %d times\n",
				len(scene.Commands), scene.DelayMs, scene.UsageCount))
			result.WriteString(fmt.Sprintf("   Created: %s\n\n", scene.CreatedAt.Format("2006-01-02 15:04:05")))
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleClearCachedScene removes a cached scene
func HandleClearCachedScene(client *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		sceneName, ok := args["scene_name"].(string)
		if !ok {
			return mcp.NewToolResultError("scene_name is required"), nil
		}

		err := globalSceneCache.DeleteScene(sceneName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to clear scene: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Scene '%s' has been cleared from cache", sceneName)), nil
	}
}

// HandleExportScene exports a cached scene as JSON
func HandleExportScene(client *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		sceneName, ok := args["scene_name"].(string)
		if !ok {
			return mcp.NewToolResultError("scene_name is required"), nil
		}

		scene, err := globalSceneCache.GetScene(sceneName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to export scene: %v", err)), nil
		}

		// Export as JSON for sharing/backup
		jsonData, err := json.MarshalIndent(scene, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize scene: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Scene export for '%s':\n\n```json\n%s\n```", sceneName, string(jsonData))), nil
	}
}