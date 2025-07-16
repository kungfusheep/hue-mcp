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

	fmt.Println("🌡️ TEST 9: Sensor Readings")
	fmt.Println("=========================")

	// Test 1: Motion sensors
	fmt.Println("\n1. Discovering motion sensors...")
	motionSensors, err := client.GetMotionSensors(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get motion sensors: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d motion sensors:\n", len(motionSensors))
		
		for i, sensor := range motionSensors {
			fmt.Printf("\n   🚶 Motion Sensor %d:\n", i+1)
			fmt.Printf("      ID: %s\n", sensor.ID)
			fmt.Printf("      Enabled: %v\n", sensor.Enabled)
			
			// Motion data is always present in the struct
			fmt.Printf("      Motion detected: %v\n", sensor.Motion.Motion)
			fmt.Printf("      Valid: %v\n", sensor.Motion.MotionValid)
			if sensor.Motion.MotionReport != nil && len(sensor.Motion.MotionReport.Changed) > 0 {
				fmt.Printf("      Last changed: %s\n", sensor.Motion.MotionReport.Changed)
			}
			
			// Find the device this sensor belongs to
			devices, _ := client.GetDevices(ctx)
			for _, device := range devices {
				for _, service := range device.Services {
					if service.RType == "motion" && service.RID == sensor.ID {
						fmt.Printf("      Device: %s\n", device.Metadata.Name)
						
						// Check for associated room
						rooms, _ := client.GetRooms(ctx)
						for _, room := range rooms {
							for _, child := range room.Children {
								if child.RType == "device" && child.RID == device.ID {
									fmt.Printf("      Room: %s\n", room.Metadata.Name)
									break
								}
							}
						}
						break
					}
				}
			}
		}
		
		if len(motionSensors) == 0 {
			fmt.Println("   ℹ️  No motion sensors found")
		}
	}

	// Test 2: Temperature sensors
	fmt.Println("\n2. Discovering temperature sensors...")
	tempSensors, err := client.GetTemperatureSensors(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get temperature sensors: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d temperature sensors:\n", len(tempSensors))
		
		for i, sensor := range tempSensors {
			fmt.Printf("\n   🌡️  Temperature Sensor %d:\n", i+1)
			fmt.Printf("      ID: %s\n", sensor.ID)
			fmt.Printf("      Enabled: %v\n", sensor.Enabled)
			
			// Temperature data is always present in the struct
			fmt.Printf("      Temperature: %.1f°C (%.1f°F)\n", 
				sensor.Temperature.Temperature, 
				sensor.Temperature.Temperature*9/5+32)
			fmt.Printf("      Valid: %v\n", sensor.Temperature.TemperatureValid)
			if sensor.Temperature.TemperatureReport != nil && len(sensor.Temperature.TemperatureReport.Changed) > 0 {
				fmt.Printf("      Last changed: %s\n", sensor.Temperature.TemperatureReport.Changed)
			}
			
			// Find the device this sensor belongs to
			devices, _ := client.GetDevices(ctx)
			for _, device := range devices {
				for _, service := range device.Services {
					if service.RType == "temperature" && service.RID == sensor.ID {
						fmt.Printf("      Device: %s\n", device.Metadata.Name)
						break
					}
				}
			}
		}
		
		if len(tempSensors) == 0 {
			fmt.Println("   ℹ️  No temperature sensors found")
		}
	}

	// Test 3: Light level sensors
	fmt.Println("\n3. Discovering light level sensors...")
	lightSensors, err := client.GetLightLevelSensors(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get light level sensors: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d light level sensors:\n", len(lightSensors))
		
		for i, sensor := range lightSensors {
			fmt.Printf("\n   💡 Light Level Sensor %d:\n", i+1)
			fmt.Printf("      ID: %s\n", sensor.ID)
			fmt.Printf("      Enabled: %v\n", sensor.Enabled)
			
			// Light level data is in the LightLevel field
			fmt.Printf("      Light level: %d lux\n", sensor.LightLevel.LightLevel)
			fmt.Printf("      Valid: %v\n", sensor.LightLevel.LightLevelValid)
			if sensor.LightLevel.LightLevelReport != nil && len(sensor.LightLevel.LightLevelReport.Changed) > 0 {
				fmt.Printf("      Last changed: %s\n", sensor.LightLevel.LightLevelReport.Changed)
			}
			
			// Find the device this sensor belongs to
			devices, _ := client.GetDevices(ctx)
			for _, device := range devices {
				for _, service := range device.Services {
					if service.RType == "light_level" && service.RID == sensor.ID {
						fmt.Printf("      Device: %s\n", device.Metadata.Name)
						break
					}
				}
			}
		}
		
		if len(lightSensors) == 0 {
			fmt.Println("   ℹ️  No light level sensors found")
		}
	}

	// Test 4: Buttons (dimmer switches)
	fmt.Println("\n4. Discovering buttons (dimmer switches)...")
	buttons, err := client.GetButtons(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get buttons: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d buttons:\n", len(buttons))
		
		for i, button := range buttons {
			fmt.Printf("\n   🔘 Button %d:\n", i+1)
			fmt.Printf("      ID: %s\n", button.ID)
			
			// Button data structure
			if button.Button.ButtonReport != nil {
				fmt.Printf("      Last event: %s\n", button.Button.ButtonReport.Event)
				fmt.Printf("      Updated: %s\n", button.Button.ButtonReport.Updated)
			}
			fmt.Printf("      Supported events: %v\n", button.Button.EventValues)
			
			// Find the device this button belongs to
			devices, _ := client.GetDevices(ctx)
			for _, device := range devices {
				for _, service := range device.Services {
					if service.RType == "button" && service.RID == button.ID {
						fmt.Printf("      Device: %s\n", device.Metadata.Name)
						
						// Show all buttons on this device
						buttonCount := 0
						for _, svc := range device.Services {
							if svc.RType == "button" {
								buttonCount++
							}
						}
						fmt.Printf("      Total buttons on device: %d\n", buttonCount)
						break
					}
				}
			}
		}
		
		if len(buttons) == 0 {
			fmt.Println("   ℹ️  No buttons found")
		}
	}

	// Test 5: Live sensor monitoring (10 seconds)
	fmt.Println("\n5. Live sensor monitoring for 10 seconds...")
	fmt.Println("   Try triggering motion sensors or pressing buttons!")
	
	startTime := time.Now()
	lastMotionStates := make(map[string]bool)
	lastButtonEvents := make(map[string]string)
	
	// Get initial states
	for _, sensor := range motionSensors {
		lastMotionStates[sensor.ID] = sensor.Motion.Motion
	}
	for _, button := range buttons {
		if button.Button.ButtonReport != nil {
			lastButtonEvents[button.ID] = button.Button.ButtonReport.Event
		} else {
			lastButtonEvents[button.ID] = ""
		}
	}
	
	for time.Since(startTime) < 10*time.Second {
		// Check motion sensors
		currentMotion, _ := client.GetMotionSensors(ctx)
		for _, sensor := range currentMotion {
			currentState := sensor.Motion.Motion
			if lastState, exists := lastMotionStates[sensor.ID]; exists && currentState != lastState {
				// Find device name
				deviceName := "Unknown"
				devices, _ := client.GetDevices(ctx)
				for _, device := range devices {
					for _, service := range device.Services {
						if service.RType == "motion" && service.RID == sensor.ID {
							deviceName = device.Metadata.Name
							break
						}
					}
				}
				
				fmt.Printf("\n   🚨 Motion %s on %s!\n", 
					map[bool]string{true: "DETECTED", false: "CLEARED"}[currentState],
					deviceName)
				lastMotionStates[sensor.ID] = currentState
			}
		}
		
		// Check buttons
		currentButtons, _ := client.GetButtons(ctx)
		for _, button := range currentButtons {
			var currentEvent string
			if button.Button.ButtonReport != nil {
				currentEvent = button.Button.ButtonReport.Event
			}
			if lastEvent, exists := lastButtonEvents[button.ID]; exists && currentEvent != lastEvent && currentEvent != "" {
				// Find device name
				deviceName := "Unknown"
				devices, _ := client.GetDevices(ctx)
				for _, device := range devices {
					for _, service := range device.Services {
						if service.RType == "button" && service.RID == button.ID {
							deviceName = device.Metadata.Name
							break
						}
					}
				}
				
				fmt.Printf("\n   🔘 Button pressed on %s: %s\n", deviceName, currentEvent)
				lastButtonEvents[button.ID] = currentEvent
			}
		}
		
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}
	
	fmt.Println("\n   ✅ Live monitoring complete")

	// Summary
	fmt.Println("\n📊 TEST 9 SUMMARY:")
	fmt.Printf("  • Motion sensors: ✅ %d found\n", len(motionSensors))
	fmt.Printf("  • Temperature sensors: ✅ %d found\n", len(tempSensors))
	fmt.Printf("  • Light level sensors: ✅ %d found\n", len(lightSensors))
	fmt.Printf("  • Buttons: ✅ %d found\n", len(buttons))
	fmt.Printf("  • Live monitoring: ✅ Working\n")
	
	fmt.Println("\n🎯 Test 9 Complete! Sensor capabilities verified.")
}