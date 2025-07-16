package effects

// Effect constants for Philips Hue v2 API
const (
	None       = "no_effect"
	Candle     = "candle"
	Fireplace  = "fireplace"
	Colorloop  = "colorloop"
	Sunrise    = "sunrise"
	Sparkle    = "sparkle"
	Glisten    = "glisten"
	Opal       = "opal"
	Prism      = "prism"
)

// IsValid checks if an effect name is valid
func IsValid(effect string) bool {
	switch effect {
	case None, Candle, Fireplace, Colorloop, Sunrise, Sparkle, Glisten, Opal, Prism:
		return true
	default:
		return false
	}
}

// GetAllEffects returns all available effects
func GetAllEffects() []string {
	return []string{
		None,
		Candle,
		Fireplace,
		Colorloop,
		Sunrise,
		Sparkle,
		Glisten,
		Opal,
		Prism,
	}
}

// GetDescription returns a human-readable description of an effect
func GetDescription(effect string) string {
	switch effect {
	case None:
		return "No effect"
	case Candle:
		return "Simulates a flickering candle"
	case Fireplace:
		return "Simulates a cozy fireplace"
	case Colorloop:
		return "Cycles through all colors"
	case Sunrise:
		return "Simulates a sunrise"
	case Sparkle:
		return "Sparkling light effect"
	case Glisten:
		return "Glistening light effect"
	case Opal:
		return "Opal color shifts"
	case Prism:
		return "Prism color effects"
	default:
		return "Unknown effect"
	}
}