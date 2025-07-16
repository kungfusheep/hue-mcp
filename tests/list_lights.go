package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/kungfusheep/hue-mcp/hue"
)

func main() {
	// Get configuration from environment variables
	bridgeIP := os.Getenv("HUE_BRIDGE_IP")
	username := os.Getenv("HUE_USERNAME")
	
	if bridgeIP == "" || username == "" {
		log.Fatal("Please set HUE_BRIDGE_IP and HUE_USERNAME environment variables")
	}

	// Create HTTP client that skips certificate verification
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 10 * time.Second,
	}

	// Create Hue client
	client := hue.NewClient(bridgeIP, username, httpClient)

	// Test connection
	ctx := context.Background()
	if err := client.TestConnection(ctx); err != nil {
		log.Fatal("Failed to connect to Hue bridge:", err)
	}

	fmt.Println("=== Philips Hue Light Discovery ===")
	fmt.Println()

	// Get all lights
	lights, err := client.GetLights(ctx)
	if err != nil {
		log.Fatal("Error getting lights:", err)
	}

	// Sort lights by name for consistent output
	sort.Slice(lights, func(i, j int) bool {
		return lights[i].Metadata.Name < lights[j].Metadata.Name
	})

	// Categorize lights
	var officeLights []hue.Light
	var playbars []hue.Light
	var otherLights []hue.Light

	fmt.Printf("Found %d total lights:\n", len(lights))
	fmt.Println(strings.Repeat("-", 80))

	for _, light := range lights {
		// Check if it's an office light
		isOfficeLight := false
		lightNameLower := strings.ToLower(light.Metadata.Name)
		
		// Check for office in the name
		if strings.Contains(lightNameLower, "office") {
			isOfficeLight = true
		}
		
		// Check for specific office light names
		officeNames := []string{"office 1", "office 2", "office 3", "office 4", "petes office lamp"}
		for _, officeName := range officeNames {
			if lightNameLower == officeName {
				isOfficeLight = true
				break
			}
		}

		// Check if it's a Playbar or Play device
		isPlaybar := false
		if strings.Contains(lightNameLower, "play") || 
		   strings.Contains(lightNameLower, "playbar") ||
		   strings.Contains(light.Metadata.Archetype, "hue_play") {
			isPlaybar = true
		}

		// Categorize the light
		if isOfficeLight {
			officeLights = append(officeLights, light)
		} else if isPlaybar {
			playbars = append(playbars, light)
		} else {
			otherLights = append(otherLights, light)
		}

		// Print light details
		status := "OFF"
		brightness := "N/A"
		if light.On.On {
			status = "ON"
			brightness = fmt.Sprintf("%.0f%%", light.Dimming.Brightness)
		}

		category := ""
		if isOfficeLight {
			category = " [OFFICE]"
		}
		if isPlaybar {
			category += " [PLAYBAR/PLAY]"
		}

		fmt.Printf("%-30s | ID: %-36s | %-3s | %6s | Archetype: %-15s%s\n",
			light.Metadata.Name,
			light.ID,
			status,
			brightness,
			light.Metadata.Archetype,
			category)
	}

	// Print summary
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Total lights: %d\n", len(lights))
	fmt.Printf("- Office lights: %d\n", len(officeLights))
	fmt.Printf("- Playbar/Play devices: %d\n", len(playbars))
	fmt.Printf("- Other lights: %d\n", len(otherLights))

	// Get room information
	fmt.Println("\n=== Room Information ===")
	rooms, err := client.GetRooms(ctx)
	if err != nil {
		log.Printf("Error getting rooms: %v", err)
	} else {
		for _, room := range rooms {
			if strings.ToLower(room.Metadata.Name) == "office" {
				fmt.Printf("\nOffice Room found (ID: %s)\n", room.ID)
				
				// Get the grouped light service for this room
				for _, service := range room.Services {
					if service.RType == "grouped_light" {
						group, err := client.GetGroup(ctx, service.RID)
						if err == nil {
							fmt.Printf("Office Group Light ID: %s\n", service.RID)
							fmt.Printf("Office Group Status: ")
							if group.On.On {
								fmt.Printf("ON (%.0f%%)\n", group.Dimming.Brightness)
							} else {
								fmt.Printf("OFF\n")
							}
						}
					}
				}
				
				// List lights in the office room
				fmt.Printf("\nLights in Office room:\n")
				for _, child := range room.Children {
					if child.RType == "device" {
						// Find the light associated with this device
						for _, light := range lights {
							if light.Owner.RID == child.RID {
								fmt.Printf("  - %s (ID: %s)\n", light.Metadata.Name, light.ID)
							}
						}
					}
				}
			}
		}
	}

	// Print detailed office lights
	if len(officeLights) > 0 {
		fmt.Println("\n=== Detailed Office Light Information ===")
		for _, light := range officeLights {
			fmt.Printf("\nLight: %s\n", light.Metadata.Name)
			fmt.Printf("  ID: %s\n", light.ID)
			fmt.Printf("  Archetype: %s\n", light.Metadata.Archetype)
			fmt.Printf("  Status: %s\n", func() string {
				if light.On.On {
					return fmt.Sprintf("ON (%.0f%%)", light.Dimming.Brightness)
				}
				return "OFF"
			}())
			
			// Check for effects support
			if light.Effects != nil && len(light.Effects.EffectValues) > 0 {
				fmt.Printf("  Supported effects: %v\n", light.Effects.EffectValues)
				fmt.Printf("  Current effect: %s\n", light.Effects.Effect)
			}
			
			// Check for color support
			if light.Color != nil {
				fmt.Printf("  Color support: Yes (Gamut type: %s)\n", light.Color.GamutType)
			}
		}
	}

	// Look for potential Playbars with non-standard names
	fmt.Println("\n=== Searching for potential Playbars/Play devices ===")
	fmt.Println("(Looking at archetype and capabilities)")
	
	for _, light := range lights {
		// Check various indicators of Playbar/Play devices
		potentialPlaybar := false
		reasons := []string{}
		
		// Check archetype
		archetype := strings.ToLower(light.Metadata.Archetype)
		if strings.Contains(archetype, "play") {
			potentialPlaybar = true
			reasons = append(reasons, fmt.Sprintf("archetype=%s", light.Metadata.Archetype))
		}
		
		// Check if it's in the office and has certain characteristics
		isInOffice := false
		for _, officeLight := range officeLights {
			if officeLight.ID == light.ID {
				isInOffice = true
				break
			}
		}
		
		// Hue Play devices typically have certain characteristics
		if isInOffice && light.Color != nil && light.Effects != nil {
			// Check if it's not already identified as office numbered lights
			lightNameLower := strings.ToLower(light.Metadata.Name)
			if !strings.HasPrefix(lightNameLower, "office ") && 
			   lightNameLower != "petes office lamp" {
				potentialPlaybar = true
				reasons = append(reasons, "in office + has color/effects")
			}
		}
		
		if potentialPlaybar && !strings.Contains(strings.ToLower(light.Metadata.Name), "play") {
			fmt.Printf("\nPotential Playbar/Play: %s\n", light.Metadata.Name)
			fmt.Printf("  ID: %s\n", light.ID)
			fmt.Printf("  Reasons: %s\n", strings.Join(reasons, ", "))
		}
	}
}