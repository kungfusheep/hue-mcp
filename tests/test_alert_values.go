package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	fmt.Println("🔍 Testing Alert Actions")
	fmt.Println("========================")

	// Get a test light
	lights, err := client.GetLights(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get lights: %v\n", err)
		return
	}

	if len(lights) == 0 {
		fmt.Println("❌ No lights found")
		return
	}

	testLight := lights[0]
	fmt.Printf("Using light: %s (ID: %s)\n\n", testLight.Metadata.Name, testLight.ID)

	// Test different alert values
	alertValues := []string{
		"select",
		"lselect",
		"breathe",
		"okay",
		"channelchange",
		"finish",
		"stop",
		"identify",
		"none",
		"blink",
		"flash",
	}

	fmt.Println("Testing different alert action values...")
	for _, alertValue := range alertValues {
		fmt.Printf("\nTesting alert action: '%s'\n", alertValue)
		
		// Build the update directly
		updateData := map[string]interface{}{
			"alert": map[string]string{
				"action": alertValue,
			},
		}
		
		jsonData, _ := json.Marshal(updateData)
		fmt.Printf("Request body: %s\n", string(jsonData))
		
		// Try the update
		err := client.UpdateLight(ctx, testLight.ID, hue.LightUpdate{
			Alert: &hue.Alert{Action: alertValue},
		})
		
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success! Light should be responding to '%s'\n", alertValue)
			time.Sleep(3 * time.Second)
		}
	}

	// Try signaling effect as an alternative
	fmt.Println("\n\nTesting signaling effect as alternative...")
	
	update := hue.LightUpdate{}
	signalingData := map[string]interface{}{
		"signal": "on",
		"duration": 5000,
		"color": []map[string]interface{}{
			{
				"xy": map[string]float64{
					"x": 0.3,
					"y": 0.3,
				},
			},
		},
	}
	
	// Marshal and unmarshal to convert to proper type
	jsonBytes, _ := json.Marshal(map[string]interface{}{"signaling": signalingData})
	json.Unmarshal(jsonBytes, &update)
	
	err = client.UpdateLight(ctx, testLight.ID, update)
	if err != nil {
		fmt.Printf("❌ Signaling failed: %v\n", err)
	} else {
		fmt.Printf("✅ Signaling might work as an alternative!\n")
	}
}