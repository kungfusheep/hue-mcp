# Philips Hue v2 MCP Server

A Model Context Protocol (MCP) server for Philips Hue v2 API, enabling native lighting effects and comprehensive control for your AI agents.

## Features

- ✅ **Native v2 Effects**: Candle, fire, sparkle, cosmos, and more!
- ✅ **Comprehensive Light Control**: On/off, brightness, color, effects
- ✅ **Group Management**: Control entire rooms at once
- ✅ **Scene Support**: Activate and manage scenes
- ✅ **Sensor Integration**: Motion, temperature, light level, buttons
- ✅ **Device Discovery**: Automatic detection of all Hue devices
- ✅ **Batch Commands**: Efficient multi-command execution
- ✅ **Real-time Identification**: Make lights blink for identification

## Prerequisites

1. Go 1.21 or later
2. Philips Hue Bridge with v2 API support
3. Hue Bridge API username (see setup below)

## Setup

### 1. Get Your Hue Bridge IP and Username

Find your bridge IP:
```bash
# On macOS/Linux
arp -a | grep -i philips

# Or visit https://discovery.meethue.com/
```

Get an API username:
```bash
# Press the link button on your Hue Bridge, then run:
curl -X POST http://<BRIDGE_IP>/api -H "Content-Type: application/json" -d '{"devicetype":"hue_mcp#claude"}'
```

### 2. Build the MCP Server

```bash
# Clone the repository
git clone https://github.com/kungfusheep/hue-mcp.git
cd hue-mcp

# Build the binary
go build -o hue-mcp
```

### 3. Set Environment Variables

```bash
export HUE_BRIDGE_IP="192.168.1.100"  # Your bridge IP
export HUE_USERNAME="your-api-username-here"
```

### 4. Configure Claude Desktop (example)

Add to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "hue": {
      "command": "/absolute/path/to/hue-mcp",
      "env": {
        "HUE_BRIDGE_IP": "YOUR_BRIDGE_IP",
        "HUE_USERNAME": "YOUR_API_USERNAME"
      }
    }
  }
}
```

### 5. Restart Claude Desktop

Quit and restart Claude Desktop to load the new configuration.

## Usage Examples

Once configured, you can ask Claude to:

- "Turn on all office lights"
- "Set the living room to candle effect"
- "Dim bedroom lights to 20%"
- "Make the kitchen lights blue"
- "Create a fire effect in the office"
- "List all motion sensors"
- "Show me all available scenes"
- "Identify which light is Office 1"
- "Turn on candle effect on all office lights at once" (uses batch commands)

## Available Tools

- `list_lights` - Discover all lights
- `light_on/off` - Control individual lights
- `light_brightness` - Set brightness (0-100%)
- `light_color` - Set color (hex or name)
- `light_effect` - Apply effects (candle, fire, sparkle, etc.)
- `list_groups` - Discover all groups/rooms
- `group_on/off` - Control entire groups
- `group_brightness` - Set group brightness
- `group_color` - Set group color
- `group_effect` - Apply effects to groups
- `list_scenes` - List available scenes
- `activate_scene` - Activate a scene
- `list_rooms` - Discover all rooms with devices
- `list_motion_sensors` - Get motion sensor states
- `list_temperature_sensors` - Get temperature readings
- `identify_light` - Make a light breathe for identification
- `batch_commands` - Execute multiple commands efficiently

## Troubleshooting

1. **"Failed to connect to Hue bridge"**
   - Verify your bridge IP is correct
   - Ensure your API username is valid
   - Check you're on the same network as the bridge

2. **"Light/group not found"**
   - Use `list_lights` or `list_groups` to see available IDs
   - Light names are case-sensitive

3. **Effects not working**
   - Not all lights support all effects
   - Use dynamic effect discovery to see supported effects

## Development

Run tests:
```bash
go test ./...
```

Run comprehensive test suite:
```bash
# Set environment variables first
go run test_comprehensive.go
```

## License

Apache 2.0 Licence
