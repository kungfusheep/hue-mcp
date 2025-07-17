package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// lightsCmd represents the lights command group
var lightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "Control individual lights",
	Long:  `Commands for controlling individual Philips Hue lights.`,
}

// listLightsCmd lists all available lights
var listLightsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available lights",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		lights, err := hueClient.GetLights(ctx)
		if err != nil {
			return fmt.Errorf("failed to get lights: %w", err)
		}

		if jsonOutput {
			printJSON(lights)
			return nil
		}

		// Human-readable output
		fmt.Printf("Found %d lights:\n\n", len(lights))
		for _, light := range lights {
			status := "off"
			if light.On.On {
				status = fmt.Sprintf("on (brightness: %.0f%%)", light.Dimming.Brightness)
			}
			fmt.Printf("%-30s %s\n", light.Metadata.Name, status)
			fmt.Printf("  ID: %s\n", light.ID)
			fmt.Printf("  Type: %s\n", light.Metadata.Archetype)
			if light.Color != nil {
				fmt.Printf("  Color: X=%.3f Y=%.3f\n", light.Color.XY.X, light.Color.XY.Y)
			}
			fmt.Println()
		}
		return nil
	},
}

// lightOnCmd turns a light on
var lightOnCmd = &cobra.Command{
	Use:   "on <light-name-or-id>",
	Short: "Turn a light on",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		
		// Resolve light name to ID
		lightID, err := resolveLightID(ctx, args[0])
		if err != nil {
			return err
		}
		
		err = hueClient.TurnOnLight(ctx, lightID)
		if err != nil {
			return fmt.Errorf("failed to turn on light: %w", err)
		}
		
		printMessage("Light %s turned on", args[0])
		return nil
	},
}

// lightOffCmd turns a light off
var lightOffCmd = &cobra.Command{
	Use:   "off <light-name-or-id>",
	Short: "Turn a light off",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		
		// Resolve light name to ID
		lightID, err := resolveLightID(ctx, args[0])
		if err != nil {
			return err
		}
		
		err = hueClient.TurnOffLight(ctx, lightID)
		if err != nil {
			return fmt.Errorf("failed to turn off light: %w", err)
		}
		
		printMessage("Light %s turned off", args[0])
		return nil
	},
}

// lightColorCmd sets light color
var lightColorCmd = &cobra.Command{
	Use:   "color <light-name-or-id> <color>",
	Short: "Set light color (hex or name)",
	Long:  `Set light color using hex code (#FF0000) or color name (red, blue, green, etc.)`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		color := args[1]
		ctx := context.Background()
		
		// Resolve light name to ID
		lightID, err := resolveLightID(ctx, args[0])
		if err != nil {
			return err
		}
		
		// Convert color name to hex if needed
		hexColor := namedColorToHex(color)
		if hexColor == "" {
			hexColor = color
		}
		
		err = hueClient.SetLightColor(ctx, lightID, hexColor)
		if err != nil {
			return fmt.Errorf("failed to set color: %w", err)
		}
		
		printMessage("Light %s color set to %s", args[0], color)
		return nil
	},
}

// lightBrightnessCmd sets light brightness
var lightBrightnessCmd = &cobra.Command{
	Use:   "brightness <light-name-or-id> <percent>",
	Short: "Set light brightness (0-100)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		brightness, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return fmt.Errorf("invalid brightness value: %w", err)
		}
		
		if brightness < 0 || brightness > 100 {
			return fmt.Errorf("brightness must be between 0 and 100")
		}
		
		ctx := context.Background()
		
		// Resolve light name to ID
		lightID, err := resolveLightID(ctx, args[0])
		if err != nil {
			return err
		}
		
		err = hueClient.SetLightBrightness(ctx, lightID, brightness)
		if err != nil {
			return fmt.Errorf("failed to set brightness: %w", err)
		}
		
		printMessage("Light %s brightness set to %.0f%%", args[0], brightness)
		return nil
	},
}

// lightStateCmd shows current state of a light
var lightStateCmd = &cobra.Command{
	Use:   "state <light-name-or-id>",
	Short: "Show current state of a light",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		
		// Resolve light name to ID
		lightID, err := resolveLightID(ctx, args[0])
		if err != nil {
			return err
		}
		
		light, err := hueClient.GetLight(ctx, lightID)
		if err != nil {
			return fmt.Errorf("failed to get light: %w", err)
		}
		
		if jsonOutput {
			printJSON(light)
			return nil
		}
		
		// Human-readable output
		fmt.Printf("Light: %s\n", light.Metadata.Name)
		fmt.Printf("Type: %s\n", light.Metadata.Archetype)
		fmt.Printf("On: %v\n", light.On.On)
		if light.On.On {
			fmt.Printf("Brightness: %.0f%%\n", light.Dimming.Brightness)
		}
		if light.Color != nil {
			fmt.Printf("Color XY: (%.3f, %.3f)\n", light.Color.XY.X, light.Color.XY.Y)
		}
		if light.ColorTemperature != nil && light.ColorTemperature.MirekValid {
			fmt.Printf("Color Temperature: %d mirek\n", light.ColorTemperature.Mirek)
		}
		if light.Effects != nil && light.Effects.Effect != "" {
			fmt.Printf("Effect: %s\n", light.Effects.Effect)
		}
		
		return nil
	},
}

func init() {
	// Add subcommands
	lightsCmd.AddCommand(listLightsCmd)
	lightsCmd.AddCommand(lightOnCmd)
	lightsCmd.AddCommand(lightOffCmd)
	lightsCmd.AddCommand(lightColorCmd)
	lightsCmd.AddCommand(lightBrightnessCmd)
	lightsCmd.AddCommand(lightStateCmd)
	
	// Add to root
	rootCmd.AddCommand(lightsCmd)
}