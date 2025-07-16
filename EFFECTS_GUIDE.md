# Hue MCP Effects & Sequences Guide

This guide showcases the lighting effects and sequences available in the Hue MCP server. All effects run asynchronously, meaning you get control back immediately while your lights perform the effects in the background.

## Quick Start Examples

### 🚨 Alert/Notification
```
"I need to get someone's attention in the office"
→ Uses alert_effect with red flashes
```

### 🌈 Party Mode
```
"Let's create a party atmosphere with the lights"
→ Combines strobe_effect and color_loop for dynamic lighting
```

### 🌅 Wake-Up Light
```
"Simulate a sunrise in my bedroom over 10 minutes"
→ Uses custom_sequence with gradual color and brightness changes
```

### 💗 Heartbeat
```
"Make the lamp pulse like a heartbeat"
→ Uses pulse_effect with specific timing
```

## Available Effects

### 1. Flash Effect (`flash_effect`)
Quick on/off flashes - perfect for notifications or emphasis.

**Use cases:**
- Doorbell or phone notifications
- Timer completion alerts
- Visual alarms
- Party lighting accents

**Example:** "Flash the office lights red 3 times when my timer goes off"

### 2. Pulse Effect (`pulse_effect`)
Smooth brightness fading up and down - creates a breathing effect.

**Use cases:**
- Meditation or relaxation lighting
- Sleep aid (slow pulsing)
- Subtle ambient effects
- Status indicators (slow pulse = standby)

**Example:** "Make the bedroom lamp pulse gently between 10% and 40% brightness"

### 3. Color Loop (`color_loop`)
Continuous cycling through colors - runs until stopped.

**Use cases:**
- Rainbow effects for kids' rooms
- Team colors for game day
- Seasonal themes (red/green for Christmas)
- Mood lighting that changes

**Example:** "Cycle the playroom lights through rainbow colors"

### 4. Strobe Effect (`strobe_effect`)
Rapid flashing for dramatic effect. ⚠️ Use responsibly!

**Use cases:**
- Dance parties
- Halloween effects
- Emergency simulation
- High-energy moments

**Example:** "Create a disco strobe in white for 10 seconds"

### 5. Alert Effect (`alert_effect`)
Pre-programmed attention-getting pattern.

**Use cases:**
- Incoming call notifications
- Security alerts
- Kitchen timers
- Meeting reminders

**Example:** "Alert me with the desk lamp when someone's at the door"

### 6. Custom Sequences (`custom_sequence`)
Build any complex lighting choreography you can imagine!

**Use cases:**
- Sunrise/sunset simulation
- Scene transitions
- Synchronized multi-room effects
- Storytelling with lights
- Complex notifications

## Combining Effects

You can run multiple effects simultaneously on different lights:

```
1. Start rainbow loop on Play bars
2. Add gentle pulse on desk lamp
3. When timer ends, flash overhead lights
4. All effects run independently!
```

## Advanced Custom Sequence Examples

### Sunrise Simulation
Gradually transition from darkness → deep red → orange → yellow → bright white

### Police Light
Alternate red and blue rapidly between two lights

### Thunderstorm
Random white flashes of varying intensity

### Romantic Dinner
Slow transitions between warm candlelight colors

### Sports Team Celebration
Flash team colors in sequence across all lights

## Tips for Best Results

1. **Start Simple**: Try individual effects before combining them
2. **Consider Timing**: Slower effects are more relaxing, faster ones energizing
3. **Use Groups**: Apply effects to room groups for coordinated lighting
4. **Mix Effects**: Different effects on different lights create depth
5. **Save Favorites**: Note down custom sequences you love

## Common Commands

- **See what's running**: "List all active light sequences"
- **Stop an effect**: "Stop sequence [ID]"
- **Quick test**: "Flash my office light blue twice"
- **Ambient mode**: "Start a slow color loop on the living room lights"

Remember: All effects are non-blocking by default, so you can stack them up and create complex lighting scenes while continuing to work!