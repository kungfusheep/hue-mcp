package cmd

import "strings"

// namedColorToHex converts common color names to hex values
func namedColorToHex(color string) string {
	colors := map[string]string{
		"red":     "#FF0000",
		"green":   "#00FF00",
		"blue":    "#0000FF",
		"yellow":  "#FFFF00",
		"cyan":    "#00FFFF",
		"magenta": "#FF00FF",
		"white":   "#FFFFFF",
		"warm":    "#FFA500",
		"cool":    "#ADD8E6",
		"orange":  "#FFA500",
		"purple":  "#800080",
		"pink":    "#FFC0CB",
	}
	
	return colors[strings.ToLower(color)]
}