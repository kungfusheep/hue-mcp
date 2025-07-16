package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kungfusheep/hue-mcp/hue"
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

	fmt.Println("🔍 Finding All Office Lights")
	fmt.Println("============================")

	// Get all lights
	lights, err := client.GetLights(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get lights: %v\n", err)
		return
	}

	// Get all rooms to find office
	rooms, err := client.GetRooms(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get rooms: %v\n", err)
		return
	}

	// Find office room
	var officeRoomID string
	for _, room := range rooms {
		if room.Metadata.Name == "Office" {
			officeRoomID = room.ID
			fmt.Printf("✅ Found Office room (ID: %s)\n", officeRoomID)
			
			// List all lights in this room
			fmt.Println("\n📋 Lights in Office room:")
			for _, child := range room.Children {
				if child.RType == "device" {
					fmt.Printf("  • Device: %s\n", child.RID)
				}
			}
			break
		}
	}

	fmt.Printf("\n🎯 Total lights found: %d\n", len(lights))
	
	// Categorize lights
	var officeLights []hue.Light
	var playbarLights []hue.Light
	var otherLights []hue.Light

	for _, light := range lights {
		name := light.Metadata.Name
		archetype := light.Metadata.Archetype
		
		// Check if it's office-related
		if strings.Contains(strings.ToLower(name), "office") || 
		   strings.Contains(strings.ToLower(name), "pete") {
			officeLights = append(officeLights, light)
		} else if strings.Contains(strings.ToLower(name), "play") || 
		          archetype == "hue_play" ||
		          strings.Contains(strings.ToLower(name), "bar") {
			playbarLights = append(playbarLights, light)
		} else {
			otherLights = append(otherLights, light)
		}
	}

	// Display office lights
	fmt.Printf("\n📍 Office Lights (%d):\n", len(officeLights))
	for _, light := range officeLights {
		status := "off"
		if light.On.On {
			status = fmt.Sprintf("on, brightness: %.0f%%", light.Dimming.Brightness)
		}
		fmt.Printf("  • %s (%s)\n", light.Metadata.Name, light.Metadata.Archetype)
		fmt.Printf("    ID: %s\n", light.ID)
		fmt.Printf("    Status: %s\n", status)
		if light.Effects != nil && len(light.Effects.EffectValues) > 0 {
			fmt.Printf("    Supports effects: %v\n", light.Effects.EffectValues)
		}
	}

	// Display playbar/play lights
	fmt.Printf("\n🎮 Playbar/Play Lights (%d):\n", len(playbarLights))
	for _, light := range playbarLights {
		status := "off"
		if light.On.On {
			status = fmt.Sprintf("on, brightness: %.0f%%", light.Dimming.Brightness)
		}
		fmt.Printf("  • %s (%s)\n", light.Metadata.Name, light.Metadata.Archetype)
		fmt.Printf("    ID: %s\n", light.ID)
		fmt.Printf("    Status: %s\n", status)
		if light.Effects != nil && len(light.Effects.EffectValues) > 0 {
			fmt.Printf("    Supports effects: %v\n", light.Effects.EffectValues)
		}
	}

	// Show a few other lights for context
	fmt.Printf("\n🏠 Other Lights (showing first 5):\n")
	for i, light := range otherLights {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(otherLights)-5)
			break
		}
		fmt.Printf("  • %s (%s)\n", light.Metadata.Name, light.Metadata.Archetype)
	}

	// Test if we're missing any lights
	fmt.Println("\n❓ Looking for potential missed office lights...")
	for _, light := range lights {
		name := strings.ToLower(light.Metadata.Name)
		// Check for any light that might be in the office but not caught by our filters
		if strings.Contains(name, "speaker") || 
		   strings.Contains(name, "desk") ||
		   strings.Contains(name, "monitor") ||
		   strings.Contains(name, "screen") {
			found := false
			for _, ol := range officeLights {
				if ol.ID == light.ID {
					found = true
					break
				}
			}
			for _, pl := range playbarLights {
				if pl.ID == light.ID {
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("  • Potential office light: %s (%s)\n", light.Metadata.Name, light.Metadata.Archetype)
			}
		}
	}
}