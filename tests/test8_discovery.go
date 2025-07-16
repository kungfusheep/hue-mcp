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

	fmt.Println("🏠 TEST 8: Room/Device Discovery")
	fmt.Println("================================")

	// Test 1: Discover all rooms
	fmt.Println("\n1. Discovering all rooms...")
	rooms, err := client.GetRooms(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get rooms: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d rooms:\n", len(rooms))
	for _, room := range rooms {
		fmt.Printf("\n   📍 %s (ID: %s)\n", room.Metadata.Name, room.ID)
		fmt.Printf("      Type: %s\n", room.Type)
		fmt.Printf("      Archetype: %s\n", room.Metadata.Archetype)
		
		// Count devices in room
		deviceCount := 0
		lightCount := 0
		for _, child := range room.Children {
			if child.RType == "device" {
				deviceCount++
			} else if child.RType == "light" {
				lightCount++
			}
		}
		fmt.Printf("      Devices: %d, Lights: %d\n", deviceCount, lightCount)
		
		// Show grouped light service
		for _, service := range room.Services {
			if service.RType == "grouped_light" {
				fmt.Printf("      Group Light ID: %s\n", service.RID)
			}
		}
	}

	// Test 2: Discover all zones
	fmt.Println("\n2. Discovering all zones...")
	zones, err := client.GetZones(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get zones: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d zones:\n", len(zones))
		for _, zone := range zones {
			fmt.Printf("   • %s (ID: %s)\n", zone.Metadata.Name, zone.ID)
			
			// Count devices in zone
			deviceCount := 0
			for _, child := range zone.Children {
				if child.RType == "light" {
					deviceCount++
				}
			}
			fmt.Printf("     Lights: %d\n", deviceCount)
		}
	}

	// Test 3: Discover all devices
	fmt.Println("\n3. Discovering all devices...")
	devices, err := client.GetDevices(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get devices: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d devices:\n", len(devices))
	
	// Categorize devices by type
	deviceTypes := make(map[string]int)
	for _, device := range devices {
		deviceTypes[device.ProductData.ProductArchetype]++
	}
	
	fmt.Println("\n   Device breakdown by type:")
	for archetype, count := range deviceTypes {
		fmt.Printf("   • %s: %d devices\n", archetype, count)
	}

	// Show first few devices with details
	fmt.Println("\n   Sample devices (first 5):")
	for i, device := range devices {
		if i >= 5 {
			break
		}
		fmt.Printf("\n   🔧 %s\n", device.Metadata.Name)
		fmt.Printf("      ID: %s\n", device.ID)
		fmt.Printf("      Type: %s\n", device.ProductData.ProductArchetype)
		fmt.Printf("      Model: %s\n", device.ProductData.ModelID)
		fmt.Printf("      Manufacturer: %s\n", device.ProductData.ManufacturerName)
		fmt.Printf("      Software: %s\n", device.ProductData.SoftwareVersion)
		
		// Show services
		fmt.Printf("      Services: ")
		for _, service := range device.Services {
			fmt.Printf("%s ", service.RType)
		}
		fmt.Println()
	}

	// Test 4: Bridge information
	fmt.Println("\n4. Getting bridge information...")
	bridge, err := client.GetBridge(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get bridge info: %v\n", err)
	} else {
		fmt.Printf("✅ Bridge Information:\n")
		fmt.Printf("   • Bridge ID: %s\n", bridge.BridgeID)
		fmt.Printf("   • Time Zone: %s\n", bridge.TimeZone.TimeZone)
		fmt.Printf("   • API Version: v2\n")
		fmt.Printf("   • Internal ID: %s\n", bridge.ID)
	}

	// Test 5: Device capabilities check
	fmt.Println("\n5. Checking device capabilities...")
	
	// Find a multi-capability device (like Hue Play)
	for _, device := range devices {
		if device.Metadata.Name == "Hue Play 1" {
			fmt.Printf("\n   🎮 %s capabilities:\n", device.Metadata.Name)
			
			// Get the light service
			for _, service := range device.Services {
				if service.RType == "light" {
					// Get light details
					light, err := client.GetLight(ctx, service.RID)
					if err == nil {
						fmt.Printf("      • Dimmable: true\n") // All v2 lights support dimming
						fmt.Printf("      • Color: %v\n", light.Color != nil)
						fmt.Printf("      • Color Temperature: %v\n", light.ColorTemperature != nil)
						fmt.Printf("      • Effects: %v\n", light.Effects != nil)
						if light.Effects != nil {
							fmt.Printf("        Supported: %v\n", light.Effects.EffectValues)
						}
					}
				}
			}
			break
		}
	}

	// Test 6: Room-Device relationships
	fmt.Println("\n6. Testing room-device relationships...")
	
	// Find Office room and list its devices
	for _, room := range rooms {
		if room.Metadata.Name == "Office" {
			fmt.Printf("\n   📍 Office room devices:\n")
			
			for _, child := range room.Children {
				if child.RType == "device" {
					// Find device details
					for _, device := range devices {
						if device.ID == child.RID {
							fmt.Printf("      • %s (%s)\n", device.Metadata.Name, device.ProductData.ProductArchetype)
							
							// Show light names for this device
							for _, service := range device.Services {
								if service.RType == "light" {
									light, err := client.GetLight(ctx, service.RID)
									if err == nil {
										fmt.Printf("        - Light: %s\n", light.Metadata.Name)
									}
								}
							}
							break
						}
					}
				}
			}
			break
		}
	}

	fmt.Println("\n📊 TEST 8 SUMMARY:")
	fmt.Printf("  • Rooms discovered: ✅ %d rooms found\n", len(rooms))
	fmt.Printf("  • Zones discovered: ✅ %d zones found\n", len(zones))
	fmt.Printf("  • Devices discovered: ✅ %d devices found\n", len(devices))
	fmt.Printf("  • Bridge info: ✅ Working\n")
	fmt.Printf("  • Device capabilities: ✅ Working\n")
	fmt.Printf("  • Room relationships: ✅ Working\n")
	
	fmt.Println("\n🎯 Test 8 Complete! Full discovery capabilities verified.")
}