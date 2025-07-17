package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/kungfusheep/hue/mcp"
)

// Effect flags
var (
	effectColor    string
	flashCount     int
	flashDuration  int
	minBrightness  float64
	maxBrightness  float64
	pulseDuration  int
	pulseCount     int
	transitionTime int
	strobeRate     int
	duration       int
)

// effectsCmd represents the effects command group
var effectsCmd = &cobra.Command{
	Use:   "effects",
	Short: "Create lighting effects",
	Long:  `Commands for creating dynamic lighting effects like flash, pulse, and strobe.`,
}

// flashCmd creates a flash effect
var flashCmd = &cobra.Command{
	Use:   "flash <light-name-or-id>",
	Short: "Create a flashing effect",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		
		// Resolve light name to ID
		targetID, err := resolveLightID(ctx, args[0])
		if err != nil {
			return err
		}
		
		// For CLI, we need to run the effect synchronously
		// Execute commands directly instead of using the scheduler
		
		// Get current state to restore later
		light, err := hueClient.GetLight(ctx, targetID)
		if err != nil {
			return fmt.Errorf("failed to get light state: %w", err)
		}
		originalOn := light.On.On
		
		// Ensure light is on first
		if !originalOn {
			err = hueClient.TurnOnLight(ctx, targetID)
			if err != nil {
				return fmt.Errorf("failed to turn on light: %w", err)
			}
		}
		
		for i := 0; i < flashCount; i++ {
			// Flash on with color at full brightness
			err = hueClient.SetLightColor(ctx, targetID, effectColor)
			if err != nil {
				return fmt.Errorf("failed to set flash color: %w", err)
			}
			err = hueClient.SetLightBrightness(ctx, targetID, 100)
			if err != nil {
				return fmt.Errorf("failed to set brightness: %w", err)
			}
			time.Sleep(time.Duration(flashDuration) * time.Millisecond)
			
			// Flash off
			err = hueClient.TurnOffLight(ctx, targetID)
			if err != nil {
				return fmt.Errorf("failed to turn off light: %w", err)
			}
			time.Sleep(time.Duration(flashDuration) * time.Millisecond)
			
			// Turn back on for next flash (except last iteration)
			if i < flashCount-1 {
				err = hueClient.TurnOnLight(ctx, targetID)
				if err != nil {
					return fmt.Errorf("failed to turn light back on: %w", err)
				}
			}
		}
		
		// Restore original state
		if originalOn {
			err = hueClient.TurnOnLight(ctx, targetID)
			if err != nil {
				return fmt.Errorf("failed to restore light state: %w", err)
			}
			// Restore original brightness
			err = hueClient.SetLightBrightness(ctx, targetID, light.Dimming.Brightness)
			if err != nil {
				return fmt.Errorf("failed to restore brightness: %w", err)
			}
		}
		
		printMessage("Flash effect completed on %s", args[0])
		printMessage("Color: %s | Flashes: %d", effectColor, flashCount)
		
		return nil
	},
}

// pulseCmd creates a pulse effect
var pulseCmd = &cobra.Command{
	Use:   "pulse <light-name-or-id>",
	Short: "Create a breathing/pulse effect",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		
		// Resolve light name to ID
		targetID, err := resolveLightID(ctx, args[0])
		if err != nil {
			return err
		}
		
		// For CLI, run the pulse effect synchronously
		
		// Get current state to restore later
		light, err := hueClient.GetLight(ctx, targetID)
		if err != nil {
			return fmt.Errorf("failed to get light state: %w", err)
		}
		
		printMessage("Pulse effect started on %s", args[0])
		printMessage("Brightness: %.0f%% - %.0f%% | Pulses: %d", minBrightness, maxBrightness, pulseCount)
		
		// Make sure light is on
		if !light.On.On {
			err = hueClient.TurnOnLight(ctx, targetID)
			if err != nil {
				return fmt.Errorf("failed to turn on light: %w", err)
			}
		}
		
		// Execute pulse cycles
		halfDuration := time.Duration(pulseDuration/2) * time.Millisecond
		for i := 0; i < pulseCount; i++ {
			// Fade down to min
			err = hueClient.SetLightBrightness(ctx, targetID, minBrightness)
			if err != nil {
				return fmt.Errorf("failed to set min brightness: %w", err)
			}
			time.Sleep(halfDuration)
			
			// Fade up to max
			err = hueClient.SetLightBrightness(ctx, targetID, maxBrightness)
			if err != nil {
				return fmt.Errorf("failed to set max brightness: %w", err)
			}
			time.Sleep(halfDuration)
		}
		
		// Restore original brightness
		if light.On.On {
			err = hueClient.SetLightBrightness(ctx, targetID, light.Dimming.Brightness)
			if err != nil {
				return fmt.Errorf("failed to restore brightness: %w", err)
			}
		}
		
		printMessage("Pulse effect completed")
		
		return nil
	},
}

// strobeCmd creates a strobe effect
var strobeCmd = &cobra.Command{
	Use:   "strobe <light-name-or-id>",
	Short: "Create a strobe effect (use responsibly!)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		
		// Resolve light name to ID
		targetID, err := resolveLightID(ctx, args[0])
		if err != nil {
			return err
		}
		
		// For CLI, run strobe effect synchronously
		
		printMessage("⚠️  Strobe effect started on %s", args[0])
		printMessage("Color: %s | Rate: %dms | Duration: %dms", effectColor, strobeRate, duration)
		
		// Calculate iterations
		iterations := duration / (strobeRate * 2)
		
		for i := 0; i < iterations; i++ {
			// Strobe on
			err := hueClient.SetLightColor(ctx, targetID, effectColor)
			if err != nil {
				return fmt.Errorf("failed to set strobe color: %w", err)
			}
			err = hueClient.TurnOnLight(ctx, targetID)
			if err != nil {
				return fmt.Errorf("failed to turn on light: %w", err)
			}
			time.Sleep(time.Duration(strobeRate) * time.Millisecond)
			
			// Strobe off
			err = hueClient.TurnOffLight(ctx, targetID)
			if err != nil {
				return fmt.Errorf("failed to turn off light: %w", err)
			}
			time.Sleep(time.Duration(strobeRate) * time.Millisecond)
		}
		
		printMessage("Strobe effect completed")
		
		return nil
	},
}

// stopCmd stops a running effect
var stopCmd = &cobra.Command{
	Use:   "stop <sequence-id>",
	Short: "Stop a running effect",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sequenceID := args[0]
		
		err := mcp.GetScheduler().StopSequence(sequenceID)
		if err != nil {
			return fmt.Errorf("failed to stop sequence: %w", err)
		}
		
		printMessage("Sequence %s stopped", sequenceID)
		return nil
	},
}

// listSequencesCmd lists all running sequences
var listSequencesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running effects",
	RunE: func(cmd *cobra.Command, args []string) error {
		sequences := mcp.GetScheduler().GetSequences()
		
		if jsonOutput {
			printJSON(sequences)
			return nil
		}
		
		if len(sequences) == 0 {
			fmt.Println("No active sequences")
			return nil
		}
		
		fmt.Printf("Active sequences (%d):\n\n", len(sequences))
		for id, seq := range sequences {
			status := "stopped"
			if seq.Running {
				status = "running"
			}
			fmt.Printf("- %s: %s [%s]\n", id, seq.Name, status)
			fmt.Printf("  Commands: %d | Loop: %v\n", len(seq.Commands), seq.Loop)
		}
		
		return nil
	},
}

func init() {
	// Flash flags
	flashCmd.Flags().StringVar(&effectColor, "color", "#FFFFFF", "Flash color (hex or name)")
	flashCmd.Flags().IntVar(&flashCount, "count", 3, "Number of flashes")
	flashCmd.Flags().IntVar(&flashDuration, "duration", 200, "Flash duration in milliseconds")
	
	// Pulse flags
	pulseCmd.Flags().Float64Var(&minBrightness, "min", 10, "Minimum brightness (0-100)")
	pulseCmd.Flags().Float64Var(&maxBrightness, "max", 100, "Maximum brightness (0-100)")
	pulseCmd.Flags().IntVar(&pulseDuration, "duration", 2000, "Pulse duration in milliseconds")
	pulseCmd.Flags().IntVar(&pulseCount, "count", 5, "Number of pulses")
	
	// Strobe flags
	strobeCmd.Flags().StringVar(&effectColor, "color", "#FFFFFF", "Strobe color (hex or name)")
	strobeCmd.Flags().IntVar(&strobeRate, "rate", 100, "Strobe rate in milliseconds")
	strobeCmd.Flags().IntVar(&duration, "duration", 5000, "Total duration in milliseconds")
	
	// Add subcommands
	effectsCmd.AddCommand(flashCmd)
	effectsCmd.AddCommand(pulseCmd)
	effectsCmd.AddCommand(strobeCmd)
	effectsCmd.AddCommand(stopCmd)
	effectsCmd.AddCommand(listSequencesCmd)
	
	// Add to root
	rootCmd.AddCommand(effectsCmd)
}