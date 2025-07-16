package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kungfusheep/hue-mcp/hue"
)

func main() {
	bridgeIP := os.Getenv("HUE_BRIDGE_IP")
	username := os.Getenv("HUE_USERNAME")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client := hue.NewClient(bridgeIP, username, httpClient)
	ctx := context.Background()

	fmt.Println("🔍 DEBUG: Group Effects Issue")
	fmt.Println("=============================")
	
	// Find the office group
	rooms, _ := client.GetRooms(ctx)
	var officeGroupID string
	for _, room := range rooms {
		if room.Metadata.Name == "Office" {
			for _, service := range room.Services {
				if service.RType == "grouped_light" {
					officeGroupID = service.RID
					break
				}
			}
			break
		}
	}
	
	fmt.Printf("Office Group ID: %s\n", officeGroupID)
	
	// Check if group supports effects
	fmt.Println("\n1. Checking if group supports effects...")
	group, err := client.GetGroup(ctx, officeGroupID)
	if err != nil {
		fmt.Printf("❌ Error getting group: %v\n", err)
		return
	}
	
	fmt.Printf("Group has effects field: %v\n", group.Effects != nil)
	if group.Effects != nil {
		fmt.Printf("Current group effect: %s\n", group.Effects.Effect)
		fmt.Printf("Group effect values: %v\n", group.Effects.EffectValues)
	}
	
	// Let's try a different approach - set effects on individual lights
	fmt.Println("\n2. Let's try setting candle effect on individual office lights...")
	
	lights, _ := client.GetLights(ctx)
	officeLights := []string{"Office 1", "Office 2", "Office 3", "Office 4", "Petes Office Lamp"}
	
	for _, lightName := range officeLights {
		for _, light := range lights {
			if light.Metadata.Name == lightName {
				fmt.Printf("   Setting candle effect on %s...\n", lightName)
				
				// Make sure it's on first
				client.TurnOnLight(ctx, light.ID)
				time.Sleep(500 * time.Millisecond)
				
				// Set effect
				err := client.SetLightEffect(ctx, light.ID, "candle", 0)
				if err != nil {
					fmt.Printf("   ❌ Failed on %s: %v\n", lightName, err)
				} else {
					fmt.Printf("   ✅ Candle effect set on %s\n", lightName)
				}
				
				// Small delay between lights
				time.Sleep(500 * time.Millisecond)
				break
			}
		}
	}
	
	fmt.Println("\n3. All individual lights should now have candle effect!")
	fmt.Println("   Watch your office lights for 10 seconds...")
	
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Printf("   %d... ", countdown)
		time.Sleep(1 * time.Second)
	}
	
	fmt.Println("\n4. Turning off effects on all individual lights...")
	for _, lightName := range officeLights {
		for _, light := range lights {
			if light.Metadata.Name == lightName {
				client.SetLightEffect(ctx, light.ID, "no_effect", 0)
				break
			}
		}
	}
	
	fmt.Println("\n🎯 Did you see the candle effect on the individual lights?")
	fmt.Println("   If yes: Individual effects work, but group effects might not be supported")
	fmt.Println("   If no: There might be a deeper issue with effects on your light models")
}