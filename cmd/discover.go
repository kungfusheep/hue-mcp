package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// DiscoveredBridge represents a bridge found via discovery
type DiscoveredBridge struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
}

// BridgeInfo represents additional bridge information
type BridgeInfo struct {
	Bridge    DiscoveredBridge
	Reachable bool
	Name      string
	APIVersion string
	SoftwareVersion string
}

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover Hue bridges on your network",
	Long: `Discover Philips Hue bridges using the official Philips discovery service.
This will find all bridges registered to your network and test their connectivity.`,
	RunE: runDiscover,
}

func runDiscover(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Discovering Hue bridges...")
	
	// Discover bridges using official API
	bridges, err := discoverBridges()
	if err != nil {
		fmt.Printf("❌ %v\n\n", err)
		fmt.Println("🔧 Manual discovery alternatives:")
		fmt.Println("   1. Check your router's admin page for connected devices")
		fmt.Println("   2. Use a network scanner: nmap -sn 192.168.1.0/24")
		fmt.Println("   3. Look for devices with 'Philips' in the manufacturer")
		fmt.Println("   4. Test connectivity: curl http://IP_ADDRESS/api/config")
		return nil
	}

	if len(bridges) == 0 {
		fmt.Println("❌ No Hue bridges found on your network")
		fmt.Println("\nTroubleshooting:")
		fmt.Println("  • Make sure your bridge is connected to the same network")
		fmt.Println("  • Check that your bridge has internet connectivity")
		fmt.Println("  • Try again in a few seconds")
		return nil
	}

	if jsonOutput {
		printJSON(bridges)
		return nil
	}

	// Test connectivity and get additional info for each bridge
	fmt.Printf("\nFound %d Hue bridge(s):\n\n", len(bridges))
	
	for i, bridge := range bridges {
		fmt.Printf("🌉 Bridge %d\n", i+1)
		fmt.Printf("   ID: %s\n", bridge.ID)
		fmt.Printf("   IP Address: %s\n", bridge.InternalIPAddress)
		
		// Test connectivity
		info := testBridgeConnectivity(bridge)
		if info.Reachable {
			fmt.Printf("   Status: ✅ Reachable\n")
			if info.Name != "" {
				fmt.Printf("   Name: %s\n", info.Name)
			}
			if info.SoftwareVersion != "" {
				fmt.Printf("   Software: %s\n", info.SoftwareVersion)
			}
		} else {
			fmt.Printf("   Status: ❌ Not reachable\n")
		}
		
		fmt.Println()
	}

	// Show usage instructions
	if len(bridges) > 0 {
		primaryBridge := bridges[0]
		fmt.Println("📋 To use this bridge:")
		fmt.Printf("   export HUE_BRIDGE_IP=\"%s\"\n", primaryBridge.InternalIPAddress)
		fmt.Println("   # Get API username by pressing the bridge button and running:")
		fmt.Printf("   curl -X POST http://%s/api -d '{\"devicetype\":\"hue#cli\"}'\n", primaryBridge.InternalIPAddress)
	}

	return nil
}

func discoverBridges() ([]DiscoveredBridge, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Retry with exponential backoff for rate limits
	maxRetries := 3
	baseDelay := 2 * time.Second
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(attempt) * baseDelay
			fmt.Printf("Rate limited, retrying in %v...\n", delay)
			time.Sleep(delay)
		}
		
		resp, err := client.Get("https://discovery.meethue.com/")
		if err != nil {
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("failed to contact discovery service after %d attempts: %w", maxRetries, err)
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("discovery service rate limited (try again in a few minutes)")
			}
			resp.Body.Close()
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("discovery service returned status %d", resp.StatusCode)
		}

		var bridges []DiscoveredBridge
		if err := json.NewDecoder(resp.Body).Decode(&bridges); err != nil {
			return nil, fmt.Errorf("failed to parse discovery response: %w", err)
		}

		return bridges, nil
	}
	
	return nil, fmt.Errorf("failed to discover bridges after %d attempts", maxRetries)
}

func testBridgeConnectivity(bridge DiscoveredBridge) BridgeInfo {
	info := BridgeInfo{
		Bridge:    bridge,
		Reachable: false,
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Test basic connectivity with bridge config endpoint
	url := fmt.Sprintf("http://%s/api/config", bridge.InternalIPAddress)
	resp, err := client.Get(url)
	if err != nil {
		return info
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return info
	}

	// Parse basic bridge info
	var config struct {
		Name           string `json:"name"`
		SoftwareVersion string `json:"swversion"`
		APIVersion     string `json:"apiversion"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&config); err == nil {
		info.Name = config.Name
		info.SoftwareVersion = config.SoftwareVersion
		info.APIVersion = config.APIVersion
	}

	info.Reachable = true
	return info
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}