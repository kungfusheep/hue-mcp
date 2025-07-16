# Philips Hue v2 MCP Server

A Model Context Protocol (MCP) server that provides comprehensive control over Philips Hue lights using the v2 API. This server unlocks native lighting effects like candle, fireplace, and other dynamic scenes not available in v1 API implementations.

## Features

- **Native Effects**: Candle, Fireplace, Colorloop, Sunrise, Sparkle, and more
- **Full Light Control**: On/off, brightness, color, effects for individual lights
- **Group Management**: Control entire rooms or zones synchronously  
- **Scene Support**: List, activate, and create scenes
- **System Tools**: Bridge info, light discovery, identification
- **90%+ v2 API Coverage**: Comprehensive implementation of Hue v2 endpoints

## Installation

### Prerequisites

- Go 1.21 or later
- Philips Hue Bridge with API access
- Your Hue bridge IP address and username

### Building from Source

```bash
git clone https://github.com/kungfusheep/hue-mcp.git
cd hue-mcp
go build -o hue-mcp .
```

## Configuration

Set the following environment variables:

```bash
export HUE_BRIDGE_IP="192.168.1.100"  # Your bridge IP
export HUE_USERNAME="your-hue-username" # Your API username
```

### Getting a Hue Username

If you don't have a username yet:

1. Press the link button on your Hue bridge
2. Within 30 seconds, run:
```bash
curl -X POST http://<bridge-ip>/api -H "Content-Type: application/json" -d '{"devicetype":"hue_mcp_server"}'
```
3. Use the returned username

## Usage with Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "hue": {
      "command": "/path/to/hue-mcp",
      "env": {
        "HUE_BRIDGE_IP": "192.168.1.100",
        "HUE_USERNAME": "your-username"
      }
    }
  }
}
```

## Available Tools

### Light Control
- `light_on` - Turn a light on
- `light_off` - Turn a light off
- `light_brightness` - Set brightness (0-100%)
- `light_color` - Set color (hex code or name)
- `light_effect` - Apply effects (candle, fireplace, etc.)

### Group Control
- `group_on` - Turn a group on
- `group_off` - Turn a group off
- `group_brightness` - Set group brightness
- `group_color` - Set group color
- `group_effect` - Apply effects to group

### Scene Management
- `list_scenes` - List all available scenes
- `activate_scene` - Activate a scene
- `create_scene` - Create a new scene

### Room & Zone Management
- `list_rooms` - List all rooms with their lights
- `list_zones` - List all zones
- `list_devices` - List all devices with details
- `get_device` - Get detailed device information

### Sensor Tools
- `list_motion_sensors` - List motion sensors and their states
- `list_temperature_sensors` - List temperature sensors with readings
- `list_light_level_sensors` - List light level sensors with lux readings
- `list_buttons` - List buttons (dimmer switches) and last events

### Entertainment
- `list_entertainment` - List entertainment configurations
- `start_entertainment` - Start entertainment mode
- `stop_entertainment` - Stop entertainment mode

### System Tools
- `list_lights` - List all lights with current states
- `list_groups` - List all groups/rooms
- `get_light_state` - Get detailed light state
- `bridge_info` - Get bridge information
- `identify_light` - Make a light blink for identification

## Examples

In Claude Desktop:

```
"Turn on the candle effect in my office"
"Set all lights to warm white at 50% brightness"
"Create a cozy fireplace atmosphere"
"List all my lights and their current states"
```

## Supported Effects

- `no_effect` - Disable effects
- `candle` - Flickering candle effect
- `fireplace` - Cozy fireplace simulation
- `colorloop` - Cycle through colors
- `sunrise` - Sunrise simulation
- `sparkle` - Sparkling effect
- `glisten` - Glistening effect
- `opal` - Opal color shifts
- `prism` - Prism color effects

## Development

### Running Tests

```bash
go test ./...
```

### Project Structure

```
hue-mcp/
├── main.go           # Entry point and MCP server setup
├── hue/
│   ├── hue.go       # Hue v2 API client
│   └── types.go     # API types and models
├── mcp/
│   └── mcp.go       # MCP tool handlers
└── effects/
    └── effects.go   # Effect definitions
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.