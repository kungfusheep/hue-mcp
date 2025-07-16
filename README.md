# Philips Hue v2 MCP Server

A Model Context Protocol (MCP) server for Philips Hue v2 API, enabling native lighting effects and comprehensive control for your AI agents.

## Features

### Core Lighting Control
- ✅ **Native v2 Effects**: Candle, fire, sparkle, cosmos, prism, opal, glisten, and more!
- ✅ **Comprehensive Light Control**: On/off, brightness, color, effects
- ✅ **Group Management**: Control entire rooms at once
- ✅ **Scene Support**: Activate and manage scenes
- ✅ **Device Discovery**: Automatic detection of all Hue devices

### Advanced Features (New!)
- 🚀 **Non-blocking Operations**: All commands execute asynchronously by default
- 🎭 **Pre-built Effects**: Flash, pulse, color loop, strobe, and alert patterns
- 🎨 **Custom Sequences**: Build complex lighting choreography with precise timing
- 🔄 **Real-time Event Streaming**: Subscribe to motion, button, and light state changes
- 📦 **Batch Commands**: Execute multiple commands with timing control
- 🎮 **Entertainment Areas**: Support for gaming and media sync (DTLS foundation ready)
- 🔍 **CRUD Operations**: Full create, read, update, delete for all resources

## Prerequisites

1. Go 1.21 or later
2. Philips Hue Bridge with v2 API support
3. Hue Bridge API username (see setup below)

## Setup

### 1. Get Your Hue Bridge IP and Username

Find your bridge IP:
```bash
curl https://discovery.meethue.com/
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

### Basic Control
- "Turn on all office lights"
- "Set the living room to candle effect"
- "Dim bedroom lights to 20%"
- "Make the kitchen lights blue"

### Effects & Sequences
- "Flash the office lights red when my timer goes off"
- "Make the lamp pulse like a heartbeat"
- "Start a rainbow color loop on the kids' room lights"
- "Create a sunrise simulation in the bedroom"
- "Alert me with the desk lamp" (rapid attention-getting flashes)

### Advanced Control
- "Create a custom sequence that fades from red to blue over 10 seconds"
- "Run a party mode with strobe effects"
- "Show me all running light effects"
- "Stop all light animations"

### Sensors & Automation
- "Subscribe to motion sensor events"
- "List all temperature sensors"
- "Show me when someone presses the Hue button"

## Available Tools

### Basic Light Control
- `list_lights` - Discover all available lights
- `light_on/off` - Control individual lights
- `light_brightness` - Set brightness (0-100%)
- `light_color` - Set color (hex or name)
- `light_effect` - Apply native effects (candle, fire, sparkle, etc.)
- `identify_light` - Make a light breathe for identification

### Group & Room Control
- `list_groups` - Discover all groups/rooms
- `group_on/off` - Control entire groups
- `group_brightness` - Set group brightness
- `group_color` - Set group color
- `group_effect` - Apply effects to groups
- `list_rooms` - Discover all rooms with devices

### Scenes & Automation
- `list_scenes` - List available scenes
- `activate_scene` - Activate a scene
- `batch_commands` - Execute multiple commands with timing (async by default!)

### Pre-built Effects 🎭
- `flash_effect` - Attention-getting flashes (notifications, alerts)
- `pulse_effect` - Smooth breathing effect (meditation, ambiance)
- `color_loop` - Continuous color cycling (parties, mood lighting)
- `strobe_effect` - Rapid disco strobe (⚠️ use responsibly!)
- `alert_effect` - Pre-programmed alert pattern

### Advanced Sequencing 🎨
- `custom_sequence` - Build complex multi-step lighting choreography
- `list_sequences` - View all running effects
- `stop_sequence` - Stop a running effect

### Sensors & Events
- `list_motion_sensors` - Get motion sensor states
- `list_temperature_sensors` - Get temperature readings
- `start_event_stream` - Subscribe to real-time events
- `stop_event_stream` - Stop event subscription

### Entertainment & CRUD
- `list_entertainment` - View entertainment areas
- `create_resource` - Create new resources (lights, groups, etc.)
- `update_resource` - Modify existing resources
- `delete_resource` - Remove resources

## Key Features Explained

### 🚀 Non-blocking Operations
All lighting commands execute asynchronously by default. This means:
- Claude responds immediately while lights change in the background
- You can stack multiple effects on different lights
- Complex sequences won't freeze the conversation
- Use `async: false` in batch commands if you need to wait

### 🎭 Effects System
The MCP includes a powerful effects engine:
- **Pre-built effects** for common scenarios (alerts, ambiance, parties)
- **Custom sequences** for precise choreography
- **Parallel execution** - run multiple effects simultaneously
- **Loop support** - effects can repeat indefinitely
- See [EFFECTS_GUIDE.md](EFFECTS_GUIDE.md) for detailed examples

### 🔄 Real-time Events
Subscribe to live updates from your Hue system:
- Motion sensor triggers
- Button presses
- Light state changes
- Temperature updates

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

## Development Status

This MCP server provides comprehensive coverage of the Philips Hue v2 API (90%+):
- ✅ Complete light, group, scene, and room control
- ✅ Full sensor integration
- ✅ Real-time event streaming
- ✅ Advanced effects and sequencing
- ✅ Non-blocking asynchronous operations
- 🚧 Entertainment streaming (DTLS foundation implemented, full streaming in progress)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Apache 2.0 Licence
