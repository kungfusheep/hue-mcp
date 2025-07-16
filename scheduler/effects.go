package scheduler

import (
	"fmt"
	"time"
)

// EffectType represents the type of effect
type EffectType string

const (
	EffectFlash     EffectType = "flash"
	EffectPulse     EffectType = "pulse"
	EffectColorLoop EffectType = "colorloop"
	EffectStrobe    EffectType = "strobe"
	EffectFade      EffectType = "fade"
	EffectRainbow   EffectType = "rainbow"
	EffectAlert     EffectType = "alert"
)

// Effect represents a lighting effect configuration
type Effect struct {
	Type     EffectType
	Target   string                 // Light or group ID
	Duration time.Duration          // Total duration of the effect
	Params   map[string]interface{} // Effect-specific parameters
}

// CreateFlashEffect creates a flash effect sequence
// Flashes the light with a specified color and returns to previous state
func CreateFlashEffect(targetID string, color string, flashCount int, flashDuration time.Duration) *Sequence {
	commands := []Command{}
	
	for i := 0; i < flashCount; i++ {
		// Flash on with color
		commands = append(commands, Command{
			Type:   "light",
			Action: "color",
			Target: targetID,
			Params: map[string]interface{}{"color": color},
			Delay:  0,
		})
		
		// Brief hold
		commands = append(commands, Command{
			Type:   "light",
			Action: "on",
			Target: targetID,
			Delay:  flashDuration,
		})
		
		// Turn off
		commands = append(commands, Command{
			Type:   "light",
			Action: "off",
			Target: targetID,
			Delay:  flashDuration,
		})
	}
	
	return &Sequence{
		Name:     fmt.Sprintf("Flash %s", targetID),
		Commands: commands,
		Loop:     false,
	}
}

// CreatePulseEffect creates a pulsing brightness effect
func CreatePulseEffect(targetID string, minBrightness, maxBrightness float64, pulseDuration time.Duration, pulseCount int) *Sequence {
	commands := []Command{}
	stepDuration := pulseDuration / 10 // 10 steps per pulse
	
	for i := 0; i < pulseCount; i++ {
		// Fade up
		for j := 0; j < 5; j++ {
			brightness := minBrightness + (maxBrightness-minBrightness)*float64(j)/5.0
			commands = append(commands, Command{
				Type:   "light",
				Action: "brightness",
				Target: targetID,
				Params: map[string]interface{}{"brightness": brightness},
				Delay:  stepDuration,
			})
		}
		
		// Fade down
		for j := 5; j > 0; j-- {
			brightness := minBrightness + (maxBrightness-minBrightness)*float64(j)/5.0
			commands = append(commands, Command{
				Type:   "light",
				Action: "brightness",
				Target: targetID,
				Params: map[string]interface{}{"brightness": brightness},
				Delay:  stepDuration,
			})
		}
	}
	
	return &Sequence{
		Name:     fmt.Sprintf("Pulse %s", targetID),
		Commands: commands,
		Loop:     false,
	}
}

// CreateColorLoopEffect creates a smooth color transition effect
func CreateColorLoopEffect(targetID string, colors []string, transitionTime time.Duration) *Sequence {
	commands := []Command{}
	
	for _, color := range colors {
		commands = append(commands, Command{
			Type:   "light",
			Action: "color",
			Target: targetID,
			Params: map[string]interface{}{"color": color},
			Delay:  transitionTime,
		})
	}
	
	return &Sequence{
		Name:     fmt.Sprintf("ColorLoop %s", targetID),
		Commands: commands,
		Loop:     true, // This effect loops by default
	}
}

// CreateStrobeEffect creates a strobe light effect
func CreateStrobeEffect(targetID string, color string, strobeRate time.Duration, duration time.Duration) *Sequence {
	commands := []Command{}
	iterations := int(duration / (strobeRate * 2))
	
	// Set color first
	commands = append(commands, Command{
		Type:   "light",
		Action: "color",
		Target: targetID,
		Params: map[string]interface{}{"color": color},
		Delay:  0,
	})
	
	for i := 0; i < iterations; i++ {
		// Turn on
		commands = append(commands, Command{
			Type:   "light",
			Action: "on",
			Target: targetID,
			Delay:  strobeRate,
		})
		
		// Turn off
		commands = append(commands, Command{
			Type:   "light",
			Action: "off",
			Target: targetID,
			Delay:  strobeRate,
		})
	}
	
	return &Sequence{
		Name:     fmt.Sprintf("Strobe %s", targetID),
		Commands: commands,
		Loop:     false,
	}
}

// CreateRainbowEffect creates a rainbow color cycle
func CreateRainbowEffect(targetID string, stepDuration time.Duration) *Sequence {
	// Rainbow colors in order
	colors := []string{
		"#FF0000", // Red
		"#FF7F00", // Orange
		"#FFFF00", // Yellow
		"#00FF00", // Green
		"#0000FF", // Blue
		"#4B0082", // Indigo
		"#9400D3", // Violet
	}
	
	return CreateColorLoopEffect(targetID, colors, stepDuration)
}

// CreateAlertEffect creates an attention-grabbing alert effect
func CreateAlertEffect(targetID string, alertColor string, normalColor string) *Sequence {
	commands := []Command{
		// Quick flashes
		{Type: "light", Action: "color", Target: targetID, Params: map[string]interface{}{"color": alertColor}, Delay: 0},
		{Type: "light", Action: "brightness", Target: targetID, Params: map[string]interface{}{"brightness": 100}, Delay: 100 * time.Millisecond},
		{Type: "light", Action: "brightness", Target: targetID, Params: map[string]interface{}{"brightness": 20}, Delay: 100 * time.Millisecond},
		{Type: "light", Action: "brightness", Target: targetID, Params: map[string]interface{}{"brightness": 100}, Delay: 100 * time.Millisecond},
		{Type: "light", Action: "brightness", Target: targetID, Params: map[string]interface{}{"brightness": 20}, Delay: 100 * time.Millisecond},
		{Type: "light", Action: "brightness", Target: targetID, Params: map[string]interface{}{"brightness": 100}, Delay: 100 * time.Millisecond},
		
		// Return to normal
		{Type: "light", Action: "color", Target: targetID, Params: map[string]interface{}{"color": normalColor}, Delay: 500 * time.Millisecond},
		{Type: "light", Action: "brightness", Target: targetID, Params: map[string]interface{}{"brightness": 50}, Delay: 0},
	}
	
	return &Sequence{
		Name:     fmt.Sprintf("Alert %s", targetID),
		Commands: commands,
		Loop:     false,
	}
}

// CreateFadeEffect creates a smooth fade between two states
func CreateFadeEffect(targetID string, startColor string, endColor string, startBrightness, endBrightness float64, duration time.Duration, steps int) *Sequence {
	commands := []Command{}
	stepDuration := duration / time.Duration(steps)
	
	// Set initial state
	commands = append(commands, Command{
		Type:   "light",
		Action: "color",
		Target: targetID,
		Params: map[string]interface{}{"color": startColor},
		Delay:  0,
	})
	commands = append(commands, Command{
		Type:   "light",
		Action: "brightness",
		Target: targetID,
		Params: map[string]interface{}{"brightness": startBrightness},
		Delay:  0,
	})
	
	// For simplicity, we'll just transition brightness and then change color
	// A more sophisticated implementation would interpolate colors
	for i := 1; i <= steps; i++ {
		progress := float64(i) / float64(steps)
		brightness := startBrightness + (endBrightness-startBrightness)*progress
		
		commands = append(commands, Command{
			Type:   "light",
			Action: "brightness",
			Target: targetID,
			Params: map[string]interface{}{"brightness": brightness},
			Delay:  stepDuration,
		})
	}
	
	// Set final color
	commands = append(commands, Command{
		Type:   "light",
		Action: "color",
		Target: targetID,
		Params: map[string]interface{}{"color": endColor},
		Delay:  0,
	})
	
	return &Sequence{
		Name:     fmt.Sprintf("Fade %s", targetID),
		Commands: commands,
		Loop:     false,
	}
}

// CreateGroupEffect applies an effect to all lights in a group
func CreateGroupEffect(effect *Sequence, groupID string) *Sequence {
	// Convert all light commands to group commands
	groupCommands := make([]Command, len(effect.Commands))
	for i, cmd := range effect.Commands {
		groupCmd := cmd
		if cmd.Type == "light" {
			groupCmd.Type = "group"
			groupCmd.Target = groupID
		}
		groupCommands[i] = groupCmd
	}
	
	return &Sequence{
		Name:     fmt.Sprintf("Group %s - %s", groupID, effect.Name),
		Commands: groupCommands,
		Loop:     effect.Loop,
	}
}