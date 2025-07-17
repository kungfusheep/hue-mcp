package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/kungfusheep/hue/client"
)

var (
	// Stream flags
	streamFilter string
	streamRaw    bool
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream real-time events from sensors and lights",
	Long: `Stream real-time events from the Hue bridge including:
- Motion sensor triggers
- Temperature changes
- Light level changes
- Button presses
- Light state changes

Use filters to focus on specific event types.`,
	RunE: runStream,
}

func runStream(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start event stream
	eventStream, err := hueClient.StreamEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to start event stream: %w", err)
	}
	defer eventStream.Close()

	fmt.Println("🔴 Streaming live events (Ctrl+C to stop)...")
	fmt.Println()

	// Parse filters
	var filters []string
	if streamFilter != "" {
		filters = strings.Split(streamFilter, ",")
		fmt.Printf("Filtering for: %s\n\n", streamFilter)
	}

	// Event loop
	for {
		select {
		case <-sigChan:
			fmt.Println("\n✋ Stopping event stream...")
			return nil

		case event := <-eventStream.Events():
			if shouldShowEvent(event, filters) {
				if streamRaw {
					printJSON(event)
				} else {
					printHumanEvent(event)
				}
			}

		case err := <-eventStream.Errors():
			printError("Stream error: %v", err)
		}
	}
}

func shouldShowEvent(event client.Event, filters []string) bool {
	if len(filters) == 0 {
		return true
	}

	// Check event type against filters
	for _, data := range event.Data {
		for _, filter := range filters {
			filter = strings.ToLower(strings.TrimSpace(filter))
			// Check various event types
			if (filter == "motion" && data.Motion != nil) ||
			   (filter == "temperature" && data.Temperature != nil) ||
			   (filter == "light" && data.Type == "light") ||
			   (filter == "button" && data.Type == "button") ||
			   strings.Contains(strings.ToLower(data.Type), filter) {
				return true
			}
		}
	}
	return false
}

func printHumanEvent(event client.Event) {
	timestamp := time.Now().Format("15:04:05")
	
	for _, data := range event.Data {
		printEventData(timestamp, data)
	}
}

func printEventData(timestamp string, data client.EventData) {
	switch data.Type {
	case "motion":
		if data.Motion != nil {
			if data.Motion.Motion {
				fmt.Printf("[%s] 🚶 Motion detected! (sensor: %s)\n", timestamp, data.ID)
			} else {
				fmt.Printf("[%s] 💤 Motion cleared (sensor: %s)\n", timestamp, data.ID)
			}
		}

	case "temperature":
		if data.Temperature != nil {
			celsius := float64(data.Temperature.Temperature) / 100.0
			fmt.Printf("[%s] 🌡️  Temperature: %.1f°C / %.1f°F (sensor: %s)\n", 
				timestamp, celsius, celsius*9/5+32, data.ID)
		}

	case "light_level":
		if data.Light != nil {
			fmt.Printf("[%s] ☀️  Light level: %d lux (sensor: %s)\n", 
				timestamp, data.Light.LightLevel, data.ID)
		}

	case "button":
		if data.Button != nil && data.Button.ButtonReport != nil {
			fmt.Printf("[%s] 🔘 Button pressed: %s (device: %s)\n", 
				timestamp, data.Button.ButtonReport.Event, data.ID)
		}

	case "light":
		// Light state changes
		changes := []string{}
		
		if data.On != nil {
			if data.On.On {
				changes = append(changes, "turned ON")
			} else {
				changes = append(changes, "turned OFF")
			}
		}
		
		if data.Dimming != nil {
			changes = append(changes, fmt.Sprintf("brightness: %.0f%%", data.Dimming.Brightness))
		}
		
		if data.Color != nil && data.Color.XY.X > 0 {
			changes = append(changes, "color changed")
		}
		
		if data.Effects != nil && data.Effects.Effect != "" {
			changes = append(changes, fmt.Sprintf("effect: %s", data.Effects.Effect))
		}
		
		if len(changes) > 0 {
			fmt.Printf("[%s] 💡 Light %s: %s\n", timestamp, data.ID, strings.Join(changes, ", "))
		}

	default:
		// Generic event
		fmt.Printf("[%s] 📡 %s event: %s\n", timestamp, data.Type, data.ID)
	}
}

func init() {
	streamCmd.Flags().StringVarP(&streamFilter, "filter", "f", "", 
		"Filter events (comma-separated: motion,temperature,light,button)")
	streamCmd.Flags().BoolVarP(&streamRaw, "raw", "r", false, 
		"Show raw JSON events")
	
	rootCmd.AddCommand(streamCmd)
}