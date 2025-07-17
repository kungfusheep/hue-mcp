package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// sensorsCmd represents the sensors command group
var sensorsCmd = &cobra.Command{
	Use:   "sensors",
	Short: "Monitor sensors",
	Long:  `Commands for viewing sensor states (motion, temperature, light level).`,
}

// listMotionSensorsCmd lists all motion sensors
var listMotionSensorsCmd = &cobra.Command{
	Use:   "motion",
	Short: "List motion sensors and their states",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		sensors, err := hueClient.GetMotionSensors(ctx)
		if err != nil {
			return fmt.Errorf("failed to get motion sensors: %w", err)
		}

		if jsonOutput {
			printJSON(sensors)
			return nil
		}

		// Human-readable output
		fmt.Printf("Found %d motion sensors:\n\n", len(sensors))
		for _, sensor := range sensors {
			motionState := "No motion"
			if sensor.Motion.Motion {
				motionState = "🚶 Motion detected"
			}
			
			fmt.Printf("👁️  Motion Sensor %s\n", sensor.ID)
			fmt.Printf("   State: %s\n", motionState)
			fmt.Printf("   Valid: %v\n", sensor.Motion.MotionValid)
			fmt.Printf("   Enabled: %v\n", sensor.Enabled)
			fmt.Println()
		}
		
		return nil
	},
}

// listTemperatureSensorsCmd lists all temperature sensors
var listTemperatureSensorsCmd = &cobra.Command{
	Use:   "temperature",
	Short: "List temperature sensors and their readings",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		sensors, err := hueClient.GetTemperatureSensors(ctx)
		if err != nil {
			return fmt.Errorf("failed to get temperature sensors: %w", err)
		}

		if jsonOutput {
			printJSON(sensors)
			return nil
		}

		// Human-readable output
		fmt.Printf("Found %d temperature sensors:\n\n", len(sensors))
		for _, sensor := range sensors {
			temp := sensor.Temperature.Temperature / 100.0 // Convert to Celsius
			
			fmt.Printf("🌡️  Temperature Sensor %s\n", sensor.ID)
			fmt.Printf("   Temperature: %.1f°C (%.1f°F)\n", temp, temp*9/5+32)
			fmt.Printf("   Valid: %v\n", sensor.Temperature.TemperatureValid)
			fmt.Printf("   Enabled: %v\n", sensor.Enabled)
			fmt.Println()
		}
		
		return nil
	},
}

// listLightSensorsCmd lists all light level sensors
var listLightSensorsCmd = &cobra.Command{
	Use:   "light",
	Short: "List light level sensors and their readings",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		sensors, err := hueClient.GetLightLevelSensors(ctx)
		if err != nil {
			return fmt.Errorf("failed to get light sensors: %w", err)
		}

		if jsonOutput {
			printJSON(sensors)
			return nil
		}

		// Human-readable output
		fmt.Printf("Found %d light level sensors:\n\n", len(sensors))
		for _, sensor := range sensors {
			fmt.Printf("☀️  Light Sensor %s\n", sensor.ID)
			fmt.Printf("   Light Level: %d lux\n", sensor.LightLevel.LightLevel)
			fmt.Printf("   Valid: %v\n", sensor.LightLevel.LightLevelValid)
			fmt.Printf("   Enabled: %v\n", sensor.Enabled)
			fmt.Println()
		}
		
		return nil
	},
}

func init() {
	// Add subcommands
	sensorsCmd.AddCommand(listMotionSensorsCmd)
	sensorsCmd.AddCommand(listTemperatureSensorsCmd)
	sensorsCmd.AddCommand(listLightSensorsCmd)
	
	// Add to root
	rootCmd.AddCommand(sensorsCmd)
}