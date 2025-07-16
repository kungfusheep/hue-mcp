# Light Discovery in the Hue MCP Codebase

## Overview

This document explains how lights are discovered and categorized in the Hue MCP codebase, with a focus on identifying office lights and Playbar devices.

## Key Findings from Code Analysis

### 1. Office Light Names
Based on the codebase search, the following office lights are consistently referenced:
- **Office 1** - Overhead spot light
- **Office 2** - Overhead spot light  
- **Office 3** - Overhead spot light
- **Office 4** - Overhead spot light
- **Petes Office Lamp** - Table/desk lamp
- **Hue Play 1** - Behind speaker (Play bar)
- **Hue Play 2** - Behind speaker (Play bar)

### 2. Light Discovery API

The Hue v2 API provides light information through:
```go
// GetLights returns all lights
func (c *Client) GetLights(ctx context.Context) ([]Light, error)
```

Each light has the following structure:
```go
type Light struct {
    ID       string    // Unique identifier
    Metadata Metadata  // Contains name and archetype
    On       OnState   // On/off state
    Dimming  Dimming   // Brightness info
    // ... other fields
}

type Metadata struct {
    Name      string // Human-readable name
    Archetype string // Device type (e.g., "hue_play", "sultan_bulb")
}
```

### 3. Office Group

The office is also defined as a room/group:
- Room name: **"Office"**
- Contains a grouped_light service for controlling all office lights together
- Used extensively in group control tests

### 4. Playbar/Play Device Detection

While "Playbar" isn't explicitly mentioned in most code, the Hue Play devices are identified by:
1. Names containing "Play" (e.g., "Hue Play 1", "Hue Play 2")
2. Archetype field potentially being "hue_play"
3. Being located in the office (behind speakers according to HAND-OVER.md)

## Light Discovery Script

Two scripts were created:

### 1. `list_lights.go` - Full Implementation
A complete light discovery script that:
- Connects to the Hue bridge via HTTPS
- Retrieves all lights using the v2 API
- Categorizes lights as office lights, Playbar/Play devices, or other
- Shows room/group associations
- Displays current state and capabilities

**Usage:**
```bash
export HUE_BRIDGE_IP="your-bridge-ip"
export HUE_USERNAME="your-hue-username"
go run list_lights.go
```

### 2. `discover_lights_demo.go` - Demonstration
A demo script that shows the categorization logic without requiring API credentials.

## Key Implementation Details

### Authentication
The Hue v2 API requires:
- Bridge IP address (found in test_server.sh: 192.168.87.51)
- Username/API key (must be obtained via bridge button press)
- HTTPS with self-signed certificates (requires `InsecureSkipVerify`)

### Light Categorization Logic

```go
// Office light detection
isOfficeLight := false
if strings.Contains(strings.ToLower(light.Name), "office") {
    isOfficeLight = true
}

// Playbar/Play detection  
isPlaybar := false
if strings.Contains(strings.ToLower(light.Name), "play") || 
   strings.Contains(strings.ToLower(light.Name), "playbar") ||
   light.Archetype == "hue_play" {
    isPlaybar = true
}
```

### Common Archetypes
From the codebase, common light archetypes include:
- `sultan_bulb` - Standard bulbs (Office 1-4)
- `table_shade` - Table lamps
- `hue_play` - Play bars
- `light_strip` - LED strips
- `pendant_round` - Ceiling pendants

## Testing and Verification

The codebase includes extensive tests that demonstrate:
1. Individual light control (test2_onoff.go, test3_brightness.go)
2. Group control for the office (test6_group.go)
3. Scene creation with office lights (test7_scenes.go)
4. Batch operations on multiple office lights (batch_advanced_demo.go)

All tests consistently reference the same set of office lights, confirming the naming convention.

## Recommendations

1. **Playbar Identification**: Since "Playbar" might be a colloquial term, look for:
   - Devices with "Play" in the name
   - Devices with `hue_play` archetype
   - Devices positioned behind speakers/TV

2. **Alternative Names**: Some Playbar devices might have custom names like:
   - "TV Light"
   - "Speaker Light"
   - "Entertainment Light"
   
3. **Room Association**: Check which lights are associated with the Office room to ensure all office lights are captured.

4. **Dynamic Discovery**: The actual light names and IDs will vary by installation, so always use the API to discover current devices rather than hardcoding.