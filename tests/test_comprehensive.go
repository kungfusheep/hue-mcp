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

	if bridgeIP == "" || username == "" {
		fmt.Println("❌ Please set HUE_BRIDGE_IP and HUE_USERNAME environment variables")
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

	fmt.Println("🔍 TEST 1: Basic Light Discovery and Status")
	fmt.Println("===========================================")
	
	// Test connection
	fmt.Println("\n1. Testing bridge connection...")
	if err := client.TestConnection(ctx); err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		return
	}
	fmt.Println("✅ Bridge connection successful!")

	// Get all lights
	fmt.Println("\n2. Discovering lights...")
	lights, err := client.GetLights(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get lights: %v\n", err)
		return
	}
	fmt.Printf("✅ Found %d lights total\n", len(lights))

	// Find office lights
	fmt.Println("\n3. Finding office lights...")
	officeLights := []string{"Office 1", "Office 2", "Office 3", "Office 4", "Petes Office Lamp", "Hue Play 1", "Hue Play 2"}
	var foundLights []hue.Light
	
	for _, light := range lights {
		for _, officeName := range officeLights {
			if light.Metadata.Name == officeName {
				foundLights = append(foundLights, light)
				break
			}
		}
	}
	
	fmt.Printf("✅ Found %d office lights:\n", len(foundLights))
	for _, light := range foundLights {
		status := "❌ OFF"
		if light.On.On {
			status = fmt.Sprintf("✅ ON (%.0f%%)", light.Dimming.Brightness)
		}
		fmt.Printf("  • %s: %s\n", light.Metadata.Name, status)
	}

	// Get office group
	fmt.Println("\n4. Finding office group...")
	rooms, err := client.GetRooms(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get rooms: %v\n", err)
		return
	}
	
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
	
	if officeGroupID == "" {
		fmt.Println("❌ Could not find office group")
		return
	}
	
	fmt.Printf("✅ Found office group: %s\n", officeGroupID)
	
	// Get group status
	group, err := client.GetGroup(ctx, officeGroupID)
	if err != nil {
		fmt.Printf("❌ Failed to get group status: %v\n", err)
		return
	}
	
	groupStatus := "❌ OFF"
	if group.On.On {
		groupStatus = fmt.Sprintf("✅ ON (%.0f%%)", group.Dimming.Brightness)
	}
	fmt.Printf("✅ Office group status: %s\n", groupStatus)

	fmt.Println("\n📊 TEST 1 SUMMARY:")
	fmt.Printf("  • Bridge connection: ✅ Working\n")
	fmt.Printf("  • Light discovery: ✅ Working (%d lights found)\n", len(foundLights))
	fmt.Printf("  • Office group: ✅ Working (%s)\n", officeGroupID)
	fmt.Printf("  • Status reading: ✅ Working\n")
	
	fmt.Println("\n🎯 Test 1 Complete!")
}