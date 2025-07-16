#!/bin/bash

# Claude Desktop configuration installer for Hue MCP

CONFIG_DIR="$HOME/Library/Application Support/Claude"
CONFIG_FILE="$CONFIG_DIR/claude_desktop_config.json"

echo "🔧 Installing Hue MCP for Claude Desktop..."

# Create config directory if it doesn't exist
mkdir -p "$CONFIG_DIR"

# Check if config file exists
if [ -f "$CONFIG_FILE" ]; then
    echo "⚠️  Existing configuration found at: $CONFIG_FILE"
    echo "📋 Current content:"
    cat "$CONFIG_FILE"
    echo ""
    echo "⚠️  Please manually merge the following configuration:"
else
    echo "✅ Creating new configuration file..."
    cp claude_desktop_config.json "$CONFIG_FILE"
    echo "✅ Configuration installed!"
fi

echo ""
echo "📝 Hue MCP Configuration:"
echo "========================"
cat claude_desktop_config.json
echo ""
echo "========================"

echo ""
echo "📌 Next steps:"
echo "1. If you had existing configuration, manually merge the above into: $CONFIG_FILE"
echo "2. Restart Claude Desktop"
echo "3. Try commands like:"
echo "   - 'Turn on the office lights'"
echo "   - 'Set living room to candle effect'"
echo "   - 'List all lights'"
echo ""
echo "🎉 Installation complete!"