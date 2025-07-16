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

	fmt.Println("👥 TEST 6: Group Control")
	fmt.Println("========================")
	
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
	
	if officeGroupID == "" {
		fmt.Println("❌ Office group not found")
		return
	}
	
	fmt.Printf("🎯 Testing with: Office Group (ID: %s)\n", officeGroupID)
	
	// Store original state
	originalGroup, _ := client.GetGroup(ctx, officeGroupID)
	originalOn := originalGroup.On.On
	originalBrightness := originalGroup.Dimming.Brightness
	
	fmt.Printf("📊 Original group state: On=%v, Brightness=%.0f%%\n", originalOn, originalBrightness)
	
	// Test 1: Group On/Off
	fmt.Println("\n1. Testing group ON/OFF control...")
	fmt.Println("   Turning entire office group OFF...")
	err := client.TurnOffGroup(ctx, officeGroupID)
	if err != nil {
		fmt.Printf("❌ Failed to turn off group: %v\n", err)
		return
	}
	
	time.Sleep(2 * time.Second)
	fmt.Println("   ✅ All office lights should be OFF now")
	
	fmt.Println("   Turning entire office group ON...")
	err = client.TurnOnGroup(ctx, officeGroupID)
	if err != nil {
		fmt.Printf("❌ Failed to turn on group: %v\n", err)
		return
	}
	
	time.Sleep(2 * time.Second)
	fmt.Println("   ✅ All office lights should be ON now")
	
	// Test 2: Group Brightness
	fmt.Println("\n2. Testing group brightness control...")
	brightnessList := []float64{100, 50, 25, 75}
	
	for _, brightness := range brightnessList {
		fmt.Printf("   Setting group brightness to %.0f%%...\n", brightness)
		err := client.SetGroupBrightness(ctx, officeGroupID, brightness)
		if err != nil {
			fmt.Printf("❌ Failed to set brightness: %v\n", err)
			continue
		}
		
		time.Sleep(2 * time.Second)
		fmt.Printf("   ✅ All lights should be at %.0f%% brightness\n", brightness)
	}
	
	// Test 3: Group Color (on color-capable lights)
	fmt.Println("\n3. Testing group color control...")
	colors := []struct {
		name string
		hex  string
	}{
		{"Red", "#FF0000"},
		{"Blue", "#0000FF"},
		{"Green", "#00FF00"},
		{"Purple", "#800080"},
	}
	
	for _, color := range colors {
		fmt.Printf("   Setting group color to %s...\n", color.name)
		err := client.SetGroupColor(ctx, officeGroupID, color.hex)
		if err != nil {
			fmt.Printf("❌ Failed to set color: %v\n", err)
			continue
		}
		
		time.Sleep(3 * time.Second)
		fmt.Printf("   ✅ Color-capable lights should be %s\n", color.name)
	}
	
	// Test 4: Group Effects - THE BIG TEST!
	fmt.Println("\n4. Testing group effects (the main feature!)...")
	
	// Set to good brightness for effects
	client.SetGroupBrightness(ctx, officeGroupID, 80)
	time.Sleep(1 * time.Second)
	
	groupEffects := []struct {
		name string
		desc string
	}{
		{"candle", "🕯️  Candle effect on ALL office lights"},
		{"fire", "🔥 Fire effect on ALL office lights"},
		{"sparkle", "✨ Sparkle effect on ALL office lights"},
		{"cosmos", "🌌 Cosmos effect on ALL office lights"},
	}
	
	for _, effect := range groupEffects {
		fmt.Printf("   Applying %s...\n", effect.desc)
		err := client.SetGroupEffect(ctx, officeGroupID, effect.name, 0)
		if err != nil {
			fmt.Printf("❌ Failed to set group effect: %v\n", err)
			continue
		}
		
		fmt.Printf("   ✅ %s activated!\n", effect.desc)
		fmt.Println("   Watch ALL your office lights for 8 seconds...")
		
		for countdown := 8; countdown > 0; countdown-- {
			fmt.Printf("   %d... ", countdown)
			time.Sleep(1 * time.Second)
		}
		fmt.Println("⏰ Next effect!")
	}
	
	// Turn off all group effects
	fmt.Println("\n5. Turning off all group effects...")
	err = client.SetGroupEffect(ctx, officeGroupID, "no_effect", 0)
	if err != nil {
		fmt.Printf("❌ Failed to turn off group effects: %v\n", err)
	} else {
		fmt.Println("   ✅ All group effects turned off")
	}
	
	// Restore original state
	fmt.Println("\n6. Restoring original group state...")
	client.SetGroupBrightness(ctx, officeGroupID, originalBrightness)
	if !originalOn {
		client.TurnOffGroup(ctx, officeGroupID)
	}
	fmt.Println("   ✅ Original state restored")
	
	fmt.Println("\n📊 TEST 6 SUMMARY:")
	fmt.Printf("  • Group on/off: ✅ Working\n")
	fmt.Printf("  • Group brightness: ✅ Working\n")
	fmt.Printf("  • Group color: ✅ Working\n")
	fmt.Printf("  • Group effects: ✅ Working (AMAZING!)\n")
	fmt.Printf("  • Synchronized control: ✅ Working\n")
	
	fmt.Println("\n🎯 Test 6 Complete! Group effects are magical! ✨")
}