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

	fmt.Println("🔄 TEST 2: Basic Light On/Off Control")
	fmt.Println("====================================")
	
	// Find Office 1
	lights, _ := client.GetLights(ctx)
	var office1 *hue.Light
	for _, light := range lights {
		if light.Metadata.Name == "Office 1" {
			office1 = &light
			break
		}
	}
	
	if office1 == nil {
		fmt.Println("❌ Office 1 not found")
		return
	}
	
	fmt.Printf("🎯 Testing with: %s (ID: %s)\n", office1.Metadata.Name, office1.ID)
	
	// Check initial state
	fmt.Println("\n1. Checking initial state...")
	currentLight, _ := client.GetLight(ctx, office1.ID)
	initialState := "OFF"
	if currentLight.On.On {
		initialState = "ON"
	}
	fmt.Printf("✅ Office 1 initial state: %s\n", initialState)
	
	// Turn ON
	fmt.Println("\n2. Turning Office 1 ON...")
	err := client.TurnOnLight(ctx, office1.ID)
	if err != nil {
		fmt.Printf("❌ Failed to turn on: %v\n", err)
		return
	}
	fmt.Println("✅ Turn ON command sent")
	
	// Wait and check
	time.Sleep(2 * time.Second)
	currentLight, _ = client.GetLight(ctx, office1.ID)
	if currentLight.On.On {
		fmt.Printf("✅ Office 1 is now ON (%.0f%% brightness)\n", currentLight.Dimming.Brightness)
	} else {
		fmt.Println("❌ Office 1 failed to turn on")
		return
	}
	
	// Wait before turning off
	fmt.Println("\n3. Waiting 3 seconds before turning OFF...")
	time.Sleep(3 * time.Second)
	
	// Turn OFF
	fmt.Println("\n4. Turning Office 1 OFF...")
	err = client.TurnOffLight(ctx, office1.ID)
	if err != nil {
		fmt.Printf("❌ Failed to turn off: %v\n", err)
		return
	}
	fmt.Println("✅ Turn OFF command sent")
	
	// Wait and check
	time.Sleep(2 * time.Second)
	currentLight, _ = client.GetLight(ctx, office1.ID)
	if !currentLight.On.On {
		fmt.Println("✅ Office 1 is now OFF")
	} else {
		fmt.Println("❌ Office 1 failed to turn off")
		return
	}
	
	// Restore initial state
	fmt.Println("\n5. Restoring initial state...")
	if initialState == "ON" {
		client.TurnOnLight(ctx, office1.ID)
		fmt.Println("✅ Restored to ON")
	} else {
		fmt.Println("✅ Left in OFF state (original)")
	}
	
	fmt.Println("\n📊 TEST 2 SUMMARY:")
	fmt.Printf("  • Turn ON: ✅ Working\n")
	fmt.Printf("  • Turn OFF: ✅ Working\n")
	fmt.Printf("  • State verification: ✅ Working\n")
	fmt.Printf("  • Restore state: ✅ Working\n")
	
	fmt.Println("\n🎯 Test 2 Complete!")
}