package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/kungfusheep/hue/client"
)

// hueScenesCmd represents the native Hue scenes command group
var hueScenesCmd = &cobra.Command{
	Use:   "hue-scenes",
	Short: "Manage native Hue scenes",
	Long:  `Commands for managing native Philips Hue scenes (not cached scenes).`,
}

// listHueScenesCmd lists all native Hue scenes
var listHueScenesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all native Hue scenes",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		scenes, err := hueClient.GetScenes(ctx)
		if err != nil {
			return fmt.Errorf("failed to list scenes: %w", err)
		}

		if jsonOutput {
			printJSON(scenes)
			return nil
		}

		// Human-readable output
		fmt.Printf("Found %d Hue scenes:\n\n", len(scenes))
		for _, scene := range scenes {
			fmt.Printf("📋 %s\n", scene.Metadata.Name)
			fmt.Printf("   ID: %s\n", scene.ID)
			if scene.IDV1 != "" {
				fmt.Printf("   V1 ID: %s\n", scene.IDV1)
			}
			if scene.Group.RID != "" {
				fmt.Printf("   Group: %s\n", scene.Group.RID)
			}
			fmt.Printf("   Actions: %d\n", len(scene.Actions))
			fmt.Println()
		}
		
		return nil
	},
}

// activateHueSceneCmd activates a native Hue scene
var activateHueSceneCmd = &cobra.Command{
	Use:   "activate <scene-name-or-id>",
	Short: "Activate a native Hue scene",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		
		// Resolve scene name to ID
		sceneID, err := resolveSceneID(ctx, args[0])
		if err != nil {
			return err
		}
		
		err = hueClient.ActivateScene(ctx, sceneID)
		if err != nil {
			return fmt.Errorf("failed to activate scene: %w", err)
		}
		
		printMessage("Scene %s activated", args[0])
		return nil
	},
}

// createHueSceneCmd creates a new Hue scene
var createHueSceneCmd = &cobra.Command{
	Use:   "create <name> <group-name-or-id>",
	Short: "Create a new Hue scene for a group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		ctx := context.Background()
		
		// Resolve group name to ID
		groupID, err := resolveGroupID(ctx, args[1])
		if err != nil {
			return err
		}
		
		// Note: This creates an empty scene. In practice, you'd want to
		// capture current light states or specify actions
		sceneCreate := client.SceneCreate{
			Type: "scene",
			Metadata: client.Metadata{
				Name: name,
			},
			Group: client.ResourceIdentifier{
				RID:   groupID,
				RType: "grouped_light",
			},
			Actions: []client.SceneAction{}, // Empty for now
		}
		
		scene, err := hueClient.CreateScene(ctx, sceneCreate)
		
		if err != nil {
			return fmt.Errorf("failed to create scene: %w", err)
		}
		
		printMessage("Scene '%s' created with ID: %s", name, scene.ID)
		return nil
	},
}

// findHueSceneCmd finds scenes by name
var findHueSceneCmd = &cobra.Command{
	Use:   "find <search-term>",
	Short: "Find scenes by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchTerm := strings.ToLower(args[0])
		ctx := context.Background()
		
		scenes, err := hueClient.GetScenes(ctx)
		if err != nil {
			return fmt.Errorf("failed to list scenes: %w", err)
		}

		// Filter scenes by name
		matchCount := 0
		fmt.Printf("Scenes matching '%s':\n\n", searchTerm)
		
		for _, scene := range scenes {
			if strings.Contains(strings.ToLower(scene.Metadata.Name), searchTerm) {
				fmt.Printf("- %s (ID: %s)\n", scene.Metadata.Name, scene.ID)
				matchCount++
			}
		}

		if matchCount == 0 {
			fmt.Printf("No scenes found matching '%s'\n", searchTerm)
			return nil
		}
		
		return nil
	},
}

func init() {
	// Add subcommands
	hueScenesCmd.AddCommand(listHueScenesCmd)
	hueScenesCmd.AddCommand(activateHueSceneCmd)
	hueScenesCmd.AddCommand(createHueSceneCmd)
	hueScenesCmd.AddCommand(findHueSceneCmd)
	
	// Add to root
	rootCmd.AddCommand(hueScenesCmd)
}