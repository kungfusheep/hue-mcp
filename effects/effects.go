package effects

// Effect constants for Philips Hue v2 API
const (
	None       = "no_effect"
	Candle     = "candle"
	Fire       = "fire"
	Prism      = "prism"
	Sparkle    = "sparkle"
	Opal       = "opal"
	Glisten    = "glisten"
	Underwater = "underwater"
	Cosmos     = "cosmos"
	Sunbeam    = "sunbeam"
	Enchant    = "enchant"
)

// IsValid checks if an effect name is valid
func IsValid(effect string) bool {
	switch effect {
	case None, Candle, Fire, Prism, Sparkle, Opal, Glisten, Underwater, Cosmos, Sunbeam, Enchant:
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
		Fire,
		Prism,
		Sparkle,
		Opal,
		Glisten,
		Underwater,
		Cosmos,
		Sunbeam,
		Enchant,
	}
}

// GetDescription returns a human-readable description of an effect
func GetDescription(effect string) string {
	switch effect {
	case None:
		return "No effect"
	case Candle:
		return "Simulates a flickering candle"
	case Fire:
		return "Simulates a cozy fireplace"
	case Prism:
		return "Prism color effects"
	case Sparkle:
		return "Sparkling light effect"
	case Opal:
		return "Opal color shifts"
	case Glisten:
		return "Glistening light effect"
	case Underwater:
		return "Underwater bubble effect"
	case Cosmos:
		return "Cosmic space effect"
	case Sunbeam:
		return "Warm sunbeam effect"
	case Enchant:
		return "Magical enchanted effect"
	default:
		return "Unknown effect"
	}
}