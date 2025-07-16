package main

import (
	"fmt"
	"strings"
)

// This is a demonstration of how to discover and categorize lights
// In a real implementation, you would get this data from the Hue API

type Light struct {
	ID       string
	Name     string
	Archetype string
	IsOn     bool
	Brightness float64
}

func main() {
	fmt.Println("=== Philips Hue Light Discovery Demo ===")
	fmt.Println("\nThis demonstrates how the light discovery script would work.")
	fmt.Println("In a real implementation, this data would come from the Hue API.\n")

	// Example lights that might be discovered
	// Based on the codebase analysis, these are typical light names found
	exampleLights := []Light{
		{ID: "abc123", Name: "Office 1", Archetype: "sultan_bulb", IsOn: true, Brightness: 80},
		{ID: "def456", Name: "Office 2", Archetype: "sultan_bulb", IsOn: true, Brightness: 80},
		{ID: "ghi789", Name: "Office 3", Archetype: "sultan_bulb", IsOn: false, Brightness: 0},
		{ID: "jkl012", Name: "Office 4", Archetype: "sultan_bulb", IsOn: false, Brightness: 0},
		{ID: "mno345", Name: "Petes Office Lamp", Archetype: "table_shade", IsOn: true, Brightness: 60},
		{ID: "pqr678", Name: "Hue Play 1", Archetype: "hue_play", IsOn: true, Brightness: 100},
		{ID: "stu901", Name: "Hue Play 2", Archetype: "hue_play", IsOn: true, Brightness: 100},
		{ID: "vwx234", Name: "Living Room Light", Archetype: "pendant_round", IsOn: false, Brightness: 0},
		{ID: "yz5678", Name: "Bedroom Lamp", Archetype: "table_shade", IsOn: false, Brightness: 0},
		{ID: "bcd901", Name: "Kitchen Strip", Archetype: "light_strip", IsOn: true, Brightness: 50},
		{ID: "efg234", Name: "TV Playbar", Archetype: "hue_play", IsOn: true, Brightness: 75},
		{ID: "hij567", Name: "Desk Light", Archetype: "desk_lamp", IsOn: true, Brightness: 90},
	}

	// Display all lights
	fmt.Printf("Found %d total lights:\n", len(exampleLights))
	fmt.Println(strings.Repeat("-", 80))
	
	// Categorize lights
	var officeLights []Light
	var playbars []Light
	var otherLights []Light

	for _, light := range exampleLights {
		// Check if it's an office light
		isOfficeLight := false
		lightNameLower := strings.ToLower(light.Name)
		
		// Check for "office" in the name
		if strings.Contains(lightNameLower, "office") {
			isOfficeLight = true
		}
		
		// Check if it's a Playbar or Play device
		isPlaybar := false
		if strings.Contains(lightNameLower, "play") || 
		   strings.Contains(lightNameLower, "playbar") ||
		   light.Archetype == "hue_play" {
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
		if light.IsOn {
			status = "ON"
			brightness = fmt.Sprintf("%.0f%%", light.Brightness)
		}

		category := ""
		if isOfficeLight {
			category = " [OFFICE]"
		}
		if isPlaybar {
			category += " [PLAYBAR/PLAY]"
		}

		fmt.Printf("%-25s | ID: %-10s | %-3s | %6s | Archetype: %-15s%s\n",
			light.Name, light.ID, status, brightness, light.Archetype, category)
	}

	// Print summary
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Total lights: %d\n", len(exampleLights))
	fmt.Printf("- Office lights: %d\n", len(officeLights))
	fmt.Printf("- Playbar/Play devices: %d (includes non-office Play devices)\n", len(playbars))
	fmt.Printf("- Other lights: %d\n", len(otherLights))

	// Office light details
	fmt.Println("\n=== Office Lights ===")
	for _, light := range officeLights {
		fmt.Printf("- %s (ID: %s, Archetype: %s)\n", light.Name, light.ID, light.Archetype)
	}

	// Playbar/Play device details
	fmt.Println("\n=== Playbar/Play Devices ===")
	for _, light := range playbars {
		location := "Other"
		if strings.Contains(strings.ToLower(light.Name), "office") {
			location = "Office"
		}
		fmt.Printf("- %s (ID: %s, Location: %s)\n", light.Name, light.ID, location)
	}

	// Detection logic explanation
	fmt.Println("\n=== Detection Logic ===")
	fmt.Println("Office lights are identified by:")
	fmt.Println("1. Name containing 'office' (case-insensitive)")
	fmt.Println("2. Specific known names: 'Office 1-4', 'Petes Office Lamp'")
	fmt.Println("\nPlaybar/Play devices are identified by:")
	fmt.Println("1. Name containing 'play' or 'playbar' (case-insensitive)")
	fmt.Println("2. Archetype being 'hue_play'")
	fmt.Println("\nNote: Some lights may be both office lights AND Playbar/Play devices")
	fmt.Println("(e.g., 'Hue Play 1' and 'Hue Play 2' in the office)")

	// Real implementation notes
	fmt.Println("\n=== Real Implementation Notes ===")
	fmt.Println("To use this with actual Hue bridge:")
	fmt.Println("1. Set environment variables:")
	fmt.Println("   export HUE_BRIDGE_IP='your-bridge-ip'")
	fmt.Println("   export HUE_USERNAME='your-hue-username'")
	fmt.Println("2. The script would then:")
	fmt.Println("   - Connect to the Hue bridge using HTTPS")
	fmt.Println("   - Call /clip/v2/resource/light to get all lights")
	fmt.Println("   - Parse the JSON response")
	fmt.Println("   - Apply the same categorization logic shown here")
}